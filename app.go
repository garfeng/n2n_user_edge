package main

import (
	"bytes"
	"changeme/lib"
	"changeme/model"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"
	"strings"
	"syscall"
)

// App struct
type App struct {
	ctx       context.Context
	cmd       *exec.Cmd
	cancelCmd context.CancelFunc
	
	messageReceiver *MessageReceiver
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{
		messageReceiver: NewMessageReceiver(20),
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}
func (a *App) LoadText(name string) (string, error) {
	buff, err := ioutil.ReadFile(name)
	if err != nil {
		return "", err
	}
	return string(buff), nil
}

func (a *App) SaveText(name string, text string) error {
	return ioutil.WriteFile(name, []byte(text), 0755)
}

func (a *App) PostHttp(url string, data string) (string, error) {
	r := bytes.NewBufferString(data)
	req, err := http.Post(url, "application/json", r)
	if err != nil {
		return "", err
	}
	defer req.Body.Close()
	buff, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return "", err
	}
	return string(buff), nil
}

func (a *App) ReadMessage() *Message {
	return a.messageReceiver.Read()
}

const (
	KAccountFile = "etc/account.json"
	KConfigFile  = "etc/config.json"
)

func (a *App) SetupN2N() error {
	Log.Info("setup N2N")
	if a.cmd != nil {
		Log.Error("edge running")
		return errors.New("edge running")
	}

	user, err := model.LoadJSON[model.User](KAccountFile)
	if err != nil {
		Log.Error(err)
		return err
	}
	param, err := model.LoadJSON[model.SetupEdgeParam](KConfigFile)
	if err != nil {
		Log.Error(err)
		return err
	}

	args := []string{}

	args = append(args, "-I", user.Username)
	args = append(args, "-J", user.Password)

	for k, v := range param.EdgeParams {
		if k != "otherParams" {
			args = append(args, k, v)
		} else {
			otherParams := strings.Split(v, " ")
			for _, op := range otherParams {
				if op != "" {
					args = append(args, op)
				}
			}
		}
	}

	Log.Info("run cmd ./edge ", strings.Join(args, " "))

	var ctx context.Context
	ctx, a.cancelCmd = context.WithCancel(context.Background())
	a.cmd = exec.CommandContext(ctx, "./edge", args...)
	a.cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}

	a.cmd.Stdout = a.messageReceiver.NewSender("stdout")
	a.cmd.Stderr = a.messageReceiver.NewSender("stderr")

	err = a.cmd.Start()
	if err != nil {
		Log.Error(err)
		a.cmd = nil
		a.cancelCmd = nil

		return err
	}

	go a.WaitForN2NFinish()

	return nil
}

func (a *App) WaitForN2NFinish() {
	err := a.cmd.Wait()
	a.cmd = nil
	a.cancelCmd = nil
	if err != nil {
		Log.Error(err)
		fmt.Println("[Error]", err)
	}
}

func (a *App) ShutdownN2N() error {
	if a.cmd == nil {
		return nil
	}
	if a.cancelCmd == nil {
		Log.Error("cmd not nil but cancelFunc is nil")
		return errors.New("cmd not nil but cancelFunc is nil")
	}

	a.cancelCmd()
	return nil
}

func (a *App) Keygen(username, password string) (string, error) {
	return lib.Keygen(username, password)
}

func (a *App) ChangePassword(data model.ChangePasswordParam) error {
	oldHash, err := a.Keygen(data.Username, data.OldPassword)
	if err != nil {
		return err
	}
	newHash, err := a.Keygen(data.Username, data.NewPassword)
	if err != nil {
		return err
	}

	req := &model.ChangePasswordParam{
		Username:    data.Username,
		OldPassword: oldHash,
		NewPassword: newHash,
	}

	cfg, err := model.LoadJSON[model.SetupEdgeParam](KConfigFile)
	if err != nil {
		return err
	}

	buff, _ := json.Marshal(req)

	resp, err := a.PostHttp(cfg.Server+"/auth/changePassword", string(buff))
	if err != nil {
		return err
	}

	r := new(model.ChangePasswordResp)
	err = json.Unmarshal([]byte(resp), r)
	if err != nil {
		return err
	}

	if r.Status != 0 {
		return errors.New(r.Message)
	}

	return nil
}

package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/garfeng/n2n_user_edge/lib"
	"github.com/garfeng/n2n_user_edge/model"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"io/ioutil"
	"net/http"
	"os/exec"
	"strings"
)

// App struct
type App struct {
	ctx       context.Context
	cmd       *exec.Cmd
	cancelCmd context.CancelFunc
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
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

const (
	KAccountFile = "etc/account.json"
	KConfigFile  = "etc/config.json"
)

func (a *App) SetupN2N() error {
	defer runtime.EventsEmit(a.ctx, EventToggleOnLineStatus)

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
	hideCmdWindow(a.cmd)

	a.cmd.Stdout = NewMessageSender(a.ctx, "stdout")
	a.cmd.Stderr = NewMessageSender(a.ctx, "stderr")

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
	defer runtime.EventsEmit(a.ctx, EventToggleOnLineStatus)

	cmd := a.cmd
	err := cmd.Wait()

	a.cmd = nil
	a.cancelCmd = nil
	if err != nil {
		Log.Error(err)
		fmt.Println("[Error]", err)
	}
}

const (
	EventToggleOnLineStatus = "toggle_online"
)

func (a *App) ShutdownN2N() error {
	defer runtime.EventsEmit(a.ctx, EventToggleOnLineStatus)
	if a.cmd == nil {
		return nil
	}
	if a.cancelCmd == nil {
		Log.Error("cmd not nil but cancelFunc is nil")
		return errors.New("cmd not nil but cancelFunc is nil")
	}
	a.cancelCmd()
	a.cmd = nil
	a.cancelCmd = nil
	return nil
}

func (a *App) IsOnline() bool {
	return a.cmd != nil
}

func (a *App) Keygen(username, password string) (string, error) {
	return lib.Keygen(username, password)
}

func (a *App) Title() string {
	const WindowTitle = "N2N User Edge"

	if a.IsOnline() {
		return WindowTitle + " - [💖 Online]"
	}
	return WindowTitle + " - [Offline]"
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

package main

import (
	"bytes"
	"changeme/model"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

// App struct
type App struct {
	ctx context.Context
	cmd *exec.Cmd
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
	if a.cmd != nil {
		return errors.New("edge running")
	}

	user, err := model.LoadJSON[model.User](KAccountFile)
	if err != nil {
		return err
	}
	param, err := model.LoadJSON[model.SetupEdgeParam](KConfigFile)
	if err != nil {
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

	a.cmd = exec.Command("./edge", args...)

	a.cmd.Stdout = os.Stdout
	a.cmd.Stderr = os.Stderr

	err = a.cmd.Start()
	if err != nil {
		return err
	}

	return nil
}

func (a *App) ShutdownN2N() error {
	// TODO: shutdown
	return nil
}

func (a *App) Keygen(username, password string) (string, error) {
	cmd := exec.Command("./n2n-keygen", username, password)
	w := bytes.NewBuffer(nil)
	cmd.Stdout = w
	cmd.Stderr = os.Stderr
	err := cmd.Run()

	if err != nil {
		return "", err
	}

	s := w.String()
	sList := strings.Split(s, username)
	if len(sList) < 2 {
		return "", errors.New("invalid result of keygen")
	}
	return strings.TrimSpace(sList[1]), nil
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

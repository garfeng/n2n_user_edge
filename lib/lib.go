package lib

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

func Keygen(username, password string) (string, error) {
	cmd := exec.Command("./n2n-keygen", username, password)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}

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

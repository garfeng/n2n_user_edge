package lib

import (
	"os/exec"
	"syscall"
)

func hideCmdWindow(c *exec.Cmd) {
	c.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
}

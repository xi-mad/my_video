//go:build windows
// +build windows

package plantform

import (
	"os/exec"
	"syscall"
)

func PrepareBackgroundCommand(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
}

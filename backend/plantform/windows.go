//go:build windows

package plantform

import (
	"os/exec"
	"syscall"
)

func PrepareBackgroundCommand(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
}

func OpenFolder(path string) error {
	cmd := exec.Command("explorer.exe", path)
	PrepareBackgroundCommand(cmd)
	return cmd.Start()
}

func OpenInBrowser(url string) error {
	cmd := exec.Command("cmd", `/c`, `start`, url)
	PrepareBackgroundCommand(cmd)
	return cmd.Start()
}

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
	return exec.Command("explorer.exe", path).Start()
}

func OpenInBrowser(url string) error {
	return exec.Command("cmd", `/c`, `start`, url).Start()
}

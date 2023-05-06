//go:build windows

package plantform

import (
	"log"
	"os/exec"
	"syscall"
)

func PrepareBackgroundCommand(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
}

func OpenFolder(path string) {
	cmd := exec.Command("explorer.exe", path)
	err := cmd.Start()
	if err != nil {
		log.Println(err)
	}
}

func OpenInBrowser(url string) {
	cmd := exec.Command("cmd", `/c`, `start`, url)
	PrepareBackgroundCommand(cmd)
	err := cmd.Start()
	if err != nil {
		log.Println(err)
	}
}

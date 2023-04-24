//go:build macOS

package plantform

import "os/exec"

func PrepareBackgroundCommand(cmd *exec.Cmd) {
}

func OpenFolder(path string) error {
	return exec.Command("open", path).Start()
}

func OpenInBrowser(url string) error {
	return exec.Command("open", url).Start()
}

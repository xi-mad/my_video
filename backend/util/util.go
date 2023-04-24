package util

import (
	"encoding/base64"
	"github.com/xi-mad/my_video/plantform"
	"io"
	"log"
	"os"
	"os/exec"
)

func Image2Base64(path string) (b64 string, err error) {
	if f, err := os.Open(path); err != nil {
		return "", err
	} else {
		defer func() {
			_ = f.Close()
		}()
		if bytes, err := io.ReadAll(f); err != nil {
			return "", err
		} else {
			b64 = base64.StdEncoding.EncodeToString(bytes)
		}
	}
	return
}

func ExecCmd(name string, arg ...string) (err error) {
	cmd := exec.Command(name, arg...)
	if name != "explorer.exe" {
		plantform.PrepareBackgroundCommand(cmd)
	}
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("exec cmd: %s %s, err: %s", name, arg, err)
	} else {
		log.Printf("exec cmd: %s %s, output: %s", name, arg, out)
	}
	return
}

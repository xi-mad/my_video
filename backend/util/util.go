package util

import (
	"encoding/base64"
	"io"
	"os"
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

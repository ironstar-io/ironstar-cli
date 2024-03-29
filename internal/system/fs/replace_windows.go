//go:build windows
// +build windows

package fs

import (
	"fmt"
	iofs "io/fs"
	"os"
)

// Replace ...
func Replace(path string, body []byte, octal iofs.FileMode) {
	var _, err = os.Stat(path)

	if os.IsNotExist(err) {
		TouchByteArray(path, body, octal)
		return
	}

	err = os.WriteFile(path, body, 0)
	if err != nil {
		fmt.Println("There was an issue replacing file contents: ", err)
	}
}

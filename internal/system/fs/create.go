package fs

import (
	"fmt"
	"io/fs"
	"os"
)

// TouchByteArray ...
func TouchByteArray(path string, body []byte, octal fs.FileMode) error {
	if err := os.WriteFile(path, body, octal); err != nil {
		fmt.Println("There was an error creating a file: ", err)
		return err
	}

	return nil
}

// TouchEmpty ...
func TouchEmpty(path string, mode fs.FileMode) error {
	n, err := os.Create(path)
	if err != nil {
		fmt.Println("There was an error creating a file: ", err)
		return err
	}
	// Change file permission bit
	err = os.Chmod(path, mode)
	if err != nil {
		fmt.Println("There was an error creating a file: ", err)
		return err
	}
	n.Close()

	return nil
}

// TouchOrReplace ...
func TouchOrReplace(path string, body []byte, octal fs.FileMode) {
	var _, errf = os.Stat(path)

	// create file if not exists
	if os.IsNotExist(errf) {
		TouchByteArray(path, body, octal)
		return
	}

	Replace(path, body, octal)
}

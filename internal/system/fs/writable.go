package fs

import (
	"github.com/ironstar-io/ironstar-cli/internal/system/fs/goos"
)

// Writable - Check if a file/folder is writable
func Writable(path string) bool {
	return goos.Writable(path)
}

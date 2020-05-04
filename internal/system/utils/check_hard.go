package utils

import (
	"os/exec"
)

// CheckCmd - Checks if an executable is available and if not, returns false
func CheckCmd(p string) bool {
	_, err := exec.LookPath(p)
	if err != nil {
		return false
	}

	return true
}

package cmd

import (
	"os/exec"
	"strings"
)

// ChildProcess - Execute a command and return the output value. No exit on stdErr
func ChildProcess(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	stdoutStderr, err := cmd.CombinedOutput()

	return strings.TrimSpace(string(stdoutStderr)), err
}

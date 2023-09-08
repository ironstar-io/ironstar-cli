package remote_command

import (
	"github.com/spf13/cobra"
)

// RemoteCommandCmd - `iron remote-command`
var RemoteCommandCmd = &cobra.Command{
	Use:   "remote-command",
	Short: "",
	Long:  "",
	Run:   list,
}

// RCCmd - `iron rc` (remote-command alias)
var RCCmd = &cobra.Command{
	Hidden: true,
	Use:    "rc",
	Short:  "",
	Long:   "",
	Run:    list,
}

// RemoteCommandsCmd - `iron remote-commands` (remote-command alias)
var RemoteCommandsCmd = &cobra.Command{
	Hidden: true,
	Use:    "remote-commands",
	Short:  "",
	Long:   "",
	Run:    list,
}

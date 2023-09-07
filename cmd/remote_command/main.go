package remote_command

import (
	"github.com/spf13/cobra"
)

// RemoteCommandCmd - `iron remote-command`
var RemoteCommandCmd = &cobra.Command{
	Hidden: true,
	Use:    "remote-command",
	Short:  "",
	Long:   "",
	Run:    list,
}

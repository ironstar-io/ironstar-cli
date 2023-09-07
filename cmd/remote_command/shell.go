package remote_command

import (
	"fmt"
	"os"

	"github.com/ironstar-io/ironstar-cli/cmd/flags"
	"github.com/ironstar-io/ironstar-cli/internal/api"
	"github.com/ironstar-io/ironstar-cli/internal/remote_command"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// ShellCmd - `iron remote-commands shell`
var ShellCmd = &cobra.Command{
	Use:   "shell [command]",
	Short: "Execute a shell command",
	Long:  "Execute a shell command",
	Run:   shell,
}

func shell(cmd *cobra.Command, args []string) {
	err := remote_command.Shell(args, flags.Acc)
	if err != nil {
		if err != api.ErrIronstarAPICall {
			fmt.Println()
			color.Red(err.Error())
		}

		os.Exit(1)
	}
}

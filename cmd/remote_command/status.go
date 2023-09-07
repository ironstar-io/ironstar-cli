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

// StatusCmd - `iron remote-command status [NameOrID]`
var StatusCmd = &cobra.Command{
	Use:   "status [nameOrID]",
	Short: "Remote command status",
	Long:  "Show the status of a remote command",
	Run:   status,
}

func status(cmd *cobra.Command, args []string) {
	err := remote_command.Status(args, flags.Acc)
	if err != nil {
		if err != api.ErrIronstarAPICall {
			fmt.Println()
			color.Red(err.Error())
		}

		os.Exit(1)
	}
}

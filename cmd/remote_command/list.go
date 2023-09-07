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

// ListCmd - `iron remote-commands list`
var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "List executed remote commands",
	Long:  "List executed remote command for an Ironstar environment",
	Run:   list,
}

func list(cmd *cobra.Command, args []string) {
	err := remote_command.List(args, flags.Acc)
	if err != nil {
		if err != api.ErrIronstarAPICall {
			fmt.Println()
			color.Red(err.Error())
		}

		os.Exit(1)
	}
}

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

// DrushCmd - `iron remote-commands drush [args]`
var DrushCmd = &cobra.Command{
	Use:   "drush [args]",
	Short: "Execute a drush command",
	Long:  "Execute a drush command",
	Run:   drush,
}

func drush(cmd *cobra.Command, args []string) {
	err := remote_command.Drush(args, flags.Acc)
	if err != nil {
		if err != api.ErrIronstarAPICall {
			fmt.Println()
			color.Red(err.Error())
		}

		os.Exit(1)
	}
}

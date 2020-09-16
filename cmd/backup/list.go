package backup

import (
	"fmt"
	"os"

	"gitlab.com/ironstar-io/ironstar-cli/cmd/flags"
	"gitlab.com/ironstar-io/ironstar-cli/internal/api"
	"gitlab.com/ironstar-io/ironstar-cli/internal/backup"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// ListCmd - `iron backup list`
var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "List available backups for a subscription",
	Long:  "List available backups for a subscription",
	Run:   list,
}

func list(cmd *cobra.Command, args []string) {
	err := backup.List(args, flags.Acc)
	if err != nil {
		if err != api.ErrIronstarAPICall {
			fmt.Println()
			color.Red(err.Error())
		}

		os.Exit(1)
	}
}

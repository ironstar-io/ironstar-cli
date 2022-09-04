package backup

import (
	"fmt"
	"os"

	"github.com/ironstar-io/ironstar-cli/cmd/flags"
	"github.com/ironstar-io/ironstar-cli/internal/api"
	"github.com/ironstar-io/ironstar-cli/internal/backup"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// InfoCmd - `iron backup info`
var InfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Backup information",
	Long:  "View information about backups",
	Run:   info,
}

func info(cmd *cobra.Command, args []string) {
	err := backup.Info(args, flags.Acc)
	if err != nil {
		if err != api.ErrIronstarAPICall {
			fmt.Println()
			color.Red(err.Error())
		}

		os.Exit(1)
	}
}

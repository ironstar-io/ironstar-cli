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

// NewCmd - `iron backup new`
var NewCmd = &cobra.Command{
	Use:   "new",
	Short: "Create new backup",
	Long:  "Backup a subscription environment",
	Run:   newBackup,
}

func newBackup(cmd *cobra.Command, args []string) {
	err := backup.New(args, flags.Acc)
	if err != nil {
		if err != api.ErrIronstarAPICall {
			fmt.Println()
			color.Red(err.Error())
		}

		os.Exit(1)
	}
}

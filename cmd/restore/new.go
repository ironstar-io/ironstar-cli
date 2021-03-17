package restore

import (
	"fmt"
	"os"

	"gitlab.com/ironstar-io/ironstar-cli/cmd/flags"
	"gitlab.com/ironstar-io/ironstar-cli/internal/api"
	"gitlab.com/ironstar-io/ironstar-cli/internal/restore"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// NewCmd - `iron restore new`
var NewCmd = &cobra.Command{
	Use:   "new",
	Short: "Create new restore",
	Long:  "Restore a backup to an environment",
	Run:   newRestore,
}

func newRestore(cmd *cobra.Command, args []string) {
	err := restore.New(args, flags.Acc)
	if err != nil {
		if err != api.ErrIronstarAPICall {
			fmt.Println()
			color.Red(err.Error())
		}

		os.Exit(1)
	}
}

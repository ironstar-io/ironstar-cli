package restore

import (
	"fmt"
	"os"

	"github.com/ironstar-io/ironstar-cli/cmd/flags"
	"github.com/ironstar-io/ironstar-cli/internal/api"
	"github.com/ironstar-io/ironstar-cli/internal/restore"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// InfoCmd - `iron restore info`
var InfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Restore information",
	Long:  "View in formation about restores",
	Run:   info,
}

func info(cmd *cobra.Command, args []string) {
	err := restore.Info(args, flags.Acc)
	if err != nil {
		if err != api.ErrIronstarAPICall {
			fmt.Println()
			color.Red(err.Error())
		}

		os.Exit(1)
	}
}

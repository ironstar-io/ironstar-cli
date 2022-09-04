package sync

import (
	"fmt"
	"os"

	"github.com/ironstar-io/ironstar-cli/cmd/flags"
	"github.com/ironstar-io/ironstar-cli/internal/api"
	"github.com/ironstar-io/ironstar-cli/internal/sync"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// InfoCmd - `iron sync info`
var InfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Sync information",
	Long:  "View information about syncs",
	Run:   info,
}

func info(cmd *cobra.Command, args []string) {
	err := sync.Info(args, flags.Acc)
	if err != nil {
		if err != api.ErrIronstarAPICall {
			fmt.Println()
			color.Red(err.Error())
		}

		os.Exit(1)
	}
}

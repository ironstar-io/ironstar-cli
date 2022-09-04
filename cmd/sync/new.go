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

// NewCmd - `iron sync new`
var NewCmd = &cobra.Command{
	Use:   "new",
	Short: "Create new environment sync",
	Long:  "Sync a source environment to a destination environment",
	Run:   newSync,
}

func newSync(cmd *cobra.Command, args []string) {
	err := sync.New(args, flags.Acc)
	if err != nil {
		if err != api.ErrIronstarAPICall {
			fmt.Println()
			color.Red(err.Error())
		}

		os.Exit(1)
	}
}

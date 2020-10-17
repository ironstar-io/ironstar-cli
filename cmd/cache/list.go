package cache

import (
	"fmt"
	"os"

	"gitlab.com/ironstar-io/ironstar-cli/cmd/flags"
	"gitlab.com/ironstar-io/ironstar-cli/internal/api"
	"gitlab.com/ironstar-io/ironstar-cli/internal/cache"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// ListInvalidationsCmd - `iron cache invalidation list`
var ListInvalidationsCmd = &cobra.Command{
	Use:   "list",
	Short: "List of cache invalidations",
	Long:  "List of cache invalidations for an environment",
	Run:   listInvalidations,
}

func listInvalidations(cmd *cobra.Command, args []string) {
	err := cache.ListInvalidations(args, flags.Acc)
	if err != nil {
		if err != api.ErrIronstarAPICall {
			fmt.Println()
			color.Red(err.Error())
		}

		os.Exit(1)
	}
}

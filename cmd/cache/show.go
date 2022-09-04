package cache

import (
	"fmt"
	"os"

	"github.com/ironstar-io/ironstar-cli/cmd/flags"
	"github.com/ironstar-io/ironstar-cli/internal/api"
	"github.com/ironstar-io/ironstar-cli/internal/cache"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// ShowInvalidationCmd - `iron cache invalidation show [name]`
var ShowInvalidationCmd = &cobra.Command{
	Use:   "show [name]",
	Short: "Show status of a single cache invalidation",
	Long:  "Show status of a single cache invalidation",
	Run:   showInvalidation,
}

func showInvalidation(cmd *cobra.Command, args []string) {
	err := cache.ShowInvalidation(args, flags.Acc)
	if err != nil {
		if err != api.ErrIronstarAPICall {
			fmt.Println()
			color.Red(err.Error())
		}

		os.Exit(1)
	}
}

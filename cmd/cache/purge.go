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

// CreateCmd - `iron cache invalidation create`
var CreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Purge the cache",
	Long:  "Purge the cache for an environment",
	Run:   create,
}

func create(cmd *cobra.Command, args []string) {
	err := cache.Purge(args, flags.Acc)
	if err != nil {
		if err != api.ErrIronstarAPICall {
			fmt.Println()
			color.Red(err.Error())
		}

		os.Exit(1)
	}
}

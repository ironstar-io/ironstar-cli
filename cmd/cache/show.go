package cache

import (
	"fmt"
	"os"
	"strings"

	"github.com/ironstar-io/ironstar-cli/cmd/flags"
	"github.com/ironstar-io/ironstar-cli/internal/api"
	"github.com/ironstar-io/ironstar-cli/internal/cache"
	"github.com/ironstar-io/ironstar-cli/internal/system/utils"

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
			if strings.ToLower(flags.Acc.Output) == "json" {
				utils.PrintErrorJSON(err)
				os.Exit(1)
			}

			fmt.Println()
			color.Red(err.Error())
		}

		os.Exit(1)
	}
}

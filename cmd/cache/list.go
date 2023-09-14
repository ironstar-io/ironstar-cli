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

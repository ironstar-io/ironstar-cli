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

// ShowInvalidationCmd - `iron cache show [name]`
var ShowInvalidationCmd = newShowInvalidationCmd(false)

// LegacyShowInvalidationCmd preserves `iron cache invalidation show [name]`.
var LegacyShowInvalidationCmd = newShowInvalidationCmd(true)

func newShowInvalidationCmd(hidden bool) *cobra.Command {
	return &cobra.Command{
		Hidden:  hidden,
		Use:     "show [name]",
		Short:   "Show a cache invalidation",
		Long:    "Show the current status and details of one cache invalidation.",
		Example: "  iron cache show <name> --subscription example --environment stage",
		Args:    cobra.ExactArgs(1),
		Run:     showInvalidation,
	}
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

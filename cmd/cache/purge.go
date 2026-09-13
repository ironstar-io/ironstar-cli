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

const invalidateLong = `Invalidate cached content for an environment.

Supply one --url to invalidate that URL. URL invalidations are soft by default:
the cached object is marked stale and revalidated on a later request. Use
--type hard to remove the cached object immediately.

Omit --url to hard-invalidate the entire environment cache. Whole-environment
invalidation requires confirmation unless --yes (-y) is supplied.

After a URL invalidation is accepted, a later FAILED status does not necessarily
mean the request was not submitted. Fastly may report failure when the URL was
not found in cache.`

const invalidateExample = `  # Soft-invalidate one URL
  iron cache invalidate --subscription example --environment dev --url https://dev.example.com/article

  # Hard-invalidate one URL
  iron cache invalidate --subscription example --environment stage --url https://stage.example.com/article --type hard

  # Hard-invalidate the entire environment cache and skip confirmation
  iron cache invalidate --subscription example --environment dev --yes`

// InvalidateCmd - `iron cache invalidate`
var InvalidateCmd = newInvalidateCmd("invalidate", false)

// CreateCmd preserves the legacy `iron cache invalidation create` command path.
var CreateCmd = newInvalidateCmd("create", true)

func newInvalidateCmd(use string, hidden bool) *cobra.Command {
	return &cobra.Command{
		Hidden:  hidden,
		Use:     use + " [flags]",
		Short:   "Invalidate cached content",
		Long:    invalidateLong,
		Example: invalidateExample,
		Args:    cobra.NoArgs,
		Run:     invalidate,
	}
}

func invalidate(cmd *cobra.Command, args []string) {
	err := cache.Invalidate(args, flags.Acc)
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

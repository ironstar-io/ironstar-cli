package cache

import (
	"github.com/spf13/cobra"
)

// CacheCmd - `iron cache`
var CacheCmd = &cobra.Command{
	Use:   "cache",
	Short: "Manage cache invalidations",
	Long: `Create and inspect cache invalidations for an Ironstar environment.

Invalidate either one HTTPS URL or the entire environment cache. URL
invalidations may be soft or hard; whole-environment invalidations are always
hard and require confirmation unless --yes (-y) is supplied.`,
	Example: `  # Soft-invalidate one URL (the default for URL invalidations)
  iron cache invalidate --subscription example --environment dev --url https://dev.example.com/article

  # Hard-invalidate one URL
  iron cache invalidate --subscription example --environment stage --url https://stage.example.com/article --type hard

  # Hard-invalidate the entire environment cache (prompts for confirmation)
  iron cache invalidate --subscription example --environment dev

  # List recent invalidations and inspect one result
  iron cache list --subscription example --environment dev
  iron cache show <name> --subscription example --environment dev`,
	Args: cobra.NoArgs,
}

// InvalidationCmd preserves the legacy `iron cache invalidation` command path.
var InvalidationCmd = &cobra.Command{
	Hidden: true,
	Use:    "invalidation",
	Short:  "Legacy cache invalidation commands",
	Run:    listInvalidations,
}

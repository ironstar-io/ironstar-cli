package cache

import (
	"github.com/spf13/cobra"
)

// CacheCmd - `iron cache`
var CacheCmd = &cobra.Command{
	Hidden: true,
	Use:    "cache",
	Short:  "",
	Long:   "",
	Run:    listInvalidations,
}

// InvalidationCmd - `iron cache invalidation`
var InvalidationCmd = &cobra.Command{
	Hidden: true,
	Use:    "invalidation",
	Short:  "",
	Long:   "",
	Run:    listInvalidations,
}

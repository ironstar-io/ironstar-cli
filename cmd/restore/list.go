package restore

import (
	"github.com/spf13/cobra"
)

// ListCmd - `iron restore list`
var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "Restore information",
	Long:  "View in formation about restores",
	Run:   info,
}

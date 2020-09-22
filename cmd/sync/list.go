package sync

import (
	"github.com/spf13/cobra"
)

// ListCmd - `iron sync list`
var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "Sync information",
	Long:  "View information about syncs",
	Run:   info,
}

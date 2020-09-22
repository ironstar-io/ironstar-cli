package backup

import (
	"github.com/spf13/cobra"
)

// ListCmd - `iron backup list`
var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "Backup information",
	Long:  "View in formation about backups",
	Run:   info,
}

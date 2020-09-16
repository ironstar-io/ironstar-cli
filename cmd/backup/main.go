package backup

import (
	"github.com/spf13/cobra"
)

// BackupCmd - `iron backup`
var BackupCmd = &cobra.Command{
	Use:   "backup",
	Short: "",
	Long:  "",
	Run:   new,
}

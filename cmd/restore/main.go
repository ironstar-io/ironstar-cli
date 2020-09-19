package restore

import (
	"github.com/spf13/cobra"
)

// RestoreCmd - `iron restore`
var RestoreCmd = &cobra.Command{
	Use:   "restore",
	Short: "",
	Long:  "",
	Run:   new,
}

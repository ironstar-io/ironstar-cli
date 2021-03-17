package sync

import (
	"github.com/spf13/cobra"
)

// SyncCmd - `iron sync`
var SyncCmd = &cobra.Command{
	Use:   "sync",
	Short: "",
	Long:  "",
	Run:   newSync,
}

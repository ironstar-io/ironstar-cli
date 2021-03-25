package antivirus

import (
	"github.com/spf13/cobra"
)

// AntivirusCmd - `iron antivirus`
var AntivirusCmd = &cobra.Command{
	Hidden: true,
	Use:    "antivirus",
	Short:  "",
	Long:   "",
	Run:    listScans,
}

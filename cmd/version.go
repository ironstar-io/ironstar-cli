package cmd

import (
	"gitlab.com/ironstar-io/ironstar-cli/internal/system/version"

	"github.com/spf13/cobra"
)

// VersionCmd - `iron version`
var VersionCmd = &cobra.Command{
	Use:   "version [version]",
	Short: "Print Ironstar CLI version information",
	Long:  "Print Ironstar CLI version information including 'Build Date', 'Compiler' and 'Platform'",
	Run: func(cmd *cobra.Command, args []string) {
		version.Display()
	},
}

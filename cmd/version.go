package cmd

import (
	"gitlab.com/ironstar-io/ironstar-cli/internal/system/version"

	"github.com/spf13/cobra"
)

// VersionCmd - `tok version`
var VersionCmd = &cobra.Command{
	Use:   "version [version]",
	Short: "Print Tokdaido version information",
	Long:  "Print Tokdaido version information including 'Build Date', 'Compiler' and 'Platform'",
	Run: func(cmd *cobra.Command, args []string) {
		version.Display()
		// if len(args) == 0 || args[0] == "" {

		// 	return
		// }

		// version.Select(args[0])
	},
}

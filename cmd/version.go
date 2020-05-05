package cmd

import (
	"fmt"
	"os"

	"gitlab.com/ironstar-io/ironstar-cli/cmd/flags"
	"gitlab.com/ironstar-io/ironstar-cli/internal/api"
	"gitlab.com/ironstar-io/ironstar-cli/internal/system/version"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// VersionCmd - `iron version`
var VersionCmd = &cobra.Command{
	Use:   "version [version]",
	Short: "Print Ironstar CLI version information",
	Long:  "Print Ironstar CLI version information including 'Build Date', 'Compiler' and 'Platform'",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 || args[0] == "" {
			version.Display()

			return
		}

		err := version.Select(args[0], flags.Acc)
		if err != nil {
			if err != api.ErrIronstarAPICall {
				fmt.Println()
				color.Red(err.Error())
			}

			os.Exit(1)
		}
	},
}

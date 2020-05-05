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

// UpgradeCmd - `iron upgrade`
var UpgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "Upgrade to the latest version of the Ironstar CLI",
	Long:  "Upgrade to the latest version of the Ironstar CLI",
	Run: func(cmd *cobra.Command, args []string) {
		err := version.Upgrade(flags.Acc)
		if err != nil {
			if err != api.ErrIronstarAPICall {
				fmt.Println()
				color.Red(err.Error())
			}

			os.Exit(1)
		}
	},
}

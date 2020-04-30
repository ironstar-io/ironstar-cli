package auth

import (
	"fmt"
	"os"

	"gitlab.com/ironstar-io/ironstar-cli/cmd/flags"
	"gitlab.com/ironstar-io/ironstar-cli/internal/auth"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// MFAEnableCmd - `tok auth mfa enable`
var MFAEnableCmd = &cobra.Command{
	Use:   "mfa enable",
	Short: "Enable MFA",
	Long:  "Enable MFA for the currently logged in user",
	Run: func(cmd *cobra.Command, args []string) {
		err := auth.MFAEnable(args, flags.Acc)
		if err != nil {
			fmt.Println()
			color.Red(err.Error())

			os.Exit(1)
		}
	},
}

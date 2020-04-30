package auth

import (
	"fmt"
	"os"

	"gitlab.com/ironstar-io/ironstar-cli/internal/auth"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// MFADisableCmd - `tok auth mfa disable`
var MFADisableCmd = &cobra.Command{
	Use:   "mfa disable",
	Short: "Disable MFA",
	Long:  "Disable MFA for the currently logged in user",
	Run: func(cmd *cobra.Command, args []string) {
		err := auth.MFADisable(args)
		if err != nil {
			fmt.Println()
			color.Red(err.Error())

			os.Exit(1)
		}
	},
}

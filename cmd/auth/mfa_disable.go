package auth

import (
	"fmt"
	"os"

	"gitlab.com/ironstar-io/ironstar-cli/cmd/flags"
	"gitlab.com/ironstar-io/ironstar-cli/internal/auth"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// MFADisableCmd - `iron auth mfa disable`
var MFADisableCmd = &cobra.Command{
	Use:   "disable",
	Short: "Disable MFA",
	Long:  "Disable MFA for the currently logged in user",
	Run: func(cmd *cobra.Command, args []string) {
		err := auth.MFADisable(args, flags.Acc)
		if err != nil {
			fmt.Println()
			color.Red(err.Error())

			os.Exit(1)
		}
	},
}

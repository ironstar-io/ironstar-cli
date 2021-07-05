package auth

import (
	"fmt"
	"os"

	"gitlab.com/ironstar-io/ironstar-cli/cmd/flags"
	"gitlab.com/ironstar-io/ironstar-cli/internal/auth"
	"gitlab.com/ironstar-io/ironstar-cli/internal/types"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// MFAEnableCmd - `iron auth mfa enable`
var MFAEnableCmd = &cobra.Command{
	Use:   "enable",
	Short: "Enable MFA",
	Long:  "Enable MFA for the currently logged in user",
	Run: func(cmd *cobra.Command, args []string) {
		_, err := auth.MFAEnable(flags.Acc, types.Keylink{})
		if err != nil {
			fmt.Println()
			color.Red(err.Error())

			os.Exit(1)
		}
	},
}

package auth

import (
	"fmt"
	"os"

	"github.com/ironstar-io/ironstar-cli/cmd/flags"
	"github.com/ironstar-io/ironstar-cli/internal/auth"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// MFARecoveryCmd - `iron auth mfa recovery`
var MFARecoveryCmd = &cobra.Command{
	Use:   "recovery",
	Short: "Recovery MFA",
	Long:  "Recovery MFA for the currently logged in user",
	Run: func(cmd *cobra.Command, args []string) {
		err := auth.MFARecovery(args, flags.Acc)
		if err != nil {
			fmt.Println()
			color.Red(err.Error())

			os.Exit(1)
		}
	},
}

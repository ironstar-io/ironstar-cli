package cmd

import (
	"fmt"
	"os"

	"github.com/ironstar-io/ironstar-cli/cmd/flags"
	"github.com/ironstar-io/ironstar-cli/internal/api"
	"github.com/ironstar-io/ironstar-cli/internal/auth"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// LogoutCmd - `iron auth logout`
var LogoutCmd = &cobra.Command{
	Use:   "logout [email]",
	Short: "Logout from your Ironstar API session",
	Long:  "Logout from your Ironstar API session. This destroys the session token in the Ironstar system",
	Run: func(cmd *cobra.Command, args []string) {
		err := auth.Logout(args, flags.Acc)
		if err != nil {
			if err != api.ErrIronstarAPICall {
				fmt.Println()
				color.Red(err.Error())
			}

			os.Exit(1)
		}
	},
}

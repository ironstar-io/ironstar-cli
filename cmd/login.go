package cmd

import (
	"fmt"
	"os"

	"gitlab.com/ironstar-io/ironstar-cli/internal/api"
	"gitlab.com/ironstar-io/ironstar-cli/internal/auth"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// LoginCmd - `tok auth login`
var LoginCmd = &cobra.Command{
	Use:   "login [email]",
	Short: "Authenticate for the Ironstar API",
	Long:  "Authenticate and store credentials for use against the Ironstar API",
	Run: func(cmd *cobra.Command, args []string) {
		err := auth.IronstarAPILogin(args, PasswordFlag)
		if err != nil {
			if err != api.ErrIronstarAPICall {
				fmt.Println()
				color.Red(err.Error())
			}

			os.Exit(1)
		}
	},
}

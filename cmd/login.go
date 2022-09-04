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

// LoginCmd - `iron auth login`
var LoginCmd = &cobra.Command{
	Use:   "login [email]",
	Short: "Authenticate for the Ironstar API",
	Long:  "Authenticate and store credentials for use against the Ironstar API",
	Run: func(cmd *cobra.Command, args []string) {
		err := auth.IronstarAPILogin(args, flags.Acc)
		if err != nil {
			if err != api.ErrIronstarAPICall {
				fmt.Println()
				color.Red(err.Error())
			}

			os.Exit(1)
		}
	},
}

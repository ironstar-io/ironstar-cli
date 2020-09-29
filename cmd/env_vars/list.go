package env_vars

import (
	"fmt"
	"os"

	"gitlab.com/ironstar-io/ironstar-cli/cmd/flags"
	"gitlab.com/ironstar-io/ironstar-cli/internal/api"
	"gitlab.com/ironstar-io/ironstar-cli/internal/env_vars"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// ListCmd - `iron env-vars list`
var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "List available environment variables for a subscription",
	Long:  "List available environment variables for a subscription",
	Run:   list,
}

func list(cmd *cobra.Command, args []string) {
	err := env_vars.List(args, flags.Acc)
	if err != nil {
		if err != api.ErrIronstarAPICall {
			fmt.Println()
			color.Red(err.Error())
		}

		os.Exit(1)
	}
}

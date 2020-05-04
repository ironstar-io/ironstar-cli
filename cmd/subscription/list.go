package subscription

import (
	"fmt"
	"os"

	"gitlab.com/ironstar-io/ironstar-cli/cmd/flags"
	"gitlab.com/ironstar-io/ironstar-cli/internal/api"
	"gitlab.com/ironstar-io/ironstar-cli/internal/subscription"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// ListCmd - `iron subscription list`
var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "List available subscriptions",
	Long:  "List all available subscriptions for the user",
	Run: func(cmd *cobra.Command, args []string) {
		err := subscription.List(args, flags.Acc)
		if err != nil {
			if err != api.ErrIronstarAPICall {
				fmt.Println()
				color.Red(err.Error())
			}

			os.Exit(1)
		}
	},
}

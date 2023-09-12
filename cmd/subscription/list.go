package subscription

import (
	"fmt"
	"os"
	"strings"

	"github.com/ironstar-io/ironstar-cli/cmd/flags"
	"github.com/ironstar-io/ironstar-cli/internal/api"
	"github.com/ironstar-io/ironstar-cli/internal/subscription"
	"github.com/ironstar-io/ironstar-cli/internal/system/utils"

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
				if strings.ToLower(flags.Acc.Output) == "json" {
					utils.PrintErrorJSON(err)
					os.Exit(1)
				}

				fmt.Println()
				color.Red(err.Error())
			}

			os.Exit(1)
		}
	},
}

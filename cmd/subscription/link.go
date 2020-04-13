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

// LinkCmd - `tok subscription link [subscription]`
var LinkCmd = &cobra.Command{
	Use:   "link [subscription]",
	Short: "Link your project to a subscription",
	Long:  "Link your project to a subscription",
	Run: func(cmd *cobra.Command, args []string) {
		err := subscription.Link(args, flags.Acc)
		if err != nil {
			if err != api.ErrIronstarAPICall {
				fmt.Println()
				color.Red(err.Error())
			}

			os.Exit(1)
		}

		color.Green("Successfully linked project!")
	},
}

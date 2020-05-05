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

// ShowCmd - `iron subscription show`
var ShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show linked subscription",
	Long:  "Show linked subscription",
	Run:   show,
}

func show(cmd *cobra.Command, args []string) {
	err := subscription.Show(args, flags.Acc)
	if err != nil {
		if err != api.ErrIronstarAPICall {
			fmt.Println()
			color.Red(err.Error())
		}

		os.Exit(1)
	}
}

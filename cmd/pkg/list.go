package pkg

import (
	"fmt"
	"os"

	"gitlab.com/ironstar-io/ironstar-cli/cmd/flags"
	"gitlab.com/ironstar-io/ironstar-cli/internal/api"
	"gitlab.com/ironstar-io/ironstar-cli/internal/pkg"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// ListCmd - `iron package list`
var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "List packages",
	Long:  "List packages for a subscription",
	Run:   list,
}

func list(cmd *cobra.Command, args []string) {
	err := pkg.List(args, flags.Login, flags.Output, flags.Subscription)
	if err != nil {
		if err != api.ErrIronstarAPICall {
			fmt.Println()
			color.Red(err.Error())
		}

		os.Exit(1)
	}
}

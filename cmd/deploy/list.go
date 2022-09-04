package deploy

import (
	"fmt"
	"os"

	"github.com/ironstar-io/ironstar-cli/cmd/flags"
	"github.com/ironstar-io/ironstar-cli/internal/api"
	"github.com/ironstar-io/ironstar-cli/internal/deploy"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// ListCmd - `iron deploy list`
var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "List deployments",
	Long:  "List all deployments for a subscription, optionally filter by environment or package",
	Run:   list,
}

func list(cmd *cobra.Command, args []string) {
	err := deploy.List(args, flags.Acc)
	if err != nil {
		if err != api.ErrIronstarAPICall {
			fmt.Println()
			color.Red(err.Error())
		}

		os.Exit(1)
	}
}

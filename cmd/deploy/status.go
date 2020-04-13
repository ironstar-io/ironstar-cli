package deploy

import (
	"fmt"
	"os"

	"gitlab.com/ironstar-io/ironstar-cli/cmd/flags"
	"gitlab.com/ironstar-io/ironstar-cli/internal/api"
	"gitlab.com/ironstar-io/ironstar-cli/internal/deploy"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// StatusCmd - `iron deploy status [deployID]`
var StatusCmd = &cobra.Command{
	Use:   "status [deployID]",
	Short: "Deployment status",
	Long:  "Show the status of a deployment",
	Run:   status,
}

func status(cmd *cobra.Command, args []string) {
	err := deploy.Status(args, flags.Acc)
	if err != nil {
		if err != api.ErrIronstarAPICall {
			fmt.Println()
			color.Red(err.Error())
		}

		os.Exit(1)
	}
}

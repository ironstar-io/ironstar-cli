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

// DeployCmd - `iron deploy`
var DeployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Create a new deployment",
	Long:  "Packages up the project as a tarball and uploads to Ironstar for remote deployment",
	Run:   create,
}

func create(cmd *cobra.Command, args []string) {
	err := deploy.Create(args, flags.Acc)
	if err != nil {
		if err != api.ErrIronstarAPICall {
			fmt.Println()
			color.Red(err.Error())
		}

		os.Exit(1)
	}
}

package pkg

import (
	"fmt"
	"os"

	"github.com/ironstar-io/ironstar-cli/cmd/flags"
	"github.com/ironstar-io/ironstar-cli/internal/api"
	"github.com/ironstar-io/ironstar-cli/internal/pkg"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// CreateCmd - `iron package create`
var CreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create package",
	Long:  "Create a new packages for a subscription",
	Run:   create,
}

func create(cmd *cobra.Command, args []string) {
	err := pkg.Create(args, flags.Acc)
	if err != nil {
		if err != api.ErrIronstarAPICall {
			fmt.Println()
			color.Red(err.Error())
		}

		os.Exit(1)
	}
}

package environment

import (
	"fmt"
	"os"

	"gitlab.com/ironstar-io/ironstar-cli/cmd/flags"
	"gitlab.com/ironstar-io/ironstar-cli/internal/api"
	"gitlab.com/ironstar-io/ironstar-cli/internal/environment"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// ShowCmd - `iron environment show`
var ShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show linked environment",
	Long:  "Show linked environment",
	Run:   show,
}

func show(cmd *cobra.Command, args []string) {
	err := environment.Show(args, flags.Acc)
	if err != nil {
		if err != api.ErrIronstarAPICall {
			fmt.Println()
			color.Red(err.Error())
		}

		os.Exit(1)
	}
}

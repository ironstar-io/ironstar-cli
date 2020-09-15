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

// LinkCmd - `iron environment link [environment]`
var LinkCmd = &cobra.Command{
	Use:   "link [environment]",
	Short: "Link your project to a environment",
	Long:  "Link your project to a environment",
	Run: func(cmd *cobra.Command, args []string) {
		err := environment.Link(args, flags.Acc)
		if err != nil {
			if err != api.ErrIronstarAPICall {
				fmt.Println()
				color.Red(err.Error())
			}

			os.Exit(1)
		}

		color.Green("Successfully linked to project!")
	},
}

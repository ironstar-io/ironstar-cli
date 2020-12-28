package new_relic

import (
	"fmt"
	"os"

	"gitlab.com/ironstar-io/ironstar-cli/cmd/flags"
	"gitlab.com/ironstar-io/ironstar-cli/internal/api"
	"gitlab.com/ironstar-io/ironstar-cli/internal/new_relic"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// ConfigureCmd - `iron new-relic configure`
var ConfigureCmd = &cobra.Command{
	Use:   "configure",
	Short: "Configure New Relic for your environment",
	Long:  "Prompts for your New Relic configuration values and applies them to your environment",
	Run:   configure,
}

func configure(cmd *cobra.Command, args []string) {
	err := new_relic.Configure(args, flags.Acc)
	if err != nil {
		if err != api.ErrIronstarAPICall {
			fmt.Println()
			color.Red(err.Error())
		}

		os.Exit(1)
	}
}

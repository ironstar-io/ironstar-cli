package newrelic

import (
	"fmt"
	"os"

	"github.com/ironstar-io/ironstar-cli/cmd/flags"
	"github.com/ironstar-io/ironstar-cli/internal/api"
	"github.com/ironstar-io/ironstar-cli/internal/newrelic"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// ConfigureCmd - `iron newrelic configure`
var ConfigureCmd = &cobra.Command{
	Use:   "configure",
	Short: "Configure New Relic for your environment",
	Long:  "Prompts for your New Relic configuration values and applies them to your environment",
	Run:   configure,
}

func configure(cmd *cobra.Command, args []string) {
	err := newrelic.Configure(args, flags.Acc)
	if err != nil {
		if err != api.ErrIronstarAPICall {
			fmt.Println()
			color.Red(err.Error())
		}

		os.Exit(1)
	}
}

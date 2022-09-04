package env_vars

import (
	"fmt"
	"os"

	"github.com/ironstar-io/ironstar-cli/cmd/flags"
	"github.com/ironstar-io/ironstar-cli/internal/api"
	"github.com/ironstar-io/ironstar-cli/internal/env_vars"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// AddCmd - `iron env-vars add`
var AddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add an environment variable",
	Long:  "Add an environment variable to an Ironstar environment",
	Run:   add,
}

func add(cmd *cobra.Command, args []string) {
	err := env_vars.Add(args, flags.Acc)
	if err != nil {
		if err != api.ErrIronstarAPICall {
			fmt.Println()
			color.Red(err.Error())
		}

		os.Exit(1)
	}
}

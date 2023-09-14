package env_vars

import (
	"fmt"
	"os"
	"strings"

	"github.com/ironstar-io/ironstar-cli/cmd/flags"
	"github.com/ironstar-io/ironstar-cli/internal/api"
	"github.com/ironstar-io/ironstar-cli/internal/env_vars"
	"github.com/ironstar-io/ironstar-cli/internal/system/utils"

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
			if strings.ToLower(flags.Acc.Output) == "json" {
				utils.PrintErrorJSON(err)
				os.Exit(1)
			}

			fmt.Println()
			color.Red(err.Error())
		}

		os.Exit(1)
	}
}

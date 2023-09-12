package hooks

import (
	"fmt"
	"os"
	"strings"

	"github.com/ironstar-io/ironstar-cli/cmd/flags"
	"github.com/ironstar-io/ironstar-cli/internal/api"
	"github.com/ironstar-io/ironstar-cli/internal/environment"
	"github.com/ironstar-io/ironstar-cli/internal/system/utils"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// AddHookCmd - `iron environment add-hook [hookname]`
var AddHookCmd = &cobra.Command{
	Use:   "add-hook [hookname]",
	Short: "Add environment hook",
	Long:  "Add an environment hook (PRE_DEPLOYMENT_BACKUP)",
	Run:   AddHook,
}

func AddHook(cmd *cobra.Command, args []string) {
	err := environment.AddHook(args, flags.Acc)
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

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

// RemoveHookCmd - `iron environment remove-hook [hookname]`
var RemoveHookCmd = &cobra.Command{
	Use:   "remove-hook [hookname]",
	Short: "Remove environment hook",
	Long:  "Remove an environment hook (PRE_DEPLOYMENT_BACKUP)",
	Run:   RemoveHook,
}

func RemoveHook(cmd *cobra.Command, args []string) {
	err := environment.RemoveHook(args, flags.Acc)
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

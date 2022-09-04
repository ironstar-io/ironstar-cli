package hooks

import (
	"fmt"
	"os"

	"github.com/ironstar-io/ironstar-cli/cmd/flags"
	"github.com/ironstar-io/ironstar-cli/internal/api"
	"github.com/ironstar-io/ironstar-cli/internal/environment"

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
			fmt.Println()
			color.Red(err.Error())
		}

		os.Exit(1)
	}
}

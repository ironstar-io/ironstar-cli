package environment

import (
	"fmt"
	"os"

	"github.com/ironstar-io/ironstar-cli/cmd/flags"
	"github.com/ironstar-io/ironstar-cli/internal/api"
	"github.com/ironstar-io/ironstar-cli/internal/constants"
	"github.com/ironstar-io/ironstar-cli/internal/environment"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// DisableRestoreCmd - `iron environment disable-restore`
var DisableRestoreCmd = &cobra.Command{
	Use:   "disable-restore",
	Short: "Disable restore for an environment",
	Long:  "Disable restore for an environment",
	Run:   DisableRestore,
}

func DisableRestore(cmd *cobra.Command, args []string) {
	err := environment.EnableDisableRestore(args, flags.Acc, constants.RestorePermissionNotAllowed)
	if err != nil {
		if err != api.ErrIronstarAPICall {
			fmt.Println()
			color.Red(err.Error())
		}

		os.Exit(1)
	}
}

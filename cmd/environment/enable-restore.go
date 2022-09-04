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

// EnableRestoreCmd - `iron environment enable-restore`
var EnableRestoreCmd = &cobra.Command{
	Use:   "enable-restore",
	Short: "Enable restore for an environment",
	Long:  "Enable restore for an environment",
	Run:   EnableRestore,
}

func EnableRestore(cmd *cobra.Command, args []string) {
	err := environment.EnableDisableRestore(args, flags.Acc, constants.RestorePermissionAllowed)
	if err != nil {
		if err != api.ErrIronstarAPICall {
			fmt.Println()
			color.Red(err.Error())
		}

		os.Exit(1)
	}
}

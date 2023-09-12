package environment

import (
	"fmt"
	"os"
	"strings"

	"github.com/ironstar-io/ironstar-cli/cmd/flags"
	"github.com/ironstar-io/ironstar-cli/internal/api"
	"github.com/ironstar-io/ironstar-cli/internal/constants"
	"github.com/ironstar-io/ironstar-cli/internal/environment"
	"github.com/ironstar-io/ironstar-cli/internal/system/utils"

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

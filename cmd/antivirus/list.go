package antivirus

import (
	"fmt"
	"os"
	"strings"

	"github.com/ironstar-io/ironstar-cli/cmd/flags"
	"github.com/ironstar-io/ironstar-cli/internal/antivirus"
	"github.com/ironstar-io/ironstar-cli/internal/api"
	"github.com/ironstar-io/ironstar-cli/internal/system/utils"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// ListScans - `iron antivirus list-scans`
var ListScans = &cobra.Command{
	Use:   "list-scans",
	Short: "List of antivirus scans",
	Long:  "List of antivirus scans for an environment",
	Run:   listScans,
}

func listScans(cmd *cobra.Command, args []string) {
	err := antivirus.ListScans(args, flags.Acc)
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

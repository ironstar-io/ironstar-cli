package auth

import (
	"fmt"
	"os"
	"strings"

	"github.com/ironstar-io/ironstar-cli/cmd/flags"
	"github.com/ironstar-io/ironstar-cli/internal/auth"
	"github.com/ironstar-io/ironstar-cli/internal/system/utils"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// SetActiveCmd - `iron auth set-active [email]`
var SetActiveCmd = &cobra.Command{
	Use:   "set-active [email]",
	Short: "Set credential as active",
	Long:  "Set credential as active",
	Run: func(cmd *cobra.Command, args []string) {
		err := auth.IronstarSetActiveCredentials(args, flags.Acc)
		if err != nil {
			if strings.ToLower(flags.Acc.Output) == "json" {
				utils.PrintErrorJSON(err)
				os.Exit(1)
			}

			fmt.Println()
			color.Red(err.Error())

			os.Exit(1)
		}
	},
}

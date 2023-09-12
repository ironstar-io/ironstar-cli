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

// ShowCmd - `iron auth show`
var ShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show available credentials",
	Long:  "Show available credentials for the Ironstar API and how they map against your projects",
	Run: func(cmd *cobra.Command, args []string) {
		err := auth.IronstarShowCredentials(args, flags.Acc)
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

package auth

import (
	"fmt"
	"os"

	"gitlab.com/ironstar-io/ironstar-cli/internal/auth"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// SetActiveCmd - `tok auth set-active [email]`
var SetActiveCmd = &cobra.Command{
	Use:   "set-active [email]",
	Short: "Set credential as active",
	Long:  "Set credential as active",
	Run: func(cmd *cobra.Command, args []string) {
		err := auth.IronstarSetActiveCredentials(args)
		if err != nil {
			fmt.Println()
			color.Red(err.Error())

			os.Exit(1)
		}
	},
}

package cmd

import (
	"fmt"
	"os"

	"github.com/ironstar-io/ironstar-cli/internal/api"
	"github.com/ironstar-io/ironstar-cli/internal/services"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// InitCmd - `iron init`
var InitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initializes an Ironstar project",
	Long:  "Initializes an Ironstar project by creating configuration files in the current working directory",
	Run: func(cmd *cobra.Command, args []string) {
		err := services.InitializeIronstarProject()
		if err != nil {
			if err != api.ErrIronstarAPICall {
				fmt.Println()
				color.Red(err.Error())
			}

			os.Exit(1)
		}

		color.Green("Successfully created the required Ironstar configuration files!")
		fmt.Println()
		color.Green("Run `iron subscription link [subscription-name]` to link this project to a remote Ironstar subscription")
	},
}

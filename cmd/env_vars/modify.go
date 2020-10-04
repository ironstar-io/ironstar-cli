package env_vars

import (
	"fmt"
	"os"

	"gitlab.com/ironstar-io/ironstar-cli/cmd/flags"
	"gitlab.com/ironstar-io/ironstar-cli/internal/api"
	"gitlab.com/ironstar-io/ironstar-cli/internal/env_vars"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// ModifyCmd - `iron env-vars modify`
var ModifyCmd = &cobra.Command{
	Use:   "modify",
	Short: "Modify an environment variable",
	Long:  "Modify an environment variable in an Ironstar environment",
	Run:   modify,
}

func modify(cmd *cobra.Command, args []string) {
	err := env_vars.Modify(args, flags.Acc)
	if err != nil {
		if err != api.ErrIronstarAPICall {
			fmt.Println()
			color.Red(err.Error())
		}

		os.Exit(1)
	}
}

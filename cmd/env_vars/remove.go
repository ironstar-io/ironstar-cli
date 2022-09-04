package env_vars

import (
	"fmt"
	"os"

	"github.com/ironstar-io/ironstar-cli/cmd/flags"
	"github.com/ironstar-io/ironstar-cli/internal/api"
	"github.com/ironstar-io/ironstar-cli/internal/env_vars"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// RemoveCmd - `iron env-vars remove`
var RemoveCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove an environment variable",
	Long:  "Remove an environment variable from an Ironstar environment",
	Run:   remove,
}

func remove(cmd *cobra.Command, args []string) {
	err := env_vars.Remove(args, flags.Acc)
	if err != nil {
		if err != api.ErrIronstarAPICall {
			fmt.Println()
			color.Red(err.Error())
		}

		os.Exit(1)
	}
}

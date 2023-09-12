package remote_command

import (
	"fmt"
	"os"
	"strings"

	"github.com/ironstar-io/ironstar-cli/cmd/flags"
	"github.com/ironstar-io/ironstar-cli/internal/api"
	"github.com/ironstar-io/ironstar-cli/internal/remote_command"
	"github.com/ironstar-io/ironstar-cli/internal/system/utils"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// DrushCacheRebuildCmd - `iron remote-commands drush-cache-rebuild`
var DrushCacheRebuildCmd = &cobra.Command{
	Use:   "drush-cache-rebuild",
	Short: "Execute drush cache-rebuild",
	Long:  "Execute drush cache-rebuild",
	Run:   drushCacheRebuild,
}

func drushCacheRebuild(cmd *cobra.Command, args []string) {
	err := remote_command.DrushCacheRebuild(args, flags.Acc)
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

package remote_command

import (
	"fmt"
	"os"

	"github.com/ironstar-io/ironstar-cli/cmd/flags"
	"github.com/ironstar-io/ironstar-cli/internal/api"
	"github.com/ironstar-io/ironstar-cli/internal/remote_command"

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
			fmt.Println()
			color.Red(err.Error())
		}

		os.Exit(1)
	}
}

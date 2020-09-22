package backup

import (
	"fmt"
	"os"

	"gitlab.com/ironstar-io/ironstar-cli/cmd/flags"
	"gitlab.com/ironstar-io/ironstar-cli/internal/api"
	"gitlab.com/ironstar-io/ironstar-cli/internal/backup"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// DeleteCmd - `iron backup delete`
var DeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a backup",
	Long:  "Delete an environment backup iteration",
	Run:   delete,
}

func delete(cmd *cobra.Command, args []string) {
	err := backup.Delete(args, flags.Acc)
	if err != nil {
		if err != api.ErrIronstarAPICall {
			fmt.Println()
			color.Red(err.Error())
		}

		os.Exit(1)
	}
}

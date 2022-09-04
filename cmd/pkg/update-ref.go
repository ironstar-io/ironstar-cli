package pkg

import (
	"fmt"
	"os"

	"github.com/ironstar-io/ironstar-cli/cmd/flags"
	"github.com/ironstar-io/ironstar-cli/internal/api"
	"github.com/ironstar-io/ironstar-cli/internal/pkg"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// UpdateRefCmd - `iron package update-ref`
var UpdateRefCmd = &cobra.Command{
	Use:   "update-ref",
	Short: "Update package ref",
	Long:  "Update the user defined 'ref' field for a package",
	Run:   updateRef,
}

func updateRef(cmd *cobra.Command, args []string) {
	err := pkg.UpdateRef(args, flags.Acc)
	if err != nil {
		if err != api.ErrIronstarAPICall {
			fmt.Println()
			color.Red(err.Error())
		}

		os.Exit(1)
	}
}

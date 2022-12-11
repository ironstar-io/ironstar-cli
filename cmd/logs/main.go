package logs

import (
	"fmt"
	"os"

	"github.com/ironstar-io/ironstar-cli/cmd/flags"
	"github.com/ironstar-io/ironstar-cli/internal/api"
	"github.com/ironstar-io/ironstar-cli/internal/logs"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// LogsCmd - `iron deploy`
var LogsCmd = &cobra.Command{
	Use:   "logs",
	Short: "Streams logs for an environment",
	Long:  "Streams logs for an environment",
	Run:   streamLogs,
}

func streamLogs(cmd *cobra.Command, args []string) {
	err := logs.Display(args, flags.Acc)
	if err != nil {
		if err != api.ErrIronstarAPICall {
			fmt.Println()
			color.Red(err.Error())
		}

		os.Exit(1)
	}
}

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

// DownloadCmd - `iron backup download`
var DownloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Backup download",
	Long:  "Download a backup",
	Run:   download,
}

func download(cmd *cobra.Command, args []string) {
	err := backup.Download(args, flags.Acc)
	if err != nil {
		if err != api.ErrIronstarAPICall {
			fmt.Println()
			color.Red(err.Error())
		}

		os.Exit(1)
	}
}

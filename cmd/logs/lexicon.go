package logs

import (
	"fmt"
	"os"

	"gitlab.com/ironstar-io/ironstar-cli/internal/api"
	"gitlab.com/ironstar-io/ironstar-cli/internal/logs"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// LexiconCmd - `iron logs lexicon`
var LexiconCmd = &cobra.Command{
	Use:   "lexicon",
	Short: "Display log terms",
	Long:  "Display a list of all log abbreviations and their meaning",
	Run:   lexicon,
}

func lexicon(cmd *cobra.Command, args []string) {
	err := logs.Lexicon()
	if err != nil {
		if err != api.ErrIronstarAPICall {
			fmt.Println()
			color.Red(err.Error())
		}

		os.Exit(1)
	}
}

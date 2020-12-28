package new_relic

import (
	"github.com/spf13/cobra"
)

// NewRelicCmd - `iron new-relic`
var NewRelicCmd = &cobra.Command{
	Use:   "new-relic",
	Short: "",
	Long:  "",
	Run:   configure,
}

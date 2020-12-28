package newrelic

import (
	"github.com/spf13/cobra"
)

// NewRelicCmd - `iron newrelic`
var NewRelicCmd = &cobra.Command{
	Use:   "newrelic",
	Short: "",
	Long:  "",
	Run:   configure,
}

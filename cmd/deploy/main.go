package deploy

import (
	"github.com/spf13/cobra"
)

// DeployCmd - `iron deploy`
var DeployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "",
	Long:  "",
	Run:   list,
}

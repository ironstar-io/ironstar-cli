package environment

import (
	"github.com/spf13/cobra"
)

// EnvironmentCmd - `iron environment`
var EnvironmentCmd = &cobra.Command{
	Use:   "environment",
	Short: "",
	Long:  "",
	Run:   list,
}

// EnvCmd - `iron env`
var EnvCmd = &cobra.Command{
	Hidden: true,
	Use:    "env",
	Short:  "",
	Long:   "",
	Run:    list,
}

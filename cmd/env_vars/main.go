package env_vars

import (
	"github.com/spf13/cobra"
)

// EnvVarsCmd - `iron env-vars`
var EnvVarsCmd = &cobra.Command{
	Hidden: true,
	Use:    "env-vars",
	Short:  "",
	Long:   "",
	Run:    list,
}

package auth

import (
	"fmt"

	"github.com/spf13/cobra"
)

// AuthCmd - `tok auth`
var AuthCmd = &cobra.Command{
	Use:   "auth",
	Short: "",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("TODO - List available 'auth' commands")
	},
}

// MFACmd - `tok auth mfa`
var MFACmd = &cobra.Command{
	Use:   "mfa",
	Short: "",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("TODO - List available 'auth mfa' commands")
	},
}

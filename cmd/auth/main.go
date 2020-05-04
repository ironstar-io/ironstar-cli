package auth

import (
	"fmt"

	"github.com/spf13/cobra"
)

// AuthCmd - `iron auth`
var AuthCmd = &cobra.Command{
	Use:   "auth",
	Short: "",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("TODO - List available 'auth' commands")
	},
}

// MFACmd - `iron auth mfa`
var MFACmd = &cobra.Command{
	Use:   "mfa",
	Short: "",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("TODO - List available 'auth mfa' commands")
	},
}

package subscription

import (
	"fmt"

	"github.com/spf13/cobra"
)

// SubscriptionCmd - `tok subscription`
var SubscriptionCmd = &cobra.Command{
	Use:   "subscription",
	Short: "",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("TODO - List available 'subscription' commands")
	},
}

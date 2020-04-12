package subscription

import (
	"github.com/spf13/cobra"
)

// SubscriptionCmd - `iron subscription`
var SubscriptionCmd = &cobra.Command{
	Use:   "subscription",
	Short: "",
	Long:  "",
	Run:   show,
}

// SubCmd - `iron sub`
var SubCmd = &cobra.Command{
	Hidden: true,
	Use:    "sub",
	Short:  "",
	Long:   "",
	Run:    show,
}

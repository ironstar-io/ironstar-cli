package auth

import (
	"fmt"
	"os"
	"time"

	"gitlab.com/ironstar-io/ironstar-cli/internal/services"

	"github.com/olekukonko/tablewriter"
)

func IronstarShowCredentials(args []string) error {
	credSet, err := services.ReadInCredentials()
	if err != nil {
		return err
	}

	fmt.Println("Active Credentials:")

	active := credSet.Active
	if active == "" {
		active = "UNSET"
	}
	fmt.Println(active)
	fmt.Println()

	fmt.Println("Available Credentials:")

	acreds := make([][]string, len(credSet.Keychain))

	for _, v := range credSet.Keychain {
		acreds = append(acreds, []string{v.Login, v.Expiry.Format(time.RFC3339)})
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Login", "Expiry"})
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.AppendBulk(acreds)
	table.Render()

	return nil
}

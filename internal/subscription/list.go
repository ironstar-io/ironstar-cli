package subscription

import (
	"fmt"
	"os"
	"strings"

	"gitlab.com/ironstar-io/ironstar-cli/cmd/flags"
	"gitlab.com/ironstar-io/ironstar-cli/internal/api"
	"gitlab.com/ironstar-io/ironstar-cli/internal/services"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
)

func List(args []string, flg flags.Accumulator) error {
	creds, err := services.ResolveUserCredentials(flg.Login)
	if err != nil {
		return err
	}

	color.Green("Using login [" + creds.Login + "]")

	uar, err := api.GetUserSubscriptions(creds)
	if err != nil {
		return err
	}

	fmt.Println()
	fmt.Println("Available Subscriptions:")

	uarRows := make([][]string, len(uar))
	for _, access := range uar {
		row := make([][]string, 1)
		row = append(row, []string{access.Subscription.HashedID, access.Subscription.Alias, access.Subscription.ApplicationType, access.Role.Name, strings.Join(access.Role.Permissions, ", ")})
		uarRows = append(row, uarRows...)
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Name", "Application Type", "Role", "Permissions"})
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.AppendBulk(uarRows)
	table.Render()

	return nil
}

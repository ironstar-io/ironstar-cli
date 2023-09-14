package subscription

import (
	"fmt"
	"os"
	"strings"

	"github.com/ironstar-io/ironstar-cli/cmd/flags"
	"github.com/ironstar-io/ironstar-cli/internal/api"
	"github.com/ironstar-io/ironstar-cli/internal/services"
	"github.com/ironstar-io/ironstar-cli/internal/system/utils"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
)

func List(args []string, flg flags.Accumulator) error {
	creds, err := services.ResolveUserCredentials(flg.Login)
	if err != nil {
		return err
	}

	if strings.ToLower(flg.Output) != "json" {
		color.Green("Using login [" + creds.Login + "]")
	}

	uar, err := api.GetUserSubscriptions(creds, flg.Output)
	if err != nil {
		return err
	}

	if strings.ToLower(flg.Output) == "json" {
		utils.PrintInterfaceAsJSON(uar)
		return nil
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

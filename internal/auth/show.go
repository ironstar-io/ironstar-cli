package auth

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/ironstar-io/ironstar-cli/cmd/flags"
	"github.com/ironstar-io/ironstar-cli/internal/services"
	"github.com/ironstar-io/ironstar-cli/internal/system/utils"

	"github.com/olekukonko/tablewriter"
)

func IronstarShowCredentials(args []string, flg flags.Accumulator) error {
	credSet, err := services.ReadInCredentials()
	if err != nil {
		return err
	}

	if strings.ToLower(flg.Output) == "json" {
		utils.PrintInterfaceAsJSON(credSet)
		return nil
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
		// Prepend rows, we want dates ordered oldest to newest
		row := make([][]string, 1)
		row = append(row, []string{v.Login, v.Expiry.Format(time.RFC3339)})
		acreds = append(row, acreds...)
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Login", "Expiry"})
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.AppendBulk(acreds)
	table.Render()

	return nil
}

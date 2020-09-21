package environment

import (
	"fmt"
	"os"

	"gitlab.com/ironstar-io/ironstar-cli/cmd/flags"
	"gitlab.com/ironstar-io/ironstar-cli/internal/api"
	"gitlab.com/ironstar-io/ironstar-cli/internal/services"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
)

func List(args []string, flg flags.Accumulator) error {
	creds, err := services.ResolveUserCredentials(flg.Login)
	if err != nil {
		return err
	}

	sub, err := api.GetSubscriptionContext(creds, flg)
	if err != nil {
		return err
	}

	if sub.Alias == "" {
		return errors.New("No Ironstar subscription has been linked to this project. Have you run `iron subscription link [subscription-name]`")
	}

	color.Green("Using login [" + creds.Login + "] for subscription '" + sub.Alias + "' (" + sub.HashedID + ")")

	envs, err := api.GetSubscriptionEnvironments(creds, sub.HashedID)
	if err != nil {
		return err
	}

	fmt.Println()
	fmt.Println("Available Environments:")

	envRows := make([][]string, len(envs))
	for _, env := range envs {
		class := "Non-Production"
		if env.Class == "cw" {
			class = "Production"
		}

		// Prepend rows, we want dates ordered oldest to newest
		row := make([][]string, 1)
		row = append(row, []string{env.Name, env.DNSName, class, env.RestorePermission})
		envRows = append(row, envRows...)
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Name", "DNS Name", "Class", "Restore Permission"})
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.AppendBulk(envRows)
	table.Render()

	return nil
}

package environment

import (
	"fmt"
	"os"
	"strings"

	"github.com/ironstar-io/ironstar-cli/cmd/flags"
	"github.com/ironstar-io/ironstar-cli/internal/api"
	"github.com/ironstar-io/ironstar-cli/internal/errs"
	"github.com/ironstar-io/ironstar-cli/internal/services"
	"github.com/ironstar-io/ironstar-cli/internal/system/utils"

	"github.com/olekukonko/tablewriter"
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
		return errs.ErrNoSubLink
	}

	utils.PrintCommandContext(flg.Output, creds.Login, sub.Alias, sub.HashedID)

	envs, err := api.GetSubscriptionEnvironments(creds, flg.Output, sub.HashedID)
	if err != nil {
		return err
	}

	if strings.ToLower(flg.Output) == "json" {
		utils.PrintInterfaceAsJSON(envs)
		return nil
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

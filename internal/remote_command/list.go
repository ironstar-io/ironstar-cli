package remote_command

import (
	"fmt"
	"os"
	"time"

	"github.com/ironstar-io/ironstar-cli/cmd/flags"
	"github.com/ironstar-io/ironstar-cli/internal/api"
	"github.com/ironstar-io/ironstar-cli/internal/services"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
)

func List(args []string, flg flags.Accumulator) error {
	creds, err := services.ResolveUserCredentials(flg.Login)
	if err != nil {
		return err
	}

	seCtx, err := api.GetSubscriptionEnvironmentContext(creds, flg)
	if err != nil {
		return err
	}

	if seCtx.Subscription.Alias == "" {
		return errors.New("No Ironstar subscription has been linked to this project. Have you run `iron subscription link [subscription-name]`")
	}

	color.Green("Using login [" + creds.Login + "] for subscription '" + seCtx.Subscription.Alias + "' (" + seCtx.Subscription.HashedID + ")")

	rcs, err := api.GetRemoteCommands(creds, seCtx.Subscription.HashedID, seCtx.Environment.HashedID)
	if err != nil {
		return err
	}

	fmt.Println()
	fmt.Println("Executed remote commands for environment [" + seCtx.Environment.Name + "]:")

	rcRows := make([][]string, len(rcs))
	for _, rc := range rcs {
		// Prepend rows, we want dates ordered oldest to newest
		row := make([][]string, 1)
		row = append(row, []string{rc.Name, rc.Command, rc.Status.Lifecycle, rc.Creator.Name, rc.CreatedAt.Format(time.RFC3339)})
		rcRows = append(row, rcRows...)
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Name", "Command", "Status", "Creator", "Created"})
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.AppendBulk(rcRows)
	table.Render()

	return nil
}

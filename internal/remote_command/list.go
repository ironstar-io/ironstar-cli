package remote_command

import (
	"fmt"
	"os"
	"strings"
	"time"

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

	seCtx, err := api.GetSubscriptionEnvironmentContext(creds, flg)
	if err != nil {
		return err
	}

	if seCtx.Subscription.Alias == "" {
		return errs.ErrNoSubLink
	}

	utils.PrintCommandContext(flg.Output, creds.Login, seCtx.Subscription.Alias, seCtx.Subscription.HashedID)

	rcs, err := api.GetRemoteCommands(creds, flg.Output, seCtx.Subscription.HashedID, seCtx.Environment.HashedID)
	if err != nil {
		return err
	}

	if strings.ToLower(flg.Output) == "json" {
		utils.PrintInterfaceAsJSON(rcs)
		return nil
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

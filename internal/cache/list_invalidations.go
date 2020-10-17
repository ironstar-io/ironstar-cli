package cache

import (
	"fmt"
	"os"
	"time"

	"gitlab.com/ironstar-io/ironstar-cli/cmd/flags"
	"gitlab.com/ironstar-io/ironstar-cli/internal/api"
	"gitlab.com/ironstar-io/ironstar-cli/internal/services"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
)

func ListInvalidations(args []string, flg flags.Accumulator) error {
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

	cis, err := api.GetEnvironmentCacheInvalidations(creds, seCtx.Subscription.HashedID, seCtx.Environment.HashedID)
	if err != nil {
		return err
	}

	if len(cis) == 0 {
		fmt.Println()
		fmt.Println("No cache invalidations found for this environment")

		return nil
	}

	fmt.Println()
	fmt.Println("Recent cache invalidations for environment [" + seCtx.Environment.Name + "]:")

	cacheInvalidationRows := make([][]string, len(cis))
	for _, ci := range cis {
		// Prepend rows, we want dates ordered oldest to newest
		row := make([][]string, 1)
		row = append(row, []string{ci.Name, ci.Status, ci.CreatedAt.Format(time.RFC3339)})
		cacheInvalidationRows = append(row, cacheInvalidationRows...)
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Name", "Status", "Created"})
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.AppendBulk(cacheInvalidationRows)
	table.Render()

	return nil
}

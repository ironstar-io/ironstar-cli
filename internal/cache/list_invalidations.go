package cache

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
		return errs.ErrNoSubLink
	}

	utils.PrintCommandContext(flg.Output, creds.Login, seCtx.Subscription.Alias, seCtx.Subscription.HashedID)

	cis, err := api.GetEnvironmentCacheInvalidations(creds, flg.Output, seCtx.Subscription.HashedID, seCtx.Environment.HashedID)
	if err != nil {
		return err
	}

	if strings.ToLower(flg.Output) == "json" {
		utils.PrintInterfaceAsJSON(cis)
		return nil
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

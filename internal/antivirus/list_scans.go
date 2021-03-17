package antivirus

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"gitlab.com/ironstar-io/ironstar-cli/cmd/flags"
	"gitlab.com/ironstar-io/ironstar-cli/internal/api"
	"gitlab.com/ironstar-io/ironstar-cli/internal/services"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
)

func ListScans(args []string, flg flags.Accumulator) error {
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

	avs, err := api.GetEnvironmentAntivirusScans(creds, seCtx.Subscription.HashedID, seCtx.Environment.HashedID)
	if err != nil {
		return err
	}

	if len(avs) == 0 {
		fmt.Println()
		fmt.Println("No antivirus scans found for this environment")

		return nil
	}

	fmt.Println()
	fmt.Println("Recent antivirus scans for environment [" + seCtx.Environment.Name + "]:")

	antivirusScanRows := make([][]string, len(avs))
	for _, av := range avs {
		// Prepend rows, we want dates ordered oldest to newest
		row := make([][]string, 1)
		row = append(row, []string{av.Result, strconv.Itoa(av.ScannedFiles), strconv.Itoa(av.InfectedFiles), strconv.Itoa(av.Duration), av.CreatedAt.Format(time.RFC3339)})
		antivirusScanRows = append(row, antivirusScanRows...)
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Result", "Scanned Files", "Infected Files", "Duration", "Created"})
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.AppendBulk(antivirusScanRows)
	table.Render()

	return nil
}

package deploy

import (
	"fmt"
	"os"
	"time"

	"gitlab.com/ironstar-io/ironstar-cli/internal/api"
	"gitlab.com/ironstar-io/ironstar-cli/internal/types"

	"github.com/olekukonko/tablewriter"
)

func DisplayDeploymentActivity(creds types.Keylink, deployID string) error {

	dac, err := api.GetDeploymentActivity(creds, deployID)
	if err != nil {
		return err
	}

	fmt.Println()
	fmt.Println("ACTIVITY: ")

	daRows := make([][]string, len(dac))
	for _, da := range dac {
		// Prepend rows, we want dates ordered oldest to newest
		row := make([][]string, 1)
		row = append(row, []string{da.CreatedAt.Format(time.RFC3339), da.Message, da.Flag})
		daRows = append(row, daRows...)
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Date", "Action", "Flag"})
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.AppendBulk(daRows)
	table.Render()

	return nil
}

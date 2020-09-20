package sync

import (
	"fmt"
	"os"
	"strings"
	"time"

	"gitlab.com/ironstar-io/ironstar-cli/cmd/flags"
	"gitlab.com/ironstar-io/ironstar-cli/internal/api"
	"gitlab.com/ironstar-io/ironstar-cli/internal/constants"
	"gitlab.com/ironstar-io/ironstar-cli/internal/services"
	"gitlab.com/ironstar-io/ironstar-cli/internal/types"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
)

func Info(args []string, flg flags.Accumulator) error {
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

	if len(args) > 0 {
		err = DisplayIndividualSyncInfo(creds, sub, args[0])
		if err != nil {
			return err
		}

		return nil
	}

	srs, err := api.GetSubscriptionSyncRequests(creds, sub.HashedID)
	if err != nil {
		return err
	}

	fmt.Println()
	fmt.Println("Recent syncs for subscription [" + sub.Alias + "]:")
	fmt.Println()

	srsRows := make([][]string, len(srs))
	for _, sr := range srs {
		tt := CalcRestoreTimeTaken(sr.Status, sr.CreatedAt, sr.CompletedAt)
		components := strings.Join(sr.BackupRequest.Components, ", ")

		srsRows = append(srsRows, []string{sr.Name, sr.SrcEnvironment.Name, sr.DestEnvironment.Name, sr.Initiator.DisplayName, sr.CreatedAt.Format(time.RFC3339), tt, sr.Status, components})
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Name", "Src Environment", "Dest Environment", "Initiator", "Start Time", "Time Taken", "Status", "Components"})
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.AppendBulk(srsRows)
	table.Render()

	return nil
}

func DisplayIndividualSyncInfo(creds types.Keylink, sub types.Subscription, syncName string) error {
	sr, err := api.GetSubscriptionSync(creds, sub.HashedID, syncName)
	if err != nil {
		return err
	}

	fmt.Println()
	fmt.Println("Name:          " + sr.Name)
	fmt.Println("Status:        " + sr.Status)
	fmt.Println("Initiator:     " + sr.Initiator.DisplayName)
	// fmt.Println("Source Backup: " + CalcBackupIterationName(rr.BackupIteration.ClientName, rr.BackupIteration.Iteration))
	fmt.Println("Started:       " + sr.CreatedAt.Format(time.RFC3339))
	if !sr.CompletedAt.IsZero() {
		fmt.Println("Completed:     " + sr.CompletedAt.Format(time.RFC3339))
		fmt.Println("Duration:      " + CalcRestoreTimeTaken(sr.Status, sr.CreatedAt, sr.CompletedAt))
	}

	fmt.Println()
	fmt.Println("Reservation:")
	fmt.Println(strings.Join(sr.Components, ", "))

	if len(sr.Results) > 0 {
		DisplayComponentInfo(sr.Results)
	}

	return nil
}

func DisplayComponentInfo(components []types.RestoreRequestResult) {
	fmt.Println()

	compRows := make([][]string, len(components))
	for _, comp := range components {
		compRows = append(compRows, []string{comp.Name, comp.Result, comp.CreatedAt.Format(time.RFC3339)})
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Component", "Result", "Time"})
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.AppendBulk(compRows)
	table.Render()
}

func CalcRestoreTimeTaken(status string, createdAt, completedAt time.Time) string {
	if status != constants.BACKUP_COMPLETE || completedAt.IsZero() {
		return time.Since(createdAt).Round(time.Second).String()
	}

	return completedAt.Sub(createdAt).Round(time.Second).String()
}

func CalcBackupIterationName(clientName, iterationName string) string {
	if clientName != "" {
		return clientName
	}

	return iterationName
}

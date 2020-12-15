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
	"gitlab.com/ironstar-io/ironstar-cli/internal/system/utils"
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
		tt := utils.CalcSyncTimeTaken(sr.Status, sr.CreatedAt, sr.CompletedAt)

		// Prepend rows, we want dates ordered oldest to newest
		row := make([][]string, 1)
		row = append(row, []string{sr.Name, sr.SrcEnvironment.Name, sr.DestEnvironment.Name, sr.Initiator.DisplayName, sr.CreatedAt.Format(time.RFC3339), tt, sr.Status})
		srsRows = append(row, srsRows...)
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Name", "Src Environment", "Dest Environment", "Initiator", "Start Time", "Time Taken", "Status"})
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
	fmt.Println("Started:       " + sr.CreatedAt.Format(time.RFC3339))
	if !sr.CompletedAt.IsZero() {
		fmt.Println("Completed:     " + sr.CompletedAt.Format(time.RFC3339))
		fmt.Println("Duration:      " + utils.CalcSyncTimeTaken(sr.Status, sr.CreatedAt, sr.CompletedAt))
	}

	fmt.Println()
	fmt.Println("Source Environment:      " + sr.SrcEnvironment.Name)
	fmt.Println("Destination Environment: " + sr.DestEnvironment.Name)
	fmt.Println()

	fmt.Println("Reservation:")
	fmt.Println(strings.Join(sr.BackupRequest.Components, ", "))

	b, _ := api.GetEnvironmentBackup(creds, sub.HashedID, sr.SrcEnvironment.HashedID, sr.Name, constants.SKIP_ERRORS)

	if b.BackupIteration.Iteration != "" {
		size := utils.CalcBackupSize(b.BackupIteration.Components)

		fmt.Println()
		fmt.Println("------------")
		fmt.Println("BACKUP PHASE")
		fmt.Println("------------")
		fmt.Println("Name: " + b.BackupIteration.Iteration)
		fmt.Println("Status:     " + b.BackupIteration.Status)
		fmt.Println("Started:    " + b.BackupIteration.CreatedAt.Format(time.RFC3339))
		if !b.BackupIteration.CompletedAt.IsZero() {
			fmt.Println("Completed:  " + b.BackupIteration.CompletedAt.Format(time.RFC3339))
			fmt.Println("Duration:   " + utils.CalcBackupTimeTaken(b.BackupIteration.Status, b.BackupIteration.CreatedAt, b.BackupIteration.CompletedAt))
		}
		fmt.Println("Size:       " + size)

		if len(b.BackupIteration.Components) > 0 {
			utils.DisplayBackupComponentInfo(b.BackupIteration.Components)
		}
	}

	if sr.RestoreRequest.Name != "" {
		fmt.Println()
		fmt.Println("------------")
		fmt.Println("RESTORE PHASE")
		fmt.Println("------------")
		fmt.Println("Name:          " + sr.RestoreRequest.Name)
		fmt.Println("Status:        " + sr.RestoreRequest.Status)
		fmt.Println("Started:       " + sr.RestoreRequest.CreatedAt.Format(time.RFC3339))
		if !sr.RestoreRequest.CompletedAt.IsZero() {
			fmt.Println("Completed:     " + sr.RestoreRequest.CompletedAt.Format(time.RFC3339))
			fmt.Println("Duration:      " + utils.CalcRestoreTimeTaken(sr.RestoreRequest.Status, sr.RestoreRequest.CreatedAt, sr.RestoreRequest.CompletedAt))
		}

		if len(sr.RestoreRequest.Results) > 0 {
			utils.DisplayRestoreComponentInfo(sr.RestoreRequest.Results)
		}

		if sr.RestoreRequest.Status == constants.RESTORE_COMPLETE {
			fmt.Println()
			color.Green("SYNC COMPLETE!")
		}
	}

	return nil
}

func CalcBackupIterationName(clientName, iterationName string) string {
	if clientName != "" {
		return clientName
	}

	return iterationName
}

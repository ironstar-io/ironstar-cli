package restore

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/ironstar-io/ironstar-cli/cmd/flags"
	"github.com/ironstar-io/ironstar-cli/internal/api"
	"github.com/ironstar-io/ironstar-cli/internal/constants"
	"github.com/ironstar-io/ironstar-cli/internal/services"
	"github.com/ironstar-io/ironstar-cli/internal/system/utils"
	"github.com/ironstar-io/ironstar-cli/internal/types"

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

	if flg.Environment != "" {
		env, err := api.GetEnvironmentContext(creds, flg, sub.HashedID)
		if err != nil {
			return err
		}

		if len(args) > 0 {
			err = DisplayIndividualRestoreInfo(creds, env, sub, args[0])
			if err != nil {
				return err
			}
		} else {
			err = DisplayEnvironmentRestoreInfo(creds, env, sub)
			if err != nil {
				return err
			}
		}

		return nil
	}

	if len(args) > 0 {
		fmt.Println()
		fmt.Println("If you're looking for information on a specific restore, please specify the environment using the flag `--env=[env-name]`. Exiting...")

		return nil
	}

	ris, err := api.GetSubscriptionRestoreIterations(creds, sub.HashedID)
	if err != nil {
		return err
	}

	fmt.Println()
	fmt.Println("Recent restores for subscription [" + sub.Alias + "]:")
	fmt.Println()

	risRows := make([][]string, len(ris))
	for _, ri := range ris {
		tt := utils.CalcRestoreTimeTaken(ri.Status, ri.CreatedAt, ri.CompletedAt)
		components := CalcRestoreResultNames(ri.Results)
		sbu := CalcBackupIterationName(ri.BackupIteration.ClientName, ri.BackupIteration.Iteration)

		// Prepend rows, we want dates ordered oldest to newest
		row := make([][]string, 1)
		row = append(row, []string{ri.Name, ri.Environment.Name, sbu, ri.Initiator.DisplayName, ri.CreatedAt.Format(time.RFC3339), tt, ri.Status, components})
		risRows = append(row, risRows...)
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Name", "Dest Environment", "Source Backup", "Initiator", "Start Time", "Time Taken", "Status", "Components"})
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.AppendBulk(risRows)
	table.Render()

	return nil
}

func DisplayEnvironmentRestoreInfo(creds types.Keylink, env types.Environment, sub types.Subscription) error {
	ris, err := api.GetEnvironmentRestoreIterations(creds, sub.HashedID, env.HashedID)
	if err != nil {
		return err
	}

	fmt.Println()
	fmt.Println("Recent restores for destination environment [" + env.Name + "]:")
	fmt.Println()

	risRows := make([][]string, len(ris))
	for _, ri := range ris {
		tt := utils.CalcRestoreTimeTaken(ri.Status, ri.CreatedAt, ri.CompletedAt)
		components := CalcRestoreResultNames(ri.Results)
		sbu := CalcBackupIterationName(ri.BackupIteration.ClientName, ri.BackupIteration.Iteration)

		// Prepend rows, we want dates ordered oldest to newest
		risRows = append(risRows, []string{ri.Name, sbu, ri.Initiator.DisplayName, ri.CreatedAt.Format(time.RFC3339), tt, ri.Status, components})
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Name", "Source Backup", "Initiator", "Start Time", "Time Taken", "Status", "Components"})
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.AppendBulk(risRows)
	table.Render()

	return nil
}

func DisplayIndividualRestoreInfo(creds types.Keylink, env types.Environment, sub types.Subscription, restoreName string) error {
	rr, err := api.GetEnvironmentRestore(creds, sub.HashedID, env.HashedID, restoreName)
	if err != nil {
		return err
	}

	fmt.Println()
	fmt.Println("Name:          " + rr.Name)
	fmt.Println("Status:        " + rr.Status)
	fmt.Println("Initiator:     " + rr.Initiator.DisplayName)
	fmt.Println("Source Backup: " + CalcBackupIterationName(rr.BackupIteration.ClientName, rr.BackupIteration.Iteration))
	fmt.Println("Started:       " + rr.CreatedAt.Format(time.RFC3339))
	if !rr.CompletedAt.IsZero() {
		fmt.Println("Completed:     " + rr.CompletedAt.Format(time.RFC3339))
		fmt.Println("Duration:      " + utils.CalcRestoreTimeTaken(rr.Status, rr.CreatedAt, rr.CompletedAt))
	}

	fmt.Println()
	fmt.Println("Reservation:")
	fmt.Println(strings.Join(rr.Components, ", "))

	if len(rr.Results) > 0 {
		DisplayComponentInfo(rr.Results)
	}

	if rr.Status == constants.RESTORE_COMPLETE {
		fmt.Println()
		color.Green("RESTORE COMPLETE!")
	}

	return nil
}

func DisplayComponentInfo(components []types.RestoreRequestResult) {
	fmt.Println()

	compRows := make([][]string, len(components))
	for _, comp := range components {
		// Prepend rows, we want dates ordered oldest to newest
		row := make([][]string, 1)
		row = append(row, []string{comp.Name, comp.Result, comp.CreatedAt.Format(time.RFC3339)})
		compRows = append(row, compRows...)
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Component", "Result", "Time"})
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.AppendBulk(compRows)
	table.Render()
}

func CalcRestoreResultNames(components []types.RestoreRequestResult) string {
	var compNames []string
	for _, comp := range components {
		compNames = append(compNames, comp.Name)
	}

	return strings.Join(compNames, ", ")
}

func CalcBackupIterationName(clientName, iterationName string) string {
	if clientName != "" {
		return clientName
	}

	return iterationName
}

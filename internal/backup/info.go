package backup

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/ironstar-io/ironstar-cli/cmd/flags"
	"github.com/ironstar-io/ironstar-cli/internal/api"
	"github.com/ironstar-io/ironstar-cli/internal/constants"
	"github.com/ironstar-io/ironstar-cli/internal/errs"
	"github.com/ironstar-io/ironstar-cli/internal/services"
	"github.com/ironstar-io/ironstar-cli/internal/system/utils"
	"github.com/ironstar-io/ironstar-cli/internal/types"

	"github.com/ironstar-io/cron"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
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
		return errs.ErrNoSubLink
	}

	utils.PrintCommandContext(flg.Output, creds.Login, sub.Alias, sub.HashedID)

	if flg.Environment != "" {
		env, err := api.GetEnvironmentContext(creds, flg, sub.HashedID)
		if err != nil {
			return err
		}

		if len(args) > 0 {
			err = DisplayIndividualBackupInfo(creds, flg.Output, env, sub, args[0])
			if err != nil {
				return err
			}
		} else {
			err = DisplayEnvironmentBackupInfo(creds, flg.Output, env, sub, flg.BackupType)
			if err != nil {
				return err
			}
		}

		return nil
	}

	if len(args) > 0 {
		fmt.Println()
		fmt.Println("If you're looking for information on a specific backup, please specify the environment using the flag `--env=[env-name]`. Exiting...")

		return nil
	}

	bis, err := api.GetSubscriptionBackupIterations(creds, flg.Output, sub.HashedID, flg.BackupType)
	if err != nil {
		return err
	}

	if strings.ToLower(flg.Output) == "json" {
		utils.PrintInterfaceAsJSON(bis)
		return nil
	}

	fmt.Println()
	fmt.Println("Recent backup runs for subscription [" + sub.Alias + "]:")
	fmt.Println()

	if flg.BackupType != "" {
		fmt.Println("Displaying results for backup type '" + flg.BackupType + "'")
	}

	bisRows := make([][]string, len(bis))
	for _, bi := range bis {
		tt := utils.CalcBackupTimeTaken(bi.Status, bi.CreatedAt, bi.CompletedAt)
		size := utils.CalcBackupSize(bi.Components)
		components := utils.CalcBackupComponentNames(bi.Components)

		// Prepend rows, we want dates ordered oldest to newest
		row := make([][]string, 1)
		row = append(row, []string{bi.BackupRequest.Kind, bi.Iteration, bi.Environment.Name, bi.CreatedAt.Format(time.RFC3339), tt, bi.Status, size, components})
		bisRows = append(row, bisRows...)
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Type", "Name", "Environment", "Start Time", "Time Taken", "Status", "Size", "Components"})
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.AppendBulk(bisRows)
	table.Render()

	return nil
}

func DisplayEnvironmentBackupInfo(creds types.Keylink, output string, env types.Environment, sub types.Subscription, backupType string) error {
	bis, err := api.GetEnvironmentBackupIterations(creds, output, sub.HashedID, env.HashedID, backupType)
	if err != nil {
		return err
	}

	if strings.ToLower(output) == "json" {
		utils.PrintInterfaceAsJSON(bis)
		return nil
	}

	fmt.Println()
	fmt.Println("Recent backup runs for environment [" + env.Name + "]:")
	fmt.Println()

	if backupType != "" {
		fmt.Println("Displaying results for backup type '" + backupType + "'")
	}

	bisRows := make([][]string, len(bis))
	for _, bi := range bis {
		tt := utils.CalcBackupTimeTaken(bi.Status, bi.CreatedAt, bi.CompletedAt)
		size := utils.CalcBackupSize(bi.Components)
		components := utils.CalcBackupComponentNames(bi.Components)

		keepUntil := ""
		if !bi.BackupRequest.KeepUntil.IsZero() {
			keepUntil = bi.BackupRequest.KeepUntil.Format(time.RFC3339)
		}

		// Prepend rows, we want dates ordered oldest to newest
		row := make([][]string, 1)
		row = append(row, []string{bi.BackupRequest.Kind, bi.Iteration, bi.CreatedAt.Format(time.RFC3339), tt, bi.Status, size, components, keepUntil})
		bisRows = append(row, bisRows...)
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Type", "Name", "Start Time", "Time Taken", "Status", "Size", "Components", "Expires After"})
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.AppendBulk(bisRows)
	table.Render()

	return nil
}

func DisplayIndividualBackupInfo(creds types.Keylink, output string, env types.Environment, sub types.Subscription, backupName string) error {
	b, err := api.GetEnvironmentBackup(creds, output, sub.HashedID, env.HashedID, backupName, constants.DISPLAY_ERRORS)
	if err != nil {
		return err
	}

	if strings.ToLower(output) == "json" {
		utils.PrintInterfaceAsJSON(b)
		return nil
	}

	if b.BackupIteration.Iteration != "" {
		size := utils.CalcBackupSize(b.BackupIteration.Components)

		fmt.Println()
		fmt.Println("Type:       " + b.BackupRequest.Kind)
		if b.BackupIteration.ClientName != "" {
			fmt.Println("Name:       " + b.BackupIteration.ClientName)
		} else {
			fmt.Println("Name:       " + b.BackupRequest.Name)
		}
		fmt.Println("Identifier: " + b.BackupIteration.Iteration)
		fmt.Println("Status:     " + b.BackupIteration.Status)
		fmt.Println("Started:    " + b.BackupIteration.CreatedAt.Format(time.RFC3339))
		if !b.BackupIteration.CompletedAt.IsZero() {
			fmt.Println("Completed:  " + b.BackupIteration.CompletedAt.Format(time.RFC3339))
			fmt.Println("Duration:   " + utils.CalcBackupTimeTaken(b.BackupIteration.Status, b.BackupIteration.CreatedAt, b.BackupIteration.CompletedAt))
		}
		fmt.Println("Size:       " + size)

		fmt.Println()
		fmt.Println("Reservation:")
		fmt.Println(strings.Join(b.BackupIteration.Reservation, ", "))

		if len(b.BackupIteration.Components) > 0 {
			utils.DisplayBackupComponentInfo(b.BackupIteration.Components)
		}

		if b.BackupIteration.Status == constants.BACKUP_COMPLETE {
			fmt.Println()
			color.Green("BACKUP COMPLETE!")
		}

		return nil
	}

	fmt.Println()
	fmt.Println("Type:       " + b.BackupRequest.Kind)
	fmt.Println("Name:       " + b.BackupRequest.Name)
	fmt.Println("Components: " + strings.Join(b.BackupRequest.Components, ", "))
	fmt.Println("Created:    " + b.BackupRequest.CreatedAt.Format(time.RFC3339))

	fmt.Println()
	if b.BackupRequest.Kind == "scheduled" {
		fmt.Println("Schedule:           " + b.BackupRequest.Schedule)

		specParser := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
		sched, _ := specParser.Parse(b.BackupRequest.Schedule)

		now := time.Now()
		next := sched.Next(now)

		fmt.Println("Next Scheduled Run: " + next.Format(time.RFC3339))
	} else {
		fmt.Println("A backup run has been registered, and will commence in the Ironstar system shortly")
	}

	return nil
}

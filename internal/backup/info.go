package backup

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"gitlab.com/ironstar-io/ironstar-cli/cmd/flags"
	"gitlab.com/ironstar-io/ironstar-cli/internal/api"
	"gitlab.com/ironstar-io/ironstar-cli/internal/constants"
	"gitlab.com/ironstar-io/ironstar-cli/internal/services"
	"gitlab.com/ironstar-io/ironstar-cli/internal/types"

	"github.com/ironstar-io/cron"

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
			err = DisplayIndividualBackupInfo(creds, env, sub, args[0])
			if err != nil {
				return err
			}
		} else {
			err = DisplayEnvironmentBackupInfo(creds, env, sub)
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

	bis, err := api.GetSubscriptionBackupIterations(creds, sub.HashedID)
	if err != nil {
		return err
	}

	fmt.Println()
	fmt.Println("Recent backup runs for subscription [" + sub.Alias + "]:")
	fmt.Println()

	bisRows := make([][]string, len(bis))
	for _, bi := range bis {
		tt := CalcBackupTimeTaken(bi.Status, bi.CreatedAt, bi.CompletedAt)
		size := CalcBackupSize(bi.Components)
		components := CalcBackupComponentNames(bi.Components)
		bisRows = append(bisRows, []string{bi.BackupRequest.Kind, bi.Iteration, bi.Environment.Name, bi.CreatedAt.Format(time.RFC3339), tt, bi.Status, size, components})
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Type", "Name", "Environment", "Start Time", "Time Taken", "Status", "Size", "Components"})
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.AppendBulk(bisRows)
	table.Render()

	return nil
}

func DisplayEnvironmentBackupInfo(creds types.Keylink, env types.Environment, sub types.Subscription) error {
	bis, err := api.GetEnvironmentBackupIterations(creds, sub.HashedID, env.HashedID)
	if err != nil {
		return err
	}

	fmt.Println()
	fmt.Println("Recent backup runs for environment [" + env.Name + "]:")
	fmt.Println()

	bisRows := make([][]string, len(bis))
	for _, bi := range bis {
		tt := CalcBackupTimeTaken(bi.Status, bi.CreatedAt, bi.CompletedAt)
		size := CalcBackupSize(bi.Components)
		components := CalcBackupComponentNames(bi.Components)
		bisRows = append(bisRows, []string{bi.BackupRequest.Kind, bi.Iteration, bi.CreatedAt.Format(time.RFC3339), tt, bi.Status, size, components})
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Type", "Name", "Start Time", "Time Taken", "Status", "Size", "Components"})
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.AppendBulk(bisRows)
	table.Render()

	return nil
}

func DisplayIndividualBackupInfo(creds types.Keylink, env types.Environment, sub types.Subscription, backupName string) error {
	b, err := api.GetEnvironmentBackup(creds, sub.HashedID, env.HashedID, backupName)
	if err != nil {
		return err
	}

	if b.BackupIteration.Iteration != "" {
		size := CalcBackupSize(b.BackupIteration.Components)

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
			fmt.Println("Duration:   " + CalcBackupTimeTaken(b.BackupIteration.Status, b.BackupIteration.CreatedAt, b.BackupIteration.CompletedAt))
		}
		fmt.Println("Size:       " + size)

		fmt.Println()
		fmt.Println("Reservation:")
		fmt.Println(strings.Join(b.BackupIteration.Reservation, ", "))

		if len(b.BackupIteration.Components) > 0 {
			DisplayComponentInfo(b.BackupIteration.Components)
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

func DisplayComponentInfo(components []types.BackupIterationComponent) {
	fmt.Println()

	compRows := make([][]string, len(components))
	for _, comp := range components {
		dur := (time.Duration(int64(comp.BackupDuration)) * time.Second).Round(time.Second).String()
		compRows = append(compRows, []string{comp.Name, strconv.Itoa(comp.BackupSize) + " MiB", dur, comp.ArchiveKey, comp.Result})
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Component", "Size", "Duration", "Filename", "Result"})
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.AppendBulk(compRows)
	table.Render()
}

func CalcBackupTimeTaken(status string, createdAt, completedAt time.Time) string {
	if status != constants.BACKUP_COMPLETE || completedAt.IsZero() {
		return time.Since(createdAt).Round(time.Second).String()
	}

	return completedAt.Sub(createdAt).Round(time.Second).String()
}

func CalcBackupSize(components []types.BackupIterationComponent) string {
	var size int
	for _, comp := range components {
		size = size + comp.BackupSize
	}

	return strconv.Itoa(size) + " MiB"
}

func CalcBackupComponentNames(components []types.BackupIterationComponent) string {
	var compNames []string
	for _, comp := range components {
		compNames = append(compNames, comp.Name)
	}

	return strings.Join(compNames, ", ")
}

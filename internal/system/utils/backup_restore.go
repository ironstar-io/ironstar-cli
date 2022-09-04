package utils

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/olekukonko/tablewriter"
	"github.com/ironstar-io/ironstar-cli/internal/constants"
	"github.com/ironstar-io/ironstar-cli/internal/types"
)

func DisplayRestoreComponentInfo(components []types.RestoreRequestResult) {
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

func DisplayBackupComponentInfo(components []types.BackupIterationComponent) {
	fmt.Println()

	compRows := make([][]string, len(components))
	for _, comp := range components {
		dur := (time.Duration(int64(comp.BackupDuration)) * time.Second).Round(time.Second).String()

		// Prepend rows, we want dates ordered oldest to newest
		row := make([][]string, 1)
		row = append(row, []string{comp.Name, strconv.Itoa(comp.BackupSize) + " MiB", dur, comp.Result})
		compRows = append(row, compRows...)
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Component", "Size", "Duration", "Result"})
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

func CalcSyncTimeTaken(status string, createdAt, completedAt time.Time) string {
	if status != constants.SYNC_COMPLETE || completedAt.IsZero() {
		return time.Since(createdAt).Round(time.Second).String()
	}

	return completedAt.Sub(createdAt).Round(time.Second).String()
}

func CalcRestoreTimeTaken(status string, createdAt, completedAt time.Time) string {
	if status != constants.RESTORE_COMPLETE || completedAt.IsZero() {
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

func CalculateRestoreComponents(ogComponents []string) []string {
	if len(ogComponents) == 0 {
		return []string{"all"}
	}

	return RemoveStringFromSlice(ogComponents, "logs")
}

func CalculateBackupComponents(ogComponents []string) []string {
	if len(ogComponents) == 0 {
		return []string{"all"}
	}

	return ogComponents
}

func CalculateRestoreStrat(strategy string) string {
	if strategy == "" {
		return "merge"
	}

	return strategy
}

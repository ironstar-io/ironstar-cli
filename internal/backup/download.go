package backup

import (
	"os"
	"path/filepath"
	"strings"

	"gitlab.com/ironstar-io/ironstar-cli/cmd/flags"
	"gitlab.com/ironstar-io/ironstar-cli/internal/api"
	"gitlab.com/ironstar-io/ironstar-cli/internal/constants"
	"gitlab.com/ironstar-io/ironstar-cli/internal/services"
	"gitlab.com/ironstar-io/ironstar-cli/internal/system/utils"
	"gitlab.com/ironstar-io/ironstar-cli/internal/types"

	"github.com/fatih/color"
	"github.com/pkg/errors"
)

func Download(args []string, flg flags.Accumulator) error {
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

	backupName, err := GetBackupName(args)
	if err != nil {
		return err
	}

	reqComps := utils.CalculateBackupComponents(flg.Component)

	b, err := api.GetEnvironmentBackup(creds, seCtx.Subscription.HashedID, seCtx.Environment.HashedID, backupName, constants.DISPLAY_ERRORS)
	if err != nil {
		return err
	}

	dlComps, err := matchDownloadComponents(reqComps, b.BackupIteration.Components)
	if err != nil {
		return err
	}

	for _, dlComp := range dlComps {
		savePath, err := calcSavePath(flg.SavePath, seCtx.Subscription.Alias, seCtx.Environment.Name, dlComp.ArchiveKey)
		if err != nil {
			return err
		}

		err = api.DownloadEnvironmentBackupComponent(creds, seCtx.Subscription.HashedID, seCtx.Environment.HashedID, backupName, savePath, dlComp)
		if err != nil {
			return err
		}
	}

	return nil
}

func matchDownloadComponents(dlComps []string, buComps []types.BackupIterationComponent) ([]types.BackupIterationComponent, error) {
	result := []types.BackupIterationComponent{}

	if utils.SliceIncludes(dlComps, "all") {
		return buComps, nil
	}

	if utils.SliceIncludes(dlComps, "database:all") {
		for _, buComp := range buComps {
			if strings.Contains(buComp.Name, "database:") {
				result = append(result, buComp)
			}
		}
	}

	for _, buComp := range buComps {
		if utils.SliceIncludes(dlComps, buComp.Name) {
			result = append(result, buComp)
		}
	}

	if len(result) == 0 {
		return nil, errors.New("No components for this backup matched your request. Please check the backup and try again.")
	}

	return result, nil
}

func calcSavePath(savePathFlag, subAlias, envName, filename string) (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	if savePathFlag == "" {
		return filepath.Join(wd, ".ironstar", "backups", subAlias, envName, filename), nil
	}

	return filepath.Join(savePathFlag, filename), nil
}

func backupComponentsToList(buComps []types.BackupIterationComponent) []string {
	result := []string{}
	for _, buComp := range buComps {
		result = append(result, buComp.Name)
	}

	return result
}

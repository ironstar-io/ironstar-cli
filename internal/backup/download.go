package backup

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gitlab.com/ironstar-io/ironstar-cli/cmd/flags"
	"gitlab.com/ironstar-io/ironstar-cli/internal/api"
	"gitlab.com/ironstar-io/ironstar-cli/internal/constants"
	"gitlab.com/ironstar-io/ironstar-cli/internal/services"
	"gitlab.com/ironstar-io/ironstar-cli/internal/system/fs"
	"gitlab.com/ironstar-io/ironstar-cli/internal/system/utils"
	"gitlab.com/ironstar-io/ironstar-cli/internal/types"

	"github.com/fatih/color"
	slugify "github.com/metal3d/go-slugify"
	"github.com/pkg/errors"
)

func Download(args []string, flg flags.Accumulator) error {
	creds, err := services.ResolveUserCredentials(flg.Login)
	if err != nil {
		return err
	}

	sCtx, err := api.GetSubscriptionContext(creds, flg)
	if err != nil {
		return err
	}

	if sCtx.Alias == "" {
		return errors.New("No Ironstar subscription has been linked to this project. Have you run `iron subscription link [subscription-name]`")
	}

	color.Green("Using login [" + creds.Login + "] for subscription '" + sCtx.Alias + "' (" + sCtx.HashedID + ")")

	backupName, err := GetBackupName(args, flg.Backup)
	if err != nil {
		return err
	}

	reqComps := utils.CalculateBackupComponents(flg.Component)

	b, err := api.GetSubscriptionBackup(creds, sCtx.HashedID, backupName, constants.DISPLAY_ERRORS)
	if err != nil {
		return err
	}

	dlComps, err := matchDownloadComponents(reqComps, b.BackupIteration.Components)
	if err != nil {
		return err
	}

	savePath, err := calcSavePath(flg.SavePath, sCtx.Alias, b.BackupIteration.EnvironmentName, backupName)
	if err != nil {
		return err
	}

	fmt.Println()

	confirmDL := services.ConfirmationPrompt("Backup downloads are included in your outbound bandwidth allocation each month.\n\nAre you sure you want to download this backup?", "y", flg.AutoAccept)
	if !confirmDL {
		fmt.Println("Exiting...")
		return nil
	}

	fmt.Println()

	for _, dlComp := range dlComps {
		file := calcFilename(savePath, dlComp.Name)

		err = api.DownloadEnvironmentBackupComponent(creds, sCtx.HashedID, b.BackupIteration.EnvironmentName, backupName, file, dlComp)
		if err != nil {
			return err
		}
	}

	fmt.Println()
	color.Green("Download completed and saved to " + savePath)

	return nil
}

func matchDownloadComponents(reqComps []string, buComps []types.BackupIterationComponent) ([]types.BackupIterationComponent, error) {
	result := []types.BackupIterationComponent{}

	if utils.SliceIncludes(reqComps, "all") {
		return buComps, nil
	}

	if utils.SliceIncludes(reqComps, "database:all") {
		for _, buComp := range buComps {
			if strings.Contains(buComp.Name, "database:") {
				result = append(result, buComp)
			}
		}
	}

	for _, reqComp := range reqComps {
		for _, buComp := range buComps {
			if buComp.Name == reqComp || buComp.Name == "path:"+reqComp {
				result = append(result, buComp)
			}
		}
	}

	if len(result) == 0 {
		return nil, errors.New("No components for this backup matched your request. Please check the backup and try again.")
	}

	return result, nil
}

func calcSavePath(savePathFlag, subAlias, envName, backupName string) (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	if savePathFlag == "" {
		path := filepath.Join(home, ".ironstar", "backups", subAlias, envName, backupName)

		fs.Mkdir(path)

		return path, nil
	}

	fs.Mkdir(savePathFlag)

	return savePathFlag, nil
}

func calcFilename(savePath, compName string) string {
	safeComp := slugify.Marshal(strings.ReplaceAll(compName, ":", "-"))

	if strings.Contains(safeComp, "database") {
		return filepath.Join(savePath, safeComp+".sql.gz")
	}

	return filepath.Join(savePath, safeComp+".tar.gz")
}

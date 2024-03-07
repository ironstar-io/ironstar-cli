package backup

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ironstar-io/ironstar-cli/cmd/flags"
	"github.com/ironstar-io/ironstar-cli/internal/api"
	"github.com/ironstar-io/ironstar-cli/internal/constants"
	"github.com/ironstar-io/ironstar-cli/internal/errs"
	"github.com/ironstar-io/ironstar-cli/internal/services"
	"github.com/ironstar-io/ironstar-cli/internal/system/fs"
	"github.com/ironstar-io/ironstar-cli/internal/system/utils"
	"github.com/ironstar-io/ironstar-cli/internal/types"

	"github.com/fatih/color"
	"github.com/jinzhu/copier"
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
		return errs.ErrNoSubLink
	}

	utils.PrintCommandContext(flg.Output, creds.Login, sCtx.Alias, sCtx.HashedID)

	reqComps := utils.CalculateBackupComponents(flg.Component)

	b, err := getTargetBackup(args, flg, creds, sCtx.HashedID)
	if err != nil {
		return err
	}

	dlComps, err := matchDownloadComponents(reqComps, b.Components)
	if err != nil {
		return err
	}

	savePath, err := calcSavePath(flg.SavePath, sCtx.Alias, b.EnvironmentName, b.Iteration)
	if err != nil {
		return err
	}

	fmt.Println()

	confirmDL := services.ConfirmationPrompt(fmt.Sprintf("Backup downloads are included in your outbound bandwidth allocation each month.\n\nAre you sure you want to download the backup '%s'?", b.Iteration), "y", flg.AutoAccept)
	if !confirmDL {
		fmt.Println("Exiting...")
		return nil
	}

	fmt.Println()

	for _, dlComp := range dlComps {
		file := calcFilename(savePath, dlComp.Name)

		err = api.DownloadEnvironmentBackupComponent(creds, flg.Output, sCtx.HashedID, b.EnvironmentName, b.Iteration, file, dlComp)
		if err != nil {
			return err
		}
	}

	fmt.Println()
	color.Green("Download completed and saved to " + savePath)

	return nil
}

func getTargetBackup(args []string, flg flags.Accumulator, creds types.Keylink, subHashedId string) (*types.BackupIteration, error) {
	if flg.Latest {
		env, err := api.GetEnvironmentContext(creds, flg, subHashedId)
		if err != nil {
			return nil, err
		}

		b, err := api.GetLatestEnvironmentBackupIteration(creds, flg.Output, subHashedId, env.HashedID)
		if err != nil {
			return nil, err
		}
		b.EnvironmentName = env.Name

		return &b, nil
	}

	backupName, err := GetBackupName(args, flg.Backup)
	if err != nil {
		return nil, err
	}

	b, err := api.GetSubscriptionBackup(creds, flg.Output, subHashedId, backupName, constants.DISPLAY_ERRORS)
	if err != nil {
		return nil, err
	}

	nbr := &types.BackupIteration{}
	err = copier.Copy(nbr, b.BackupIteration)
	if err != nil {
		return nil, err
	}
	nbr.BackupRequest = b.BackupRequest

	return nbr, nil
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

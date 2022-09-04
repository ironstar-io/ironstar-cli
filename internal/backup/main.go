package backup

import (
	"errors"
	"strings"

	"github.com/ironstar-io/ironstar-cli/internal/services"
	"github.com/ironstar-io/ironstar-cli/internal/system/utils"
	"github.com/ironstar-io/ironstar-cli/internal/types"
)

func GetBackupName(args []string, backupName string) (string, error) {
	if backupName != "" {
		return backupName, nil
	}

	if len(args) > 0 {
		return args[0], nil
	}

	input, err := services.StdinPrompt("Backup Name: ")
	if err != nil {
		return "", err
	}

	return input, nil

}

func MatchBackupComponents(reqComps []string, buComps []types.BackupIterationComponent) ([]string, error) {
	buCompList := BackupComponentsToList(buComps)

	result := []string{}

	if utils.SliceIncludes(reqComps, "all") {
		return buCompList, nil
	}

	if utils.SliceIncludes(reqComps, "database:all") {
		for _, buComp := range buComps {
			if strings.Contains(buComp.Name, "database:") {
				result = append(result, buComp.Name)
			}
		}
	}

	for _, reqComp := range reqComps {
		if utils.SliceIncludes(buCompList, reqComp) {
			result = append(result, reqComp)
		}
		if utils.SliceIncludes(buCompList, "path:"+reqComp) {
			result = append(result, "path:"+reqComp)
		}
	}

	if len(result) == 0 {
		return nil, errors.New("No components for this backup matched your request. Please check the backup and try again.")
	}

	return result, nil
}

func BackupComponentsToList(buComps []types.BackupIterationComponent) []string {
	result := []string{}
	for _, buComp := range buComps {
		result = append(result, buComp.Name)
	}

	return result
}

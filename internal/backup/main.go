package backup

import (
	"errors"
	"strings"

	"gitlab.com/ironstar-io/ironstar-cli/internal/services"
	"gitlab.com/ironstar-io/ironstar-cli/internal/system/utils"
	"gitlab.com/ironstar-io/ironstar-cli/internal/types"
)

func GetBackupName(args []string) (string, error) {
	var name string
	if len(args) == 0 {
		input, err := services.StdinPrompt("Backup Name: ")
		if err != nil {
			return "", err
		}
		name = input
	} else {
		name = args[0]
	}

	return name, nil
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

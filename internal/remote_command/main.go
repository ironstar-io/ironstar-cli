package remote_command

import (
	"strings"

	"github.com/ironstar-io/ironstar-cli/internal/types"
)

func envVarKeyValue(envVars []string) []types.RemoteCommandEnvironmentVariable {
	if len(envVars) == 0 {
		return nil
	}

	result := []types.RemoteCommandEnvironmentVariable{}
	for _, ev := range envVars {
		kv := strings.Split(ev, "=")
		if len(kv) <= 1 {
			continue
		}

		result = append(result, types.RemoteCommandEnvironmentVariable{
			Key:   kv[0],
			Value: kv[1],
		})
	}

	return result
}

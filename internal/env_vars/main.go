package env_vars

import (
	"github.com/ironstar-io/ironstar-cli/cmd/flags"
	"github.com/ironstar-io/ironstar-cli/internal/services"

	"github.com/fatih/color"
)

func PullEnvVarKey(flg flags.Accumulator) string {
	if flg.Key == "" {
		key, err := services.StdinPrompt("Environment variable key: ")
		if err != nil {
			return ""
		}

		return key
	}

	return flg.Key
}

func PullEnvVarValue(flg flags.Accumulator) string {
	if flg.Value == "" {
		val, err := services.StdinPrompt("Environment variable value: ")
		if err != nil {
			return ""
		}

		return val
	}

	color.Yellow("Warning: Supplying a secret via a command line flag is potentially insecure")

	return flg.Value
}

func PullEnvVarVarType(flg flags.Accumulator) string {
	if flg.VarType == "" {
		return "PROTECTED"
	}

	return flg.VarType
}

package env_vars

import (
	"fmt"

	"github.com/ironstar-io/ironstar-cli/cmd/flags"
	"github.com/ironstar-io/ironstar-cli/internal/api"
	"github.com/ironstar-io/ironstar-cli/internal/services"

	"github.com/fatih/color"
	"github.com/pkg/errors"
)

func Add(args []string, flg flags.Accumulator) error {
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

	key := PullEnvVarKey(flg)
	if key == "" {
		return errors.New("No key provided for the new environment variable. This can be specified with `--key={KEY}` or through the prompt.")
	}

	value := PullEnvVarValue(flg)
	if value == "" {
		return errors.New("No value provided for the new environment variable. This can be specified with `--value={VALUE}` or through the prompt.")
	}

	varType := PullEnvVarVarType(flg)
	if varType != "PROTECTED" && varType != "VISIBLE" {
		return errors.New("Invalid variable type provided for the new environment variable. This can be specified with `--var-type={VALUE}` and must be one of 'PROTECTED' or 'VISIBLE'.")
	}

	err = api.PostEnvironmentEnvVars(creds, seCtx.Subscription.HashedID, seCtx.Environment.HashedID, key, value, varType)
	if err != nil {
		return err
	}

	fmt.Println()
	fmt.Println("Added environment variable for environment [" + seCtx.Environment.Name + "]")
	fmt.Println()
	color.Green("Key: " + key + " (" + varType + ")")

	return nil
}

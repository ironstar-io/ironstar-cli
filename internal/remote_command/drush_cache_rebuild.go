package remote_command

import (
	"fmt"
	"time"

	"github.com/ironstar-io/ironstar-cli/cmd/flags"
	"github.com/ironstar-io/ironstar-cli/internal/api"
	"github.com/ironstar-io/ironstar-cli/internal/services"

	"github.com/fatih/color"
	"github.com/pkg/errors"
)

func DrushCacheRebuild(args []string, flg flags.Accumulator) error {
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

	var timeout int
	if flg.Timeout != 0 {
		timeout = flg.Timeout
	}

	rc, err := api.PostRemoteCommandDrushCacheRebuild(creds, seCtx.Subscription.HashedID, seCtx.Environment.HashedID, envVarKeyValue(flg.EnvironmentVars), timeout)
	if err != nil {
		return err
	}

	fmt.Println()
	fmt.Println("Triggered remote command for environment [" + seCtx.Environment.Name + "]")
	fmt.Println()
	color.Green("ID: " + rc.HashedId)
	color.Green("Name: " + rc.Name)
	color.Green("Command: " + rc.Command)
	color.Green(fmt.Sprintf("Timeout: %s", time.Duration(rc.Timeout)*time.Second))
	color.Green(fmt.Sprintf("Environment Variables: %s", rc.EnvironmentVariables))

	return nil
}

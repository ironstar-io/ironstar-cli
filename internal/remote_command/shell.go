package remote_command

import (
	"fmt"
	"strings"
	"time"

	"github.com/ironstar-io/ironstar-cli/cmd/flags"
	"github.com/ironstar-io/ironstar-cli/internal/api"
	"github.com/ironstar-io/ironstar-cli/internal/errs"
	"github.com/ironstar-io/ironstar-cli/internal/services"
	"github.com/ironstar-io/ironstar-cli/internal/system/utils"

	"github.com/fatih/color"
)

func Shell(args []string, flg flags.Accumulator) error {
	creds, err := services.ResolveUserCredentials(flg.Login)
	if err != nil {
		return err
	}

	seCtx, err := api.GetSubscriptionEnvironmentContext(creds, flg)
	if err != nil {
		return err
	}

	if seCtx.Subscription.Alias == "" {
		return errs.ErrNoSubLink
	}

	utils.PrintCommandContext(flg.Output, creds.Login, seCtx.Subscription.Alias, seCtx.Subscription.HashedID)

	var timeout int
	if flg.Timeout != 0 {
		timeout = flg.Timeout
	}

	var workDir string
	if flg.WorkDir != "" {
		workDir = flg.WorkDir
	}

	rc, err := api.PostRemoteCommandShell(creds, flg.Output, seCtx.Subscription.HashedID, seCtx.Environment.HashedID, workDir, strings.Join(args, " "), envVarKeyValue(flg.EnvironmentVars), timeout)
	if err != nil {
		return err
	}

	if strings.ToLower(flg.Output) == "json" {
		utils.PrintInterfaceAsJSON(rc)
		return nil
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

package environment

import (
	"fmt"
	"strings"

	"github.com/ironstar-io/ironstar-cli/cmd/flags"
	"github.com/ironstar-io/ironstar-cli/internal/api"
	"github.com/ironstar-io/ironstar-cli/internal/constants"
	"github.com/ironstar-io/ironstar-cli/internal/errs"
	"github.com/ironstar-io/ironstar-cli/internal/services"
	"github.com/ironstar-io/ironstar-cli/internal/system/utils"

	"github.com/fatih/color"
)

func EnableDisableRestore(args []string, flg flags.Accumulator, restoreState string) error {
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

	err = api.PatchEnvironment(creds, flg.Output, seCtx.Subscription.HashedID, seCtx.Environment.HashedID, restoreState)
	if err != nil {
		return err
	}

	if restoreState == constants.RestorePermissionAllowed {
		if strings.ToLower(flg.Output) == "json" {
			utils.PrintInterfaceAsJSON(map[string]string{
				"result": "success",
				"info":   "Successfully enabled restores to the environment '" + seCtx.Environment.Name + "'",
			})
			return nil
		}

		fmt.Println()
		color.Green("Successfully enabled restores to the environment '" + seCtx.Environment.Name + "'")
	}
	if restoreState == constants.RestorePermissionNotAllowed {
		if strings.ToLower(flg.Output) == "json" {
			utils.PrintInterfaceAsJSON(map[string]string{
				"result": "success",
				"info":   "Successfully disabled restores to the environment '" + seCtx.Environment.Name + "'",
			})
			return nil
		}

		fmt.Println()
		color.Green("Successfully disabled restores to the environment '" + seCtx.Environment.Name + "'")
	}

	return nil
}

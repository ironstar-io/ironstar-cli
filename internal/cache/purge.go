package cache

import (
	"fmt"
	"strings"

	"github.com/ironstar-io/ironstar-cli/cmd/flags"
	"github.com/ironstar-io/ironstar-cli/internal/api"
	"github.com/ironstar-io/ironstar-cli/internal/errs"
	"github.com/ironstar-io/ironstar-cli/internal/services"
	"github.com/ironstar-io/ironstar-cli/internal/system/utils"

	"github.com/fatih/color"
)

func Purge(args []string, flg flags.Accumulator) error {
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

	ci, err := api.PostEnvironmentCacheInvalidation(creds, flg.Output, seCtx.Subscription.HashedID, seCtx.Environment.HashedID)
	if err != nil {
		return err
	}

	if strings.ToLower(flg.Output) == "json" {
		utils.PrintInterfaceAsJSON(ci)
		return nil
	}

	fmt.Println()
	color.Green("Cache purge has commenced. To see an up-to-date status please run `iron cache invalidation show " + ci.Name + " --subscription=" + seCtx.Subscription.Alias + " --environment=" + seCtx.Environment.Name + "`")

	return nil
}

package cache

import (
	"fmt"
	"strings"

	"github.com/ironstar-io/ironstar-cli/cmd/flags"
	"github.com/ironstar-io/ironstar-cli/internal/api"
	"github.com/ironstar-io/ironstar-cli/internal/errs"
	"github.com/ironstar-io/ironstar-cli/internal/services"
	"github.com/ironstar-io/ironstar-cli/internal/system/utils"
)

func ShowInvalidation(args []string, flg flags.Accumulator) error {
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

	name, err := GetCacheInvalidationName(args)
	if err != nil {
		return err
	}

	ci, err := api.GetEnvironmentCacheInvalidation(creds, flg.Output, seCtx.Subscription.HashedID, seCtx.Environment.HashedID, name)
	if err != nil {
		return err
	}

	if strings.ToLower(flg.Output) == "json" {
		utils.PrintInterfaceAsJSON(ci)
		return nil
	}

	fmt.Println()
	fmt.Println("Cache invalidation '" + name + "' for environment [" + seCtx.Environment.Name + "]: " + ci.Status)

	return nil
}

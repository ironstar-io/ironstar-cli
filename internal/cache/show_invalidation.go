package cache

import (
	"fmt"

	"github.com/ironstar-io/ironstar-cli/cmd/flags"
	"github.com/ironstar-io/ironstar-cli/internal/api"
	"github.com/ironstar-io/ironstar-cli/internal/services"

	"github.com/fatih/color"
	"github.com/pkg/errors"
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
		return errors.New("No Ironstar subscription has been linked to this project. Have you run `iron subscription link [subscription-name]`")
	}

	color.Green("Using login [" + creds.Login + "] for subscription '" + seCtx.Subscription.Alias + "' (" + seCtx.Subscription.HashedID + ")")

	name, err := GetCacheInvalidationName(args)
	if err != nil {
		return err
	}

	ci, err := api.GetEnvironmentCacheInvalidation(creds, seCtx.Subscription.HashedID, seCtx.Environment.HashedID, name)
	if err != nil {
		return err
	}

	fmt.Println()
	fmt.Println("Cache invalidation '" + name + "' for environment [" + seCtx.Environment.Name + "]: " + ci.Status)

	return nil
}

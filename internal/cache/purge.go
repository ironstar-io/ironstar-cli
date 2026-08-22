package cache

import (
	"fmt"
	"strings"

	"github.com/ironstar-io/ironstar-cli/cmd/flags"
	"github.com/ironstar-io/ironstar-cli/internal/api"
	"github.com/ironstar-io/ironstar-cli/internal/errs"
	"github.com/ironstar-io/ironstar-cli/internal/services"
	"github.com/ironstar-io/ironstar-cli/internal/system/utils"
	"github.com/ironstar-io/ironstar-cli/internal/types"

	"github.com/fatih/color"
	"github.com/pkg/errors"
)

func Purge(args []string, flg flags.Accumulator) error {
	urls, err := normalizeCachePurgeURLs(flg.URLs)
	if err != nil {
		return err
	}

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

	targets := urls
	if len(targets) == 0 {
		targets = []string{""}
	}

	invalidations := make([]types.CacheInvalidation, 0, len(targets))
	for _, url := range targets {
		ci, err := api.PostEnvironmentCacheInvalidation(creds, flg.Output, seCtx.Subscription.HashedID, seCtx.Environment.HashedID, url)
		if err != nil {
			if url != "" {
				if err == api.ErrIronstarAPICall {
					return err
				}
				return errors.Wrapf(err, "Failed to purge URL %q", url)
			}
			return err
		}
		invalidations = append(invalidations, ci)
	}

	if strings.ToLower(flg.Output) == "json" {
		if len(invalidations) == 1 {
			utils.PrintInterfaceAsJSON(invalidations[0])
		} else {
			utils.PrintInterfaceAsJSON(invalidations)
		}
		return nil
	}

	for index, ci := range invalidations {
		fmt.Println()
		if len(urls) > 0 {
			color.Green("Cache purge has commenced for " + urls[index] + ". To see an up-to-date status please run `iron cache invalidation show " + ci.Name + " --subscription=" + seCtx.Subscription.Alias + " --environment=" + seCtx.Environment.Name + "`")
		} else {
			color.Green("Cache purge has commenced. To see an up-to-date status please run `iron cache invalidation show " + ci.Name + " --subscription=" + seCtx.Subscription.Alias + " --environment=" + seCtx.Environment.Name + "`")
		}
	}

	return nil
}

func normalizeCachePurgeURLs(values []string) ([]string, error) {
	urls := make([]string, 0, len(values))
	for _, value := range values {
		url := strings.TrimSpace(value)
		if url == "" {
			return nil, errors.New("URL must not be empty")
		}
		urls = append(urls, url)
	}
	return urls, nil
}

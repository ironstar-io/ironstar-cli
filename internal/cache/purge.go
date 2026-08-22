package cache

import (
	"fmt"
	"net/url"
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
	purgeURL, err := normalizeCachePurgeURL(flg.URLs)
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

	request := makeCacheInvalidationRequest(purgeURL)
	ci, err := api.PostEnvironmentCacheInvalidation(creds, flg.Output, seCtx.Subscription.HashedID, seCtx.Environment.HashedID, request)
	if err != nil {
		if purgeURL != "" {
			if err == api.ErrIronstarAPICall {
				return err
			}
			return errors.Wrapf(err, "Failed to purge URL %q", purgeURL)
		}
		return err
	}

	if strings.ToLower(flg.Output) == "json" {
		utils.PrintInterfaceAsJSON(ci)
		return nil
	}

	fmt.Println()
	if purgeURL != "" {
		color.Green("Cache purge has commenced for " + purgeURL + ". To see an up-to-date status please run `iron cache invalidation show " + ci.Name + " --subscription=" + seCtx.Subscription.Alias + " --environment=" + seCtx.Environment.Name + "`")
	} else {
		color.Green("Cache purge has commenced. To see an up-to-date status please run `iron cache invalidation show " + ci.Name + " --subscription=" + seCtx.Subscription.Alias + " --environment=" + seCtx.Environment.Name + "`")
	}

	return nil
}

func normalizeCachePurgeURL(values []string) (string, error) {
	if len(values) == 0 {
		return "", nil
	}
	if len(values) > 1 {
		return "", errors.New("only one URL may be purged at a time")
	}

	rawURL := strings.TrimSpace(values[0])
	if rawURL == "" {
		return "", errors.New("URL must not be empty")
	}

	parsedURL, err := url.Parse(rawURL)
	if err != nil || !parsedURL.IsAbs() || !strings.EqualFold(parsedURL.Scheme, "https") || parsedURL.Hostname() == "" {
		return "", errors.Errorf("URL %q must be an absolute HTTPS URL", rawURL)
	}
	if parsedURL.User != nil {
		return "", errors.Errorf("URL %q must not include user information", rawURL)
	}
	if parsedURL.Fragment != "" {
		return "", errors.Errorf("URL %q must not include a fragment", rawURL)
	}

	return rawURL, nil
}

func makeCacheInvalidationRequest(purgeURL string) types.PostCacheInvalidationRequestParams {
	if purgeURL == "" {
		return types.PostCacheInvalidationRequestParams{
			Kind:             types.CacheInvalidationKindEnvironment,
			InvalidationType: types.CacheInvalidationTypeHard,
		}
	}

	return types.PostCacheInvalidationRequestParams{
		Kind:             types.CacheInvalidationKindURL,
		InvalidationType: types.CacheInvalidationTypeSoft,
		URL:              purgeURL,
	}
}

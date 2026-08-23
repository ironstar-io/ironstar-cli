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

func Invalidate(args []string, flg flags.Accumulator) error {
	invalidationURL, err := normalizeCacheInvalidationURL(flg.URLs)
	if err != nil {
		return err
	}

	request, err := makeCacheInvalidationRequest(invalidationURL, flg.Type)
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

	if invalidationURL == "" {
		confirmed := services.ConfirmationPrompt(
			"No URL was supplied. This will hard-invalidate the entire cache for environment '"+seCtx.Environment.Name+"'. Continue?",
			"n",
			flg.AutoAccept,
		)
		if !confirmed {
			fmt.Println("No cache invalidation requested.")
			return nil
		}
	}

	ci, err := api.PostEnvironmentCacheInvalidation(creds, flg.Output, seCtx.Subscription.HashedID, seCtx.Environment.HashedID, request)
	if err != nil {
		if invalidationURL != "" {
			if err == api.ErrIronstarAPICall {
				return err
			}
			return errors.Wrapf(err, "Failed to invalidate URL %q", invalidationURL)
		}
		return err
	}

	if strings.ToLower(flg.Output) == "json" {
		utils.PrintInterfaceAsJSON(ci)
		return nil
	}

	fmt.Println()
	if invalidationURL != "" {
		color.Green("Cache invalidation has commenced for " + invalidationURL + ". To see its current status, run `iron cache show " + ci.Name + " --subscription=" + seCtx.Subscription.Alias + " --environment=" + seCtx.Environment.Name + "`")
		fmt.Println()
		color.Yellow("Note: A later FAILED status does not necessarily mean this request was not submitted. Fastly may report failure when the URL was not found in cache.")
	} else {
		color.Green("Environment cache invalidation has commenced. To see its current status, run `iron cache show " + ci.Name + " --subscription=" + seCtx.Subscription.Alias + " --environment=" + seCtx.Environment.Name + "`")
	}

	return nil
}

func normalizeCacheInvalidationURL(values []string) (string, error) {
	if len(values) == 0 {
		return "", nil
	}
	if len(values) > 1 {
		return "", errors.New("only one URL may be invalidated at a time")
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

func makeCacheInvalidationRequest(invalidationURL, requestedType string) (types.PostCacheInvalidationRequestParams, error) {
	invalidationType := strings.ToLower(strings.TrimSpace(requestedType))
	if invalidationType == "" {
		if invalidationURL == "" {
			invalidationType = string(types.CacheInvalidationTypeHard)
		} else {
			invalidationType = string(types.CacheInvalidationTypeSoft)
		}
	}

	if invalidationType != string(types.CacheInvalidationTypeSoft) && invalidationType != string(types.CacheInvalidationTypeHard) {
		return types.PostCacheInvalidationRequestParams{}, errors.New("invalidation type must be 'soft' or 'hard'")
	}

	if invalidationURL == "" {
		if invalidationType == string(types.CacheInvalidationTypeSoft) {
			return types.PostCacheInvalidationRequestParams{}, errors.New("whole-environment invalidations must use type 'hard'")
		}
		return types.PostCacheInvalidationRequestParams{
			Kind:             types.CacheInvalidationKindEnvironment,
			InvalidationType: types.CacheInvalidationTypeHard,
		}, nil
	}

	return types.PostCacheInvalidationRequestParams{
		Kind:             types.CacheInvalidationKindURL,
		InvalidationType: types.CacheInvalidationType(invalidationType),
		URL:              invalidationURL,
	}, nil
}

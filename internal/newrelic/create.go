package newrelic

import (
	"fmt"

	"gitlab.com/ironstar-io/ironstar-cli/cmd/flags"
	"gitlab.com/ironstar-io/ironstar-cli/internal/api"
	"gitlab.com/ironstar-io/ironstar-cli/internal/services"
	"gitlab.com/ironstar-io/ironstar-cli/internal/types"

	"github.com/fatih/color"
	"github.com/pkg/errors"
)

func Configure(args []string, flg flags.Accumulator) error {
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

	fmt.Println()

	appName, err := services.StdinPrompt("App Name: ")
	if err != nil {
		return errors.New("App name must be supplied. Exiting...")
	}

	appID, err := services.StdinPrompt("App ID: ")
	if err != nil {
		return errors.New("App ID must be supplied. Exiting...")
	}

	apiKey, err := services.StdinPrompt("API Key: ")
	if err != nil {
		return errors.New("API Key must be supplied. Exiting...")
	}

	licenseKey, err := services.StdinPrompt("License Key: ")
	if err != nil {
		return errors.New("License Key must be supplied. Exiting...")
	}

	err = api.PutNewRelicApplicationConfig(creds, seCtx.Subscription.HashedID, seCtx.Environment.HashedID, types.PutNewRelicParams{
		LicenseKey:  licenseKey,
		AppName:     appName,
		AppID:       appID,
		APIKeyValue: apiKey,
		APIKeyType:  "",
	})
	if err != nil {
		return err
	}

	fmt.Println()
	fmt.Println("It may take several minutes for these changes to reflect in your New Relic dashboard")
	fmt.Println()
	color.Green("Completed successfully!")

	return nil
}

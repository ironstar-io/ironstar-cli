package environment

import (
	"fmt"

	"github.com/ironstar-io/ironstar-cli/cmd/flags"
	"github.com/ironstar-io/ironstar-cli/internal/api"
	"github.com/ironstar-io/ironstar-cli/internal/constants"
	"github.com/ironstar-io/ironstar-cli/internal/services"

	"github.com/fatih/color"
	"github.com/pkg/errors"
)

func AddHook(args []string, flg flags.Accumulator) error {
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

	hookName, err := GetHookName(args)
	if err != nil {
		return err
	}

	err = api.PostEnvironmentHook(creds, seCtx.Subscription.HashedID, seCtx.Environment.HashedID, hookName)
	if err != nil {
		return err
	}

	fmt.Println()
	color.Green("The hook '" + hookName + "' has successfully been added to your environment")

	if hookName == constants.PRE_DEPLOYMENT_BACKUP {
		fmt.Println()
		color.Yellow("Please note: There is limit of 20 backups per subscription that applies. When deploying to your environments with this hook enabled, automatically provisioned backups will contribute to this total.")
	}

	return nil
}

func RemoveHook(args []string, flg flags.Accumulator) error {
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

	hookName, err := GetHookName(args)
	if err != nil {
		return err
	}

	confirmRemove := services.ConfirmationPrompt("Are you sure you would like to remove the hook '"+hookName+"' from your environment?", "y", flg.AutoAccept)
	if !confirmRemove {
		fmt.Println("Exiting...")
		return nil
	}

	err = api.DeleteEnvironmentHook(creds, seCtx.Subscription.HashedID, seCtx.Environment.HashedID, hookName)
	if err != nil {
		return err
	}

	fmt.Println()
	color.Green("The hook '" + hookName + "' has successfully been removed from your environment")

	return nil
}

func GetHookName(args []string) (string, error) {
	var name string
	if len(args) == 0 {
		input, err := services.StdinPrompt("Hook Name: ")
		if err != nil {
			return "", err
		}
		name = input
	} else {
		name = args[0]
	}

	return name, nil
}

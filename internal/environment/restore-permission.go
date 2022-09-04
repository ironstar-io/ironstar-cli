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
		return errors.New("No Ironstar subscription has been linked to this project. Have you run `iron subscription link [subscription-name]`")
	}

	color.Green("Using login [" + creds.Login + "] for subscription '" + seCtx.Subscription.Alias + "' (" + seCtx.Subscription.HashedID + ")")

	err = api.PatchEnvironment(creds, seCtx.Subscription.HashedID, seCtx.Environment.HashedID, restoreState)
	if err != nil {
		return err
	}

	if restoreState == constants.RestorePermissionAllowed {
		fmt.Println()
		color.Green("Successfully enabled restores to the environment '" + seCtx.Environment.Name + "'")
	}
	if restoreState == constants.RestorePermissionNotAllowed {
		fmt.Println()
		color.Green("Successfully disabled restores to the environment '" + seCtx.Environment.Name + "'")
	}

	return nil
}

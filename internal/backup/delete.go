package backup

import (
	"fmt"

	"gitlab.com/ironstar-io/ironstar-cli/cmd/flags"
	"gitlab.com/ironstar-io/ironstar-cli/internal/api"
	"gitlab.com/ironstar-io/ironstar-cli/internal/services"
	"gitlab.com/ironstar-io/ironstar-cli/internal/types"

	"github.com/fatih/color"
	"github.com/pkg/errors"
)

func Delete(args []string, flg flags.Accumulator) error {
	if len(args) == 0 && flg.Name == "" {
		return errors.New("Please specify the name of the backup you'd like to delete with an argument, or the --name flag ie. `iron backup delete [backup-name]` or `iron backup delete --name=[backup-name]`")
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
		return errors.New("No Ironstar subscription has been linked to this project. Have you run `iron subscription link [subscription-name]`")
	}

	color.Green("Using login [" + creds.Login + "] for subscription '" + seCtx.Subscription.Alias + "' (" + seCtx.Subscription.HashedID + ")")

	name := flg.Name
	if name == "" {
		name = args[0]
	}

	err = api.DeleteBackupIteration(creds, types.DeleteBackupIterationParams{
		SubscriptionID: seCtx.Subscription.HashedID,
		EnvironmentID:  seCtx.Environment.HashedID,
		Name:           name,
	})
	if err != nil {
		return err
	}

	fmt.Println()
	color.Green("Successfully deleted backup of the environment '" + seCtx.Environment.Name + "' named [" + name + "]")

	return nil
}

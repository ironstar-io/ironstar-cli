package backup

import (
	"fmt"

	"github.com/ironstar-io/ironstar-cli/cmd/flags"
	"github.com/ironstar-io/ironstar-cli/internal/api"
	"github.com/ironstar-io/ironstar-cli/internal/services"
	"github.com/ironstar-io/ironstar-cli/internal/types"

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

	seCtx, err := api.GetSubscriptionContext(creds, flg)
	if err != nil {
		return err
	}

	if seCtx.Alias == "" {
		return errors.New("No Ironstar subscription has been linked to this project. Have you run `iron subscription link [subscription-name]`")
	}

	color.Green("Using login [" + creds.Login + "] for subscription '" + seCtx.Alias + "' (" + seCtx.HashedID + ")")

	name := flg.Name
	if name == "" {
		name = args[0]
	}

	err = api.DeleteBackup(creds, types.DeleteBackupParams{
		SubscriptionID: seCtx.HashedID,
		Name:           name,
	})
	if err != nil {
		return err
	}

	fmt.Println()
	color.Green("Successfully deleted backup [" + name + "]")

	return nil
}

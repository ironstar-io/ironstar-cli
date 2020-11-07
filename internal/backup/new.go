package backup

import (
	"fmt"

	"gitlab.com/ironstar-io/ironstar-cli/cmd/flags"
	"gitlab.com/ironstar-io/ironstar-cli/internal/api"
	"gitlab.com/ironstar-io/ironstar-cli/internal/services"
	"gitlab.com/ironstar-io/ironstar-cli/internal/system/utils"
	"gitlab.com/ironstar-io/ironstar-cli/internal/types"

	"github.com/fatih/color"
	"github.com/pkg/errors"
)

func New(args []string, flg flags.Accumulator) error {
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
	components := utils.CalculateBackupComponents(flg.Component)

	br, err := api.PostBackupRequest(creds, types.PostBackupRequestParams{
		SubscriptionID: seCtx.Subscription.HashedID,
		EnvironmentID:  seCtx.Environment.HashedID,
		Name:           name,
		Kind:           "manual",
		LockTables:     flg.LockTables,
		Components:     components,
	})
	if err != nil {
		return err
	}

	fmt.Println()
	fmt.Println("Creating a manual backup of environment [" + seCtx.Environment.Name + "] for subscription [" + seCtx.Subscription.Alias + "] named [" + br.Name + "]")
	fmt.Println()
	fmt.Println("The following components will be backed up:")
	for _, comp := range br.Components {
		fmt.Println("- " + comp)
	}

	if br.ETA != 0 {
		fETA := utils.CalculateFriendlyETA(br.ETA)
		fmt.Println()
		fmt.Println("This backup will take approximately " + fETA + " (based on previous similar backups) to complete")
	}

	fmt.Println()
	fmt.Println("You can check the status at any time by running `iron backup info " + br.Name + " --env=" + seCtx.Environment.Name + " --subscription=" + seCtx.Subscription.Alias + "`")
	fmt.Println()

	color.Green("Successfully commenced backup of the environment '" + seCtx.Environment.Name + "'")

	return nil
}

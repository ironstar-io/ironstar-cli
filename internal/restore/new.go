package restore

import (
	"fmt"

	"github.com/ironstar-io/ironstar-cli/cmd/flags"
	"github.com/ironstar-io/ironstar-cli/internal/api"
	"github.com/ironstar-io/ironstar-cli/internal/backup"
	"github.com/ironstar-io/ironstar-cli/internal/constants"
	"github.com/ironstar-io/ironstar-cli/internal/services"
	"github.com/ironstar-io/ironstar-cli/internal/system/utils"
	"github.com/ironstar-io/ironstar-cli/internal/types"

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

	if flg.Backup == "" {
		return errors.New("A source backup must be specified with the --backup=[backup-name] flag")
	}

	b, err := api.GetSubscriptionBackup(creds, seCtx.Subscription.HashedID, flg.Backup, constants.DISPLAY_ERRORS)
	if err != nil {
		return err
	}

	reqComps := utils.CalculateRestoreComponents(flg.Component)
	if len(reqComps) == 0 {
		return errors.New("At least one component must be specified with the --component=[component-name] flag")
	}
	components, err := backup.MatchBackupComponents(reqComps, b.BackupIteration.Components)
	if err != nil {
		return err
	}

	color.Green("Using login [" + creds.Login + "] for subscription '" + seCtx.Subscription.Alias + "' (" + seCtx.Subscription.HashedID + ")")

	name := flg.Name
	strategy := utils.CalculateRestoreStrat(flg.Strategy)

	rr, err := api.PostRestoreRequest(creds, types.PostRestoreRequestParams{
		SubscriptionID: seCtx.Subscription.HashedID,
		EnvironmentID:  seCtx.Environment.HashedID,
		Name:           name,
		Strategy:       strategy,
		Backup:         flg.Backup,
		Components:     components,
	})
	if err != nil {
		return err
	}

	fmt.Println()
	fmt.Println("Creating a restore to environment [" + seCtx.Environment.Name + "] from backup [" + flg.Backup + "] named [" + rr.Name + "]")
	fmt.Println()
	fmt.Println("The following components will be restored:")
	for _, comp := range rr.Components {
		fmt.Println("- " + comp)
	}

	if rr.ETA != 0 {
		fETA := utils.CalculateFriendlyETA(rr.ETA)
		fmt.Println()
		fmt.Println("This restore will take approximately " + fETA + " to complete")
	}

	fmt.Println()
	fmt.Println("You can check the status at any time by running `iron restore info " + rr.Name + " --env=" + seCtx.Environment.Name + " --subscription=" + seCtx.Subscription.Alias + "`")
	fmt.Println()

	color.Green("Successfully commenced restore")

	return nil
}

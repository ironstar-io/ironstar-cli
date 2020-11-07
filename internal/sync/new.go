package sync

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

	sub, err := api.GetSubscriptionContext(creds, flg)
	if err != nil {
		return err
	}

	if sub.Alias == "" {
		return errors.New("No Ironstar subscription has been linked to this project. Have you run `iron subscription link [subscription-name]`")
	}

	// Check supplied components
	components := utils.CalculateRestoreComponents(flg.Component)
	if len(components) == 0 {
		return errors.New("At least one component must be specified with the --component=[component-name] flag")
	}

	// Check and pull source/destination environments
	if flg.SrcEnvironment == "" {
		return errors.New("The source environment for the sync must be specified with the --src-env=[environment-name] flag")
	}
	srcEnv, err := api.GetSubscriptionEnvironment(creds, sub.HashedID, flg.SrcEnvironment)
	if err != nil {
		return err
	}
	if flg.DestEnvironment == "" {
		return errors.New("The destination environment for the sync must be specified with the --dest-env=[environment-name] flag")
	}
	destEnv, err := api.GetSubscriptionEnvironment(creds, sub.HashedID, flg.DestEnvironment)
	if err != nil {
		return err
	}

	if flg.SrcEnvironment == flg.DestEnvironment {
		return errors.New("Cannot sync between the same environment.")
	}

	color.Green("Using login [" + creds.Login + "] for subscription '" + sub.Alias + "' (" + sub.HashedID + ")")

	strategy := utils.CalculateRestoreStrat(flg.Strategy)

	if flg.UseLatestBackup == true {
		err = RestoreFromLatestBackup(creds, flg, sub, srcEnv, destEnv, components, strategy)
		if err != nil {
			return err
		}

		return nil
	}

	sr, err := api.PostSyncRequest(creds, types.PostSyncRequestParams{
		SubscriptionID:  sub.HashedID,
		SrcEnvironment:  srcEnv.HashedID,
		DestEnvironment: destEnv.HashedID,
		RestoreStrategy: strategy,
		Components:      components,
	})
	if err != nil {
		return err
	}

	fmt.Println()
	fmt.Println("Creating a sync from environment [" + srcEnv.Name + "] to environment [" + destEnv.Name + "] named [" + sr.Name + "]")
	fmt.Println()

	if sr.ETA != 0 {
		fETA := utils.CalculateFriendlyETA(sr.ETA)
		fmt.Println()
		fmt.Println("This sync will take approximately " + fETA + " to complete")
	}

	fmt.Println()
	fmt.Println("You can check the status at any time by running `iron sync info " + sr.Name + "`")

	return nil
}

func RestoreFromLatestBackup(creds types.Keylink, flg flags.Accumulator, sub types.Subscription, srcEnv, destEnv types.Environment, components []string, strategy string) error {
	rr, err := api.PostSyncRequestUseLatestBackup(creds, types.PostSyncRequestParams{
		SubscriptionID:  sub.HashedID,
		RestoreStrategy: strategy,
		SrcEnvironment:  srcEnv.HashedID,
		DestEnvironment: destEnv.HashedID,
		Components:      components,
	})
	if err != nil {
		return err
	}

	fmt.Println()
	fmt.Println("Creating a restore to environment [" + destEnv.Name + "] from backup run [" + rr.BackupIteration.Iteration + "] named [" + rr.Name + "]")
	fmt.Println()
	fmt.Println("The backup portion of the sync was skipped due to the `--use-latest-backup` flag being set")
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
	fmt.Println("You can check the status at any time by running `iron restore info " + rr.Name + " --env=" + destEnv.Name + "`")
	fmt.Println()

	color.Green("Successfully commenced restore")

	return nil
}

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
	components := CalculateSyncRestoreComponents(flg.Component)
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

	strategy := CalculateRestoreStrat(flg.Strategy)

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
		fmt.Println("This sync will take approximately " + fETA + "to complete")
	}

	fmt.Println()
	fmt.Println("You can check the status at any time by running `iron sync info " + sr.Name + "`")

	return nil
}

func CalculateSyncRestoreComponents(ogComponents []string) []string {
	if len(ogComponents) == 0 {
		return []string{"all"}
	}

	return utils.RemoveStringFromSlice(ogComponents, "logs")
}

func CalculateRestoreStrat(strategy string) string {
	if strategy == "" {
		return "merge"
	}

	return strategy
}

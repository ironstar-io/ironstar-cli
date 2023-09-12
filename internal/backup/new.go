package backup

import (
	"fmt"
	"strings"

	"github.com/ironstar-io/ironstar-cli/cmd/flags"
	"github.com/ironstar-io/ironstar-cli/internal/api"
	"github.com/ironstar-io/ironstar-cli/internal/errs"
	"github.com/ironstar-io/ironstar-cli/internal/services"
	"github.com/ironstar-io/ironstar-cli/internal/system/utils"
	"github.com/ironstar-io/ironstar-cli/internal/types"

	"github.com/fatih/color"
	slugify "github.com/metal3d/go-slugify"
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
		return errs.ErrNoSubLink
	}

	utils.PrintCommandContext(flg.Output, creds.Login, seCtx.Subscription.Alias, seCtx.Subscription.HashedID)

	name := slugify.Marshal(flg.Name, false)
	if name != flg.Name && strings.ToLower(flg.Output) != "json" {
		fmt.Println()
		color.Yellow("NOTE: '" + flg.Name + "' is not a valid name for a backup in the Ironstar system. We have slugified this value for you to be '" + name + "'")
		fmt.Println()
	}

	components := utils.CalculateBackupComponents(flg.Component)

	br, err := api.PostBackupRequest(creds, flg.Output, types.PostBackupRequestParams{
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

	if strings.ToLower(flg.Output) == "json" {
		utils.PrintInterfaceAsJSON(br)
		return nil
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

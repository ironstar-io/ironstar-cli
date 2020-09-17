package backup

import (
	"fmt"
	"strconv"

	"gitlab.com/ironstar-io/ironstar-cli/cmd/flags"
	"gitlab.com/ironstar-io/ironstar-cli/internal/api"
	"gitlab.com/ironstar-io/ironstar-cli/internal/services"
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
	components := CalculatePostBackupComponents(flg.Component)
	kind := CalculateBackupRequestKind(flg.Type)

	br, err := api.PostBackupRequest(creds, types.PostBackupRequestParams{
		SubscriptionID: seCtx.Subscription.HashedID,
		EnvironmentID:  seCtx.Environment.HashedID,
		Name:           name,
		Kind:           kind,
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

	fmt.Println()

	if br.ETA != 0 {
		fmt.Println("This backup will take approximately " + strconv.Itoa(br.ETA) + " minutes (based on previous similar backups) to complete")
	}

	fmt.Println()
	fmt.Println("You can check the status at any time by running `iron backup status`")
	fmt.Println()

	color.Green("Successfully commenced backup of the environment '" + seCtx.Environment.Name + "'")

	return nil
}

func CalculatePostBackupComponents(ogComponents []string) []string {
	if len(ogComponents) == 0 {
		return []string{"all"}
	}

	return ogComponents
}

func CalculateBackupRequestKind(ogKind string) string {
	if ogKind == "scheduled" {
		return "scheduled"
	}

	return "manual"
}

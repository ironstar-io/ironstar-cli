package remote_command

import (
	"fmt"
	"time"

	"github.com/ironstar-io/ironstar-cli/cmd/flags"
	"github.com/ironstar-io/ironstar-cli/internal/api"
	"github.com/ironstar-io/ironstar-cli/internal/services"

	"github.com/fatih/color"
	"github.com/pkg/errors"
)

func Status(args []string, flg flags.Accumulator) error {
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

	rcni, err := services.GetRemoteCommandID(args)
	if err != nil {
		return err
	}

	rc, err := api.GetRemoteCommand(creds, seCtx.Subscription.HashedID, seCtx.Environment.HashedID, rcni)
	if err != nil {
		return err
	}

	fmt.Println()
	fmt.Println("Remote command for environment [" + seCtx.Environment.Name + "]")
	fmt.Println()
	fmt.Println("ID: " + rc.HashedId)
	fmt.Println("Name: " + rc.Name)
	fmt.Println("Command: " + rc.Command)

	fmt.Println("Status: " + rc.Status.Lifecycle)
	if rc.Status.Error != "" {
		color.Red("Error: " + rc.Status.Error)
	}
	fmt.Printf("Timeout: %s\n", time.Duration(rc.Timeout)*time.Second)
	fmt.Printf("Environment Variables: %s\n", rc.EnvironmentVariables)

	fmt.Printf("Creator: %s (%s)\n", rc.Creator.Name, rc.Creator.Kind)
	fmt.Printf("Created At: %s\n", rc.CreatedAt.Format(time.RFC3339))

	return nil
}

package remote_command

import (
	"fmt"
	"strings"
	"time"

	"github.com/ironstar-io/ironstar-cli/cmd/flags"
	"github.com/ironstar-io/ironstar-cli/internal/api"
	"github.com/ironstar-io/ironstar-cli/internal/errs"
	"github.com/ironstar-io/ironstar-cli/internal/services"
	"github.com/ironstar-io/ironstar-cli/internal/system/utils"

	"github.com/fatih/color"
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
		return errs.ErrNoSubLink
	}

	utils.PrintCommandContext(flg.Output, creds.Login, seCtx.Subscription.Alias, seCtx.Subscription.HashedID)

	rcni, err := services.GetRemoteCommandID(args)
	if err != nil {
		return err
	}

	rc, err := api.GetRemoteCommand(creds, flg.Output, seCtx.Subscription.HashedID, seCtx.Environment.HashedID, rcni)
	if err != nil {
		return err
	}

	if strings.ToLower(flg.Output) == "json" {
		utils.PrintInterfaceAsJSON(rc)
		return nil
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
	fmt.Println()
	fmt.Printf("Creator: %s (%s)\n", rc.Creator.Name, rc.Creator.Kind)
	fmt.Printf("Created At: %s\n", rc.CreatedAt.Format(time.RFC3339))

	return nil
}

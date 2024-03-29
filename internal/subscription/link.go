package subscription

import (
	"strings"

	"github.com/ironstar-io/ironstar-cli/cmd/flags"
	"github.com/ironstar-io/ironstar-cli/internal/api"
	"github.com/ironstar-io/ironstar-cli/internal/services"

	"github.com/fatih/color"
)

func Link(args []string, flg flags.Accumulator) error {
	proj, err := services.GetProjectData(flg.AutoAccept)
	if err != nil {
		return err
	}

	creds, err := services.ResolveUserCredentials(flg.Login)
	if err != nil {
		return err
	}

	if strings.ToLower(flg.Output) != "json" {
		color.Green("Using login [" + creds.Login + "]")
	}

	var hashOrAlias string
	if len(args) == 0 {
		ha, err := services.StdinPrompt("Subscription ID or Alias: ")
		if err != nil {
			return err
		}

		hashOrAlias = ha
	} else {
		hashOrAlias = args[0]
	}

	sub, err := api.GetSubscription(creds, flg.Output, hashOrAlias)
	if err != nil {
		return err
	}

	return services.LinkSubscriptionToProject(proj, sub)
}

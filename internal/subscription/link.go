package subscription

import (
	"gitlab.com/ironstar-io/ironstar-cli/internal/api"
	"gitlab.com/ironstar-io/ironstar-cli/internal/services"

	"github.com/fatih/color"
)

func Link(args []string, loginFlag string) error {
	proj, err := services.GetProjectData()
	if err != nil {
		return err
	}

	creds, err := services.ResolveUserCredentials(loginFlag)
	if err != nil {
		return err
	}

	color.Green("Using login [" + creds.Login + "]")

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

	sub, err := api.GetSubscription(creds, hashOrAlias)
	if err != nil {
		return err
	}

	return services.LinkSubscriptionToProject(proj, sub)
}

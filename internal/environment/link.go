package environment

import (
	"errors"

	"gitlab.com/ironstar-io/ironstar-cli/cmd/flags"
	"gitlab.com/ironstar-io/ironstar-cli/internal/api"
	"gitlab.com/ironstar-io/ironstar-cli/internal/services"

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

	sub, err := api.GetSubscriptionContext(creds, flg)
	if err != nil {
		return err
	}

	if sub.Alias == "" {
		return errors.New("No Ironstar subscription has been linked to this project. Have you run `iron subscription link [subscription-name]`")
	}

	color.Green("Using login [" + creds.Login + "] for subscription '" + sub.Alias + "' (" + sub.HashedID + ")")

	var envHashOrAlias string
	if len(args) == 0 {
		ha, err := services.StdinPrompt("Environment Name: ")
		if err != nil {
			return err
		}

		envHashOrAlias = ha
	} else {
		envHashOrAlias = args[0]
	}

	env, err := api.GetSubscriptionEnvironment(creds, sub.HashedID, envHashOrAlias)
	if err != nil {
		return err
	}

	return services.LinkEnvironmentToProject(proj, env)
}

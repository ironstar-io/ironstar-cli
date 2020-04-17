package deploy

import (
	"fmt"

	"gitlab.com/ironstar-io/ironstar-cli/cmd/flags"
	"gitlab.com/ironstar-io/ironstar-cli/internal/api"
	"gitlab.com/ironstar-io/ironstar-cli/internal/errs"
	"gitlab.com/ironstar-io/ironstar-cli/internal/services"
	"gitlab.com/ironstar-io/ironstar-cli/internal/types"

	"github.com/fatih/color"
	"github.com/pkg/errors"
	yaml "gopkg.in/yaml.v2"
)

func Create(args []string, flg flags.Accumulator) error {
	creds, err := services.ResolveUserCredentials(flg.Login)
	if err != nil {
		return err
	}

	sub, err := api.GetSubscriptionContext(creds, flg.Subscription)
	if err != nil {
		return err
	}

	if sub.Alias == "" {
		return errors.New("No Ironstar subscription has been linked to this project. Have you run `iron subscription link [subscription-name]`")
	}

	if flg.Output == "" {
		color.Green("Using login [" + creds.Login + "] for subscription " + sub.Alias + " (" + sub.HashedID + ")")
	}

	var envID string
	if flg.Environment == "" {
		ei, err := services.StdinPrompt("Environment ID: ")
		if err != nil {
			return errors.New("No environment ID argument supplied")
		}
		envID = ei
	}

	var packageID string
	if flg.Package == "" {
		pi, err := services.StdinPrompt("Package ID: ")
		if err != nil {
			return errors.New("No package ID argument supplied")
		}
		packageID = pi
	}

	req := &api.Request{
		RunTokenRefresh:  true,
		Credentials:      creds,
		Method:           "POST",
		Path:             "/build/" + packageID + "/deploy",
		MapStringPayload: map[string]string{"environmentName": envID},
	}

	res, err := req.Send()
	if err != nil {
		return errors.Wrap(err, errs.APISubListErrorMsg)
	}

	if res.StatusCode != 201 {
		return res.HandleFailure()
	}

	if flg.Output == "json" {
		err = services.OutputJSON(res.Body)
		if err != nil {
			return errors.Wrap(err, errs.APISubListErrorMsg)
		}

		return nil
	}

	var d types.Deployment
	err = yaml.Unmarshal(res.Body, &d)
	if err != nil {
		return err
	}

	fmt.Println()
	color.Green("Completed successfully!")
	fmt.Println()

	fmt.Println("DEPLOYMENT ID: " + d.HashedID)
	fmt.Println("PACKAGE ID: " + d.BuildID)
	fmt.Println("ENVIRONMENT: " + envID)
	fmt.Println("APPLICATION STATUS: " + d.AppStatus)
	fmt.Println("ADMIN SERVICE STATUS: " + d.AdminSvcStatus)
	fmt.Println("CREATED: " + d.CreatedAt.String())

	return nil
}

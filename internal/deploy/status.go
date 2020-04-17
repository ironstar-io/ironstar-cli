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

func Status(args []string, flg flags.Accumulator) error {
	creds, err := services.ResolveUserCredentials(flg.Login)
	if err != nil {
		return err
	}

	deployID, err := services.GetDeployID(args, flg.Deploy)
	if err != nil {
		return err
	}

	if flg.Output == "" {
		color.Green("Using login [" + creds.Login + "]")
	}

	err = DisplayDeploymentInfo(creds, deployID)
	if err != nil {
		return err
	}

	return DisplayDeploymentActivity(creds, deployID)
}

func DisplayDeploymentInfo(creds types.Keylink, deployID string) error {
	req := &api.Request{
		RunTokenRefresh:  true,
		Credentials:      creds,
		Method:           "GET",
		Path:             "/deployment/" + deployID,
		MapStringPayload: map[string]string{},
	}

	res, err := req.Send()
	if err != nil {
		return errors.Wrap(err, errs.APISubListErrorMsg)
	}

	if res.StatusCode != 200 {
		return res.HandleFailure()
	}

	var d types.Deployment
	err = yaml.Unmarshal(res.Body, &d)
	if err != nil {
		return err
	}

	fmt.Println()
	fmt.Println("DEPLOYMENT ID: " + d.HashedID)
	fmt.Println("PACKAGE ID: " + d.BuildID)
	fmt.Println("ENVIRONMENT: " + d.Environment.Name)
	fmt.Println("APPLICATION STATUS: " + d.AppStatus)
	fmt.Println("ADMIN SERVICE STATUS: " + d.AdminSvcStatus)
	fmt.Println("CREATED: " + d.CreatedAt.String())

	return nil
}

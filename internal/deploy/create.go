package deploy

import (
	"fmt"
	"time"

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

	sub, err := api.GetSubscriptionContext(creds, flg)
	if err != nil {
		return err
	}

	if sub.Alias == "" {
		return errors.New("No Ironstar subscription has been linked to this project. Have you run `iron subscription link [subscription-name]`")
	}

	color.Green("Using login [" + creds.Login + "] for subscription '" + sub.Alias + "' (" + sub.HashedID + ")")

	var envID string
	if flg.Environment == "" {
		ei, err := services.StdinPrompt("Environment ID: ")
		if err != nil {
			return errors.New("No environment ID argument supplied")
		}
		envID = ei
	} else {
		envID = flg.Environment
	}

	err = checkOperatingEnvironment(flg, creds, sub.HashedID, envID)
	if err != nil {
		return err
	}

	packageID, err := determinePackageSelection(args, flg, creds, sub.HashedID)
	if err != nil {
		return err
	}

	req := &api.Request{
		RunTokenRefresh:  true,
		Credentials:      creds,
		Method:           "POST",
		Path:             "/build/" + packageID + "/deploy",
		MapStringPayload: map[string]string{"environmentName": envID},
	}

	res, err := req.NankaiSend()
	if err != nil {
		return errors.Wrap(err, errs.APISubListErrorMsg)
	}

	if res.StatusCode != 201 {
		return res.HandleFailure()
	}

	var d types.Deployment
	err = yaml.Unmarshal(res.Body, &d)
	if err != nil {
		return err
	}

	fmt.Println()
	color.Green("Completed successfully!")
	fmt.Println()
	fmt.Println("DEPLOYMENT: " + d.Name)
	fmt.Println()
	fmt.Println("ENVIRONMENT: " + envID)
	fmt.Println("PACKAGE ID: " + d.Build.Name)
	fmt.Println("APPLICATION STATUS: " + d.AppStatus)
	fmt.Println("ADMIN SERVICE STATUS: " + d.AdminSvcStatus)
	fmt.Println("CREATED: " + d.CreatedAt.Format(time.RFC3339))
	fmt.Println()
	color.Green("You can now run 'iron deploy status " + d.Name + "' to check deployment status")

	return nil
}

func checkOperatingEnvironment(flg flags.Accumulator, creds types.Keylink, subID, envID string) error {
	req := &api.Request{
		RunTokenRefresh: true,
		Credentials:     creds,
		Method:          "GET",
		Path:            "/subscription/" + subID + "/environment/" + envID,
	}

	res, err := req.NankaiSend()
	if err != nil {
		return errors.Wrap(err, errs.APISubListErrorMsg)
	}

	if res.StatusCode != 200 {
		return res.HandleFailure()
	}

	var env types.Environment
	err = yaml.Unmarshal(res.Body, &env)
	if err != nil {
		return err
	}

	if env.Class == "cw" {
		deployToProd := services.ConfirmationPrompt("Environment '"+env.Name+"' is a production grade environment. Are you sure you would like to continue?", "n", flg.ApproveProdDeploy)
		if deployToProd {
			if !flg.ApproveProdDeploy {
				fmt.Println("")
				color.Yellow("This confirmation prompt can be skipped with the flag '--approve-prod-deploy'")
				fmt.Println("")
			}
			return nil
		}

		return errors.New("Deployment rejected by user")
	}

	return nil
}

func determinePackageSelection(args []string, flg flags.Accumulator, creds types.Keylink, subHash string) (string, error) {
	var empty string
	if flg.Package != "" {
		return flg.Package, nil
	}

	if len(args) != 0 {
		return args[0], nil
	}

	createNew := services.ConfirmationPrompt("No package specified. Would you like to create one?", "y", flg.AutoAccept)
	if createNew {
		tarpath, err := services.CreateProjectTar(flg.Exclude)
		if err != nil {
			return empty, err
		}

		res, err := api.UploadPackage(creds, subHash, tarpath, flg.Ref)
		if err != nil {
			return empty, err
		}

		var ur types.UploadResponse
		err = yaml.Unmarshal(res.Body, &ur)
		if err != nil {
			return empty, err
		}

		fmt.Println("PACKAGE ID: " + ur.BuildID)
		fmt.Println("PACKAGE NAME: " + ur.BuildName)
		fmt.Println()
		color.Green("Continuing to deployment...")

		return ur.BuildID, nil
	}

	pi, err := services.StdinPrompt("Package ID: ")
	if err != nil {
		return empty, errors.New("No package idenitifer supplied")
	}

	return pi, nil
}

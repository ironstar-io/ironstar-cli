package pkg

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
		color.Green("Using login [" + creds.Login + "] for subscription <" + sub.Alias + ">")
	}

	tarpath, err := services.CreateProjectTar(flg.Exclude)
	if err != nil {
		return err
	}

	res, err := api.UploadPackage(creds, sub.HashedID, tarpath)
	if err != nil {
		return err
	}

	if flg.Output == "json" {
		err = services.OutputJSON(res.Body)
		if err != nil {
			return errors.Wrap(err, errs.APISubListErrorMsg)
		}

		return nil
	}

	var ur types.UploadResponse
	err = yaml.Unmarshal(res.Body, &ur)
	if err != nil {
		return err
	}

	fmt.Println()
	color.Green("Package Successfully Created!")
	fmt.Println()

	fmt.Println("PACKAGE ID: " + ur.BuildID)
	fmt.Println("PACKAGE NAME: " + ur.BuildName)

	return nil
}

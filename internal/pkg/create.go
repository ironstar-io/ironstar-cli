package pkg

import (
	"fmt"

	"gitlab.com/ironstar-io/ironstar-cli/cmd/flags"
	"gitlab.com/ironstar-io/ironstar-cli/internal/api"
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

	color.Green("Using login [" + creds.Login + "] for subscription <" + sub.Alias + ">")

	tarpath, err := services.CreateProjectTar(flg)
	if err != nil {
		return err
	}

	res, err := api.UploadPackage(creds, sub.HashedID, tarpath, flg)
	if err != nil {
		return err
	}

	var ur types.UploadResponse
	err = yaml.Unmarshal(res.Body, &ur)
	if err != nil {
		return err
	}

	fmt.Println("PACKAGE ID: " + ur.BuildID)
	fmt.Println("PACKAGE NAME: " + ur.BuildName)

	return nil
}

package pkg

import (
	"fmt"
	"strings"

	"github.com/ironstar-io/ironstar-cli/cmd/flags"
	"github.com/ironstar-io/ironstar-cli/internal/api"
	"github.com/ironstar-io/ironstar-cli/internal/errs"
	"github.com/ironstar-io/ironstar-cli/internal/services"
	"github.com/ironstar-io/ironstar-cli/internal/system/utils"

	"github.com/fatih/color"
	"github.com/pkg/errors"
)

func UpdateRef(args []string, flg flags.Accumulator) error {
	creds, err := services.ResolveUserCredentials(flg.Login)
	if err != nil {
		return err
	}

	sub, err := api.GetSubscriptionContext(creds, flg)
	if err != nil {
		return err
	}

	if sub.Alias == "" {
		return errs.ErrNoSubLink
	}

	utils.PrintCommandContext(flg.Output, creds.Login, sub.Alias, sub.HashedID)

	pi, err := getBuildName(args)
	if err != nil {
		return err
	}

	ref, err := getRef(flg)
	if err != nil {
		return err
	}

	req := &api.Request{
		Retries:         3,
		RunTokenRefresh: true,
		Credentials:     creds,
		Method:          "PUT",
		Path:            "/build/" + pi,
		MapStringPayload: map[string]interface{}{
			"ref": ref,
		},
	}

	res, err := req.NankaiSend()
	if err != nil {
		return errors.Wrap(err, errs.APISubListErrorMsg)
	}

	if res.StatusCode != 204 {
		return res.HandleFailure(flg.Output)
	}

	if strings.ToLower(flg.Output) == "json" {
		utils.PrintInterfaceAsJSON(res)
		return nil
	}

	fmt.Println()
	color.Green("Completed successfully!")

	return nil
}

func getBuildName(args []string) (string, error) {
	if len(args) != 0 {
		return args[0], nil
	}

	pi, err := services.StdinPrompt("Package ID: ")
	if err != nil {
		return "", errors.New("No package idenitifer supplied")
	}

	return pi, nil
}

func getRef(flg flags.Accumulator) (string, error) {
	if flg.Ref != "" {
		return flg.Ref, nil
	}

	r, err := services.StdinPrompt("New Ref: ")
	if err != nil {
		return "", errors.New("A new ref for the package was not supplied")
	}

	return r, nil
}

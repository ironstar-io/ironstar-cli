package subscription

import (
	"fmt"
	"os"
	"strings"

	"gitlab.com/ironstar-io/ironstar-cli/cmd/flags"
	"gitlab.com/ironstar-io/ironstar-cli/internal/api"
	"gitlab.com/ironstar-io/ironstar-cli/internal/errs"
	"gitlab.com/ironstar-io/ironstar-cli/internal/services"
	"gitlab.com/ironstar-io/ironstar-cli/internal/types"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
	yaml "gopkg.in/yaml.v2"
)

func List(args []string, flg flags.Accumulator) error {
	creds, err := services.ResolveUserCredentials(flg.Login)
	if err != nil {
		return err
	}

	if flg.Output == "" {
		color.Green("Using login [" + creds.Login + "]")
	}

	req := &api.Request{
		RunTokenRefresh:  true,
		Credentials:      creds,
		Method:           "GET",
		Path:             "/user/subscriptions",
		MapStringPayload: map[string]string{},
	}

	res, err := req.Send()
	if err != nil {
		return errors.Wrap(err, errs.APISubListErrorMsg)
	}

	if res.StatusCode != 200 {
		return res.HandleFailure()
	}

	if flg.Output == "json" {
		err = services.OutputJSON(res.Body)
		if err != nil {
			return errors.Wrap(err, errs.APISubListErrorMsg)
		}

		return nil
	}

	var uar []types.UserAccessResponse
	err = yaml.Unmarshal(res.Body, &uar)
	if err != nil {
		return err
	}

	fmt.Println()
	fmt.Println("Available Subscriptions:")

	uarRows := make([][]string, len(uar))
	for _, access := range uar {
		uarRows = append(uarRows, []string{access.Subscription.HashedID, access.Subscription.Alias, access.Subscription.ApplicationType, access.Role.Name, strings.Join(access.Role.Permissions, ", ")})
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Name", "Application Type", "Role", "Permissions"})
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.AppendBulk(uarRows)
	table.Render()

	return nil
}

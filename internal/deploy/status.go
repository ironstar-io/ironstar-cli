package deploy

import (
	"fmt"
	"os"

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

func Status(args []string, flg flags.Accumulator) error {
	creds, err := services.ResolveUserCredentials(flg.Login)
	if err != nil {
		return err
	}

	deployID, err := getDeployID(args, flg.Deploy)
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

	if flg.Output == "json" {
		err = services.OutputJSON(res.Body)
		if err != nil {
			return errors.Wrap(err, errs.APISubListErrorMsg)
		}

		return nil
	}

	var bs []types.DeploymentResponse
	err = yaml.Unmarshal(res.Body, &bs)
	if err != nil {
		return err
	}

	fmt.Println("DEPLOYMENT ID: ")
	fmt.Println("PACKAGE ID: ")
	fmt.Println("ENVIRONMENT: ")
	fmt.Println("APPLICATION STATUS: ")
	fmt.Println("ADMIN SERVICE STATUS: ")
	fmt.Println("CREATED: ")

	req2 := &api.Request{
		RunTokenRefresh:  true,
		Credentials:      creds,
		Method:           "GET",
		Path:             "/deployment/" + deployID + "/activity",
		MapStringPayload: map[string]string{},
	}

	res2, err := req2.Send()
	if err != nil {
		return errors.Wrap(err, errs.APISubListErrorMsg)
	}

	if res2.StatusCode != 200 {
		return res2.HandleFailure()
	}

	var da []types.DeploymentActivityResponse
	err = yaml.Unmarshal(res.Body, &da)
	if err != nil {
		return err
	}

	fmt.Println()
	fmt.Println("ACTIVITY: ")

	daRows := make([][]string, len(da))
	for _, d := range da {
		bsRows = append(daRows, []string{b.CreatedAt.String(), b.HashedID, b.CreatedBy, b.RunningIn})
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Date", "Action"})
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.AppendBulk(daRows)
	table.Render()

	return nil
}

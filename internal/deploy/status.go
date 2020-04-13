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

	var d types.DeploymentResponse
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

	var dac []types.DeploymentActivityResponse
	err = yaml.Unmarshal(res2.Body, &dac)
	if err != nil {
		return err
	}

	fmt.Println()
	fmt.Println("ACTIVITY: ")

	daRows := make([][]string, len(dac))
	for _, da := range dac {
		daRows = append(daRows, []string{da.CreatedAt.String(), da.Message, da.Flag})
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Date", "Action", "Flag"})
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.AppendBulk(daRows)
	table.Render()

	return nil
}

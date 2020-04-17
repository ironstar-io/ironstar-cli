package deploy

import (
	"fmt"
	"os"

	"gitlab.com/ironstar-io/ironstar-cli/internal/api"
	"gitlab.com/ironstar-io/ironstar-cli/internal/errs"
	"gitlab.com/ironstar-io/ironstar-cli/internal/types"

	"github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
	yaml "gopkg.in/yaml.v2"
)

func DisplayDeploymentActivity(creds types.Keylink, deployID string) error {
	req := &api.Request{
		RunTokenRefresh:  true,
		Credentials:      creds,
		Method:           "GET",
		Path:             "/deployment/" + deployID + "/activity",
		MapStringPayload: map[string]string{},
	}

	res, err := req.Send()
	if err != nil {
		return errors.Wrap(err, errs.APISubListErrorMsg)
	}

	if res.StatusCode != 200 {
		return res.HandleFailure()
	}

	var dac []types.DeploymentActivityResponse
	err = yaml.Unmarshal(res.Body, &dac)
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

package pkg

import (
	"fmt"
	"os"
	"time"

	"gitlab.com/ironstar-io/ironstar-cli/cmd/flags"
	"gitlab.com/ironstar-io/ironstar-cli/internal/api"
	"gitlab.com/ironstar-io/ironstar-cli/internal/errs"
	"gitlab.com/ironstar-io/ironstar-cli/internal/services"
	"gitlab.com/ironstar-io/ironstar-cli/internal/system/utils"
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

	sub, err := api.GetSubscriptionContext(creds, flg.Subscription)
	if err != nil {
		return err
	}

	if sub.Alias == "" {
		return errors.New("No Ironstar subscription has been linked to this project. Have you run `iron subscription link [subscription-name]`")
	}

	color.Green("Using login [" + creds.Login + "] for subscription <" + sub.Alias + ">")

	req := &api.Request{
		RunTokenRefresh:  true,
		Credentials:      creds,
		Method:           "GET",
		Path:             "/subscription/" + sub.HashedID + "/builds",
		MapStringPayload: map[string]string{},
	}

	res, err := req.Send()
	if err != nil {
		return errors.Wrap(err, errs.APISubListErrorMsg)
	}

	if res.StatusCode != 200 {
		return res.HandleFailure()
	}

	var bs []types.Build
	err = yaml.Unmarshal(res.Body, &bs)
	if err != nil {
		return err
	}

	fmt.Println()
	fmt.Println("Available Packages:")

	var envRefs []string
	bsRows := make([][]string, len(bs))
	for _, b := range bs {
		var runningIn string
		if len(b.Deployment) > 0 {
			for _, d := range b.Deployment {
				if d != (types.Deployment{}) && d.Environment != (types.Environment{}) && d.AppStatus == "FINISHED" && !utils.StringSliceContains(envRefs, d.Environment.HashedID) {
					runningIn = d.Environment.Name
					envRefs = append(envRefs, d.Environment.HashedID)
				}
			}
		}

		// Prepend rows, we want dates ordered oldest to newest
		row := make([][]string, 1)
		row = append(row, []string{b.CreatedAt.Format(time.RFC3339), b.Name, b.CreatedBy, runningIn})
		bsRows = append(row, bsRows...)
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Date Created", "Name", "Created By", "Running In"})
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.AppendBulk(bsRows)
	table.Render()

	return nil
}

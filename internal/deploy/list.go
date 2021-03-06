package deploy

import (
	"fmt"
	"os"
	"time"

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

	sub, err := api.GetSubscriptionContext(creds, flg)
	if err != nil {
		return err
	}

	if sub.Alias == "" {
		return errors.New("No Ironstar subscription has been linked to this project. Have you run `iron subscription link [subscription-name]`")
	}

	color.Green("Using login [" + creds.Login + "] for subscription '" + sub.Alias + "' (" + sub.HashedID + ")")

	qs := services.BuildQSFilters(flg, "10")
	req := &api.Request{
		RunTokenRefresh:  true,
		Credentials:      creds,
		Method:           "GET",
		Path:             "/subscription/" + sub.HashedID + "/deployments" + qs,
		MapStringPayload: map[string]interface{}{},
	}

	res, err := req.NankaiSend()
	if err != nil {
		return errors.Wrap(err, errs.APISubListErrorMsg)
	}

	if res.StatusCode != 200 {
		return res.HandleFailure()
	}

	var ds []types.Deployment
	err = yaml.Unmarshal(res.Body, &ds)
	if err != nil {
		return err
	}

	fmt.Println()
	fmt.Println("Deployments:")

	dsRows := make([][]string, len(ds))
	for _, d := range ds {
		// Prepend rows, we want dates ordered oldest to newest
		row := make([][]string, 1)

		status := d.AppStatus
		if d.Status.Lifecycle != "" {
			status = d.Status.Lifecycle
		}

		row = append(row, []string{d.CreatedAt.Format(time.RFC3339), d.Environment.Name, d.Name, d.Build.Name, status, d.Build.Branch, d.Build.Tag})
		dsRows = append(row, dsRows...)
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Date Created", "Environment", "Deployment", "Package", "Status", "Branch", "Tag"})
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.AppendBulk(dsRows)
	table.Render()

	return nil
}

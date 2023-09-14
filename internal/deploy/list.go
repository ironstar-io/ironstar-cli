package deploy

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/ironstar-io/ironstar-cli/cmd/flags"
	"github.com/ironstar-io/ironstar-cli/internal/api"
	"github.com/ironstar-io/ironstar-cli/internal/errs"
	"github.com/ironstar-io/ironstar-cli/internal/services"
	"github.com/ironstar-io/ironstar-cli/internal/system/utils"
	"github.com/ironstar-io/ironstar-cli/internal/types"

	"github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
	yaml "gopkg.in/yaml.v2"
)

func List(args []string, flg flags.Accumulator) error {
	creds, err := services.ResolveUserCredentials(flg.Login)
	if err != nil {
		return err
	}

	if flg.Environment != "" {
		return retrieveAndDisplayEnvDeployments(creds, flg)
	}

	return retrieveAndDisplaySubDeployments(creds, flg)
}

func retrieveAndDisplayEnvDeployments(creds types.Keylink, flg flags.Accumulator) error {
	seCtx, err := api.GetSubscriptionEnvironmentContext(creds, flg)
	if err != nil {
		return err
	}

	if seCtx.Subscription.Alias == "" {
		return errs.ErrNoSubLink
	}

	utils.PrintCommandContext(flg.Output, creds.Login, seCtx.Subscription.Alias, seCtx.Subscription.HashedID)

	qs := services.BuildQSFilters(flg, "10")
	req := &api.Request{
		Retries:          3,
		RunTokenRefresh:  true,
		Credentials:      creds,
		Method:           "GET",
		Path:             "/subscription/" + seCtx.Subscription.HashedID + "/environment/" + seCtx.Environment.HashedID + "/deployments" + qs,
		MapStringPayload: map[string]interface{}{},
	}

	res, err := req.NankaiSend()
	if err != nil {
		return errors.Wrap(err, errs.APISubListErrorMsg)
	}

	if res.StatusCode != 200 {
		return res.HandleFailure(flg.Output)
	}

	var ds []types.Deployment
	err = yaml.Unmarshal(res.Body, &ds)
	if err != nil {
		return err
	}

	return displayDeployments(flg.Output, ds)
}

func retrieveAndDisplaySubDeployments(creds types.Keylink, flg flags.Accumulator) error {
	sub, err := api.GetSubscriptionContext(creds, flg)
	if err != nil {
		return err
	}

	if sub.Alias == "" {
		return errs.ErrNoSubLink
	}

	utils.PrintCommandContext(flg.Output, creds.Login, sub.Alias, sub.HashedID)

	qs := services.BuildQSFilters(flg, "10")
	req := &api.Request{
		Retries:          3,
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
		return res.HandleFailure(flg.Output)
	}

	var ds []types.Deployment
	err = yaml.Unmarshal(res.Body, &ds)
	if err != nil {
		return err
	}

	return displayDeployments(flg.Output, ds)
}

func displayDeployments(output string, ds []types.Deployment) error {
	if strings.ToLower(output) == "json" {
		utils.PrintInterfaceAsJSON(ds)
		return nil
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

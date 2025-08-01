package pkg

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

	sub, err := api.GetSubscriptionContext(creds, flg)
	if err != nil {
		return err
	}

	if sub.Alias == "" {
		return errs.ErrNoSubLink
	}

	utils.PrintCommandContext(flg.Output, creds.Login, sub.Alias, sub.HashedID)

	req := &api.Request{
		Retries:          3,
		RunTokenRefresh:  true,
		Credentials:      creds,
		Method:           "GET",
		Path:             "/subscription/" + sub.HashedID + "/builds",
		MapStringPayload: nil,
	}

	res, err := req.NankaiSend()
	if err != nil {
		return errors.Wrap(err, errs.APISubListErrorMsg)
	}

	if res.StatusCode != 200 {
		return res.HandleFailure(flg.Output)
	}

	var bs []types.Build
	err = yaml.Unmarshal(res.Body, &bs)
	if err != nil {
		return err
	}

	if strings.ToLower(flg.Output) == "json" {
		utils.PrintInterfaceAsJSON(bs)
		return nil
	}

	fmt.Println()
	fmt.Println("Available Packages:")

	mbm := calcRunningBuildEnv(bs)

	bsRows := make([][]string, len(bs))
	for _, b := range bs {
		runningIn := ""

		for _, v := range mbm {
			if b.Name == v.BuildName {
				if runningIn == "" {
					runningIn = v.EnvironmentName

					continue
				}

				runningIn = runningIn + ", " + v.EnvironmentName
			}
		}

		// Prepend rows, we want dates ordered oldest to newest
		row := make([][]string, 1)
		row = append(row, []string{b.CreatedAt.Format(time.RFC3339), b.Name, b.Branch, b.Tag, b.CreatedBy, runningIn, b.CommitSHA})
		bsRows = append(row, bsRows...)
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Date Created", "Name", "Branch", "Tag", "Created By", "Running In", "Commit SHA"})
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.AppendBulk(bsRows)
	table.Render()

	return nil
}

type MinimalBuildMatcher struct {
	BuildName       string
	EnvironmentName string
	ActiveDate      time.Time
}

func calcRunningBuildEnv(bs []types.Build) map[string]MinimalBuildMatcher {
	buildRefs := make(map[string]MinimalBuildMatcher)

	for _, b := range bs {
		if len(b.Deployment) == 0 {
			continue
		}

		for _, d := range b.Deployment {
			if !d.IsStructureEmpty() && d.Environment != (types.Environment{}) {
				if d.Status != (types.DeploymentStatus{}) && strings.ToLower(d.Status.Lifecycle) == "active" {
					if buildRefs[d.Environment.Name] == (MinimalBuildMatcher{}) {
						buildRefs[d.Environment.Name] = MinimalBuildMatcher{
							BuildName:       b.Name,
							EnvironmentName: d.Environment.Name,
							ActiveDate:      d.Status.ActiveDate,
						}
					}

					if d.Status.ActiveDate.After(buildRefs[d.Environment.Name].ActiveDate) {
						buildRefs[d.Environment.Name] = MinimalBuildMatcher{
							BuildName:       b.Name,
							EnvironmentName: d.Environment.Name,
							ActiveDate:      d.Status.ActiveDate,
						}
					}
				}
			}
		}
	}

	return buildRefs
}

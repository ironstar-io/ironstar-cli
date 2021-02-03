package deploy

import (
	"fmt"
	"strconv"
	"time"

	"gitlab.com/ironstar-io/ironstar-cli/cmd/flags"
	"gitlab.com/ironstar-io/ironstar-cli/internal/api"
	"gitlab.com/ironstar-io/ironstar-cli/internal/logs"
	"gitlab.com/ironstar-io/ironstar-cli/internal/services"
	"gitlab.com/ironstar-io/ironstar-cli/internal/system/utils"
	"gitlab.com/ironstar-io/ironstar-cli/internal/types"

	"github.com/fatih/color"
)

func Status(args []string, flg flags.Accumulator) error {
	creds, err := services.ResolveUserCredentials(flg.Login)
	if err != nil {
		return err
	}

	deployID, err := services.GetDeployID(args, flg.Deploy)
	if err != nil {
		return err
	}

	color.Green("Using login [" + creds.Login + "]")

	deployment, err := api.GetDeployment(creds, deployID)
	if err != nil {
		return err
	}

	flg.Environment = deployment.Environment.Name
	seCtx, err := api.GetSubscriptionEnvironmentContext(creds, flg)
	if err != nil {
		return err
	}

	err = DisplayDeploymentInfo(creds, deployment)
	if err != nil {
		return err
	}

	err = DisplayDeploymentActivity(creds, deployID)
	if err != nil {
		return err
	}

	cwLogs, err := logs.RetrieveArimaLogs(creds, seCtx.Subscription.Alias, seCtx.Environment.Name, deployment.Name, utils.UnixMilliseconds(deployment.CreatedAt), utils.UnixMilliseconds(time.Now()), []string{"deploy.log"})
	if err != nil {
		return err
	}

	if len(cwLogs) == 0 {
		if seCtx.Environment.LogRetention != 0 && deployment.CreatedAt.Before(time.Now().UTC().Add(time.Duration(time.Duration(-seCtx.Environment.LogRetention*24)*time.Hour)).UTC()) {
			fmt.Println()
			fmt.Println("The log for this deployment is no longer available. Logs in your " + seCtx.Environment.Name + " environment are retained for " + strconv.Itoa(int(seCtx.Environment.LogRetention)) + " days.")
			fmt.Println()
		} else {
			fmt.Println()
			fmt.Println("Logs not available for this deployment")
			fmt.Println()
		}

		return nil
	}

	fmt.Println()
	logs.StdoutArimaLogs(cwLogs)

	return nil
}

func DisplayDeploymentInfo(creds types.Keylink, d types.Deployment) error {
	fmt.Println()
	fmt.Println("DEPLOYMENT: " + d.Name)
	fmt.Println()
	fmt.Println("ENVIRONMENT: " + d.Environment.Name)
	fmt.Println("PACKAGE: " + d.Build.Name)
	fmt.Println("APPLICATION STATUS: " + d.AppStatus)
	fmt.Println("ADMIN SERVICE STATUS: " + d.AdminSvcStatus)
	fmt.Println("CREATED: " + d.CreatedAt.Format(time.RFC3339))

	return nil
}

func DisplayDeploymentLogs(creds types.Keylink, d types.Deployment) error {
	fmt.Println()
	fmt.Println("DEPLOYMENT: " + d.Name)
	fmt.Println()
	fmt.Println("ENVIRONMENT: " + d.Environment.Name)
	fmt.Println("PACKAGE: " + d.Build.Name)
	fmt.Println("APPLICATION STATUS: " + d.AppStatus)
	fmt.Println("ADMIN SERVICE STATUS: " + d.AdminSvcStatus)
	fmt.Println("CREATED: " + d.CreatedAt.Format(time.RFC3339))

	return nil
}

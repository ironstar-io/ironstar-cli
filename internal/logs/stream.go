package logs

import (
	"fmt"
	"strings"
	"time"

	"github.com/ironstar-io/ironstar-cli/cmd/flags"
	"github.com/ironstar-io/ironstar-cli/internal/api"
	"github.com/ironstar-io/ironstar-cli/internal/errs"
	"github.com/ironstar-io/ironstar-cli/internal/services"
	"github.com/ironstar-io/ironstar-cli/internal/system/utils"

	"github.com/fatih/color"
	"github.com/pkg/errors"
)

func Display(args []string, flg flags.Accumulator) error {
	creds, err := services.ResolveUserCredentials(flg.Login)
	if err != nil {
		return err
	}

	seCtx, err := api.GetSubscriptionEnvironmentContext(creds, flg)
	if err != nil {
		return err
	}

	if seCtx.Subscription.Alias == "" {
		return errs.ErrNoSubLink
	}

	utils.PrintCommandContext(flg.Output, creds.Login, seCtx.Subscription.Alias, seCtx.Subscription.HashedID)

	err = checkLogLabelFlags(flags.Acc.Filenames, flags.Acc.Sources)
	if err != nil {
		return err
	}

	startTime := flg.Start
	endTime := flg.End

	custLogs, err := RetrieveEnvironmentLogs(
		RetrieveEnvironmentLogsParams{
			Creds:     creds,
			SubAlias:  seCtx.Subscription.Alias,
			EnvName:   seCtx.Environment.Name,
			Search:    flags.Acc.Search,
			Start:     startTime,
			End:       endTime,
			Filenames: flags.Acc.Filenames,
			Sources:   flags.Acc.Sources,
		})
	if err != nil {
		return err
	}

	if strings.ToLower(flg.Output) == "json" {
		utils.PrintInterfaceAsJSON(custLogs)
		return nil
	}

	fmt.Println()
	fmt.Println("Retrieving logs for environment '" + seCtx.Environment.Name + "'")
	fmt.Println()

	if startTime != 0 {
		fmt.Println("Start Time: " + time.Unix(0, startTime*int64(time.Millisecond)).UTC().Format(time.RFC3339))
	}
	if endTime != 0 {
		fmt.Println("End Time: " + time.Unix(0, endTime*int64(time.Millisecond)).UTC().Format(time.RFC3339))
	}

	if len(flags.Acc.Filenames) > 0 {
		fmt.Println("Filenames: " + strings.Join(flags.Acc.Filenames, ", "))
		fmt.Println()
	}
	if len(flags.Acc.Sources) > 0 {
		fmt.Println("Sources: " + strings.Join(flags.Acc.Filenames, ", "))
		fmt.Println()
	}

	if custLogs == nil || len(custLogs.Results) == 0 {
		color.Yellow("No logs found in the specified time period for this environment")
		return nil
	}

	StdoutEnvironmentLogs(custLogs.Results)

	return nil
}

func checkLogLabelFlags(filenames, sources []string) error {
	if len(filenames) > 0 && len(sources) > 0 {
		return errors.New("The flags 'filenames' and 'sources' must not be set simultaneously")
	}

	return nil
}

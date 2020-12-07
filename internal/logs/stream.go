package logs

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
	"unicode/utf8"

	"gitlab.com/ironstar-io/ironstar-cli/cmd/flags"
	"gitlab.com/ironstar-io/ironstar-cli/internal/api"
	"gitlab.com/ironstar-io/ironstar-cli/internal/services"
	"gitlab.com/ironstar-io/ironstar-cli/internal/types"

	"github.com/fatih/color"
	"github.com/pkg/errors"
)

func Stream(args []string, flg flags.Accumulator) error {
	creds, err := services.ResolveUserCredentials(flg.Login)
	if err != nil {
		return err
	}

	seCtx, err := api.GetSubscriptionEnvironmentContext(creds, flg)
	if err != nil {
		return err
	}

	if seCtx.Subscription.Alias == "" {
		return errors.New("No Ironstar subscription has been linked to this project. Have you run `iron subscription link [subscription-name]`")
	}

	color.Green("Using login [" + creds.Login + "] for subscription '" + seCtx.Subscription.Alias + "' (" + seCtx.Subscription.HashedID + ")")

	fmt.Println()
	fmt.Println("Logs for your environment")
	fmt.Println()

	last, err := printArimaLogs(creds, seCtx, flg, time.Now().UTC().Add(time.Duration(-15*time.Minute)).UTC().UnixNano()/int64(time.Millisecond))
	if err != nil {
		return err
	}

	for range time.Tick(3 * time.Second) {
		go func() {
			newLast, err := printArimaLogs(creds, seCtx, flg, last)
			if err != nil {
				os.Exit(1)
			}
			last = newLast
		}()
	}

	return nil
}

func printArimaLogs(creds types.Keylink, seCtx types.SubscriptionEnvironment, flags flags.Accumulator, start int64) (int64, error) {
	payload := map[string]interface{}{
		"logStreamNames": []string{"fpm.access.log", "cache.access.log", "cron.log"},
		"start":          start,
		"end":            time.Now().UTC().UnixNano() / int64(time.Millisecond),
	}

	cwLogs, err := api.QueryEnvironmentLogs(creds, seCtx.Subscription.Alias, seCtx.Environment.Name, payload)
	if err != nil {
		return 0, err
	}

	if len(cwLogs) == 0 {
		return start, nil
	}

	for _, cwLog := range cwLogs {
		var logMsg string
		b, err := json.Marshal(cwLog.Log)
		if err != nil {
			logMsg = fmt.Sprint(cwLog.Log)
		} else {
			logMsg = string(b)
		}

		streamColor := streamNameColour(cwLog.LogStreamName)

		fmt.Printf("%s%s\n", color.New(streamColor).SprintFunc()(streamNameWithPadding(cwLog.LogStreamName)), logMsg)
	}

	s := cwLogs[len(cwLogs)-1]

	return s.IngestionTime + 1, nil
}

func streamNameWithPadding(logStreamName string) string {
	const colWidth = 18

	length := utf8.RuneCountInString(logStreamName)

	if length >= 17 {
		return logStreamName + " | "
	}

	padding := colWidth - length

	base := logStreamName
	for i := 0; i < padding; i++ {
		base = base + " "
	}

	return base + " | "
}

func streamNameColour(logStreamName string) color.Attribute {
	switch logStreamName {
	case "fpm.access.log":
		return color.FgCyan
	case "nginx.access.log":
		return color.FgGreen
	case "cache.access.log":
		return color.FgYellow
	case "cron.log":
		return color.FgMagenta
	case "deploy.log":
		return color.FgRed
	default:
		return color.FgWhite
	}
}

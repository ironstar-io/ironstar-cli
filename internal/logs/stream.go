package logs

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/ironstar-io/ironstar-cli/cmd/flags"
	"github.com/ironstar-io/ironstar-cli/internal/api"
	"github.com/ironstar-io/ironstar-cli/internal/services"
	"github.com/ironstar-io/ironstar-cli/internal/types"

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

	// labelValues, err := api.GetEnvironmentLogLabelValues(creds, seCtx.Subscription.Alias, seCtx.Environment.Name, "filename")
	// if err != nil {
	// 	return err
	// }

	// if len(labelValues) == 0 {
	// 	return errors.New("There are no logs for this environment")
	// }

	// logStreams := calcTargetLogStreams(flg.LogStreams)

	// if len(logStreams) == 0 {
	// 	return errors.New("There were no matching logs available for this environment")
	// }

	// logStreamNames := calcLogStreamNames(logStreams)
	var logStreamNames []string // TODO

	startTime := calcStartTime(flg.Start)
	endTime := calcEndTime(flg.End)

	if !flags.Acc.Stream {
		fmt.Println()
		fmt.Println("Printing logs for environment '" + seCtx.Environment.Name + "'")
		fmt.Println()
		fmt.Println("You can stream logs by passing the '-f' flag")
	} else {
		fmt.Println()
		fmt.Println("Streaming logs for environment '" + seCtx.Environment.Name + "'")
	}

	fmt.Println()
	fmt.Println("Start Time: " + time.Unix(0, startTime*int64(time.Millisecond)).UTC().String())
	if endTime != 0 && !flags.Acc.Stream {
		fmt.Println("End Time: " + time.Unix(0, endTime*int64(time.Millisecond)).UTC().String())
	}
	// fmt.Println("Log Streams: " + strings.Join(logStreamNames, ", "))
	fmt.Println()

	last, err := PrintEnvironmentLogs(creds, seCtx, flg, startTime, endTime, logStreamNames)
	if err != nil {
		return err
	}

	if flags.Acc.Stream {
		streamEnd := time.Now().Add(time.Duration(15 * time.Minute))

		for range time.Tick(3 * time.Second) {
			if time.Now().After(streamEnd) {
				fmt.Println()
				fmt.Println("Stream runnning for longer than 15 minutes. Exiting...")

				return nil
			}

			go func() {
				newLast, err := PrintEnvironmentLogs(creds, seCtx, flg, last, 0, logStreamNames)
				if err != nil {
					os.Exit(1)
				}

				last = newLast
			}()
		}
	}

	return nil
}

func calcStartTime(start int) int64 {
	if start != 0 {
		return int64(start)
	}

	return 0

	// if flags.Search != "" {
	// 	return 1
	// }

	// sort.SliceStable(availableLogStreams, func(i, j int) bool {
	// 	return availableLogStreams[i].LastEventTimestamp > availableLogStreams[j].LastEventTimestamp
	// })

	// s := availableLogStreams[0]

	// if s.LastEventTimestamp == 0 {
	// 	return time.Now().UTC().Add(time.Duration(-15*time.Minute)).UTC().UnixNano() / int64(time.Millisecond)
	// }

	// return time.Unix(0, s.LastEventTimestamp*int64(time.Millisecond)).Add(-2*time.Minute).UnixNano() / int64(time.Millisecond)
}

func formatLogTimestamp(ogTimestamp string) string {
	split1 := strings.ReplaceAll(ogTimestamp, `"`, "")
	split2 := strings.ReplaceAll(split1, `time=`, "")
	split3 := strings.ReplaceAll(split2, `+0000`, "+00:00")

	t, err := time.Parse(time.RFC3339, split3)
	if err != nil {
		return ogTimestamp
	}

	return t.Format(time.RFC3339Nano)
}

func calcEndTime(endFlag int) int64 {
	if endFlag == 0 {
		return 0
	}

	return int64(endFlag)
}

func calcTargetLogStreams(logStreamFlag []string) []string {
	if len(logStreamFlag) == 0 {
		return nil
	}

	return logStreamFlag
	// if len(logStreamFlag) == 0 {
	// 	return availableLogStreams
	// }

	// var targetStreams []types.CWLogStreamsResponse
	// for _, ls := range availableLogStreams {
	// 	if utils.SliceIncludes(logStreamFlag, ls.LogStreamName) {
	// 		targetStreams = append(targetStreams, ls)
	// 	}
	// }

	// return targetStreams
}

func calcLogStreamNames(availableLogStreams []types.CWLogStreamsResponse) []string {
	var streamNames []string
	for _, ls := range availableLogStreams {
		streamNames = append(streamNames, ls.LogStreamName)
	}

	return streamNames
}

func RetrieveEnvironmentLogs(creds types.Keylink, subAlias, envName, search string, start, end int64, logStreams []string) (*types.CustomerLogsResponse, error) {
	payload := map[string]interface{}{
		"logStreamNames": logStreams,
		"start":          start,
		"end":            end,
		"pattern":        search,
	}

	return api.QueryEnvironmentLogs(creds, subAlias, envName, payload)
}

func PrintEnvironmentLogs(creds types.Keylink, seCtx types.SubscriptionEnvironment, flags flags.Accumulator, start, end int64, logStreams []string) (int64, error) {
	custLogs, err := RetrieveEnvironmentLogs(creds, seCtx.Subscription.Alias, seCtx.Environment.Name, flags.Search, start, end, logStreams)
	if err != nil {
		return 0, err
	}

	if custLogs == nil || len(custLogs.Results) == 0 {
		return start, nil
	}

	StdoutEnvironmentLogs(custLogs.Results)

	// s := custLogs.Results[len(custLogs.Results)-1]

	// return s.Values + 1, nil
	return 1, nil // TODO
}

func StdoutEnvironmentLogs(custLogs []types.CustomerLogsResult) {
	for _, custLog := range custLogs[0].Values {
		streamColor := streamNameColour(custLog[2])
		// logMsg := stringifyLog(cwLog.Log)

		fmt.Printf("%s%s%s\n", color.New(streamColor).SprintFunc()(streamNameWithPadding(custLog[4])), color.New(color.Faint).SprintFunc()(formatLogTimestamp(fmt.Sprintf("%v", custLog[0]))+" | "), custLog[1])
	}
}

func stringifyLog(logMsg map[string]interface{}) string {
	var keys []string
	for key := range logMsg {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	var msg string
	for _, key := range keys {
		if key == "tmst" {
			continue
		}

		msg = msg + key + "=" + fmt.Sprintf("%v", logMsg[key]) + " "
	}

	if len(msg) == 0 {
		return ""
	}

	msg = msg[:len(msg)-1]

	return msg
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

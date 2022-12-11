package logs

import (
	"fmt"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/ironstar-io/ironstar-cli/internal/api"
	"github.com/ironstar-io/ironstar-cli/internal/types"

	"github.com/fatih/color"
)

type RetrieveEnvironmentLogsParams struct {
	Creds     types.Keylink
	SubAlias  string
	EnvName   string
	Search    string
	Limit     int64
	Start     int64
	End       int64
	Filenames []string
	Sources   []string
}

func RetrieveEnvironmentLogs(params RetrieveEnvironmentLogsParams) (*types.CustomerLogsResponse, error) {
	payload := map[string]interface{}{
		"output":    "merge",
		"direction": "forward",
		"limit":     100,
		"search":    params.Search,
	}

	if params.Start != 0 {
		payload["start"] = time.Unix(0, params.Start*int64(time.Millisecond)).UTC().Format(time.RFC3339)
	}
	if params.End != 0 {
		payload["end"] = time.Unix(0, params.End*int64(time.Millisecond)).UTC().Format(time.RFC3339)
	}

	return api.QueryEnvironmentLogs(params.Creds, params.SubAlias, params.EnvName, payload)
}

func StdoutEnvironmentLogs(custLogs []types.CustomerLogsResult) {
	for _, custLog := range custLogs[0].Values {
		streamColor := streamNameColour(custLog[2])

		fmt.Printf("%s%s%s\n",
			color.New(streamColor).SprintFunc()(streamNameWithPadding(custLog[3])),
			color.New(color.Faint).SprintFunc()(formatLogTimestamp(fmt.Sprintf("%v", custLog[0]))+" | "),
			custLog[1])
	}
}

func streamNameWithPadding(logStreamName string) string {
	const colWidth = 18

	length := utf8.RuneCountInString(logStreamName)

	if length >= 18 {
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
	if strings.Contains(logStreamName, "fpm") {
		return color.FgCyan
	}
	if strings.Contains(logStreamName, "nginx") {
		return color.FgGreen
	}
	if strings.Contains(logStreamName, "memcache") {
		return color.BgBlue
	}
	if strings.Contains(logStreamName, "cache") {
		return color.FgYellow
	}
	if strings.Contains(logStreamName, "cron") {
		return color.FgMagenta
	}
	if strings.Contains(logStreamName, "deploy") {
		return color.FgRed
	}
	if strings.Contains(logStreamName, "antivirus") {
		return color.FgHiRed
	}
	if strings.Contains(logStreamName, "newrelic") {
		return color.FgBlue
	}
	if strings.Contains(logStreamName, "redis") {
		return color.BgRed
	}

	return color.FgWhite
}

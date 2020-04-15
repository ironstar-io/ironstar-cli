package services

import (
	"strings"

	"gitlab.com/ironstar-io/ironstar-cli/cmd/flags"
)

func BuildQSFilters(flg flags.Accumulator) string {
	var acc []string

	if flg.Deploy != "" {
		acc = append(acc, "deployment="+flg.Deploy)
	}

	if flg.Environment != "" {
		acc = append(acc, "environment="+flg.Environment)
	}

	if flg.Subscription != "" {
		acc = append(acc, "subscription="+flg.Subscription)
	}

	if flg.Package != "" {
		acc = append(acc, "build="+flg.Package)
	}

	if len(acc) == 0 {
		return ""
	}

	return "?" + strings.Join(acc, "&")
}

package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/pkg/errors"
)

func PrintInterfaceAsJSON(data interface{}) error {
	src, err := json.Marshal(data)
	if err != nil {
		return errors.New("Failed to marshal struct: " + err.Error())
	}

	return OutputJSON(src)
}

func OutputJSON(src []byte) error {
	dest := &bytes.Buffer{}
	if err := json.Indent(dest, src, "", "  "); err != nil {
		return errors.New("Unexpected error occurred")
	}

	fmt.Println(dest.String())

	return nil
}

func PrintCommandContext(output, login, subAlias, subHash string) {
	if strings.ToLower(output) != "json" {
		color.Green("Using login [" + login + "] for subscription '" + subAlias + "' (" + subHash + ")")
	}
}

func PrintErrorJSON(err error) {
	PrintInterfaceAsJSON(map[string]string{
		"result": "error",
		"error":  err.Error(),
	})
}

package services

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
)

// OutputJSON ...
func OutputJSON(src []byte) error {
	dest := &bytes.Buffer{}
	if err := json.Indent(dest, src, "", "  "); err != nil {
		return errors.New("Unexpected error occurred")
	}

	fmt.Println(dest.String())

	return nil
}

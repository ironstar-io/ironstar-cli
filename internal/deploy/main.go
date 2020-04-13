package deploy

import (
	"errors"
)

func getDeployID(args []string, deployFlag string) (string, error) {
	if deployFlag != "" {
		return deployFlag, nil
	}

	if len(args) == 0 {
		return "", errors.New("No deployment ID argument supplied")
	}

	return args[0], nil
}

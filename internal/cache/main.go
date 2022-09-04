package cache

import (
	"github.com/ironstar-io/ironstar-cli/internal/services"
)

func GetCacheInvalidationName(args []string) (string, error) {
	var name string
	if len(args) == 0 {
		input, err := services.StdinPrompt("Cache Invalidation Name: ")
		if err != nil {
			return "", err
		}
		name = input
	} else {
		name = args[0]
	}

	return name, nil
}

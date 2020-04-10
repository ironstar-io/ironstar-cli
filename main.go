package main

import (
	"log"

	"gitlab.com/ironstar-io/ironstar-cli/cmd"
)

func main() {
	if err := cmd.RootCmd().Execute(); err != nil {
		log.Fatal(err)
	}
}

package main

import (
	"log"

	"github.com/ironstar-io/ironstar-cli/cmd"
)

func main() {
	if err := cmd.RootCmd().Execute(); err != nil {
		log.Fatal(err)
	}
}

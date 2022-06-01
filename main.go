package main

import (
	"log"
	"os"

	"github.com/csothen/tmdei-project/cmd"
)

func main() {
	// Uncomment if we need to reload the plugins and sonarqube versions
	// sonarqube.NewFetcher().Reload()
	if err := cmd.RootCmd.Execute(); err != nil {
		log.Printf("%+v\n", err)
		os.Exit(1)
	}
}

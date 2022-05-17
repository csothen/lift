package main

import (
	"log"
	"os"

	"github.com/csothen/tmdei-project/cmd"
	"github.com/csothen/tmdei-project/internal/fetcher/sonarqube"
)

func main() {
	sonarqube.NewFetcher().Reload()
	if err := cmd.RootCmd.Execute(); err != nil {
		log.Printf("%+v\n", err)
		os.Exit(1)
	}
}

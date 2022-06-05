package main

import (
	"log"
	"os"

	"github.com/csothen/tmdei-project/cmd"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		log.Printf("%+v\n", err)
		os.Exit(1)
	}
}

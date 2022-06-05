package cmd

import (
	"github.com/csothen/tmdei-project/internal/fetcher/sonarqube"
	"github.com/spf13/cobra"
)

var fetchCmd = &cobra.Command{
	Use:   "fetch",
	Short: "Fetches static data",
	Long:  "Fetches all static data necessary like sonarqube versions and plugins",
	RunE: func(cmd *cobra.Command, args []string) error {
		return sonarqube.NewFetcher().Reload()
	},
}

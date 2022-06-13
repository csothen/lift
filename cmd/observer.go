package cmd

import (
	"log"

	"github.com/csothen/env"
	"github.com/csothen/lift/internal/config"
	"github.com/csothen/lift/internal/db"
	"github.com/csothen/lift/internal/observer"
	"github.com/csothen/lift/internal/services"
	"github.com/csothen/lift/internal/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var observerCmd = &cobra.Command{
	Use:   "observer",
	Short: "Starts the observer worker",
	Long:  "Starts the observer worker responsible for notifying the user when instances are up",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Load the environment variables
		err := env.Load(".env")
		if err != nil {
			log.Println("no .env file found: %w", err)
		}

		cfg := config.New(viper.GetViper())
		utils.InstallTerraform(cfg)
		repo := db.New(cfg)
		s := services.New(repo, cfg)
		ow := observer.NewWorker(s, cfg)

		log.Printf("observer worker started")
		return ow.Start()
	},
}

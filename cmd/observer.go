package cmd

import (
	"github.com/csothen/tmdei-project/internal/config"
	"github.com/csothen/tmdei-project/internal/db"
	"github.com/csothen/tmdei-project/internal/observer"
	"github.com/csothen/tmdei-project/internal/services"
	"github.com/csothen/tmdei-project/internal/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var observerCmd = &cobra.Command{
	Use:   "observer",
	Short: "Starts the observer worker",
	Long:  "Starts the observer worker responsible for notifying the user when instances are up",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := config.New(viper.GetViper())
		utils.InstallTerraform(cfg)
		repo := db.New(cfg)
		s := services.New(repo, cfg)
		ow := observer.NewWorker(s, cfg)
		return ow.Start()
	},
}

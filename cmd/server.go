package cmd

import (
	"fmt"
	"log"
	"net/http"

	"github.com/csothen/env"
	"github.com/csothen/lift/internal/db"
	"github.com/csothen/lift/internal/graph"
	"github.com/csothen/lift/internal/services"
	"github.com/csothen/lift/internal/utils"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/csothen/lift/internal/config"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Starts the server",
	Long:  "Starts the GraphQL server API with the given configuration",
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

		router := mux.NewRouter()

		router.Handle("/", graph.NewPlaygroundHandler("/query"))
		router.Handle("/query", graph.NewHandler(s))

		log.Printf("connect to http://localhost:%d/ for GraphQL playground", cfg.Port)
		return http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), router)
	},
}

func init() {
	// Server flags
	serverCmd.PersistentFlags().Int("port", 8080, "The port for the server to listen on")
	viper.BindPFlag("port", serverCmd.PersistentFlags().Lookup("port"))
}

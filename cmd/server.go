package cmd

import (
	"fmt"
	"log"
	"net/http"

	"github.com/csothen/tmdei-project/internal/db"
	"github.com/csothen/tmdei-project/internal/graph"
	"github.com/csothen/tmdei-project/internal/middlewares"
	"github.com/csothen/tmdei-project/internal/services"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/csothen/tmdei-project/internal/config"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Starts the server",
	Long:  "Starts the GraphQL server API with the given configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := config.New(viper.GetViper())
		repo := db.New(cfg)
		s := services.New(repo)

		router := mux.NewRouter()
		router.Use(middlewares.Json)
		router.Use(middlewares.Auth(s))

		router.Handle("/", graph.NewPlaygroundHandler("/query"))
		router.Handle("/query", graph.NewHandler(s))

		log.Printf("connect to http://localhost:%d/ for GraphQL playground", cfg.Port)
		return http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), router)
	},
}

func init() {
	serverCmd.PersistentFlags().Int("port", 8080, "The port to listen on")
	viper.BindPFlag("port", serverCmd.PersistentFlags().Lookup("port"))
}

package cmd

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/csothen/tmdei-project/internal/db"
	"github.com/csothen/tmdei-project/internal/graph"
	"github.com/csothen/tmdei-project/internal/middlewares"
	"github.com/csothen/tmdei-project/internal/services"
	"github.com/gorilla/mux"
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hc-install/product"
	"github.com/hashicorp/hc-install/releases"
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
		installTerraform(cfg)
		repo := db.New(cfg)
		s := services.New(repo, cfg)

		router := mux.NewRouter()
		router.Use(middlewares.Json)
		router.Use(middlewares.Auth(s))

		router.Handle("/", graph.NewPlaygroundHandler("/query"))
		router.Handle("/query", graph.NewHandler(s))

		log.Printf("connect to http://localhost:%d/ for GraphQL playground", cfg.Port)
		return http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), router)
	},
}

// in case no terraform executable path was provided we will install it
// and pass the new executable path to the config
func installTerraform(cfg *config.Config) {
	if cfg.TerraformExecPath == "" {
		// we will install version 1.2.1 of Terraform
		installer := &releases.ExactVersion{
			Product: product.Terraform,
			Version: version.Must(version.NewVersion("1.2.1")),
		}

		execPath, err := installer.Install(context.Background())
		if err != nil {
			log.Fatal("error installing Terraform: %w", err)
		}
		cfg.TerraformExecPath = execPath
	}
}

func init() {
	// Server flags
	serverCmd.PersistentFlags().Int("port", 8080, "The port for the server to listen on")
	viper.BindPFlag("port", serverCmd.PersistentFlags().Lookup("port"))

	// Terraform flags
	serverCmd.PersistentFlags().String("terraform_exec_path", "/usr/local/bin/", "Path to the terraform executable")
	viper.BindPFlag("terraform_exec_path", serverCmd.PersistentFlags().Lookup("terraform_exec_path"))

	// Database configuration flags
	serverCmd.PersistentFlags().String("db_host", "localhost", "The host where the database lives")
	viper.BindPFlag("db_host", serverCmd.PersistentFlags().Lookup("db_host"))
	serverCmd.PersistentFlags().String("port", "user", "The database user to connect with to the database")
	viper.BindPFlag("port", serverCmd.PersistentFlags().Lookup("db_user"))
	serverCmd.PersistentFlags().String("db_password", "", "The database user's password")
	viper.BindPFlag("db_password", serverCmd.PersistentFlags().Lookup("db_password"))
	serverCmd.PersistentFlags().String("db_name", "lift", "The database name")
	viper.BindPFlag("db_name", serverCmd.PersistentFlags().Lookup("db_name"))
	serverCmd.PersistentFlags().Int("db_port", 5432, "The port in which the database is listening on")
	viper.BindPFlag("db_port", serverCmd.PersistentFlags().Lookup("db_port"))
}

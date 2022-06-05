package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// RootCmd is the root command of the application where
// all the other subcomnads belong to
var RootCmd = &cobra.Command{
	Use:   "lift",
	Short: "Lift's CLI",
	Long:  "Lift's full featured Command Line Interface",
}

func init() {
	// Debug flag
	RootCmd.PersistentFlags().Bool("debug", false, "Run the service in debug mode")
	viper.BindPFlag("debug", RootCmd.PersistentFlags().Lookup("debug"))

	// Terraform flags
	RootCmd.PersistentFlags().String("terraform_exec_path", "/usr/bin/terraform", "Path to the terraform executable")
	viper.BindPFlag("terraform_exec_path", RootCmd.PersistentFlags().Lookup("terraform_exec_path"))

	// Database configuration flags
	RootCmd.PersistentFlags().String("db_host", "localhost", "The host where the database lives")
	viper.BindPFlag("db_host", RootCmd.PersistentFlags().Lookup("db_host"))
	RootCmd.PersistentFlags().String("db_user", "user", "The database user to connect with to the database")
	viper.BindPFlag("db_user", RootCmd.PersistentFlags().Lookup("db_user"))
	RootCmd.PersistentFlags().String("db_password", "password", "The database user's password")
	viper.BindPFlag("db_password", RootCmd.PersistentFlags().Lookup("db_password"))
	RootCmd.PersistentFlags().String("db_name", "lift", "The database name")
	viper.BindPFlag("db_name", RootCmd.PersistentFlags().Lookup("db_name"))
	RootCmd.PersistentFlags().Int("db_port", 5432, "The port in which the database is listening on")
	viper.BindPFlag("db_port", RootCmd.PersistentFlags().Lookup("db_port"))

	RootCmd.AddCommand(
		serverCmd,
		observerCmd,
		fetchCmd,
	)
}

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
	RootCmd.PersistentFlags().Bool("debug", false, "Run the service in debug mode")
	viper.BindPFlag("debug", RootCmd.PersistentFlags().Lookup("debug"))

	RootCmd.AddCommand(
		serverCmd,
	)
}

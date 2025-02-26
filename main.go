package main

import (
	"fmt"
	"os"

	"github.com/galihrivanto/kotak/cli"
	"github.com/galihrivanto/kotak/config"
	"github.com/galihrivanto/kotak/db"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "kotak",
	Short: "Kotak is simple temporary email service",
}

func main() {
	rootCmd.AddCommand(cli.ServerCmd)

	rootCmd.PersistentFlags().StringP("config", "c", "./config.yaml", "Config file (default is ./config.yaml)")

	rootCmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		// load config
		c := config.Load(rootCmd.Flag("config").Value.String())
		cmd.SetContext(config.WithContext(cmd.Context(), c))

		// setup database
		dbInstance, err := db.New(c.Database)
		if err != nil {
			rootCmd.Print(err)
			os.Exit(1)
		}
		cmd.SetContext(db.WithContext(cmd.Context(), dbInstance))
	}

	rootCmd.PersistentPostRun = func(cmd *cobra.Command, args []string) {
		// teardown database
		fmt.Println("Closing database...")
		db.FromContext(cmd.Context()).Close()
	}

	if err := rootCmd.Execute(); err != nil {
		rootCmd.Print(err)
		os.Exit(1)
	}
}

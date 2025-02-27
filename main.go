package main

import (
	"fmt"
	"os"

	"github.com/galihrivanto/kotak/cli"
	"github.com/galihrivanto/kotak/config"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "kotak",
	Short: "Kotak is simple temporary email service",
}

func main() {
	rootCmd.AddCommand(cli.ServerCmd)
	rootCmd.AddCommand(cli.SendEmailCmd)

	rootCmd.PersistentFlags().StringP("config", "c", "./config.yaml", "Config file (default is ./config.yaml)")

	rootCmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		// load config
		fmt.Println("Loading configuration...")
		c := config.Load(rootCmd.Flag("config").Value.String())
		cmd.SetContext(config.WithContext(cmd.Context(), c))

	}

	if err := rootCmd.Execute(); err != nil {
		rootCmd.Print(err)
		os.Exit(1)
	}
}

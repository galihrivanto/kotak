package cli

import (
	"context"
	"fmt"
	"os"

	"github.com/galihrivanto/kotak/config"
	"github.com/galihrivanto/kotak/db"
	"github.com/galihrivanto/kotak/module"
	"github.com/galihrivanto/runner"
	"github.com/spf13/cobra"

	_ "github.com/galihrivanto/kotak/module/http"
	_ "github.com/galihrivanto/kotak/module/smtp"
)

var ServerCmd = &cobra.Command{
	Use:   "server",
	Short: "Start the server",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Server started")

		c := config.FromContext(cmd.Context())
		db := db.FromContext(cmd.Context())

		runner.Run(cmd.Context(), func(ctx context.Context) error {
			return module.Start(ctx, c, db)
		}).Handle(func(sig os.Signal) {
			if sig == os.Interrupt {
				fmt.Println("\nStopping server")
				module.Stop()
			}
		})
	},
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		c := config.FromContext(cmd.Context())

		// setup database
		fmt.Println("Setting up database...")
		dbInstance, err := db.New(c.Database)
		if err != nil {
			cmd.Print(err)
			os.Exit(1)
		}
		cmd.SetContext(db.WithContext(cmd.Context(), dbInstance))
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		// teardown database
		fmt.Println("Closing database...")
		db.FromContext(cmd.Context()).Close()
	},
}

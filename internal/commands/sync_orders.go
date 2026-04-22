package commands

import (
	"fmt"
	"stationhub-api/config"
	"stationhub-api/deps"

	"github.com/spf13/cobra"
)

func NewSyncOrdersCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "sync:orders",
		Short: "Synchronize orders from external API",
		Long:  "Fetches and synchronizes order data from the external API into the database",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := config.Load()
			db := config.ConnectDatabase(cfg)
			d := deps.New(db, cfg)

			fmt.Println("🔄 Starting orders synchronization...")
			
			// TODO: Implement your sync logic here
			// Example: return d.OrderService.Sync()
			_ = d

			fmt.Println("✅ Orders synchronized successfully")
			return nil
		},
	}
}

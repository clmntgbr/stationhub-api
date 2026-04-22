package main

import (
	"fmt"
	"os"
	"stationhub-api/command"

	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "cli",
		Short: "StationHub CLI - internal commands for background tasks",
		Long:  "StationHub CLI provides commands for synchronization, cleanup, and maintenance tasks",
	}

	rootCmd.AddCommand(
		command.NewGasPricesUpdateCommand(),
	)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

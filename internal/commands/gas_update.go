package commands

import (
	"fmt"
	"stationhub-api/service"

	"github.com/spf13/cobra"
)

func NewGasUpdateCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "gas:update",
		Short: "Update gas prices",
		Long:  "Fetches and updates gas prices",
		RunE: func(cmd *cobra.Command, args []string) error {
			gasFileService := service.NewGasFileService()

			zipFilePath, err := gasFileService.DownloadGasFile()
			if err != nil {
				return fmt.Errorf("failed to download gas file: %w", err)
			}

			extractedPath, err := gasFileService.Extract(zipFilePath)
			if err != nil {
				return fmt.Errorf("failed to extract gas file: %w", err)
			}

			fmt.Println("✅ Gas file extracted to:", extractedPath)

			return nil
		},
	}
}

package command

import (
	"fmt"
	"stationhub-api/config"
	"stationhub-api/repository"
	"stationhub-api/service"

	"github.com/spf13/cobra"
)

func NewGasPricesUpdateCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "gas:prices:update",
		Short: "Update gas prices",
		Long:  "Fetches and updates gas prices",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := config.Load()
			db := config.ConnectDatabase(cfg)

			gasFileService := service.NewGasFileService()

			stationRepository := repository.NewStationRepository(db)
			addressRepository := repository.NewAddressRepository(db)
			currentPriceRepository := repository.NewCurrentPriceRepository(db)
			priceHistoryRepository := repository.NewPriceHistoryRepository(db)

			gasPricesUpdateService := service.NewGasPricesUpdateService(
				stationRepository,
				addressRepository,
				currentPriceRepository,
				priceHistoryRepository,
			)

			zipFilePath, err := gasFileService.DownloadGasFile()
			if err != nil {
				return fmt.Errorf("failed to download gas file: %w", err)
			}

			extractedPath, err := gasFileService.Extract(zipFilePath)
			if err != nil {
				return fmt.Errorf("failed to extract gas file: %w", err)
			}

			err = gasPricesUpdateService.UpdateGasPrices(extractedPath)
			if err != nil {
				return fmt.Errorf("failed to update gas prices: %w", err)
			}

			err = gasFileService.Delete(zipFilePath)
			if err != nil {
				return fmt.Errorf("failed to delete zip file: %w", err)
			}

			err = gasFileService.Delete(extractedPath)
			if err != nil {
				return fmt.Errorf("failed to delete extracted file: %w", err)
			}

			return nil
		},
	}
}

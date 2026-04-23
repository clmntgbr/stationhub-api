package config

import (
	"fmt"
	"log"
	"time"

	"stationhub-api/domain"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func ConnectDatabase(cfg *Config) *gorm.DB {
	logLevel := logger.Warn
	if cfg.Environment == "development" {
		logLevel = logger.Info
	}

	db, err := gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})
	if err != nil {
		log.Fatal("failed to connect to database: ", err)
	}

	return db
}

func AutoMigrate(db *gorm.DB) {

	InitDB(db)

	err := db.AutoMigrate(
		&domain.User{},
		&domain.Station{},
		&domain.Address{},
		&domain.GooglePlace{},
		&domain.CurrentPrice{},
	)

	CreatePriceHistoryPartitionBase(db)
	InitPartitions(db)

	if err != nil {
		log.Fatal("failed to migrate database: ", err)
	}
}

func InitDB(db *gorm.DB) {
	err := db.Exec(`CREATE EXTENSION IF NOT EXISTS postgis;`).Error

	if err != nil {
		log.Fatal("failed to create postgis extension: ", err)
	}
}

func CreatePriceHistoryPartitionBase(db *gorm.DB) {
	err := db.Exec(`
		CREATE TABLE IF NOT EXISTS price_histories (
			id uuid DEFAULT gen_random_uuid(),
			value decimal NOT NULL,
			currency text NOT NULL,
			type text NOT NULL,
			type_id bigint NOT NULL,
			date timestamptz NOT NULL,
			station_id uuid NOT NULL,
			created_at timestamptz,
			updated_at timestamptz,

			PRIMARY KEY (id, date),

			CONSTRAINT fk_price_histories_station
			FOREIGN KEY (station_id) REFERENCES stations(id)
		)
		PARTITION BY RANGE (date);
		CREATE INDEX IF NOT EXISTS idx_price_histories_date ON price_histories (date);
		CREATE INDEX IF NOT EXISTS idx_price_histories_station_id_date ON price_histories (station_id, date);
	`).Error

	if err != nil {
		log.Fatal("failed to create partitioned table: ", err)
	}
}

func InitPartitions(db *gorm.DB) {
	now := time.Now().UTC()

	_ = CreatePartitionsRange(db, time.Date(now.Year()-1, 1, 1, 0, 0, 0, 0, time.UTC), 12)

	_ = CreatePartitionsRange(db, time.Date(now.Year(), 1, 1, 0, 0, 0, 0, time.UTC), 24)
}

func CreateInitialPartitions(db *gorm.DB) error {
	now := time.Now().UTC()

	start := time.Date(now.Year()-1, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(now.Year()+5, 1, 1, 0, 0, 0, 0, time.UTC)

	for t := start; t.Before(end); t = t.AddDate(0, 1, 0) {
		if err := CreatePartitionIfNotExists(db, t.Year(), t.Month()); err != nil {
			return err
		}
	}

	return nil
}

func CreatePartitionsRange(db *gorm.DB, start time.Time, months int) error {
	for i := 0; i < months; i++ {
		t := start.AddDate(0, i, 0)

		year := t.Year()
		month := t.Month()

		if err := CreatePartitionIfNotExists(db, year, month); err != nil {
			return err
		}
	}
	return nil
}

func CreatePartitionIfNotExists(db *gorm.DB, year int, month time.Month) error {
	start := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(0, 1, 0)

	name := fmt.Sprintf("price_histories_%04d_%02d", year, month)

	sql := fmt.Sprintf(`
CREATE TABLE IF NOT EXISTS "%s"
PARTITION OF price_histories
FOR VALUES FROM ('%s') TO ('%s');
CREATE INDEX IF NOT EXISTS idx_%s ON %s (station_id, date);
`, name, start.Format("2006-01-02"), end.Format("2006-01-02"), name, name)

	err := db.Exec(sql).Error
	if err != nil {
		return fmt.Errorf("failed to create partition: %w", err)
	}

	return nil
}

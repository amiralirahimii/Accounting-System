package db

import (
	"accountingsystem/config"
	"accountingsystem/internal/constants"
	"database/sql"
	"fmt"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func loadVars() (string, string, string, string, string, string, error) {
	host, err := config.GetEnv("DB_HOST")
	if err != nil {
		return "", "", "", "", "", "", fmt.Errorf("%w: DB_HOST", constants.ErrEnvNotFound)
	}

	user, err := config.GetEnv("DB_USER")
	if err != nil {
		return "", "", "", "", "", "", fmt.Errorf("%w: DB_USER", constants.ErrEnvNotFound)
	}

	password, err := config.GetEnv("DB_PASSWORD")
	if err != nil {
		return "", "", "", "", "", "", fmt.Errorf("%w: DB_PASSWORD", constants.ErrEnvNotFound)
	}

	dbName, err := config.GetEnv("DB_NAME")
	if err != nil {
		return "", "", "", "", "", "", fmt.Errorf("%w: DB_NAME", constants.ErrEnvNotFound)
	}

	port, err := config.GetEnv("DB_PORT")
	if err != nil {
		return "", "", "", "", "", "", fmt.Errorf("%w: DB_PORT", constants.ErrEnvNotFound)
	}

	sslMode, err := config.GetEnv("DB_SSLMODE")
	if err != nil {
		return "", "", "", "", "", "", fmt.Errorf("%w: DB_SSLMODE", constants.ErrEnvNotFound)
	}

	return host, user, password, dbName, port, sslMode, nil
}

func configureConnectionPool(sqlDB *sql.DB) {
	sqlDB.SetMaxOpenConns(20)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)
}

func Init() error {
	host, user, password, dbName, port, sslMode, err := loadVars()
	if err != nil {
		log.Fatalf("Failed to load database variables: %v", err)
		return err
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		host, user, password, dbName, port, sslMode)

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
		return err
	}

	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatalf("Failed to get sql.DB from GORM: %v", err)
		return err
	}

	configureConnectionPool(sqlDB)

	return nil
}

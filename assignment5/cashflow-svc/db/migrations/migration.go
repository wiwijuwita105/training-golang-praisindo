package migrations

import (
	"assignment5/cashflow-svc/internal/config"
	"assignment5/cashflow-svc/internal/entity"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

func CreateDB() {
	// DSN (Data Source Name)
	dsn := "host=" + config.PostgresHost + " user=" + config.PostgresUser + " password=" + config.PostgresPassword + " port=" + config.PostgresPort + " sslmode=" + config.PostgresSSLMode

	// Opening a DB connection
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{SkipDefaultTransaction: true})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
	log.Println("Database connection successfully opened")

	// Check if the database exists
	var exists bool
	err = db.Raw("SELECT EXISTS(SELECT datname FROM pg_catalog.pg_database WHERE datname = ?)", config.PostgresDB).Scan(&exists).Error
	if err != nil {
		log.Fatalln(err)
	} else {
		log.Println("Database exists:", exists, config.PostgresDB)
	}

	// Create the database if it does not exist
	if !exists {
		err = db.Exec("CREATE DATABASE " + config.PostgresDB).Error
		if err != nil {
			log.Fatalln(err)
		} else {
			log.Println("Database created successfully")
		}
	}
}

func MigrationTable(gormDB *gorm.DB) {
	// Migrate the schema
	err := gormDB.AutoMigrate(entity.User{}, entity.Wallet{}, entity.TransactionCategory{}, entity.Transaction{})
	if err != nil {
		fmt.Println("Failed to migrate database schema:", err)
	} else {
		fmt.Println("Database schema migrated successfully")
	}
}

package migrations

import (
	"assignment6/config"
	"assignment6/entity"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

func CreateDBUser() {
	// DSN (Data Source Name)
	dsn := "host=" + config.PostgresUserHost + " user=" + config.PostgresUser + " password=" + config.PostgresPassword + " port=" + config.PostgresUserPort + " sslmode=" + config.PostgresSSLMode

	// Opening a DB connection
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{SkipDefaultTransaction: true})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
	log.Println("Database connection successfully opened")

	// Check if the database exists
	var exists bool
	err = db.Raw("SELECT EXISTS(SELECT datname FROM pg_catalog.pg_database WHERE datname = ?)", config.PostgresDBUser).Scan(&exists).Error
	if err != nil {
		log.Fatalln(err)
	} else {
		log.Println("Database exists:", exists, config.PostgresDBUser)
	}

	// Create the database if it does not exist
	if !exists {
		err = db.Exec("CREATE DATABASE " + config.PostgresUser).Error
		if err != nil {
			log.Fatalln(err)
		} else {
			log.Println("Database created successfully")
		}
	}
}

func MigrationTableUser(gormDB *gorm.DB) {
	// Migrate the schema
	err := gormDB.AutoMigrate(entity.User{})
	if err != nil {
		fmt.Println("Failed to migrate database schema:", err)
	} else {
		fmt.Println("Database schema migrated successfully")
	}
}

func CreateDBWallet() {
	// DSN (Data Source Name)
	dsn := "host=" + config.PostgresWalletHost + " user=" + config.PostgresUser + " password=" + config.PostgresPassword + " port=" + config.PostgresWalletPort + " sslmode=" + config.PostgresSSLMode

	// Opening a DB connection
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{SkipDefaultTransaction: true})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
	log.Println("Database connection successfully opened")

	// Check if the database exists
	var exists bool
	err = db.Raw("SELECT EXISTS(SELECT datname FROM pg_catalog.pg_database WHERE datname = ?)", config.PostgresDBWallet).Scan(&exists).Error
	if err != nil {
		log.Fatalln(err)
	} else {
		log.Println("Database exists:", exists, config.PostgresDBWallet)
	}

	// Create the database if it does not exist
	if !exists {
		err = db.Exec("CREATE DATABASE " + config.PostgresUser).Error
		if err != nil {
			log.Fatalln(err)
		} else {
			log.Println("Database created successfully")
		}
	}
}

func MigrationTableWallet(gormDB *gorm.DB) {
	// Migrate the schema
	err := gormDB.AutoMigrate(entity.Wallet{}, entity.TransactionCategory{}, entity.Transaction{})
	if err != nil {
		fmt.Println("Failed to migrate database schema:", err)
	} else {
		fmt.Println("Database schema migrated successfully")
	}
}

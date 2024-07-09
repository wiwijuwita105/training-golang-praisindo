package main

import (
	"fmt"
	"log"
	"wallet_svc/entity"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// setup gorm connection
	dsn := "postgresql://postgres:postgres@localhost:5432/db_wallet_svc"
	gormDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{SkipDefaultTransaction: true})
	if err != nil {
		log.Fatalln(err)
	}

	// Migrate the schema
	err = gormDB.AutoMigrate(entity.Wallet{}, entity.Transaction{})
	if err != nil {
		fmt.Println("Failed to migrate database schema:", err)
	} else {
		fmt.Println("Database schema migrated successfully")
	}

	log.Println("Running grpc server in port :50052")
}

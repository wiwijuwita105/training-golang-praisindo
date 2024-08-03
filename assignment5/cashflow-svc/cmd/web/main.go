package main

import (
	"assignment5/cashflow-svc/internal/config"
	"assignment5/cashflow-svc/internal/delivery/http/route"
	"assignment5/cashflow-svc/internal/entity"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

func main() {
	//Migration Create DB
	CreateDB()

	//Reconnection to DB
	dsn := "host=" + config.PostgresHost + " user=" + config.PostgresUser + " password=" + config.PostgresPassword + " port=" + config.PostgresPort + " dbname=" + config.PostgresDB + " sslmode=" + config.PostgresSSLMode
	// Opening a DB connection
	gormDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{SkipDefaultTransaction: true})
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("Database connection successfully opened")

	// Migrate the schema
	err = gormDB.AutoMigrate(entity.User{}, entity.Wallet{}, entity.TransactionCategory{}, entity.Transaction{})
	if err != nil {
		fmt.Println("Failed to migrate database schema:", err)
	} else {
		fmt.Println("Database schema migrated successfully")
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     config.RedisHost + ":" + config.RedisiPort,
		Password: config.RedisPassword, // no password set
		DB:       config.RedisDb,       // use default DB
	})
	log.Println(rdb)
	//handler gateway
	gwServer := gin.Default()
	// Routes Gateway
	route.SetupRouter(gwServer)
	//gwServer.Group("v1/*{grpc_gateway}").Any("", gin.WrapH(gwmux))
	//gwServer.GET("/:shortlink", redirectHandler)
	log.Println("Running grpc gateway server in port :8080")
	_ = gwServer.Run()
}

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

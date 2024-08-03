package main

import (
	"assignment5/cashflow-svc/db/migrations"
	"assignment5/cashflow-svc/internal/config"
	"assignment5/cashflow-svc/internal/delivery/http/route"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

func main() {
	//Migration DB
	migrations.CreateDB()

	//Reconnection to DB
	dsn := "host=" + config.PostgresHost + " user=" + config.PostgresUser + " password=" + config.PostgresPassword + " port=" + config.PostgresPort + " dbname=" + config.PostgresDB + " sslmode=" + config.PostgresSSLMode
	// Opening a DB connection
	gormDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{SkipDefaultTransaction: true})
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("Database connection successfully opened")

	//Migration Table
	migrations.MigrationTable(gormDB)

	//Redis connection
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

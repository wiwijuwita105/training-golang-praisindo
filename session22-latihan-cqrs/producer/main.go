package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"session22-latihan-cqrs/config"
	"session22-latihan-cqrs/producer/handler"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	//Koneksion db using gorm
	gormDB, err := gorm.Open(postgres.Open(config.DBReadConn), &gorm.Config{SkipDefaultTransaction: true})
	if err != nil {
		log.Println(gormDB)
		log.Fatalln(err)

	}

	userHandler := handler.NewUserHandler(gormDB)

	r.POST("/users", userHandler.CreateUser)
	r.GET("/users", userHandler.GetAllUsers)

	fmt.Printf("Kafka PRODUCER started at http://localhost%s\n", config.KafkaProducerPort)

	if err := r.Run(config.KafkaProducerPort); err != nil {
		log.Printf("failed to run the server: %v", err)
	}
}

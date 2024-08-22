package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"latihan-cqrs/entity"
	"latihan-cqrs/handler"
	"latihan-cqrs/repository/postgres_gorm"
	"latihan-cqrs/router"
	"latihan-cqrs/service"
	"log"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	//Koneksion db using gorm
	dsn := "postgresql://postgres:postgres@localhost:5434/db_cqrs_write?sslmode=disable"
	gormDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{SkipDefaultTransaction: true})
	if err != nil {
		log.Println(gormDB)
		log.Fatalln(err)

	}

	userRepo := postgres_gorm.NewUserRepository(gormDB)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	// Routes
	router.SetupRouter(r, userHandler)

	fmt.Printf("Kafka PRODUCER started at http://localhost%s\n", entity.KafkaProducerPort)

	if err := r.Run(entity.KafkaProducerPort); err != nil {
		log.Printf("failed to run the server: %v", err)
	}

}

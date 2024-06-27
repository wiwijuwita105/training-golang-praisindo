package main

import (
	"assignment1/entity"
	"assignment1/repository/postgres_gorm"
	"assignment1/router"
	"fmt"
	"log"

	"assignment1/handler"
	"assignment1/service"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	//Koneksion db using gorm
	dsn := "postgresql://postgres:postgres@localhost:5432/db_training_golang_assignment1"
	gormDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{SkipDefaultTransaction: true})
	if err != nil {
		log.Println(gormDB)
		log.Fatalln(err)
	}

	// Migrate the schema
	err = gormDB.AutoMigrate(entity.User{}, entity.Submission{})
	if err != nil {
		fmt.Println("Failed to migrate database schema:", err)
	} else {
		fmt.Println("Database schema migrated successfully")
	}

	// uncomment to use postgres gorm
	userRepo := postgres_gorm.NewUserRepository(gormDB)

	// service and handler declaration
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	// uncomment to use postgres gorm
	submissionRepo := postgres_gorm.NewSubmissionRepository(gormDB)

	// service and handler declaration
	submissionService := service.NewSubmissionService(submissionRepo)
	submissionHandler := handler.NewSubmissionHandler(submissionService)

	// Routes
	router.SetupRouter(r, userHandler, submissionHandler)

	// Run the server
	log.Println("Running server on port 8080")
	r.Run(":8080")
}

package main

import (
	"log"
	"session4-unit-testing-crud/entity"
	"session4-unit-testing-crud/handler"
	"session4-unit-testing-crud/repository/slice"
	"session4-unit-testing-crud/router"
	"session4-unit-testing-crud/service"

	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	//setup service
	var mockUserDBInSlice []entity.User
	userRepo := slice.NewUserRepository(mockUserDBInSlice)
	userService := service.NewUserService(userRepo)
	UserHandler := handler.NewUserHandler(userService)

	//Routes
	router.SetupRouter(r, UserHandler)

	//run server
	log.Println("Running server on port 8080")
	r.Run(":8080")
}

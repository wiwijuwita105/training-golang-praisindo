package main

import (
	"log"
	"session5-validator/entity"
	"session5-validator/handler"
	"session5-validator/repository/slice"
	"session5-validator/router"
	"session5-validator/service"

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

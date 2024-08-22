package router

import (
	"github.com/gin-gonic/gin"
	"latihan-cqrs/handler"
)

func SetupRouter(r *gin.Engine, userHandler handler.IUserHandler) {
	userEndpoint := r.Group("/users")
	userEndpoint.POST("/", userHandler.CreateUser)
	userEndpoint.GET("/", userHandler.GetAllUsers)
}

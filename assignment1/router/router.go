package router

import (
	"assignment1/handler"
	"assignment1/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter(r *gin.Engine, userHandler handler.IUserHandler, submissionHandler handler.ISubmissionHandler) {
	//endpoint users group
	//mengutur endpoint public
	usersPublicEndpoint := r.Group("/users")
	usersPublicEndpoint.GET("/:id", userHandler.GetUser)
	usersPublicEndpoint.GET("/", userHandler.GetAllUsersWithRIsk)

	// Mengatur endpoint privat untuk pengguna dengan middleware autentikasi
	usersPrivateEndpoint := r.Group("/users")
	usersPrivateEndpoint.Use(middleware.AuthMiddleware())
	usersPrivateEndpoint.POST("", userHandler.CreateUser)
	usersPrivateEndpoint.POST("/", userHandler.CreateUser)
	usersPrivateEndpoint.PUT("/:id", userHandler.UpdateUser)
	usersPrivateEndpoint.DELETE("/:id", userHandler.DeleteUser)

	//endpoint submisson group
	submissionPublicEndpoint := r.Group("/submissions")
	submissionPublicEndpoint.GET("/:id", submissionHandler.GetSubmission)
	submissionPublicEndpoint.GET("/", submissionHandler.GetAllSubmissions)

	submissionPrivateEndpoint := r.Group("/submissions")
	submissionPrivateEndpoint.Use(middleware.AuthMiddleware())
	submissionPrivateEndpoint.POST("", submissionHandler.CreateSubmission)
	submissionPrivateEndpoint.POST("/", submissionHandler.CreateSubmission)
	submissionPrivateEndpoint.DELETE("/:id", submissionHandler.DeleteSubmission)
}

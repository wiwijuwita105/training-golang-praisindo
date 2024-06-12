package router

import (
	"session1-intorduction-gin/handler"
	"session1-intorduction-gin/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter(r *gin.Engine) {
	r.GET("/", handler.RootHandler)

	// r.POST("/", handler.PostHandler)

	// Tambahkan middleware AuthMiddleware ke rute yang memerlukan autentikasi
	privateEndpoint := r.Group("/private")
	privateEndpoint.Use(middleware.AuthMiddleware())
	{
		privateEndpoint.POST("/post", handler.PostHandler)
	}

}

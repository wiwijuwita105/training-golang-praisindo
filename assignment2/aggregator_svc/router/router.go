package router

import (
	"aggregator_svc/handler"
	"aggregator_svc/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRouter(r *gin.Engine, aggregator handler.IAggregatorHandler) {
	//r.GET("/health-check", healthcheck.HealthCheck)

	userEndpoint := r.Group("/user")
	userEndpoint.Use(middleware.AuthMiddleware())
	userEndpoint.GET("/:userId", aggregator.GetUser)

}

package route

import (
	"assignment5/cashflow-svc/internal/delivery/http/handler"
	"github.com/gin-gonic/gin"
)

func SetupRouter(r *gin.Engine) {
	// Register health check endpoint
	r.GET("/", handler.HealthCheckHandler)

	// Register other routes here...
}

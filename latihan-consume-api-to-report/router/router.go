package router

import (
	"github.com/gin-gonic/gin"
	"latihan-consume-api-to-report/handler"
)

func SetupRouter(r *gin.Engine, handler handler.IWeatherHandler) {
	weatherEndpoint := r.Group("/weather")
	weatherEndpoint.GET("/generateToXLSX", handler.GenerateToXLSX)
	weatherEndpoint.GET("/generateToPDF", handler.GenerateToPDF)
}

package main

import (
	"github.com/gin-gonic/gin"
	"latihan-consume-api-to-report/handler"
	"latihan-consume-api-to-report/router"
	"latihan-consume-api-to-report/service"
	"log"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	serviceWeather := service.NewWeatherService()
	handlerWeather := handler.NewWeatherHandler(serviceWeather)
	// Routes
	router.SetupRouter(r, handlerWeather)

	// Run the server
	log.Println("Running server on port 8080")
	r.Run(":8080")
}

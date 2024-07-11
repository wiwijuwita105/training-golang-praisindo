package main

import (
	"aggregator_svc/config"
	"aggregator_svc/handler"
	"aggregator_svc/router"
	"aggregator_svc/service"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	// Inisialisasi router Gin
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	// Initialize service user
	userClient := config.InitUserSvc()
	if userClient == nil {
		log.Fatal("Failed to initialize User Service Client")
	}

	// Initialize wallet service
	walletClient := config.InitWalletSvc()
	if walletClient == nil {
		log.Fatal("Failed to initialize User Service Client")
	}

	aggregationService := service.NewAggregatorService(userClient, walletClient)
	aggregatorHandler := handler.NewAggregatorHandler(*aggregationService)

	// Definisikan route
	router.SetupRouter(r, aggregatorHandler)

	// Jalankan server pada port 8080
	log.Println("Running server on port 8080")
	r.Run(":8080")
}

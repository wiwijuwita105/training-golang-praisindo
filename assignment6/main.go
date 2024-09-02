package main

import (
	"assignment6/config"
	"assignment6/entity"
	"assignment6/handler"
	"assignment6/repository/postgres_gorm"
	"assignment6/router"
	"assignment6/service"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"log/slog"
	"os"
	"time"
)

func main() {
	// Inisialisasi GIN
	r := gin.Default()

	time.Sleep(5 * time.Second)

	//Connection to BD Wallet
	walletDB, err := gorm.Open(postgres.Open(config.DBWalletDSN), &gorm.Config{})
	if err != nil {
		slog.Error("error opening read database", slog.String("dsn", config.DBWalletDSN), slog.Any("err", err))
		os.Exit(1)
	}

	//Migration Table
	log.Println("starting migration")
	if err := walletDB.AutoMigrate(entity.User{}, entity.Wallet{}, entity.TransactionCategory{}, entity.Transaction{}); err != nil {
		log.Println(err)
	}
	log.Println("finish migration")

	//assign repo dan service
	repoUser := postgres_gorm.NewUserRepository(walletDB)
	serviceUser := service.NewUserService(repoUser)
	userHandler := handler.NewUserHandler(serviceUser)

	walletRepo := postgres_gorm.NewWalletRepository(walletDB)
	walletService := service.NewWalletService(walletRepo)

	categoryRepo := postgres_gorm.NewCategoryRepository(walletDB)
	categoryService := service.NewCategoryService(categoryRepo)

	transactionRepo := postgres_gorm.NewTransactionRepository(walletDB)
	transactionService := service.NewTransactionService(walletDB, walletRepo, categoryRepo, transactionRepo)

	walletHandler := handler.NewWalletHandler(walletService, categoryService, transactionService)

	router.SetupRouter(r, userHandler, walletHandler)

	// Definisikan route

	// Jalankan server pada port 8080
	log.Println("Running server on port 8080")
	r.Run(":8080")
}

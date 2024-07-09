package main

import (
	"fmt"
	"log"
	"net"
	"wallet_svc/entity"
	grpcHandler "wallet_svc/handler/grpc"
	pb "wallet_svc/proto/wallet_service/v1"
	"wallet_svc/repository/postgres_gorm"
	"wallet_svc/service"

	"google.golang.org/grpc"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// setup gorm connection
	dsn := "postgresql://postgres:postgres@localhost:5432/db_wallet_svc"
	gormDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{SkipDefaultTransaction: true})
	if err != nil {
		log.Fatalln(err)
	}

	// Migrate the schema
	err = gormDB.AutoMigrate(entity.Wallet{}, entity.Transaction{})
	if err != nil {
		fmt.Println("Failed to migrate database schema:", err)
	} else {
		fmt.Println("Database schema migrated successfully")
	}

	// setup service
	walletRepo := postgres_gorm.NewWalletRepository(gormDB)
	walletService := service.NewWalletService(walletRepo)
	walletHandler := grpcHandler.NewWalletHandler(walletService)

	// Run the grpc server
	grpcServer := grpc.NewServer()
	pb.RegisterWalletServiceServer(grpcServer, walletHandler)
	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Println("Running grpc server in port :50052")
	_ = grpcServer.Serve(lis)
}

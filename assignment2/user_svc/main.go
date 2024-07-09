package main

import (
	"fmt"
	"log"
	"net"
	"user_svc/entity"
	grpcHandler "user_svc/handler/grpc"
	pb "user_svc/proto/user_service/v1"
	"user_svc/repository/postgres_gorm"
	"user_svc/service"

	"google.golang.org/grpc"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// setup gorm connection
	dsn := "postgresql://postgres:postgres@localhost:5432/db_user_svc"
	gormDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{SkipDefaultTransaction: true})
	if err != nil {
		log.Fatalln(err)
	}

	// Migrate the schema
	err = gormDB.AutoMigrate(entity.User{})
	if err != nil {
		fmt.Println("Failed to migrate database schema:", err)
	} else {
		fmt.Println("Database schema migrated successfully")
	}

	// setup service
	userRepo := postgres_gorm.NewUserRepository(gormDB)
	userService := service.NewUserService(userRepo)
	userHandler := grpcHandler.NewUserHandler(userService)

	// Run the grpc server
	grpcServer := grpc.NewServer()
	pb.RegisterUserServiceServer(grpcServer, userHandler)
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Println("Running grpc server in port :50051")
	_ = grpcServer.Serve(lis)
}

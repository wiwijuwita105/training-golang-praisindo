package main

import (
	"assignment5/cashflow-svc/db/migrations"
	"assignment5/cashflow-svc/internal/config"
	grpcHandler "assignment5/cashflow-svc/internal/delivery/grpc/handler"
	"assignment5/cashflow-svc/internal/delivery/http/route"
	pb "assignment5/cashflow-svc/internal/proto/user_service/v1"
	pb2 "assignment5/cashflow-svc/internal/proto/wallet_service/v1"
	"assignment5/cashflow-svc/internal/repository/postgres_gorm"
	"assignment5/cashflow-svc/internal/service"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net"
	"net/http"
	"time"
)

func main() {
	//Migration DB
	migrations.CreateDB()

	//Reconnection to DB
	dsn := "host=" + config.PostgresHost + " user=" + config.PostgresUser + " password=" + config.PostgresPassword + " port=" + config.PostgresPort + " dbname=" + config.PostgresDB + " sslmode=" + config.PostgresSSLMode
	// Opening a DB connection
	gormDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{SkipDefaultTransaction: true})
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("Database connection successfully opened")

	//Migration Table
	migrations.MigrationTable(gormDB)

	//Redis connection
	rdb := redis.NewClient(&redis.Options{
		Addr:     config.RedisHost + ":" + config.RedisiPort,
		Password: config.RedisPassword, // no password set
		DB:       config.RedisDb,       // use default DB
	})
	log.Println(rdb)

	userRepo := postgres_gorm.NewUserRepository(gormDB)
	userService := service.NewUserService(userRepo)
	userHandler := grpcHandler.NewUserHandler(userService)

	walletRepo := postgres_gorm.NewWalletRepository(gormDB)
	walletService := service.NewWalletService(walletRepo)

	categoryRepo := postgres_gorm.NewCategoryRepository(gormDB)
	categoryService := service.NewCategoryService(categoryRepo)

	transactionRepo := postgres_gorm.NewTransactionRepository(gormDB)
	transactionService := service.NewTransactionService(gormDB, walletRepo, categoryRepo, transactionRepo)

	walletHandler := grpcHandler.NewWalletHandler(walletService, categoryService, transactionService)

	// Run the grpc server
	grpcServer := grpc.NewServer()

	//define grpc user service
	pb.RegisterUserServiceServer(grpcServer, userHandler)
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	go func() {
		log.Println("Running grpc server in port :50051")
		_ = grpcServer.Serve(lis)
	}()

	// Run the grpc gateway user service
	conn, err := grpc.NewClient(
		"0.0.0.0:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalln("Failed to dial server:", err)
	}
	gwmux := runtime.NewServeMux()
	if err = pb.RegisterUserServiceHandler(context.Background(), gwmux, conn); err != nil {
		log.Fatalln("Failed to register gateway:", err)
	}

	//define grpc wallet service
	pb2.RegisterWalletServiceServer(grpcServer, walletHandler)
	lisWallet, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	go func() {
		log.Println("Running grpc server in port :50052")
		_ = grpcServer.Serve(lisWallet)
	}()
	time.Sleep(3 * time.Second)

	// Run the grpc gateway wallet service
	connWallet, err := grpc.NewClient(
		"0.0.0.0:50052",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalln("Failed to dial server:", err)
	}
	gwmuxWallet := runtime.NewServeMux()
	if err = pb.RegisterUserServiceHandler(context.Background(), gwmuxWallet, connWallet); err != nil {
		log.Fatalln("Failed to register gateway:", err)
	}

	//using GIN
	gwServer := gin.Default()

	// Routes Gateway using GIN
	route.SetupRouter(gwServer)
	//Route using GRPC
	gwServer.Group("v1/*{grpc_gateway}").Any("", gin.WrapH(gwmux), gin.WrapH(gwServer))
	log.Println("Running grpc gateway server in port :8080")
	if err := http.ListenAndServe(":8080", gwmux); err != nil {
		log.Fatal(err)
	}

	//run gwserver
	_ = gwServer.Run(":8080")
}

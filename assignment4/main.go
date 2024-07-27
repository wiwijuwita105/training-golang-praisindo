package main

import (
	"context"
	"fmt"
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
	"os"
	"shortlink/entity"
	grpcHandler "shortlink/handler/grpc"
	pb "shortlink/proto/shortlink/v1"
	"shortlink/repository/postgres_gorm"
	"shortlink/service"
	"time"
)

func main() {
	// Read environment variables
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")
	redisAddr := os.Getenv("REDIS_ADDR")

	if host == "" || user == "" || password == "" || dbname == "" || port == "" || redisAddr == "" {
		log.Fatalln("Missing one or more environment variables: DB_HOST, DB_USER, DB_PASSWORD, DB_NAME, DB_PORT, REDIS_ADDR")
	}

	// Setup GORM connection
	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", user, password, host, port, dbname)
	gormDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{SkipDefaultTransaction: true})
	if err != nil {
		log.Fatalln("Failed to initialize database:", err)
	}
	// Migrate the schema
	err = gormDB.AutoMigrate(entity.Shortlink{})
	if err != nil {
		fmt.Println("Failed to migrate database schema:", err)
	} else {
		fmt.Println("Database schema migrated successfully")
	}

	//setup connection to redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	shortlinkRepo := postgres_gorm.NewShortlinkRepository(gormDB)
	shortlinkService := service.NewShortlinkService(shortlinkRepo, rdb)
	shortlinkHandler := grpcHandler.NewShortlinkHandler(shortlinkService)

	// Run the grpc server
	grpcServer := grpc.NewServer()
	pb.RegisterShortlinkServiceServer(grpcServer, shortlinkHandler)
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	go func() {
		log.Println("Running grpc server in port :50051")
		_ = grpcServer.Serve(lis)
	}()
	time.Sleep(1 * time.Second)

	// Run the grpc gateway
	conn, err := grpc.NewClient(
		"0.0.0.0:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalln("Failed to dial server:", err)
	}
	gwmux := runtime.NewServeMux()
	if err = pb.RegisterShortlinkServiceHandler(context.Background(), gwmux, conn); err != nil {
		log.Fatalln("Failed to register gateway:", err)
	}

	// Custom handler for redirecting shortlinks
	redirectHandler := func(c *gin.Context) {
		shortlink := c.Param("shortlink")

		resp, err := shortlinkHandler.GetLongUrl(context.Background(), &pb.GetLongUrlRequest{Shortlink: shortlink})
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Shortlink not found"})
			return
		}

		longURL := resp.GetLongUrl()
		c.Redirect(http.StatusFound, longURL)
	}

	// dengan GIN
	gwServer := gin.Default()
	gwServer.Group("v1/*{grpc_gateway}").Any("", gin.WrapH(gwmux))
	gwServer.GET("/:shortlink", redirectHandler)
	log.Println("Running grpc gateway server in port :8080")
	_ = gwServer.Run()

}

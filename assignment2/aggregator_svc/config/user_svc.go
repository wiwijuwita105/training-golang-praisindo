package config

import (
	user_service "aggregator_svc/proto/user_service/v1"
	"google.golang.org/grpc"
	"log"
)

func InitUserSvc() user_service.UserServiceClient {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	return user_service.NewUserServiceClient(conn)
}

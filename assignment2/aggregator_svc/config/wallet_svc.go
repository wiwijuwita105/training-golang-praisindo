package config

import (
	transaction_service "aggregator_svc/proto/transaction_service/v1"
	wallet_service "aggregator_svc/proto/wallet_service/v1"
	"google.golang.org/grpc"
	"log"
)

func InitWalletSvc() wallet_service.WalletServiceClient {
	conn, err := grpc.Dial("localhost:50052", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	return wallet_service.NewWalletServiceClient(conn)
}

func InitTransactionSvc() transaction_service.TransactionServiceClient {
	conn, err := grpc.Dial("localhost:50052", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	return transaction_service.NewTransactionServiceClient(conn)
}

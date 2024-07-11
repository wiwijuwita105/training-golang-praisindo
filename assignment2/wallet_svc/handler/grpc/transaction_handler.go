package grpc

import (
	"context"
	"fmt"
	"log"
	"wallet_svc/entity"
	pb "wallet_svc/proto/transaction_service/v1"
	"wallet_svc/service"
)

type TransactionHandler struct {
	pb.UnimplementedTransactionServiceServer
	transactionService service.ITransactionService
}

func NewTransactionHandler(transactionService service.ITransactionService) *TransactionHandler {
	return &TransactionHandler{
		transactionService: transactionService,
	}
}

func (u *TransactionHandler) CreateTransaction(ctx context.Context, req *pb.CreateTransactionRequest) (*pb.MutationTransResponse, error) {
	createdTransaction, err := u.transactionService.CreateTransaction(ctx, &entity.TransactionRequest{
		FromID: int(req.GetFromID()),
		ToID:   int(req.GetToID()),
		Type:   string(req.GetType()),
		Amount: float64(req.GetAmount()),
	})
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &pb.MutationTransResponse{
		Message: fmt.Sprintf("Success created Transaction with ID %d", createdTransaction.ID),
	}, nil
}

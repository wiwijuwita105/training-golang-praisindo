package grpc

import (
	"context"
	"fmt"
	"google.golang.org/protobuf/types/known/timestamppb"
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
	log.Println(req)
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

func (u *TransactionHandler) GetTransactions(ctx context.Context, req *pb.GetTransactionRequest) (*pb.GetTransactionResponse, error) {
	transactions, err := u.transactionService.GetAllTransactions(ctx, entity.TransactionGetRequest{
		Type:   req.GetType(),
		UserID: int(req.GetUserID()),
		Size:   int(req.GetPageSize()),
		Page:   int(req.GetPage()),
	})
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var transactionProto []*pb.Transaction
	for _, transaction := range transactions.Transaction {
		transactionProto = append(transactionProto, &pb.Transaction{
			Id:        int32(transaction.ID),
			UserID:    int32(transaction.UserID),
			Amount:    float32(transaction.Amount),
			Type:      string(transaction.Type),
			Category:  string(transaction.Category),
			CreatedAt: timestamppb.New(transaction.CreatedAt),
			UpdatedAt: timestamppb.New(transaction.UpdatedAt),
		})
	}

	return &pb.GetTransactionResponse{
		Transactions: transactionProto,
		TotalCount:   int32(transactions.CountData),
	}, nil
}

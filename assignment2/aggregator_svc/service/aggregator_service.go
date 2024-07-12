package service

import (
	"aggregator_svc/entity"
	transaction_service "aggregator_svc/proto/transaction_service/v1"
	user_service "aggregator_svc/proto/user_service/v1"
	wallet_service "aggregator_svc/proto/wallet_service/v1"
	"context"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
	"time"
)

type IAggregatorService interface {
	GetUser(ctx context.Context, id int) (entity.UserResponse, error)
	TopupTransaction(ctx context.Context, request entity.TransactionRequest) (entity.TransactionResponse, error)
	TransferTransaction(ctx context.Context, request entity.TransactionRequest) (entity.TransactionResponse, error)
}

type AggregatorService struct {
	userService        user_service.UserServiceClient
	walletService      wallet_service.WalletServiceClient
	transactionService transaction_service.TransactionServiceClient
}

func NewAggregatorService(userService user_service.UserServiceClient, walletService wallet_service.WalletServiceClient, transactionService transaction_service.TransactionServiceClient) *AggregatorService {
	return &AggregatorService{
		userService:        userService,
		walletService:      walletService,
		transactionService: transactionService,
	}
}

func (svc *AggregatorService) TopupTransaction(ctx context.Context, request entity.TransactionRequest) (entity.TransactionResponse, error) {
	transactionResp, err := svc.transactionService.CreateTransaction(ctx, &transaction_service.CreateTransactionRequest{
		ToID:   int32(request.ToID),
		Type:   "topup",
		Amount: float32(request.Amount),
	})
	if err != nil {
		return entity.TransactionResponse{}, err
	}
	log.Println(transactionResp)
	return entity.TransactionResponse{Message: transactionResp.Message}, nil
}

func (svc *AggregatorService) TransferTransaction(ctx context.Context, request entity.TransactionRequest) (entity.TransactionResponse, error) {
	log.Println(request)
	transactionResp, err := svc.transactionService.CreateTransaction(ctx, &transaction_service.CreateTransactionRequest{
		FromID: request.FromID,
		ToID:   int32(request.ToID),
		Type:   "transfer",
		Amount: float32(request.Amount),
	})
	log.Println(transactionResp)
	log.Println(err)
	if err != nil {
		return entity.TransactionResponse{}, err
	}
	log.Println(transactionResp)
	return entity.TransactionResponse{Message: transactionResp.Message}, nil
}

func (svc *AggregatorService) GetUser(ctx context.Context, id int) (entity.UserResponse, error) {
	userId := int32(id)
	userResp, err := svc.userService.GetUserByID(ctx, &user_service.GetUserByIDRequest{Id: userId})
	if err != nil {
		return entity.UserResponse{}, err
	}
	log.Println(userId)
	walletResp, err := svc.walletService.GetWalletByUserID(ctx, &wallet_service.GetWalletByUserIDRequest{UserID: userId})
	if err != nil {
		return entity.UserResponse{}, err
	}

	user := entity.UserResponse{
		ID:        userResp.User.Id,
		Name:      userResp.User.Name,
		Email:     userResp.User.Email,
		Balance:   float64(walletResp.Wallet.Balance),
		CreatedAt: convertTimestampToTime(userResp.User.CreatedAt),
		UpdatedAt: convertTimestampToTime(userResp.User.UpdatedAt),
	}
	return user, nil
}

func convertTimestampToTime(timestamp *timestamppb.Timestamp) *time.Time {
	if timestamp == nil {
		return nil
	}
	t := timestamp.AsTime()
	return &t
}

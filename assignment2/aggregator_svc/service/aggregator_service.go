package service

import (
	"aggregator_svc/entity"
	transaction_service "aggregator_svc/proto/transaction_service/v1"
	user_service "aggregator_svc/proto/user_service/v1"
	wallet_service "aggregator_svc/proto/wallet_service/v1"
	"context"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
	"strconv"
	"time"
)

type IAggregatorService interface {
	GetUser(ctx context.Context, id int) (entity.UserResponse, error)
	CreateUser(ctx context.Context, request entity.UserCreateRequest) (entity.UserResponse, error)
	TopupTransaction(ctx context.Context, request entity.TransactionRequest) (entity.TransactionResponse, error)
	TransferTransaction(ctx context.Context, request entity.TransactionRequest) (entity.TransactionResponse, error)
	GetTransactions(ctx context.Context, param entity.TransactionGetRequest) (entity.TransactionGetResponseWithPagination, error)
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

	return entity.TransactionResponse{Message: transactionResp.Message}, nil
}

func (svc *AggregatorService) TransferTransaction(ctx context.Context, request entity.TransactionRequest) (entity.TransactionResponse, error) {
	//Get Balance User sender
	walletResp, err := svc.walletService.GetWalletByUserID(ctx, &wallet_service.GetWalletByUserIDRequest{UserID: request.FromID})
	if err != nil {
		return entity.TransactionResponse{Message: "Wallet Not Fount"}, err
	}

	if walletResp.Wallet.Balance < float32(request.Amount) {
		return entity.TransactionResponse{Message: "The balance is not enough. Your balance: " + strconv.FormatFloat(float64(walletResp.Wallet.Balance), 'f', 2, 64)}, err
	}

	//create transaction
	transactionResp, err := svc.transactionService.CreateTransaction(ctx, &transaction_service.CreateTransactionRequest{
		FromID: request.FromID,
		ToID:   int32(request.ToID),
		Type:   "transfer",
		Amount: float32(request.Amount),
	})

	if err != nil {
		return entity.TransactionResponse{}, err
	}

	return entity.TransactionResponse{Message: transactionResp.Message}, nil
}

func (svc *AggregatorService) GetTransactions(ctx context.Context, request entity.TransactionGetRequest) (entity.TransactionGetResponseWithPagination, error) {
	transactionResp, err := svc.transactionService.GetTransactions(ctx, &transaction_service.GetTransactionRequest{
		Type:     request.Type,
		UserID:   int32(request.UserID),
		PageSize: int32(request.Size),
		Page:     int32(request.Page),
	})

	if err != nil {
		return entity.TransactionGetResponseWithPagination{}, err
	}
	var transactionsWithUser []entity.TransactionGetResponse
	for _, tx := range transactionResp.Transactions {
		name := ""
		userResp, err := svc.userService.GetUserByID(ctx, &user_service.GetUserByIDRequest{Id: tx.UserID})
		if err != nil {
			return entity.TransactionGetResponseWithPagination{}, err
		}
		name = userResp.User.Name

		var createdAt, updatedAt time.Time
		if tx.CreatedAt != nil {
			createdAt = tx.CreatedAt.AsTime()
		}
		if tx.UpdatedAt != nil {
			updatedAt = tx.UpdatedAt.AsTime()
		}

		transactionsWithUser = append(transactionsWithUser, entity.TransactionGetResponse{
			ID:        int(tx.Id),
			UserID:    int(tx.UserID),
			Name:      name,
			Amount:    float64(tx.Amount),
			Type:      tx.Type,
			Category:  tx.Category,
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
		})
	}

	return entity.TransactionGetResponseWithPagination{
		Data: transactionsWithUser,
		Pagination: entity.Pagination{
			TotalData: int(transactionResp.TotalCount),
			TotalPage: (int(transactionResp.TotalCount) + request.Size - 1) / request.Size,
			PageSize:  request.Size,
			Page:      request.Page,
		},
	}, nil
}

func convertTimestampToTime(timestamp *timestamppb.Timestamp) *time.Time {
	if timestamp == nil {
		return nil
	}
	t := timestamp.AsTime()
	return &t
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

func (svc *AggregatorService) CreateUser(ctx context.Context, request entity.UserCreateRequest) (entity.UserResponse, error) {
	//create user
	createUser, err := svc.userService.CreateUser(ctx, &user_service.CreateUserRequest{
		Name:  request.Name,
		Email: request.Email,
	})
	if err != nil {
		return entity.UserResponse{}, err
	}

	//create wallet for balance
	_, err = svc.walletService.CreateWallet(ctx, &wallet_service.CreateWalletRequest{
		UserID:  createUser.Id,
		Balance: 0,
	})
	if err != nil {
		return entity.UserResponse{}, err
	}

	return entity.UserResponse{
		ID: createUser.Id,
	}, nil
}

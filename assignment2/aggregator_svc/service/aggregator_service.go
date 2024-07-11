package service

import (
	"aggregator_svc/entity"
	user_service "aggregator_svc/proto/user_service/v1"
	wallet_service "aggregator_svc/proto/wallet_service/v1"
	"context"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
	"time"
)

type IAggregatorService interface {
	GetUser(ctx context.Context, id int) (entity.UserResponse, error)
}

type AggregatorService struct {
	userService   user_service.UserServiceClient
	walletService wallet_service.WalletServiceClient
}

func NewAggregatorService(userService user_service.UserServiceClient, walletService wallet_service.WalletServiceClient) *AggregatorService {
	return &AggregatorService{
		userService:   userService,
		walletService: walletService,
	}
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

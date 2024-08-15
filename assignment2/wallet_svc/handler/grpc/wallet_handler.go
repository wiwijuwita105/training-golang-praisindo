package grpc

import (
	"assignment5/cashflow-svc/internal/proto/wallet_service/v1"
	"context"
	"fmt"
	"log"
	"wallet_svc/entity"
	"wallet_svc/service"

	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// WalletHandler is used to implement UnimplementedWalletServiceServer
type WalletHandler struct {
	v1.UnimplementedWalletServiceServer
	walletService service.IWalletService
}

// membuat instance baru dari WalletHandler
func NewWalletHandler(walletService service.IWalletService) *WalletHandler {
	return &WalletHandler{
		walletService: walletService,
	}
}

func (u *WalletHandler) GetWallets(ctx context.Context, _ *emptypb.Empty) (*v1.GetWalletsResponse, error) {
	wallets, err := u.walletService.GetAllWallets(ctx)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var walletsProto []*v1.Wallet
	for _, wallet := range wallets {
		walletsProto = append(walletsProto, &v1.Wallet{
			Id:        int32(wallet.ID),
			UserID:    int32(wallet.UserID),
			Balance:   float32(wallet.Balance),
			CreatedAt: timestamppb.New(wallet.CreatedAt),
			UpdatedAt: timestamppb.New(wallet.UpdatedAt),
		})
	}

	return &v1.GetWalletsResponse{
		Wallets: walletsProto,
	}, nil
}

func (u *WalletHandler) GetWalletByUserID(ctx context.Context, req *v1.GetWalletByUserIDRequest) (*v1.GetWalletByUserIDResponse, error) {
	wallet, err := u.walletService.GetWalletByUserID(ctx, int(req.GetUserID()))
	if err != nil {
		log.Println(err)
		return nil, err
	}
	res := &v1.GetWalletByUserIDResponse{
		Wallet: &v1.Wallet{
			Id:        int32(wallet.ID),
			UserID:    int32(wallet.UserID),
			Balance:   float32(wallet.Balance),
			CreatedAt: timestamppb.New(wallet.CreatedAt),
			UpdatedAt: timestamppb.New(wallet.UpdatedAt),
		},
	}
	return res, nil
}

func (u *WalletHandler) CreateWallet(ctx context.Context, req *v1.CreateWalletRequest) (*v1.MutationResponse, error) {
	createdWallet, err := u.walletService.CreateWallet(ctx, &entity.Wallet{
		UserID:  int(req.GetUserID()),
		Balance: float64(req.GetBalance()),
	})
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &v1.MutationResponse{
		Message: fmt.Sprintf("Success created wallet with ID %d", createdWallet.ID),
	}, nil
}

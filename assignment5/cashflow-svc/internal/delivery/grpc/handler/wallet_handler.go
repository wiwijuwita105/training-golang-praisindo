package handler

import (
	"assignment5/cashflow-svc/internal/entity"
	"assignment5/cashflow-svc/internal/model"
	pbWallet "assignment5/cashflow-svc/internal/proto/wallet_service/v1"
	"assignment5/cashflow-svc/internal/service"
	"context"
	"errors"
	"fmt"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
	"time"
)

// WalletHandler is used to implement UnimplementedWalletServiceServer
type WalletHandler struct {
	pbWallet.UnimplementedWalletServiceServer
	walletService      service.IWalletService
	categoryService    service.ICategoryService
	transactionService service.ITransactionService
}

// membuat instance baru dari WalletHandler
func NewWalletHandler(
	walletService service.IWalletService,
	categoryService service.ICategoryService,
	transactionService service.ITransactionService,
) *WalletHandler {
	return &WalletHandler{
		walletService:      walletService,
		categoryService:    categoryService,
		transactionService: transactionService,
	}
}

func (u *WalletHandler) GetWallets(ctx context.Context, _ *emptypb.Empty) (*pbWallet.GetWalletsResponse, error) {
	wallets, err := u.walletService.GetAllWallets(ctx)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var walletsProto []*pbWallet.Wallet
	for _, wallet := range wallets {
		walletsProto = append(walletsProto, &pbWallet.Wallet{
			Id:        int32(wallet.ID),
			UserID:    int32(wallet.UserID),
			Balance:   float32(wallet.Balance),
			Name:      wallet.Name,
			CreatedAt: timestamppb.New(wallet.CreatedAt),
			UpdatedAt: timestamppb.New(wallet.UpdatedAt),
		})
	}

	return &pbWallet.GetWalletsResponse{
		Wallets: walletsProto,
	}, nil
}

func (u *WalletHandler) GetWalletByID(ctx context.Context, req *pbWallet.GetWalletByIDRequest) (*pbWallet.GetWalletByIDResponse, error) {
	wallet, err := u.walletService.GetWalletByID(ctx, int(req.GetId()))
	if err != nil {
		log.Println(err)
		return nil, err
	}
	res := &pbWallet.GetWalletByIDResponse{
		Wallet: &pbWallet.Wallet{
			Id:        int32(wallet.ID),
			UserID:    int32(wallet.UserID),
			Name:      wallet.Name,
			Balance:   float32(wallet.Balance),
			CreatedAt: timestamppb.New(wallet.CreatedAt),
			UpdatedAt: timestamppb.New(wallet.UpdatedAt),
		},
	}
	return res, nil
}

func (u *WalletHandler) CreateWallet(ctx context.Context, req *pbWallet.CreateWalletRequest) (*pbWallet.MutationResponse, error) {
	createdWallet, err := u.walletService.CreateWallet(ctx, &entity.Wallet{
		UserID:  int(req.GetUserID()),
		Name:    req.GetName(),
		Balance: float64(req.GetBalance()),
	})
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &pbWallet.MutationResponse{
		Message: fmt.Sprintf("Success created wallet with ID %d", createdWallet.ID),
	}, nil
}

func (u *WalletHandler) DeleteWallet(ctx context.Context, req *pbWallet.DeleteWalletRequest) (*pbWallet.MutationResponse, error) {
	if err := u.walletService.DeleteWallet(ctx, int(req.GetId())); err != nil {
		log.Println(err)
		return nil, err
	}

	return &pbWallet.MutationResponse{
		Message: fmt.Sprintf("Success deleted wallet with ID %d", req.GetId()),
	}, nil
}

// HANDLER CATEGORY
func (u *WalletHandler) GetCategories(ctx context.Context, _ *emptypb.Empty) (*pbWallet.GetCategoriesResponse, error) {
	categories, err := u.categoryService.GetAllCategories(ctx)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var categoriesProto []*pbWallet.Category
	for _, category := range categories {
		categoriesProto = append(categoriesProto, &pbWallet.Category{
			Id:        int32(category.ID),
			Name:      category.Name,
			Type:      category.Type,
			CreatedAt: timestamppb.New(category.CreatedAt),
			UpdatedAt: timestamppb.New(category.UpdatedAt),
		})
	}

	return &pbWallet.GetCategoriesResponse{
		Categories: categoriesProto,
	}, nil
}

func (u *WalletHandler) GetCategoryByID(ctx context.Context, req *pbWallet.GetCategoryByIDRequest) (*pbWallet.GetCategoryByIDResponse, error) {
	category, err := u.categoryService.GetCategoryByID(ctx, int(req.GetId()))
	if err != nil {
		log.Println(err)
		return nil, err
	}
	res := &pbWallet.GetCategoryByIDResponse{
		Category: &pbWallet.Category{
			Id:        int32(category.ID),
			Name:      category.Name,
			Type:      category.Type,
			CreatedAt: timestamppb.New(category.CreatedAt),
			UpdatedAt: timestamppb.New(category.UpdatedAt),
		},
	}
	return res, nil
}

func (u *WalletHandler) CreateCategory(ctx context.Context, req *pbWallet.CreateCategoryRequest) (*pbWallet.MutationCategoryResponse, error) {
	createdCategory, err := u.categoryService.CreateCategory(ctx, &entity.TransactionCategory{
		Name: req.GetName(),
		Type: req.GetType(),
	})
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &pbWallet.MutationCategoryResponse{
		Message: fmt.Sprintf("Success created wallet with ID %d", createdCategory.ID),
	}, nil
}

func (u *WalletHandler) DeleteCategory(ctx context.Context, req *pbWallet.DeleteCategoryRequest) (*pbWallet.MutationCategoryResponse, error) {
	if err := u.categoryService.DeleteCategory(ctx, int(req.GetId())); err != nil {
		log.Println(err)
		return nil, err
	}

	return &pbWallet.MutationCategoryResponse{
		Message: fmt.Sprintf("Success deleted category with ID %d", req.GetId()),
	}, nil
}

// TRANSACTION HANDLER
func (u *WalletHandler) TransferWallet(ctx context.Context, req *pbWallet.TransferRequest) (*pbWallet.MutationTransferResponse, error) {
	_, err := u.transactionService.TransferWallet(ctx, model.TransferWalletRequest{
		ToID:            req.ToID,
		FromID:          req.FromID,
		Nominal:         float64(req.Nominal),
		TransactionDate: time.Now(),
	})
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &pbWallet.MutationTransferResponse{
		Message: fmt.Sprintf("Success transfer from wallet ID %d to wallet ID %d", req.FromID, req.ToID),
	}, nil
}

func (u *WalletHandler) CreateTransaction(ctx context.Context, req *pbWallet.CreateTransactionRequest) (*pbWallet.MutationTransactionResponse, error) {
	transaction, err := u.transactionService.CreateTransaction(ctx, model.TransactionRequest{
		WalletID:   req.WalletID,
		CategoryID: req.CategoryID,
		Nominal:    float64(req.Nominal),
	})
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &pbWallet.MutationTransactionResponse{
		Message: fmt.Sprintf("Success transfer with ID: %d", transaction.ID),
	}, nil
}

func (u *WalletHandler) GetTransactions(ctx context.Context, req *pbWallet.GetTransactionRequest) (*pbWallet.GetTransactionResponse, error) {
	var requestFilter model.FilterTransactionRequest
	if req.WalletID == 0 {
		if req.UserID == 0 {
			return &pbWallet.GetTransactionResponse{}, errors.New("no user id or wallet id provided")
		}
	}
	requestFilter.UserID = req.UserID
	requestFilter.WalletID = req.WalletID
	if req.StartDate != "" {
		startDate, err := time.Parse("2006-01-02", req.StartDate)
		if err != nil {
			return &pbWallet.GetTransactionResponse{}, err
		}
		requestFilter.StartTime = &startDate
	} else {
		stratDate := time.Now().AddDate(0, 0, -30)
		requestFilter.StartTime = &stratDate
	}

	if req.EndDate != "" {
		endDate, err := time.Parse("2006-01-02", req.EndDate)
		if err != nil {
			return &pbWallet.GetTransactionResponse{}, err
		}
		requestFilter.EndTime = &endDate
	} else {
		endDate := time.Now()
		requestFilter.EndTime = &endDate
	}

	transactions, err := u.transactionService.GetTransactions(ctx, requestFilter)
	if err != nil {
		return nil, err
	}

	var transactionProto []*pbWallet.Transaction
	for _, record := range transactions {
		transactionProto = append(transactionProto, &pbWallet.Transaction{
			Id:              record.ID,
			TransactionDate: timestamppb.New(record.TransactionDate),
			Type:            record.Type,
			Nominal:         float32(record.Nominal),
			WalletID:        record.WalletID,
			WalletName:      record.WalletName,
			CategoryId:      record.CategoryID,
			CategoryName:    record.CategoryName,
		})
	}
	return &pbWallet.GetTransactionResponse{
		Transactions: transactionProto,
	}, nil
}

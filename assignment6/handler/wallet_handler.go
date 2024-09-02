package handler

import (
	"assignment6/entity"
	"assignment6/model"
	"assignment6/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

type IWalletHandler interface {
	//wallet
	GetWallets(ctx *gin.Context)
	GetWalletByID(ctx *gin.Context)
	CreateWallet(ctx *gin.Context)
	DeleteWallet(ctx *gin.Context)

	//category
	GetCategories(ctx *gin.Context)
	GetCategoryByID(ctx *gin.Context)
	CreateCategory(ctx *gin.Context)
	DeleteCategory(ctx *gin.Context)

	//transaction
	TransferWallet(ctx *gin.Context)
	CreateTransaction(ctx *gin.Context)
	GetTransactions(ctx *gin.Context)
	GetLastTransactions(ctx *gin.Context)
	GetCashflowReport(ctx *gin.Context)
	GetSummaryCategory(ctx *gin.Context)
}

// WalletHandler is used to implement UnimplementedWalletServiceServer
type WalletHandler struct {
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

func (u *WalletHandler) GetWallets(c *gin.Context) {
	wallets, err := u.walletService.GetAllWallets(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, gin.H{"wallets": wallets})
}

func (u *WalletHandler) GetWalletByID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	wallet, err := u.walletService.GetWalletByID(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"wallet": wallet})
	return
}

func (u *WalletHandler) CreateWallet(ctx *gin.Context) {
	var wallet entity.Wallet
	if err := ctx.ShouldBindJSON(&wallet); err != nil {
		errMsg := err.Error()
		errMsg = convertUserMandatoryFieldErrorString(errMsg)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": errMsg})
		return
	}
	createdWallet, err := u.walletService.CreateWallet(ctx, &wallet)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"wallet": createdWallet})
	return
}

func (u *WalletHandler) DeleteWallet(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	if err := u.walletService.DeleteWallet(ctx, id); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "wallet deleted"})
	return
}

// HANDLER CATEGORY
func (u *WalletHandler) GetCategories(ctx *gin.Context) {
	categories, err := u.categoryService.GetAllCategories(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"categories": categories})
	return
}

func (u *WalletHandler) GetCategoryByID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	category, err := u.categoryService.GetCategoryByID(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"category": category})
	return
}

func (u *WalletHandler) CreateCategory(ctx *gin.Context) {
	var transactionCategory entity.TransactionCategory
	if err := ctx.ShouldBindJSON(&transactionCategory); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createdCategory, err := u.categoryService.CreateCategory(ctx, &transactionCategory)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"category": createdCategory})
	return
}

func (u *WalletHandler) DeleteCategory(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := u.categoryService.DeleteCategory(ctx, id); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "category deleted"})
	return
}

// TRANSACTION HANDLER
func (u *WalletHandler) TransferWallet(ctx *gin.Context) {
	var transaction model.TransferWalletRequest
	if err := ctx.ShouldBindJSON(&transaction); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	transfer, err := u.transactionService.TransferWallet(ctx, transaction)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	ctx.JSON(http.StatusOK, gin.H{"transfer": transfer})
	return
}

func (u *WalletHandler) CreateTransaction(ctx *gin.Context) {
	var request model.TransactionRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	transaction, err := u.transactionService.CreateTransaction(ctx, request)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"transaction": transaction})
	return
}

func (u *WalletHandler) GetTransactions(ctx *gin.Context) {
	reqWalletID := ctx.Query("wallet_id")
	reqUserID := ctx.Query("user_id")
	reqStartDate := ctx.Query("start_date")
	reqEndDate := ctx.Query("end_date")

	var requestFilter model.FilterTransactionRequest
	if reqWalletID == "0" || reqWalletID == "" {
		if reqUserID == "0" || reqUserID == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "no user id or wallet id provided"})
			return
		} else {
			userID, err := strconv.Atoi(reqUserID)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			requestFilter.UserID = int32(userID)
			requestFilter.WalletID = 0
		}
	} else {
		walletID, err := strconv.Atoi(reqWalletID)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		requestFilter.WalletID = int32(walletID)
		requestFilter.UserID = 0
	}

	if reqStartDate != "" {
		startDate, err := time.Parse("2006-01-02", reqStartDate)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		requestFilter.StartTime = &startDate
	} else {
		startDate := time.Now().AddDate(0, 0, -30)
		requestFilter.StartTime = &startDate
	}

	if reqEndDate != "" {
		endDate, err := time.Parse("2006-01-02", reqEndDate)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		requestFilter.EndTime = &endDate
	} else {
		endDate := time.Now()
		requestFilter.EndTime = &endDate
	}

	transactions, err := u.transactionService.GetTransactions(ctx, requestFilter)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"transactions": transactions})
	return
}

func (u *WalletHandler) GetLastTransactions(ctx *gin.Context) {
	var requestFilter model.LastTransactionRequest
	reqWalletID := ctx.Query("wallet_id")
	if reqWalletID == "" || reqWalletID == "0" {
		reqUserID := ctx.Query("user_id")
		if reqUserID == "0" || reqUserID == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "no user id or user id provided"})
			return
		} else {
			userID, err := strconv.Atoi(reqUserID)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			requestFilter.UserID = userID
			requestFilter.WalletID = 0
		}
	} else {
		walletID, err := strconv.Atoi(ctx.Query("wallet_id"))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Wallet ID"})
			return
		}
		requestFilter.WalletID = walletID
		requestFilter.UserID = 0
	}

	transactions, err := u.transactionService.GetLastTransactions(ctx, requestFilter)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"transactions": transactions})
	return
}

func (u *WalletHandler) GetCashflowReport(ctx *gin.Context) {
	var requestFilter model.CashFlowReportRequest
	reqWalletID := ctx.Query("wallet_id")
	if reqWalletID == "" || reqWalletID == "0" {
		reqUserID := ctx.Query("user_id")
		if reqUserID == "0" || reqUserID == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "no user id or wallet id provided"})
			return
		}
		userID, err := strconv.Atoi(reqUserID)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		requestFilter.UserID = int32(userID)
		requestFilter.WalletID = 0
	}

	walletID, err := strconv.Atoi(ctx.Query("wallet_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	requestFilter.WalletID = int32(walletID)
	requestFilter.UserID = 0

	reqStartDate := ctx.Query("start_date")
	reqEndDate := ctx.Query("end_date")

	if reqStartDate != "" {
		startDate, err := time.Parse("2006-01-02", reqStartDate)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		requestFilter.StartTime = &startDate
	} else {
		startDate := time.Now().AddDate(0, 0, -30)
		requestFilter.StartTime = &startDate
	}

	if reqEndDate != "" {
		endDate, err := time.Parse("2006-01-02", reqEndDate)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		requestFilter.EndTime = &endDate
	} else {
		endDate := time.Now()
		requestFilter.EndTime = &endDate
	}

	cashflowReport, err := u.transactionService.GetCashflowReport(ctx, requestFilter)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"cashflowReport": cashflowReport})

	return
}

func (u *WalletHandler) GetSummaryCategory(ctx *gin.Context) {
	var requestFilter model.SummaryCategoryRequest
	reqWalletID := ctx.Query("wallet_id")
	if reqWalletID == "" || reqWalletID == "0" {
		reqUserID := ctx.Query("user_id")
		if reqUserID == "0" || reqUserID == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "no user id or user id provided"})
			return
		}
		userID, err := strconv.Atoi(reqUserID)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		requestFilter.UserID = int32(userID)
		requestFilter.WalletID = 0
	}

	walletID, err := strconv.Atoi(ctx.Query("wallet_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	requestFilter.WalletID = int32(walletID)
	requestFilter.UserID = 0

	reqStartDate := ctx.Query("start_date")
	reqEndDate := ctx.Query("end_date")
	if reqStartDate != "" {
		startDate, err := time.Parse("2006-01-02", reqStartDate)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		requestFilter.StartTime = &startDate
	} else {
		stratDate := time.Now().AddDate(0, 0, -30)
		requestFilter.StartTime = &stratDate
	}

	if reqEndDate != "" {
		endDate, err := time.Parse("2006-01-02", reqEndDate)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		requestFilter.EndTime = &endDate
	} else {
		endDate := time.Now()
		requestFilter.EndTime = &endDate
	}

	summaryCategories, err := u.transactionService.GetSummaryCategory(ctx, requestFilter)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"summaryCategories": summaryCategories})
	return
}

package router

import (
	"assignment6/handler"
	"github.com/gin-gonic/gin"
)

func SetupRouter(r *gin.Engine, userHandler handler.IUserHandler, walletHandler handler.IWalletHandler) {
	///endpoint user
	usersEndpoint := r.Group("/users")
	usersEndpoint.GET("/:id", userHandler.GetUserByID)
	usersEndpoint.GET("/", userHandler.GetUsers)
	usersEndpoint.POST("", userHandler.CreateUser)

	//endpoint wallet
	walletsEndpoint := r.Group("/wallets")
	walletsEndpoint.GET("/", walletHandler.GetWallets)
	walletsEndpoint.GET("/:id", walletHandler.GetWalletByID)
	walletsEndpoint.POST("/", walletHandler.CreateWallet)
	walletsEndpoint.DELETE("/:id", walletHandler.DeleteWallet)

	//endpointCategory
	categoryEndpoint := r.Group("/category")
	categoryEndpoint.GET("/", walletHandler.GetCategories)
	categoryEndpoint.GET("/:id", walletHandler.GetCategoryByID)
	categoryEndpoint.POST("", walletHandler.CreateCategory)
	categoryEndpoint.DELETE("/:id", walletHandler.DeleteCategory)

	//transaction
	transactionsEndpoint := r.Group("/transactions")
	transactionsEndpoint.POST("/transfer", walletHandler.TransferWallet)
	transactionsEndpoint.POST("/insert", walletHandler.CreateTransaction)
	transactionsEndpoint.GET("/", walletHandler.GetTransactions)
	transactionsEndpoint.GET("/last-transaction", walletHandler.GetLastTransactions)
	transactionsEndpoint.GET("/cashflow-report", walletHandler.GetCashflowReport)
	transactionsEndpoint.GET("/summary-by-categories", walletHandler.GetSummaryCategory)
}

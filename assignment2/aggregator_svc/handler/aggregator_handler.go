package handler

import (
	"aggregator_svc/entity"
	"aggregator_svc/service"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

type IAggregatorHandler interface {
	GetUser(c *gin.Context)
	TopupTransaction(c *gin.Context)
	TransferTransaction(c *gin.Context)
	GetTransactions(c *gin.Context)
}

type AggregatorHandler struct {
	aggregatorService service.AggregatorService
}

func NewAggregatorHandler(aggregatorService service.AggregatorService) *AggregatorHandler {
	return &AggregatorHandler{
		aggregatorService: aggregatorService,
	}
}

func (h *AggregatorHandler) GetUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("userId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	user, err := h.aggregatorService.GetUser(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	c.JSON(http.StatusOK, user)
}

func (h *AggregatorHandler) TopupTransaction(c *gin.Context) {
	var request entity.TransactionRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		errMsg := err.Error()
		c.JSON(http.StatusBadRequest, gin.H{"error": errMsg})
		return
	}
	transaction, err := h.aggregatorService.TopupTransaction(c.Request.Context(), request)
	log.Println(transaction)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	c.JSON(http.StatusOK, transaction)
}

func (h *AggregatorHandler) TransferTransaction(c *gin.Context) {
	var request entity.TransactionRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		errMsg := err.Error()
		c.JSON(http.StatusBadRequest, gin.H{"error": errMsg})
		return
	}
	transaction, err := h.aggregatorService.TransferTransaction(c.Request.Context(), request)
	log.Println(transaction)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	c.JSON(http.StatusOK, transaction)
}

func (h *AggregatorHandler) GetTransactions(c *gin.Context) {
	userID, err := strconv.Atoi(c.Query("userID"))
	if err != nil {
		userID = 0 // set default value or handle error as per requirement
	}

	page, err := strconv.Atoi(c.Query("page"))
	if err != nil {
		page = 1 // set default value or handle error as per requirement
	}

	size, err := strconv.Atoi(c.Query("size"))
	if err != nil {
		size = 10 // set default value or handle error as per requirement
	}

	paramRequest := entity.TransactionGetRequest{
		Type:   c.Query("type"),
		UserID: userID,
		Page:   page,
		Size:   size,
	}
	transaction, err := h.aggregatorService.GetTransactions(c.Request.Context(), paramRequest)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	c.JSON(http.StatusOK, transaction)
}

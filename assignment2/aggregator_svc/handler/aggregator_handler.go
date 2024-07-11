package handler

import (
	"aggregator_svc/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type IAggregatorHandler interface {
	GetUser(c *gin.Context)
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

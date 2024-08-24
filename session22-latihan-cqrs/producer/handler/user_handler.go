package handler

import (
	"encoding/json"
	"github.com/IBM/sarama"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"session22-latihan-cqrs/config"
	"session22-latihan-cqrs/entity"
)

type IUserHandler interface {
	CreateUser(c *gin.Context)
	GetAllUsers(c *gin.Context)
}

type UserHandler struct {
	dbRead *gorm.DB
}

func NewUserHandler(dbRead *gorm.DB) *UserHandler {
	return &UserHandler{
		dbRead: dbRead,
	}
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	//call config producer kafka
	producer, err := config.SetupProducer()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer producer.Close()

	var user entity.UserInput
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	jsonUser, err := json.Marshal(user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	msg := &sarama.ProducerMessage{
		Topic: config.UserTopic,
		Value: sarama.StringEncoder(jsonUser),
	}

	_, _, err = producer.SendMessage(msg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusCreated, user)
}

func (h *UserHandler) GetAllUsers(c *gin.Context) {
	var users []entity.User
	if err := h.dbRead.Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, users)
}

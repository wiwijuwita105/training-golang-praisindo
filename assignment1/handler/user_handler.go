package handler

import (
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"assignment1/config"
	"assignment1/entity"
	"assignment1/service"

	"github.com/gin-gonic/gin"
)

type UserWithRiskResponse struct {
	ID             uint   `json:"id"`
	Name           string `json:"name"`
	Email          string `json:"email"`
	RiskScore      int    `json:"risk_score"`
	RiskCategory   string `json:"risk_category"`
	RiskDefinition string `json:"risk_definition"`
	CreatedAt      string `json:"created_at"`
	UpdatedAt      string `json:"updated_at"`
}

// IUserHandler mendefinisikan interface untuk handler user
type IUserHandler interface {
	CreateUser(c *gin.Context)
	GetUser(c *gin.Context)
	UpdateUser(c *gin.Context)
	DeleteUser(c *gin.Context)
	GetAllUsers(c *gin.Context)
	GetAllUsersWithRIsk(c *gin.Context)
}

type UserHandler struct {
	userService service.IUserService
}

// NewUserHandler membuat instance baru dari UserHandler
func NewUserHandler(userService service.IUserService) IUserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// UserResponse defines the structure of user response
type UserResponse struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// CreateUser menghandle permintaan untuk membuat user baru
func (h *UserHandler) CreateUser(c *gin.Context) {
	var user entity.User
	if err := c.ShouldBindJSON(&user); err != nil {
		errMsg := err.Error()
		errMsg = convertUserMandatoryFieldErrorString(errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"error": errMsg})
		return
	}

	createdUser, err := h.userService.CreateUser(c.Request.Context(), &user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Format the response message
	responseMessage := gin.H{
		"message": "User ID " + strconv.Itoa(int(createdUser.ID)) + " created successfully",
	}

	c.JSON(http.StatusCreated, responseMessage)
}

// GetUser menghandle permintaan untuk mendapatkan user berdasarkan ID
func (h *UserHandler) GetUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	user, err := h.userService.GetUserByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// UpdateUser menghandle permintaan untuk mengupdate informasi user
func (h *UserHandler) UpdateUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var user entity.User
	if err := c.ShouldBindJSON(&user); err != nil {
		errMsg := err.Error()
		errMsg = convertUserMandatoryFieldErrorString(errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"error": errMsg})
		return
	}

	updatedUser, err := h.userService.UpdateUser(c.Request.Context(), id, user)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// Format the response message
	responseMessage := gin.H{
		"message": "User ID " + strconv.Itoa(int(updatedUser.ID)) + " updated successfully",
	}

	c.JSON(http.StatusOK, responseMessage)
}

// DeleteUser menghandle permintaan untuk menghapus user
func (h *UserHandler) DeleteUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := h.userService.DeleteUser(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

// GetAllUsers menghandle permintaan untuk mendapatkan semua user
func (h *UserHandler) GetAllUsers(c *gin.Context) {
	// Default values for page and pageSize
	const defaultPage = 1
	const defaultPageSize = 10

	page, err := strconv.Atoi(c.DefaultQuery("page", strconv.Itoa(defaultPage)))
	if err != nil || page < 1 {
		page = defaultPage
	}

	pageSize, err := strconv.Atoi(c.DefaultQuery("limit", strconv.Itoa(defaultPageSize)))
	if err != nil || pageSize < 1 {
		pageSize = defaultPageSize
	}

	// ctx := context.Background()
	users, total, err := h.userService.GetAllUsers(c.Request.Context(), page, pageSize)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// users, total, err := h.userService.GetAllUsers(c.Request.Context(), page, pageSize)
	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	// 	return
	// }

	// Transform users to UserResponse
	var userResponses []UserResponse
	for _, user := range users {
		log.Println(user)
		userResponse := UserResponse{
			ID:        user.ID,
			Name:      user.Name,
			Email:     user.Email,
			CreatedAt: user.CreatedAt.Format(time.RFC3339),
			UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
		}
		userResponses = append(userResponses, userResponse)
	}

	// Calculate total pages
	totalPages := (total + int64(pageSize) - 1) / int64(pageSize)

	c.JSON(http.StatusOK, gin.H{"users": userResponses, "total_pages": totalPages, "current_page": page, "total_data": total, "total_perpage": pageSize})
}

func (h *UserHandler) GetAllUsersWithRIsk(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil || pageSize < 1 {
		pageSize = 10
	}

	users, total, err := h.userService.GetAllUsersWithRisk(c.Request.Context(), page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get users"})
		return
	}

	var responseUser []UserWithRiskResponse
	for _, sub := range users {
		var riskScore int
		var riskCategory string
		var riskDefinition string

		if len(sub.Submissions) > 0 {
			riskScore = sub.Submissions[0].RiskScore
			riskCategory = sub.Submissions[0].RiskCategory
			riskDefinition = getRiskDefinition(config.ProfileRiskCategory(sub.Submissions[0].RiskCategory))
		}

		responseUser = append(responseUser, UserWithRiskResponse{
			ID:             sub.ID,
			RiskScore:      riskScore,
			RiskCategory:   riskCategory,
			RiskDefinition: riskDefinition,
			CreatedAt:      sub.CreatedAt.Format("2006-01-02T15:04:05Z"),
			UpdatedAt:      sub.UpdatedAt.Format("2006-01-02T15:04:05Z"),
		})
	}

	response := gin.H{
		"users":        responseUser,
		"total_pages":  (total + int64(pageSize) - 1) / int64(pageSize),
		"current_page": page,
	}

	c.JSON(http.StatusOK, response)
}

func convertUserMandatoryFieldErrorString(oldErrorMsg string) string {
	switch {
	case strings.Contains(oldErrorMsg, "'Name' failed on the 'required' tag"):
		return "name is mandatory"
	case strings.Contains(oldErrorMsg, "'Email' failed on the 'required' tag"):
		return "email is mandatory"
	}
	return oldErrorMsg
}

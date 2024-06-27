package handler

import (
	"assignment1/config"
	"assignment1/entity"
	"assignment1/service"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Define your response struct
type GetSubmissionsResponse struct {
	Total       int64            `json:"total"`
	Page        int              `json:"page"`
	Limit       int              `json:"limit"`
	Submissions []SubmissionData `json:"submissions"`
}

type SubmissionData struct {
	ID             uint            `json:"id"`
	UserID         uint            `json:"user_id"`
	RiskScore      int             `json:"risk_score"`
	RiskCategory   string          `json:"risk_category"`
	RiskDefinition string          `json:"risk_definition"`
	Answers        []entity.Answer `json:"answers"`
	CreatedAt      string          `json:"created_at"`
	UpdatedAt      string          `json:"updated_at"`
}

// mendefinisikan interface untuk handler submission
type ISubmissionHandler interface {
	CreateSubmission(c *gin.Context)
	GetSubmission(c *gin.Context)
	DeleteSubmission(c *gin.Context)
	GetAllSubmissions(c *gin.Context)
	GetAllSubmissionsWithUser(c *gin.Context)
}

type SubmissionHandler struct {
	submissionService service.ISubmissionService
}

// NewUserHandler membuat instance baru dari UserHandler
func NewSubmissionHandler(submissionService service.ISubmissionService) ISubmissionHandler {
	return &SubmissionHandler{
		submissionService: submissionService,
	}
}

type CreateSubmissionRequest struct {
	UserID  uint            `json:"user_id"`
	Answers []entity.Answer `json:"answers"`
}

// menghandle permintaan untuk membuat submission baru
func (h *SubmissionHandler) CreateSubmission(c *gin.Context) {
	//defind request
	var submissionRequest CreateSubmissionRequest

	// Bind JSON dari body request ke struct submissionRequest
	if err := c.ShouldBindJSON(&submissionRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Convert answers to JSON type
	answersJSON, err := json.Marshal(submissionRequest.Answers)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process answers"})
		return
	}

	totalWeight := calculateRiskScore(submissionRequest.Answers)
	// Simpan submission ke dalam database menggunakan GORM
	submission := entity.Submission{
		UserID:       submissionRequest.UserID,
		Answer:       answersJSON,
		RiskScore:    totalWeight,
		RiskCategory: calculateRiskCategory(totalWeight),
	}

	createdSubmission, err := h.submissionService.CreateSubmission(c.Request.Context(), &submission)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Format the response message
	responseMessage := gin.H{
		"message": "Submission ID " + strconv.Itoa(int(createdSubmission.ID)) + " created successfully",
	}

	c.JSON(http.StatusCreated, responseMessage)
}

// menghandle permintaan untuk mendapatkan submission berdasarkan ID
func (h *SubmissionHandler) GetSubmission(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	submission, err := h.submissionService.GetSubmissionByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Submission not found"})
		return
	}

	// Unmarshal the answers from the submission
	var answers []entity.Answer
	err = json.Unmarshal(submission.Answer, &answers)
	if err != nil {
		log.Printf("Error unmarshalling answers: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process answers"})
		return
	}

	// Create the response struct
	response := SubmissionData{
		ID:             submission.ID,
		UserID:         submission.UserID,
		RiskScore:      submission.RiskScore,
		RiskCategory:   submission.RiskCategory,
		RiskDefinition: getRiskDefinition(config.ProfileRiskCategory(submission.RiskCategory)),
		Answers:        answers,
		CreatedAt:      submission.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:      submission.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}

	c.JSON(http.StatusOK, response)
}

// menghandle permintaan untuk menghapus user
func (h *SubmissionHandler) DeleteSubmission(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := h.submissionService.DeleteSubmission(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Submission deleted successfully"})
}

// GetAllUsers menghandle permintaan untuk mendapatkan semua user
func (h *SubmissionHandler) GetAllSubmissions(c *gin.Context) {
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

	// Get user_id from query parameters
	var userID *uint
	if userIDStr := c.Query("user_id"); userIDStr != "" {
		uid, err := strconv.Atoi(userIDStr)
		if err == nil {
			uidUint := uint(uid)
			userID = &uidUint
		}
	}

	// ctx := context.Background()

	submissions, total, err := h.submissionService.GetAllSubmissions(c.Request.Context(), page, pageSize, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var submissionData []SubmissionData
	for _, sub := range submissions {
		var answers []entity.Answer
		err := json.Unmarshal(sub.Answer, &answers)
		if err != nil {
			log.Printf("Error unmarshalling answers: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process answers"})
			return
		}
		submissionData = append(submissionData, SubmissionData{
			ID:             sub.ID,
			UserID:         sub.UserID,
			RiskScore:      sub.RiskScore,
			RiskCategory:   sub.RiskCategory,
			RiskDefinition: getRiskDefinition(config.ProfileRiskCategory(sub.RiskCategory)),
			Answers:        answers,
			CreatedAt:      sub.CreatedAt.Format("2006-01-02T15:04:05Z"),
			UpdatedAt:      sub.UpdatedAt.Format("2006-01-02T15:04:05Z"),
		})
	}

	response := GetSubmissionsResponse{
		Total:       total,
		Page:        page,
		Limit:       pageSize,
		Submissions: submissionData,
	}

	c.JSON(http.StatusOK, response)
}

func (h *SubmissionHandler) GetAllSubmissionsWithUser(c *gin.Context) {
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

	// Get user_id from query parameters
	var userID *uint
	if userIDStr := c.Query("user_id"); userIDStr != "" {
		uid, err := strconv.Atoi(userIDStr)
		if err == nil {
			uidUint := uint(uid)
			userID = &uidUint
		}
	}

	// ctx := context.Background()

	submissions, total, err := h.submissionService.GetAllSubmissionsWithUser(c.Request.Context(), page, pageSize, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var submissionData []SubmissionData
	for _, sub := range submissions {
		var answers []entity.Answer
		err := json.Unmarshal(sub.Answer, &answers)
		if err != nil {
			log.Printf("Error unmarshalling answers: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process answers"})
			return
		}
		submissionData = append(submissionData, SubmissionData{
			ID:             sub.ID,
			UserID:         sub.UserID,
			RiskScore:      sub.RiskScore,
			RiskCategory:   sub.RiskCategory,
			RiskDefinition: getRiskDefinition(config.ProfileRiskCategory(sub.RiskCategory)),
			Answers:        answers,
			CreatedAt:      sub.CreatedAt.Format("2006-01-02T15:04:05Z"),
			UpdatedAt:      sub.UpdatedAt.Format("2006-01-02T15:04:05Z"),
		})
	}

	response := GetSubmissionsResponse{
		Total:       total,
		Page:        page,
		Limit:       pageSize,
		Submissions: submissionData,
	}

	c.JSON(http.StatusOK, response)
}

// Fungsi untuk mendapatkan bobot jawaban berdasarkan question_id dan answer
func calculateRiskScore(answers []entity.Answer) int {
	totalScore := 0
	for _, answer := range answers {
		for _, question := range config.Questions { //ID QUESION
			if question.ID == answer.QuestionID {
				for _, option := range question.Options {
					if option.Answer == answer.Answer {
						totalScore += option.Weight
						break
					}
				}
				break
			}
		}
	}
	return totalScore
}

func calculateRiskCategory(totalWeight int) string {
	for _, risk := range config.RiskMapping {
		if totalWeight >= risk.MinScore && totalWeight <= risk.MaxScore {
			return string(risk.Category)
		}
	}
	return ""
}

func getRiskDefinition(category config.ProfileRiskCategory) string {
	for _, profileRisk := range config.RiskMapping {
		if profileRisk.Category == category {
			return profileRisk.Definition
		}
	}
	return ""
}

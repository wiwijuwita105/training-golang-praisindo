package handler

import (
	"github.com/gin-gonic/gin"
	"latihan-consume-api-to-report/entity"
	"latihan-consume-api-to-report/service"
	"net/http"
	"strings"
)

type IWeatherHandler interface {
	GenerateToXLSX(c *gin.Context)
	GenerateToPDF(c *gin.Context)
}

type WeatherHandler struct {
	weatherService service.IWeatherService
}

func NewWeatherHandler(weatherService service.IWeatherService) IWeatherHandler {
	return &WeatherHandler{
		weatherService: weatherService,
	}
}

func (h *WeatherHandler) GenerateToXLSX(c *gin.Context) {
	var request entity.LocationRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		errMsg := err.Error()
		errMsg = convertUserMandatoryFieldErrorString(errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"error": errMsg})
		return
	}

	weather, err := h.weatherService.ReportWeatherToXLSX(c.Request.Context(), request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, weather)
}

func (h *WeatherHandler) GenerateToPDF(c *gin.Context) {
	var request entity.LocationRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		errMsg := err.Error()
		errMsg = convertUserMandatoryFieldErrorString(errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"error": errMsg})
		return
	}

	weather, err := h.weatherService.ReportWeatherToPDF(c.Request.Context(), request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, weather)
}

func convertUserMandatoryFieldErrorString(oldErrorMsg string) string {
	switch {
	case strings.Contains(oldErrorMsg, "'Latitude' failed on the 'required' tag"):
		return "latitude is mandatory"
	case strings.Contains(oldErrorMsg, "'Longitude' failed on the 'required' tag"):
		return "longitude is mandatory"
	}
	return oldErrorMsg
}

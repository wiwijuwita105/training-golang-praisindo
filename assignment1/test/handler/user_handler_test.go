package handler_test

import (
	"assignment1/entity"
	"assignment1/handler"
	mock_service "assignment1/test/mock/service"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestUserHandler_CreateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock_service.NewMockIUserService(ctrl)
	userHandler := handler.NewUserHandler(mockService)

	gin.SetMode(gin.TestMode)

	t.Run("ValidRequest", func(t *testing.T) {
		mockService.EXPECT().CreateUser(gomock.Any(), &entity.User{
			Name:  "John Doe",
			Email: "john@example.com",
		}).Return(entity.User{
			Name:  "John Doe",
			Email: "john@example.com",
		}, nil)

		req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(`{"name":"John Doe","email":"john@example.com","password":"password"}`))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()
		router := gin.Default()
		router.POST("/users", userHandler.CreateUser)

		router.ServeHTTP(resp, req)

		require.Equal(t, http.StatusCreated, resp.Code)
		require.JSONEq(t, `{"message":"User ID 0 created successfully"}`, resp.Body.String())
	})

	t.Run("InvalidPayload_MissingName", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(`{"email":"john@example.com","password":"password"}`))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()
		router := gin.Default()
		router.POST("/users", userHandler.CreateUser)

		router.ServeHTTP(resp, req)

		require.Equal(t, http.StatusBadRequest, resp.Code)
		require.JSONEq(t, `{"error":"name is mandatory"}`, resp.Body.String())
	})

	t.Run("InvalidPayload_MissingEmail", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(`{"name":"john","password":"password"}`))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()
		router := gin.Default()
		router.POST("/users", userHandler.CreateUser)

		router.ServeHTTP(resp, req)

		require.Equal(t, http.StatusBadRequest, resp.Code)
		require.JSONEq(t, `{"error":"email is mandatory"}`, resp.Body.String())
	})

	t.Run("ServiceError", func(t *testing.T) {
		mockService.EXPECT().CreateUser(gomock.Any(), &entity.User{
			Name:  "John Doe",
			Email: "john@example.com",
		}).Return(entity.User{}, errors.New("some service error"))

		req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(`{"name":"John Doe","email":"john@example.com","password":"password"}`))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()
		router := gin.Default()
		router.POST("/users", userHandler.CreateUser)

		router.ServeHTTP(resp, req)

		require.Equal(t, http.StatusBadRequest, resp.Code)
		require.JSONEq(t, `{"error":"some service error"}`, resp.Body.String())
	})
}

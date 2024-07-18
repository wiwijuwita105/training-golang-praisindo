// package router mengatur rute untuk aplikasi
package router

import (
	gin2 "session11-crud-user-grpc-gateway-cache/handler/gin"
	"session11-crud-user-grpc-gateway-cache/middleware"

	"github.com/gin-gonic/gin"
)

// SetupRouter menginisialisasi dan mengatur rute untuk aplikasi
func SetupRouter(r *gin.Engine, userHandler gin2.IUserHandler) {
	// Mengatur endpoint publik untuk pengguna
	usersPublicEndpoint := r.Group("/users")
	// Rute untuk mendapatkan pengguna berdasarkan ID
	usersPublicEndpoint.GET("/:id", userHandler.GetUser)
	// Rute untuk mendapatkan semua pengguna
	usersPublicEndpoint.GET("", userHandler.GetAllUsers)
	usersPublicEndpoint.GET("/", userHandler.GetAllUsers)

	// Mengatur endpoint privat untuk pengguna dengan middleware autentikasi
	usersPrivateEndpoint := r.Group("/users")
	// Menambahkan middleware autentikasi untuk endpoint privat
	usersPrivateEndpoint.Use(middleware.AuthMiddleware())
	// Rute untuk membuat pengguna baru
	usersPrivateEndpoint.POST("", userHandler.CreateUser)
	usersPrivateEndpoint.POST("/", userHandler.CreateUser)
	// Rute untuk memperbarui pengguna berdasarkan ID
	usersPrivateEndpoint.PUT("/:id", userHandler.UpdateUser)
	// Rute untuk menghapus pengguna berdasarkan ID
	usersPrivateEndpoint.DELETE("/:id", userHandler.DeleteUser)
}

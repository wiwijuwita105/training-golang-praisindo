package main

import (
	"context"
	"log"
	"session6-db-pgx/handler"
	"session6-db-pgx/repository/postgrespgx"
	"session6-db-pgx/router"
	"session6-db-pgx/service"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	pgxPool, err := connectDB("postgresql://postgres:postgres@localhost:5432/db_training_golang")
	if err != nil {
		log.Fatalln(err)
	}

	// setup service

	// slice db is disabled. uncomment to enabled
	//var mockUserDBInSlice []entity.User
	//_ = slice.NewUserRepository(mockUserDBInSlice)

	// pgx db is enabled. comment to disabled
	userRepo := postgrespgx.NewUserRepository(pgxPool)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	// Routes
	router.SetupRouter(r, userHandler)

	// Run the server
	log.Println("Running server on port 8080")
	r.Run(":8080")
}

func connectDB(dbURL string) (postgrespgx.PgxPoolIface, error) {
	return pgxpool.New(context.Background(), dbURL)
}

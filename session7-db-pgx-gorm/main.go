package main

import (
	"context"
	"log"
	"session7-db-pgx-gorm/entity"
	"session7-db-pgx-gorm/handler"
	"session7-db-pgx-gorm/repository/postgres_gorm_raw"
	"session7-db-pgx-gorm/repository/postgrespgx"
	"session7-db-pgx-gorm/router"
	"session7-db-pgx-gorm/service"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	//koneksion DB
	// pgxPool, err := connectDB("postgresql://postgres:postgres@localhost:5432/db_training_golang")
	// if err != nil {
	// 	log.Fatalln(err)
	// }

	//Koneksion db using gorm
	dsn := "postgresql://postgres:postgres@localhost:5432/db_training_golang"
	gormDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{SkipDefaultTransaction: true})
	if err != nil {
		log.Println(gormDB)
		log.Fatalln(err)

	}

	log.Println("starting migration")
	if err := gormDB.AutoMigrate(entity.User{}, entity.Customer{}); err != nil {
		log.Println(err)
	}
	log.Println("finish migration")
	// setup service

	// slice db is disabled. uncomment to enabled
	// var mockUserDBInSlice []entity.User
	// _ = slice.NewUserRepository(mockUserDBInSlice)

	// uncomment to use postgres pgx
	// userRepo := postgres_pgx.NewUserRepository(pgxPool)

	// uncomment to use postgres gorm
	// userRepo := postgres_gorm.NewUserRepository(gormDB)

	// service and handler declaration
	// userService := service.NewUserService(userRepo)
	// userHandler := handler.NewUserHandler(userService)

	// uncomment to use postgres gorm raw
	userRepo := postgres_gorm_raw.NewUserRepository(gormDB)

	// service and handler declaration
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

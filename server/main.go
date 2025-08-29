package main

import (
	"context"
	"log"

	"net/http"

	"github.com/dorianneto/bugfy/db"
	userHandler "github.com/dorianneto/bugfy/internal/api/handler/user"
	repository "github.com/dorianneto/bugfy/internal/repository/user"
	service "github.com/dorianneto/bugfy/internal/service/user"
	"github.com/dorianneto/bugfy/router"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}

	dbConn, err := db.NewDatabase()
	if err != nil {
		log.Fatalf("Could not initialize DB connection: %s", err)
	}
	defer dbConn.Disconnect(context.TODO())

	if err := dbConn.Ping(context.TODO(), readpref.Primary()); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}
	log.Println("Pinged your deployment. You are successfully connected!")

	userRepo := repository.NewUserRepository(dbConn)

	userService := service.NewUserService(userRepo)

	userHandler := userHandler.NewUserHandler(userService)

	router := router.SetupRouter(userHandler)
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

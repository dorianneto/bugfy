package main

import (
	"context"
	"log"

	"net/http"

	"github.com/dorianneto/bugfy/db"
	handler "github.com/dorianneto/bugfy/internal/api/handler"
	repo "github.com/dorianneto/bugfy/internal/repository"
	service "github.com/dorianneto/bugfy/internal/service"
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

	userRepo := repo.NewUserRepository(dbConn)
	projectRepo := repo.NewProjectRepository(dbConn)
	errorRepo := repo.NewErrorRepository(dbConn)

	userService := service.NewUserService(userRepo)
	projectService := service.NewProjectService(projectRepo)
	errorService := service.NewErrorService(errorRepo)

	userHandler := handler.NewUserHandler(userService)
	projectHandler := handler.NewProjectHandler(projectService)
	errorHandler := handler.NewErrorHandler(errorService)

	router := router.SetupRouter(userHandler, projectHandler, errorHandler)
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

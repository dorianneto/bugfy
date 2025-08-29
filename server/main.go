package main

import (
	"context"
	"log"

	"net/http"

	"github.com/dorianneto/bugfy/db"
	projectH "github.com/dorianneto/bugfy/internal/api/handler/project"
	userH "github.com/dorianneto/bugfy/internal/api/handler/user"
	projectRepo "github.com/dorianneto/bugfy/internal/repository/project"
	userRepo "github.com/dorianneto/bugfy/internal/repository/user"
	projectServ "github.com/dorianneto/bugfy/internal/service/project"
	userServ "github.com/dorianneto/bugfy/internal/service/user"
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

	userRepo := userRepo.NewUserRepository(dbConn)
	projectRepo := projectRepo.NewProjectRepository(dbConn)

	userService := userServ.NewUserService(userRepo)
	projectService := projectServ.NewProjectService(projectRepo)

	userHandler := userH.NewUserHandler(userService)
	projectHandler := projectH.NewProjectHandler(projectService)

	router := router.SetupRouter(userHandler, projectHandler)
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

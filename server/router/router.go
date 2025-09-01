package router

import (
	"net/http"

	errorH "github.com/dorianneto/bugfy/internal/api/handler/error"
	projectH "github.com/dorianneto/bugfy/internal/api/handler/project"
	userH "github.com/dorianneto/bugfy/internal/api/handler/user"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func SetupRouter(userHandler *userH.UserHandler, projectHandler *projectH.ProjectHandler, errorHandler *errorH.ErrorHandler) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "http://127.0.0.1:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Route("/api/users", func(u chi.Router) {
		u.Post("/signup", userHandler.CreateUser)
		// u.Post("/login", userHandler.Login)
		// u.Get("/logout", userHandler.Logout)

		// u.Group(func(r chi.Router) {
		// 	r.Use(internalMiddleware.JWTAuth)
		// 	r.Put("/username", userHandler.UpdateUsername)
		// })
	})

	r.Route("/api/projects", func(u chi.Router) {
		u.Group(func(r chi.Router) {
			r.Post("/", projectHandler.CreateProject)
		})
	})

	r.Route("/api/errors", func(u chi.Router) {
		u.Group(func(r chi.Router) {
			r.Post("/", errorHandler.CreateError)
		})
	})

	r.Get("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})

	return r
}

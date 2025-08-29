package router

import (
	"net/http"

	userhandler "github.com/dorianneto/bugfy/internal/api/handler/user"
	authmiddleware "github.com/dorianneto/bugfy/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func SetupRouter(userH *userhandler.UserHandler) http.Handler {
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
		u.Post("/signup", userH.CreateUser)
		// u.Post("/login", userH.Login)
		// u.Get("/logout", userH.Logout)

		u.Group(func(r chi.Router) {
			r.Use(authmiddleware.JWTAuth)
			// r.Put("/username", userH.UpdateUsername)
		})
	})

	r.Get("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})

	return r
}

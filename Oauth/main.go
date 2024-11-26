// main.go
package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sidgupt12/learning/Oauth/handlers"
	// Replace with your actual import path
)

func main() {
	r := chi.NewRouter()

	// r.Use(cors.Handler(cors.Options{
	// 	AllowedOrigins:   []string{"http://localhost:8080"},
	// 	AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
	// 	AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
	// 	AllowCredentials: true,
	// }))

	// Basic middlewares
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Create handler
	h := handlers.NewHandler()

	// Serve static files
	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// OAuth routes
	r.Get("/login/google", h.HandleGoogleLogin)
	r.Get("/callback", h.HandleGoogleCallback)

	// Serve index.html
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/index.html")
	})

	// Start server
	http.ListenAndServe(":8080", r)
}

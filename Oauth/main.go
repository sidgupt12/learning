// main.go
package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sidgupt12/learning/Oauth/handlers"
	// Replace with your actual import path
)

func main() {
	r := chi.NewRouter()

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
	r.Get("/logout", h.HandleLogout)

	// Serve index.html
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/index.html")
	})

	fmt.Println("Server running on http://localhost:8080")

	// Start server
	http.ListenAndServe(":8080", r)
}

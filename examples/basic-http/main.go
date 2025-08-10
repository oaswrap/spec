package main

import (
	"log"
	"net/http"

	"github.com/oaswrap/spec"
	"github.com/oaswrap/spec/option"
)

func main() {
	// Create a new OpenAPI router
	r := spec.NewRouter(
		option.WithTitle("My API"),
		option.WithVersion("1.0.0"),
		option.WithServer("https://api.example.com"),
	)

	// Add routes
	v1 := r.Group("/api/v1")

	v1.Post("/login",
		option.Summary("User login"),
		option.Request(new(LoginRequest)),
		option.Response(200, new(LoginResponse)),
	)

	v1.Get("/users/{id}",
		option.Summary("Get user by ID"),
		option.Request(new(GetUserRequest)),
		option.Response(200, new(User)),
	)

	// Get configuration for doc paths
	cfg := r.Config()

	// Setup HTTP handlers
	mux := http.NewServeMux()

	// Add built-in documentation handlers
	mux.HandleFunc("GET "+cfg.DocsPath, r.DocsHandlerFunc())
	mux.HandleFunc("GET "+cfg.SpecPath, r.SpecHandlerFunc())

	// Add your actual API routes here
	// mux.HandleFunc("POST /api/v1/login", loginHandler)
	// mux.HandleFunc("GET /api/v1/users/{id}", getUserHandler)

	log.Printf("🚀 OpenAPI docs available at: http://localhost:3000%s", cfg.DocsPath)
	log.Printf("📄 OpenAPI spec available at: http://localhost:3000%s", cfg.SpecPath)

	http.ListenAndServe(":3000", mux)
}

type LoginRequest struct {
	Username string `json:"username" required:"true"`
	Password string `json:"password" required:"true"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type GetUserRequest struct {
	ID string `path:"id" required:"true"`
}

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

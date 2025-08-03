package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/oaswrap/spec/adapters/chiopenapi"
	"github.com/oaswrap/spec/option"
)

func main() {
	c := chi.NewRouter()
	// Create a new OpenAPI router
	r := chiopenapi.NewRouter(c,
		option.WithTitle("My API"),
		option.WithVersion("1.0.0"),
	)
	// Add routes
	r.Route("/api/v1", func(r chiopenapi.Router) {
		r.Post("/login", LoginHandler).With(
			option.Summary("User login"),
			option.Request(new(LoginRequest)),
			option.Response(200, new(LoginResponse)),
		)

		r.Get("/users/{id}", GetUserHandler).With(
			option.Summary("Get user by ID"),
			option.Request(new(GetUserRequest)),
			option.Response(200, new(User)),
		)
	})

	// Generate OpenAPI spec
	if err := r.WriteSchemaTo("openapi.yaml"); err != nil {
		log.Fatal(err)
	}

	log.Printf("🚀 OpenAPI docs available at: %s", "http://localhost:3000/docs")

	// Start the server
	if err := http.ListenAndServe(":3000", c); err != nil {
		log.Fatal(err)
	}
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

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	// Simulate login logic
	json.NewEncoder(w).Encode(LoginResponse{Token: "example-token"})
}

func GetUserHandler(w http.ResponseWriter, r *http.Request) {
	var req GetUserRequest
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}
	req.ID = id
	// Simulate fetching user by ID
	user := User{ID: req.ID, Name: "John Doe"}
	json.NewEncoder(w).Encode(user)
}

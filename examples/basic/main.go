package main

import (
	"log"

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

	// Generate OpenAPI spec
	if err := r.WriteSchemaTo("openapi.yaml"); err != nil {
		log.Fatal(err)
	}

	log.Println("✅ OpenAPI spec generated at openapi.yaml")
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

package main

import (
	"log"

	"github.com/oaswrap/spec"
	"github.com/oaswrap/spec/option"
)

func main() {
	// Create a new OpenAPI router with basic info and security scheme
	r := spec.NewRouter(
		option.WithTitle("My API"),
		option.WithVersion("1.0.0"),
		option.WithDescription("Example API"),
		option.WithServer("https://api.example.com"),
		option.WithSecurity("bearerAuth", option.SecurityHTTPBearer("Bearer")),
	)

	// Versioned API group
	v1 := r.Group("/api/v1")

	// Auth routes
	v1.Route("/auth", func(r spec.Router) {
		r.Post("/login",
			option.Summary("User Login"),
			option.Request(new(LoginRequest)),
			option.Response(200, new(Response[Token])),
		)

		r.Get("/me",
			option.Summary("Get Profile"),
			option.Security("bearerAuth"),
			option.Response(200, new(Response[User])),
		)
	}, option.GroupTags("Authentication"))

	// Generate the OpenAPI file
	if err := r.WriteSchemaTo("openapi.yaml"); err != nil {
		log.Fatal(err)
	}

	log.Println("✅ OpenAPI schema generated at openapi.yaml")
}

// Example request & response structs

type LoginRequest struct {
	Username string `json:"username" required:"true"`
	Password string `json:"password" required:"true"`
}

type Token struct {
	AccessToken string `json:"access_token"`
}

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Response[T any] struct {
	Status int `json:"status"`
	Data   T   `json:"data"`
}

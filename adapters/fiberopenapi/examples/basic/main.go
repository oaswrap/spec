package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/oaswrap/spec/adapters/fiberopenapi"
	"github.com/oaswrap/spec/option"
)

func main() {
	app := fiber.New()

	// Create a new OpenAPI router
	r := fiberopenapi.NewRouter(app,
		option.WithTitle("My API"),
		option.WithVersion("1.0.0"),
	)
	// Add routes
	v1 := r.Group("/api/v1")

	v1.Post("/login", LoginHandler).With(
		option.Summary("User login"),
		option.Request(new(LoginRequest)),
		option.Response(200, new(LoginResponse)),
	)

	v1.Get("/users/{id}", GetUserHandler).With(
		option.Summary("Get user by ID"),
		option.Request(new(GetUserRequest)),
		option.Response(200, new(User)),
	)

	// Generate OpenAPI spec
	if err := r.WriteSchemaTo("openapi.yaml"); err != nil {
		log.Fatal(err)
	}
	log.Println("✅ OpenAPI schema written to: openapi.yaml")

	log.Printf("🚀 OpenAPI docs available at: %s", "http://localhost:3000/docs")

	app.Listen(":3000")
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

func LoginHandler(c *fiber.Ctx) error {
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(map[string]string{"error": "Invalid request"})
	}
	// Simulate login logic
	return c.Status(200).JSON(LoginResponse{Token: "example-token"})
}

func GetUserHandler(c *fiber.Ctx) error {
	var req GetUserRequest
	if err := c.ParamsParser(&req); err != nil {
		return c.Status(400).JSON(map[string]string{"error": "Invalid request"})
	}
	// Simulate fetching user by ID
	user := User{ID: req.ID, Name: "John Doe"}
	return c.Status(200).JSON(user)
}

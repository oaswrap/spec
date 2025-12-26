package main

import (
	"log"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/static"
	fiberopenapi "github.com/oaswrap/spec/adapter/fiberopenapi/v3"
	"github.com/oaswrap/spec/option"
)

func main() {
	app := fiber.New()

	// Create a new OpenAPI router
	r := fiberopenapi.NewRouter(app,
		option.WithTitle("My API"),
		option.WithVersion("1.0.0"),
		option.WithSecurity("bearerAuth", option.SecurityHTTPBearer("Bearer")),
	)
	r.Get("/static/*", static.New("../../testdata"))
	// Add routes
	v1 := r.Group("/api/v1")
	v1.Post("/login", LoginHandler).With(
		option.Summary("User login"),
		option.Request(new(LoginRequest)),
		option.Response(fiber.StatusOK, new(LoginResponse)),
	)
	auth := v1.Group("/", AuthMiddleware).With(
		option.GroupSecurity("bearerAuth"),
	)
	auth.Get("/users/:id", GetUserHandler).With(
		option.Summary("Get user by ID"),
		option.Request(new(GetUserRequest)),
		option.Response(fiber.StatusOK, new(User)),
	)

	log.Printf("🚀 OpenAPI docs available at: %s", "http://localhost:3000/docs")

	if err := app.Listen(":3000"); err != nil {
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
	ID string `uri:"id" required:"true"`
}

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func AuthMiddleware(c fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader != "" && authHeader == "Bearer example-token" {
		return c.Next()
	}
	return c.Status(fiber.StatusUnauthorized).JSON(map[string]string{"error": "Unauthorized"})
}

func LoginHandler(c fiber.Ctx) error {
	var req LoginRequest
	if err := c.Bind().JSON(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(map[string]string{"error": "Invalid request"})
	}
	// Simulate login logic
	return c.Status(fiber.StatusOK).JSON(LoginResponse{Token: "example-token"})
}

func GetUserHandler(c fiber.Ctx) error {
	var req GetUserRequest
	if err := c.Bind().URI(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(map[string]string{"error": "Invalid request"})
	}
	// Simulate fetching user by ID
	user := User{ID: req.ID, Name: "John Doe"}
	return c.Status(fiber.StatusOK).JSON(user)
}

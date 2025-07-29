package main

import (
	"fmt"
	"log"

	"github.com/faizlabs/openapi-wrapper/fiberopenapi"
	"github.com/faizlabs/openapi-wrapper/option"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	// Initialize OpenAPI router with configuration
	r := fiberopenapi.NewRouter(app,
		option.WithTitle("My API"),
		option.WithVersion("1.0.0"),
		option.WithDescription("This is a sample API"),
		option.WithDocsPath("/docs"),
		option.WithServer("http://localhost:3000", "Local server"),
		option.WithSecurity("bearerAuth", option.SecurityHTTPBearer()),
		option.WithDebug(true),
	)

	api := r.Group("/api")
	v1 := api.Group("/v1")

	v1.Route("/auth", func(r fiberopenapi.Router) {
		r.Post("/login", dummyHandler).With(
			option.Summary("User Login"),
			option.Description("Endpoint for user login"),
			option.Request(new(LoginRequest)),
			option.Response(200, new(Response[Token])),
			option.Response(400, new(ErrorResponse)),
			option.Response(422, new(ValidationResponse)),
		)
		r.Post("/register", dummyHandler).With(
			option.Summary("User Registration"),
			option.Description("Endpoint for user registration"),
			option.Request(new(RegisterRequest)),
			option.Response(201, new(Response[Token])),
			option.Response(400, new(ErrorResponse)),
			option.Response(422, new(ValidationResponse)),
		)
		r.Post("/refresh-token", dummyHandler).With(
			option.Hide(),
			option.Summary("Refresh Access Token"),
			option.Description("Endpoint to refresh access token using refresh token"),
			option.Request(new(RefreshTokenRequest)),
			option.Response(200, new(Response[Token])),
			option.Response(400, new(ErrorResponse)),
			option.Response(401, new(ErrorResponse)),
			option.Response(422, new(ValidationResponse)),
		)
		r.Get("/me", dummyHandler).With(
			option.Summary("Get User Profile"),
			option.Description("Endpoint to get the authenticated user's profile"),
			option.Security("bearerAuth"),
			option.Response(200, new(Response[User])),
			option.Response(401, new(ErrorResponse)),
		)
	}).With(
		option.RouteTags("Authentication"),
	)

	v1.Route("/profile", func(r fiberopenapi.Router) {
		r.Put("/update", dummyHandler).With(
			option.Summary("Update User Profile"),
			option.Description("Endpoint to update the user's profile"),
			option.Request(new(UpdateProfileRequest)),
			option.Response(200, new(Response[User])),
			option.Response(400, new(ErrorResponse)),
			option.Response(422, new(ValidationResponse)),
		)
		r.Put("/password", dummyHandler).With(
			option.Summary("Update Password"),
			option.Description("Endpoint to update the user's password"),
			option.Request(new(UpdatePasswordRequest)),
			option.Response(200, new(MessageResponse)),
			option.Response(400, new(ErrorResponse)),
			option.Response(422, new(ValidationResponse)),
		)
		r.Post("/delete-account", dummyHandler).With(
			option.Summary("Delete User Account"),
			option.Description("Endpoint to delete the user's account"),
			option.Request(new(DeleteAccountRequest)),
			option.Response(200, new(MessageResponse)),
			option.Response(400, new(ErrorResponse)),
			option.Response(422, new(ValidationResponse)),
		)
	}).With(
		option.RouteTags("Profile"),
		option.RouteSecurity("bearerAuth"),
	)

	// Validate the OpenAPI configuration
	if err := r.Validate(); err != nil {
		log.Fatalf("OpenAPI validation failed: %v", err)
	}

	// Write the OpenAPI schema to a file (Optional)
	if err := r.WriteSchemaTo("openapi.yaml"); err != nil {
		log.Fatalf("Failed to write OpenAPI schema: %v", err)
	}
	if err := r.WriteSchemaTo("openapi.json"); err != nil {
		log.Fatalf("Failed to write OpenAPI schema: %v", err)
	}

	fmt.Println("Open http://localhost:3000/docs to view the OpenAPI documentation")

	app.Listen(":3000")
}

func dummyHandler(c *fiber.Ctx) error {
	// Dummy handler for demonstration purposes
	return c.SendString("This is a dummy handler")
}

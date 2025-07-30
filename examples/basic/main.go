package main

import (
	"log"

	"github.com/oaswrap/spec"
	"github.com/oaswrap/spec/option"
)

func main() {
	r := spec.NewGenerator(
		option.WithTitle("My API"),
		option.WithVersion("1.0.0"),
		option.WithDescription("This is my API"),
		option.WithServer("https://api.example.com", option.ServerDescription("Main server")),
		option.WithSecurity("bearerAuth", option.SecurityHTTPBearer("Bearer")),
	)
	r.Post("/api/v1/login",
		option.OperationID("loginUser"),
		option.Summary("Login User"),
		option.Description("Logs in a user and returns a token"),
		option.Tags("Authentication"),
		option.Request(new(LoginRequest)),
		option.Response(200, new(Response[Token])),
		option.Response(401, new(ErrorResponse)),
		option.Response(422, new(ValidationResponse)),
	)
	r.Get("/api/v1/users/{user_id}",
		option.OperationID("getUserDetail"),
		option.Summary("Get User Detail"),
		option.Description("Retrieves details of a user by ID"),
		option.Tags("Users"),
		option.Security("bearerAuth"),
		option.Request(new(GetUserDetailRequest)),
		option.Response(200, new(Response[User])),
		option.Response(401, new(ErrorResponse)),
	)

	if err := r.Validate(); err != nil {
		log.Fatal("Failed to validate OpenAPI schema: ", err)
	}

	if err := r.WriteSchemaTo("openapi.yaml"); err != nil {
		log.Fatal("Failed to write OpenAPI schema: ", err)
	}

	log.Println("OpenAPI schema generated successfully at openapi.yaml")
}

type LoginRequest struct {
	Username   string `json:"username" example:"john_doe" validate:"required"`
	Password   string `json:"password" example:"password123" validate:"required"`
	RememberMe bool   `json:"remember_me" example:"true"`
}

type GetUserDetailRequest struct {
	UserID string `path:"user_id" example:"12345"`
}

type Token struct {
	AccessToken  string `json:"access_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	RefreshToken string `json:"refresh_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

type User struct {
	ID       string `json:"id" example:"12345"`
	Username string `json:"username" example:"john_doe"`
	Email    string `json:"email" example:"john_doe@example.com"`
	Name     string `json:"name" example:"John Doe"`
}

type Response[T any] struct {
	Status int `json:"status" example:"200"`
	Data   T   `json:"data,omitempty"`
}

type MessageResponse struct {
	Status  int    `json:"status" example:"200"`
	Message string `json:"message" example:"Operation successful"`
}

type ErrorResponse struct {
	Status int    `json:"status" example:"400"`
	Title  string `json:"title" example:"Bad Request"`
	Detail string `json:"detail" example:"Invalid input data"`
}

type ValidationResponse struct {
	Status int          `json:"status" example:"422"`
	Title  string       `json:"title" example:"Validation Error"`
	Detail string       `json:"detail" example:"Input data does not meet validation criteria"`
	Errors []FieldError `json:"errors" example:"[{\"field\":\"email\",\"message\":\"Email is required\"}]"`
}

type FieldError struct {
	Field   string `json:"field" example:"email"`
	Message string `json:"message" example:"Email is required"`
}

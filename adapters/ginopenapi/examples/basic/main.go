package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/oaswrap/spec/adapters/ginopenapi"
	"github.com/oaswrap/spec/option"
)

func main() {
	e := gin.Default()

	// Create a new OpenAPI router
	r := ginopenapi.NewRouter(e,
		option.WithTitle("My API"),
		option.WithVersion("1.0.0"),
	)
	// Add routes
	v1 := r.Group("/api/v1")

	v1.POST("/login", LoginHandler).With(
		option.Summary("User login"),
		option.Request(new(LoginRequest)),
		option.Response(200, new(LoginResponse)),
	)

	v1.GET("/users/{id}", GetUserHandler).With(
		option.Summary("Get user by ID"),
		option.Request(new(GetUserRequest)),
		option.Response(200, new(User)),
	)

	// Generate OpenAPI spec
	if err := r.WriteSchemaTo("openapi.yaml"); err != nil {
		log.Fatal(err)
	}

	log.Printf("🚀 OpenAPI docs available at: %s", "http://localhost:3000/docs")

	if err := e.Run(":3000"); err != nil {
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

func LoginHandler(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, map[string]string{"error": "Invalid request"})
		return
	}
	// Simulate login logic
	c.JSON(200, LoginResponse{Token: "example-token"})
}

func GetUserHandler(c *gin.Context) {
	var req GetUserRequest
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(400, map[string]string{"error": "Invalid request"})
		return
	}
	// Simulate fetching user by ID
	user := User{ID: req.ID, Name: "John Doe"}
	c.JSON(200, user)
}

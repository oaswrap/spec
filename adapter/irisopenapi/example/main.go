package main

import (
	"log"

	"github.com/kataras/iris/v12"

	"github.com/oaswrap/spec/option"

	"github.com/oaswrap/spec/adapter/irisopenapi"
)

func main() {
	app := iris.New()

	r := irisopenapi.NewRouter(app,
		option.WithTitle("My API"),
		option.WithVersion("1.0.0"),
		option.WithSecurity("bearerAuth", option.SecurityHTTPBearer("Bearer")),
	)

	v1 := r.Party("/api/v1")
	v1.Post("/login", LoginHandler).With(
		option.Summary("User login"),
		option.Request(new(LoginRequest)),
		option.Response(200, new(LoginResponse)),
	)

	auth := v1.Party("", AuthMiddleware).With(
		option.GroupSecurity("bearerAuth"),
	)
	auth.Get("/users/{id}", GetUserHandler).With(
		option.Summary("Get user by ID"),
		option.Request(new(GetUserRequest)),
		option.Response(200, new(User)),
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
	ID string `param:"id" required:"true"`
}

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func AuthMiddleware(ctx iris.Context) {
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "Bearer example-token" {
		ctx.Next()
		return
	}

	ctx.StatusCode(iris.StatusUnauthorized)
	_ = ctx.JSON(map[string]string{"error": "Unauthorized"})
	ctx.StopExecution()
}

func LoginHandler(ctx iris.Context) {
	var req LoginRequest
	if err := ctx.ReadJSON(&req); err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		_ = ctx.JSON(map[string]string{"error": "Invalid request"})
		return
	}
	_ = ctx.JSON(LoginResponse{Token: "example-token"})
}

func GetUserHandler(ctx iris.Context) {
	id := ctx.Params().Get("id")
	_ = ctx.JSON(User{ID: id, Name: "John Doe"})
}

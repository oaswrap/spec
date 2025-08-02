package echoopenapi_test

import (
	"flag"
	"fmt"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/oaswrap/spec/adapters/echoopenapi"
	"github.com/oaswrap/spec/option"
	"github.com/oaswrap/spec/pkg/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var update = flag.Bool("update", false, "update golden files")

type HelloRequest struct {
	Name string `json:"name" query:"name"`
}

type HelloResponse struct {
	Response string `json:"response"`
}

func HelloHandler(c echo.Context) error {
	var req HelloRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(400, map[string]string{"error": "Invalid request"})
	}
	return c.JSON(200, map[string]string{"response": "Hello " + req.Name})
}

type EchoOpenAPISingleRouteFunc func(path string, handler echo.HandlerFunc, m ...echo.MiddlewareFunc) echoopenapi.Route

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

type Response[T any] struct {
	Status int `json:"status"`
	Data   T   `json:"data"`
}

type Token struct {
	Token string `json:"token"`
}

type User struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ErrorResponse struct {
	Status int    `json:"status"`
	Title  string `json:"title"`
	Detail string `json:"detail,omitempty"`
}

type ValidationResponse struct {
	ErrorResponse
	Errors []struct {
		Field   string `json:"field"`
		Message string `json:"message"`
	} `json:"errors"`
}

func DummyHandler(c echo.Context) error {
	return c.JSON(200, map[string]string{"message": "Dummy handler"})
}

func TestRouter_Spec(t *testing.T) {
	tests := []struct {
		name      string
		golden    string
		opts      []option.OpenAPIOption
		setup     func(r echoopenapi.Router)
		shouldErr bool
	}{
		{
			name:   "Authentication Routes",
			golden: "auth_routes",
			opts: []option.OpenAPIOption{
				option.WithTitle("Authentication API"),
				option.WithDescription("API for user authentication"),
				option.WithVersion("1.0.0"),
				option.WithSecurity("bearerAuth", option.SecurityHTTPBearer("Bearer")),
			},
			setup: func(r echoopenapi.Router) {
				api := r.Group("/api")
				v1 := api.Group("/v1")
				auth := v1.Group("/").With(option.GroupTags("Authentication"))
				auth.POST("/login", DummyHandler).With(
					option.Summary("User Login"),
					option.Description("Endpoint for user login"),
					option.OperationID("userLogin"),
					option.Request(new(LoginRequest)),
					option.Response(200, new(Response[Token])),
					option.Response(400, new(ErrorResponse)),
					option.Response(422, new(ValidationResponse)),
				)
				auth.POST("/register", DummyHandler).With(
					option.Summary("User Registration"),
					option.Description("Endpoint for user registration"),
					option.OperationID("userRegister"),
					option.Request(new(RegisterRequest)),
					option.Response(201, new(Response[User])),
					option.Response(400, new(ErrorResponse)),
					option.Response(422, new(ValidationResponse)),
				)

				v1.GET("/profile", DummyHandler).With(
					option.Summary("Get User Profile"),
					option.Description("Endpoint to get user profile information"),
					option.OperationID("getUserProfile"),
					option.Tags("User"),
					option.Security("bearerAuth"),
					option.Response(200, new(Response[User])),
					option.Response(401, new(ErrorResponse)),
				)
			},
		},
	}

	versions := map[string]string{
		"3.0.3": "3",
		"3.1.0": "31",
	}

	for _, tt := range tests {
		for version, suffix := range versions {
			t.Run(fmt.Sprintf("%s_%s", tt.name, suffix), func(t *testing.T) {
				e := echo.New()
				defaultOpts := []option.OpenAPIOption{
					option.WithOpenAPIVersion(version),
				}
				tt.opts = append(defaultOpts, tt.opts...)
				r := echoopenapi.NewRouter(e, tt.opts...)
				tt.setup(r)

				err := r.Validate()
				if tt.shouldErr {
					assert.Error(t, err, "Expected error for test: %s", tt.name)
					return
				}
				assert.NoError(t, err, "Expected no error for test: %s", tt.name)

				// Test the OpenAPI schema generation
				schema, err := r.GenerateSchema()
				require.NoError(t, err, "failed to generate schema")

				golden := filepath.Join("testdata", tt.golden+"_"+suffix+".yaml")
				if *update {
					err = r.WriteSchemaTo(golden)
					require.NoError(t, err, "failed to write golden file")
					t.Logf("Updated golden file: %s", golden)
				}

				want, err := os.ReadFile(golden)
				require.NoError(t, err, "failed to read golden file %s", golden)

				testutil.EqualYAML(t, want, schema)
			})
		}
	}
}

func TestRouter_Single(t *testing.T) {
	tests := []struct {
		method     string
		path       string
		methodFunc func(r echoopenapi.Router) EchoOpenAPISingleRouteFunc
	}{
		{
			method:     "GET",
			path:       "/hello",
			methodFunc: func(r echoopenapi.Router) EchoOpenAPISingleRouteFunc { return r.GET },
		},
		{
			method:     "POST",
			path:       "/hello",
			methodFunc: func(r echoopenapi.Router) EchoOpenAPISingleRouteFunc { return r.POST },
		},
		{
			method:     "PUT",
			path:       "/hello",
			methodFunc: func(r echoopenapi.Router) EchoOpenAPISingleRouteFunc { return r.PUT },
		},
		{
			method:     "DELETE",
			path:       "/hello",
			methodFunc: func(r echoopenapi.Router) EchoOpenAPISingleRouteFunc { return r.DELETE },
		},
		{
			method:     "PATCH",
			path:       "/hello",
			methodFunc: func(r echoopenapi.Router) EchoOpenAPISingleRouteFunc { return r.PATCH },
		},
		{
			method:     "HEAD",
			path:       "/hello",
			methodFunc: func(r echoopenapi.Router) EchoOpenAPISingleRouteFunc { return r.HEAD },
		},
		{
			method:     "OPTIONS",
			path:       "/hello",
			methodFunc: func(r echoopenapi.Router) EchoOpenAPISingleRouteFunc { return r.OPTIONS },
		},
		{
			method:     "TRACE",
			path:       "/hello",
			methodFunc: func(r echoopenapi.Router) EchoOpenAPISingleRouteFunc { return r.TRACE },
		},
		{
			method:     "CONNECT",
			path:       "/hello",
			methodFunc: func(r echoopenapi.Router) EchoOpenAPISingleRouteFunc { return r.CONNECT },
		},
	}
	for _, tt := range tests {
		t.Run(tt.method, func(t *testing.T) {
			e := echo.New()
			r := echoopenapi.NewGenerator(e,
				option.WithTitle("Test API Single"),
				option.WithVersion("1.0.0"),
			)
			// Setup the route
			route := tt.methodFunc(r)(tt.path, HelloHandler).With(
				option.Summary("Hello Handler"),
				option.Description("Handles hello requests"),
				option.OperationID(fmt.Sprintf("hello%s", tt.method)),
				option.Tags("greeting"),
				option.Request(new(HelloRequest)),
				option.Response(200, new(HelloResponse)),
			)

			// Verify the route is registered
			assert.Equal(t, tt.method, route.Method(), "Expected method to be %s", tt.method)
			assert.Equal(t, tt.path, route.Path(), "Expected path to be %s", tt.path)
			assert.NotEmpty(t, route.Name(), "Expected route name to be set")
			req := httptest.NewRequest(tt.method, tt.path, nil)
			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, req)

			assert.Equal(t, 200, rec.Code, "Expected status code 200 for %s %s", tt.method, tt.path)
			assert.Contains(t, rec.Body.String(), "Hello", "Expected response to contain 'Hello'")

			// Verify the OpenAPI schema
			if tt.method == "CONNECT" {
				schema, err := r.MarshalYAML()
				assert.NoError(t, err, "Expected no error while generating OpenAPI schema")
				assert.NotEmpty(t, schema, "Expected OpenAPI schema to be generated")
				assert.NotContains(t, string(schema), fmt.Sprintf("operationId: hello%s", tt.method))
				return
			}
			schema, err := r.MarshalYAML()
			assert.NoError(t, err, "Expected no error while generating OpenAPI schema")
			assert.NotEmpty(t, schema, "Expected OpenAPI schema to be generated")
			assert.Contains(t, string(schema), fmt.Sprintf("operationId: hello%s", tt.method))
			assert.Contains(t, string(schema), "summary: Hello Handler", "Expected OpenAPI schema to contain the summary")
		})
	}
}

func TestRouter_Use(t *testing.T) {
	t.Run("should call middleware", func(t *testing.T) {
		totalCalled := 0
		middleware := func(next echo.HandlerFunc) echo.HandlerFunc {
			return func(c echo.Context) error {
				totalCalled++
				return next(c)
			}
		}
		e := echo.New()
		r := echoopenapi.NewGenerator(e,
			option.WithTitle("Test API Middleware"),
			option.WithVersion("1.0.0"),
		)
		r.Use(middleware)

		r.GET("/test", func(c echo.Context) error {
			return c.String(200, "Hello Middleware")
		})
		req := httptest.NewRequest(echo.GET, "/test", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)
		assert.Equal(t, 200, rec.Code, "Expected status code 200")
		assert.Equal(t, "Hello Middleware", rec.Body.String(), "Expected response body to be 'Hello Middleware'")
		assert.Equal(t, 1, totalCalled, "Expected middleware to be called once")
	})
}

func TestRouter_Group(t *testing.T) {
	e := echo.New()
	r := echoopenapi.NewGenerator(e,
		option.WithTitle("Test API Group"),
		option.WithVersion("1.0.0"),
	)

	v1 := r.Group("/v1")
	v1.GET("/hello", HelloHandler).With(
		option.Summary("Hello Handler V1"),
		option.Description("Handles hello requests for V1"),
		option.OperationID("helloV1"),
		option.Tags("greeting"),
		option.Request(new(HelloRequest)),
		option.Response(200, new(HelloResponse)),
	)

	v2 := r.Group("/v2")
	v2.GET("/hello", HelloHandler).With(
		option.Summary("Hello Handler V2"),
		option.Description("Handles hello requests for V2"),
		option.OperationID("helloV2"),
		option.Tags("greeting"),
		option.Request(new(HelloRequest)),
		option.Response(200, new(HelloResponse)),
	)

	req := httptest.NewRequest(echo.GET, "/v1/hello?name=World", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	assert.Equal(t, 200, rec.Code, "Expected status code 200 for /v1/hello")
	assert.Contains(t, rec.Body.String(), "Hello World", "Expected response to contain 'Hello World'")

	req = httptest.NewRequest(echo.GET, "/v2/hello?name=Echo", nil)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	assert.Equal(t, 200, rec.Code, "Expected status code 200 for /v2/hello")
	assert.Contains(t, rec.Body.String(), "Hello Echo", "Expected response to contain 'Hello Echo'")
}

func TestRouter_StaticFS(t *testing.T) {
	e := echo.New()
	r := echoopenapi.NewGenerator(e,
		option.WithTitle("Test API StaticFS"),
		option.WithVersion("1.0.0"),
	)
	tempDir := t.TempDir()
	// Create a test file in the temporary directory
	testFilePath := fmt.Sprintf("%s/test.txt", tempDir)
	if err := os.WriteFile(testFilePath, []byte("This is a test file."), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	// Serve static files from the temporary directory
	r.StaticFS("/static", os.DirFS(tempDir))

	req := httptest.NewRequest(echo.GET, "/static/test.txt", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, 200, rec.Code, "Expected status code 200 for /static/test.txt")
	assert.Equal(t, "This is a test file.", rec.Body.String(), "Expected response body to match test file content")
}

func TestRouter_Static(t *testing.T) {
	e := echo.New()
	r := echoopenapi.NewGenerator(e,
		option.WithTitle("Test API Static"),
		option.WithVersion("1.0.0"),
	)
	tempDir := t.TempDir()
	// Create a test file in the temporary directory
	testFilePath := fmt.Sprintf("%s/test.txt", tempDir)
	if err := os.WriteFile(testFilePath, []byte("This is a test file."), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	// Serve static files from the temporary directory
	r.Static("/static", tempDir)

	req := httptest.NewRequest(echo.GET, "/static/test.txt", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, 200, rec.Code, "Expected status code 200 for /static/test.txt")
	assert.Equal(t, "This is a test file.", rec.Body.String(), "Expected response body to match test file content")
}

func TestRouter_File(t *testing.T) {
	e := echo.New()
	r := echoopenapi.NewGenerator(e,
		option.WithTitle("Test API File"),
		option.WithVersion("1.0.0"),
	)
	tempDir := t.TempDir()
	// Create a test file in the temporary directory
	testFilePath := fmt.Sprintf("%s/test.txt", tempDir)
	if err := os.WriteFile(testFilePath, []byte("This is a test file."), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	// Serve a static file
	r.File("/test.txt", testFilePath)

	req := httptest.NewRequest(echo.GET, "/test.txt", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, 200, rec.Code, "Expected status code 200 for /test.txt")
	assert.Equal(t, "This is a test file.", rec.Body.String(), "Expected response body to match test file content")
}

func TestRouter_FileFS(t *testing.T) {
	e := echo.New()
	r := echoopenapi.NewGenerator(e,
		option.WithTitle("Test API FileFS"),
		option.WithVersion("1.0.0"),
	)
	tempDir := t.TempDir()
	// Create a test file in the temporary directory
	testFilePath := fmt.Sprintf("%s/test.txt", tempDir)
	if err := os.WriteFile(testFilePath, []byte("This is a test file."), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	// Serve a static file from the filesystem
	r.FileFS("/test.txt", "test.txt", os.DirFS(tempDir))

	req := httptest.NewRequest(echo.GET, "/test.txt", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, 200, rec.Code, "Expected status code 200 for /test.txt")
	assert.Equal(t, "This is a test file.", rec.Body.String(), "Expected response body to match test file content")
}

func TestGenerator_WriteSchemaTo(t *testing.T) {
	e := echo.New()
	r := echoopenapi.NewGenerator(e,
		option.WithTitle("Test API WriteSchemaTo"),
		option.WithVersion("1.0.0"),
	)

	// Define a route
	r.GET("/hello", HelloHandler).With(
		option.Summary("Hello Handler"),
		option.Description("Handles hello requests"),
		option.OperationID("hello"),
		option.Tags("greeting"),
		option.Request(new(HelloRequest)),
		option.Response(200, new(HelloResponse)),
	)

	// Write the OpenAPI schema to a file
	tempFile := t.TempDir() + "/openapi.yaml"
	err := r.WriteSchemaTo(tempFile)
	assert.NoError(t, err, "Expected no error while writing OpenAPI schema to file")

	// Verify the file exists and is not empty
	info, err := os.Stat(tempFile)
	assert.NoError(t, err, "Expected no error while checking file stats")
	assert.False(t, info.IsDir(), "Expected file to not be a directory")
	assert.Greater(t, info.Size(), int64(0), "Expected file size to be greater than 0")
}

func TestGenerator_MarshalJSON(t *testing.T) {
	e := echo.New()
	r := echoopenapi.NewGenerator(e,
		option.WithTitle("Test API MarshalJSON"),
		option.WithVersion("1.0.0"),
	)

	// Define a route
	r.GET("/hello", HelloHandler).With(
		option.Summary("Hello Handler"),
		option.Description("Handles hello requests"),
		option.OperationID("hello"),
		option.Tags("greeting"),
		option.Request(new(HelloRequest)),
		option.Response(200, new(HelloResponse)),
	)

	// Marshal the OpenAPI schema to JSON
	schema, err := r.MarshalJSON()
	assert.NoError(t, err, "Expected no error while marshaling OpenAPI schema to JSON")
	assert.NotEmpty(t, schema, "Expected OpenAPI schema JSON to not be empty")
}

func TestGenerator_Docs(t *testing.T) {
	e := echo.New()
	r := echoopenapi.NewGenerator(e,
		option.WithTitle("Test API Docs"),
		option.WithVersion("1.0.0"),
	)

	// Define a route
	r.GET("/hello", HelloHandler).With(
		option.Summary("Hello Handler"),
		option.Description("Handles hello requests"),
		option.OperationID("hello"),
		option.Tags("greeting"),
		option.Request(new(HelloRequest)),
		option.Response(200, new(HelloResponse)),
	)

	req := httptest.NewRequest(echo.GET, "/docs", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, 200, rec.Code, "Expected status code 200 for /docs")
	assert.Contains(t, rec.Body.String(), "Test API Docs", "Expected response to contain API title")
}

func TestGenerator_DisableDocs(t *testing.T) {
	e := echo.New()
	r := echoopenapi.NewGenerator(e,
		option.WithTitle("Test API Disable Docs"),
		option.WithVersion("1.0.0"),
		option.WithDisableDocs(true),
	)

	// Define a route
	r.GET("/hello", HelloHandler).With(
		option.Summary("Hello Handler"),
		option.Description("Handles hello requests"),
		option.OperationID("hello"),
		option.Tags("greeting"),
		option.Request(new(HelloRequest)),
		option.Response(200, new(HelloResponse)),
	)

	req := httptest.NewRequest(echo.GET, "/docs", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, 404, rec.Code, "Expected status code 404 for /docs when docs are disabled")
}

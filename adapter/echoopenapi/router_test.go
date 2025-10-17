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
	"github.com/oaswrap/spec/adapter/echoopenapi"
	"github.com/oaswrap/spec/openapi"
	"github.com/oaswrap/spec/option"
	"github.com/oaswrap/spec/pkg/dto"
	"github.com/oaswrap/spec/pkg/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//nolint:gochecknoglobals // test flag for golden file updates
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
			name:   "Petstore API",
			golden: "petstore",
			opts: []option.OpenAPIOption{
				option.WithDescription("This is a sample Petstore server."),
				option.WithVersion("1.0.0"),
				option.WithTermsOfService("https://swagger.io/terms/"),
				option.WithContact(openapi.Contact{
					Email: "apiteam@swagger.io",
				}),
				option.WithLicense(openapi.License{
					Name: "Apache 2.0",
					URL:  "https://www.apache.org/licenses/LICENSE-2.0.html",
				}),
				option.WithExternalDocs("https://swagger.io", "Find more info here about swagger"),
				option.WithServer("https://petstore3.swagger.io/api/v3"),
				option.WithTags(
					openapi.Tag{
						Name:        "pet",
						Description: "Everything about your Pets",
						ExternalDocs: &openapi.ExternalDocs{
							Description: "Find out more about our Pets",
							URL:         "https://swagger.io",
						},
					},
					openapi.Tag{
						Name:        "store",
						Description: "Access to Petstore orders",
						ExternalDocs: &openapi.ExternalDocs{
							Description: "Find out more about our Store",
							URL:         "https://swagger.io",
						},
					},
					openapi.Tag{
						Name:        "user",
						Description: "Operations about user",
					},
				),
				option.WithSecurity("petstore_auth", option.SecurityOAuth2(
					openapi.OAuthFlows{
						Implicit: &openapi.OAuthFlowsImplicit{
							AuthorizationURL: "https://petstore3.swagger.io/oauth/authorize",
							Scopes: map[string]string{
								"write:pets": "modify pets in your account",
								"read:pets":  "read your pets",
							},
						},
					}),
				),
				option.WithSecurity("apiKey", option.SecurityAPIKey("api_key", openapi.SecuritySchemeAPIKeyInHeader)),
			},
			setup: func(r echoopenapi.Router) {
				pet := r.Group("/pet").With(
					option.GroupTags("pet"),
					option.GroupSecurity("petstore_auth", "write:pets", "read:pets"),
				)
				pet.PUT("/", nil).With(
					option.OperationID("updatePet"),
					option.Summary("Update an existing pet"),
					option.Description("Update the details of an existing pet in the store."),
					option.Request(new(dto.Pet)),
					option.Response(200, new(dto.Pet)),
				)
				pet.POST("/", nil).With(
					option.OperationID("addPet"),
					option.Summary("Add a new pet"),
					option.Description("Add a new pet to the store."),
					option.Request(new(dto.Pet)),
					option.Response(201, new(dto.Pet)),
				)
				pet.GET("/findByStatus", nil).With(
					option.OperationID("findPetsByStatus"),
					option.Summary("Find pets by status"),
					option.Description("Finds Pets by status. Multiple status values can be provided with comma separated strings."),
					option.Request(new(struct {
						Status string `query:"status" enum:"available,pending,sold"`
					})),
					option.Response(200, new([]dto.Pet)),
				)
				pet.GET("/findByTags", nil).With(
					option.OperationID("findPetsByTags"),
					option.Summary("Find pets by tags"),
					option.Description("Finds Pets by tags. Multiple tags can be provided with comma separated strings."),
					option.Request(new(struct {
						Tags []string `query:"tags"`
					})),
					option.Response(200, new([]dto.Pet)),
				)
				pet.POST("/{petId}/uploadImage", nil).With(
					option.OperationID("uploadFile"),
					option.Summary("Upload an image for a pet"),
					option.Description("Uploads an image for a pet."),
					option.Request(new(dto.UploadImageRequest)),
					option.Response(200, new(dto.APIResponse)),
				)
				pet.GET("/{petId}", nil).With(
					option.OperationID("getPetById"),
					option.Summary("Get pet by ID"),
					option.Description("Retrieve a pet by its ID."),
					option.Request(new(struct {
						ID int `path:"petId" required:"true"`
					})),
					option.Response(200, new(dto.Pet)),
				)
				pet.POST("/{petId}", nil).With(
					option.OperationID("updatePetWithForm"),
					option.Summary("Update pet with form"),
					option.Description("Updates a pet in the store with form data."),
					option.Request(new(dto.UpdatePetWithFormRequest)),
					option.Response(200, nil),
				)
				pet.DELETE("/{petId}", nil).With(
					option.OperationID("deletePet"),
					option.Summary("Delete a pet"),
					option.Description("Delete a pet from the store by its ID."),
					option.Request(new(dto.DeletePetRequest)),
					option.Response(204, nil),
				)
				store := r.Group("/store").With(
					option.GroupTags("store"),
				)
				store.POST("/order", nil).With(
					option.OperationID("placeOrder"),
					option.Summary("Place an order"),
					option.Description("Place a new order for a pet."),
					option.Request(new(dto.Order)),
					option.Response(201, new(dto.Order)),
				)
				store.GET("/order/{orderId}", nil).With(
					option.OperationID("getOrderById"),
					option.Summary("Get order by ID"),
					option.Description("Retrieve an order by its ID."),
					option.Request(new(struct {
						ID int `path:"orderId" required:"true"`
					})),
					option.Response(200, new(dto.Order)),
					option.Response(404, nil),
				)
				store.DELETE("/order/{orderId}", nil).With(
					option.OperationID("deleteOrder"),
					option.Summary("Delete an order"),
					option.Description("Delete an order by its ID."),
					option.Request(new(struct {
						ID int `path:"orderId" required:"true"`
					})),
					option.Response(204, nil),
				)

				user := r.Group("/user").With(
					option.GroupTags("user"),
				)
				user.POST("/createWithList", nil).With(
					option.OperationID("createUsersWithList"),
					option.Summary("Create users with list"),
					option.Description("Create multiple users in the store with a list."),
					option.Request(new([]dto.PetUser)),
					option.Response(201, nil),
				)
				user.POST("/", nil).With(
					option.OperationID("createUser"),
					option.Summary("Create a new user"),
					option.Description("Create a new user in the store."),
					option.Request(new(dto.PetUser)),
					option.Response(201, new(dto.PetUser)),
				)
				user.GET("/{username}", nil).With(
					option.OperationID("getUserByName"),
					option.Summary("Get user by username"),
					option.Description("Retrieve a user by their username."),
					option.Request(new(struct {
						Username string `path:"username" required:"true"`
					})),
					option.Response(200, new(dto.PetUser)),
					option.Response(404, nil),
				)
				user.PUT("/{username}", nil).With(
					option.OperationID("updateUser"),
					option.Summary("Update an existing user"),
					option.Description("Update the details of an existing user."),
					option.Request(new(struct {
						dto.PetUser

						Username string `path:"username" required:"true"`
					})),
					option.Response(200, new(dto.PetUser)),
					option.Response(404, nil),
				)
				user.DELETE("/{username}", nil).With(
					option.OperationID("deleteUser"),
					option.Summary("Delete a user"),
					option.Description("Delete a user from the store by their username."),
					option.Request(new(struct {
						Username string `path:"username" required:"true"`
					})),
					option.Response(204, nil),
				)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			r := echoopenapi.NewRouter(e, tt.opts...)
			tt.setup(r)

			err := r.Validate()
			if tt.shouldErr {
				require.Error(t, err, "Expected error for test: %s", tt.name)
				return
			}
			require.NoError(t, err, "Expected no error for test: %s", tt.name)

			// Test the OpenAPI schema generation
			schema, err := r.GenerateSchema()
			require.NoError(t, err, "failed to generate schema")

			golden := filepath.Join("testdata", tt.golden+".yaml")
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

type SingleRouteFunc func(path string, handler echo.HandlerFunc, m ...echo.MiddlewareFunc) echoopenapi.Route

func TestRouter_Single(t *testing.T) {
	tests := []struct {
		method     string
		path       string
		methodFunc func(r echoopenapi.Router) SingleRouteFunc
	}{
		{"GET", "/hello", func(r echoopenapi.Router) SingleRouteFunc { return r.GET }},
		{"POST", "/hello", func(r echoopenapi.Router) SingleRouteFunc { return r.POST }},
		{"PUT", "/hello", func(r echoopenapi.Router) SingleRouteFunc { return r.PUT }},
		{"PATCH", "/hello", func(r echoopenapi.Router) SingleRouteFunc { return r.PATCH }},
		{"DELETE", "/hello", func(r echoopenapi.Router) SingleRouteFunc { return r.DELETE }},
		{"HEAD", "/hello", func(r echoopenapi.Router) SingleRouteFunc { return r.HEAD }},
		{"OPTIONS", "/hello", func(r echoopenapi.Router) SingleRouteFunc { return r.OPTIONS }},
		{"TRACE", "/hello", func(r echoopenapi.Router) SingleRouteFunc { return r.TRACE }},
		{"CONNECT", "/hello", func(r echoopenapi.Router) SingleRouteFunc { return r.CONNECT }},
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
				require.NoError(t, err, "Expected no error while generating OpenAPI schema")
				assert.NotEmpty(t, schema, "Expected OpenAPI schema to be generated")
				assert.NotContains(t, string(schema), fmt.Sprintf("operationId: hello%s", tt.method))
				return
			}
			schema, err := r.MarshalYAML()
			require.NoError(t, err, "Expected no error while generating OpenAPI schema")
			assert.NotEmpty(t, schema, "Expected OpenAPI schema to be generated")
			assert.Contains(t, string(schema), fmt.Sprintf("operationId: hello%s", tt.method))
			assert.Contains(
				t,
				string(schema),
				"summary: Hello Handler",
				"Expected OpenAPI schema to contain the summary",
			)
		})
	}
}

func TestRouter(t *testing.T) {
	t.Run("Use", func(t *testing.T) {
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
	require.NoError(t, err, "Expected no error while writing OpenAPI schema to file")

	// Verify the file exists and is not empty
	info, err := os.Stat(tempFile)
	require.NoError(t, err, "Expected no error while checking file stats")
	assert.False(t, info.IsDir(), "Expected file to not be a directory")
	assert.Positive(t, info.Size(), "Expected file size to be greater than 0")
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
	require.NoError(t, err, "Expected no error while marshaling OpenAPI schema to JSON")
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

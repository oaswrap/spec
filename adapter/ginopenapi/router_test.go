package ginopenapi_test

import (
	"flag"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/oaswrap/spec/adapter/ginopenapi"
	"github.com/oaswrap/spec/openapi"
	"github.com/oaswrap/spec/option"
	"github.com/oaswrap/spec/pkg/dto"
	"github.com/oaswrap/spec/pkg/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var update = flag.Bool("update", false, "update golden files")

func TestRouter_Spec(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name      string
		golden    string
		opts      []option.OpenAPIOption
		setup     func(r ginopenapi.Router)
		shouldErr bool
	}{
		{
			name:   "Pet Store API",
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
			setup: func(r ginopenapi.Router) {
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
					option.Response(200, new(dto.ApiResponse)),
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
						Username string `path:"username" required:"true"`
						dto.PetUser
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
		{
			name: "Invalid Open API Version",
			opts: []option.OpenAPIOption{
				option.WithOpenAPIVersion("2.0.0"), // Invalid version for this test
			},
			shouldErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := gin.Default()
			opts := []option.OpenAPIOption{
				option.WithOpenAPIVersion("3.0.3"),
				option.WithTitle("Test API " + tt.name),
				option.WithVersion("1.0.0"),
				option.WithDescription("This is a test API for " + tt.name),
				option.WithReflectorConfig(
					option.RequiredPropByValidateTag(),
					option.StripDefNamePrefix("GinopenapiTest"),
				),
			}
			if len(tt.opts) > 0 {
				opts = append(opts, tt.opts...)
			}
			r := ginopenapi.NewRouter(app, opts...)

			if tt.setup != nil {
				tt.setup(r)
			}

			if tt.shouldErr {
				err := r.Validate()
				assert.Error(t, err, "expected error for invalid OpenAPI configuration")
				return
			}
			err := r.Validate()
			assert.NoError(t, err, "failed to validate OpenAPI configuration")

			// Test the OpenAPI schema generation
			schema, err := r.GenerateSchema()

			require.NoError(t, err, "failed to generate OpenAPI schema")
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

type SingleRouteFunc func(path string, handlers ...gin.HandlerFunc) ginopenapi.Route

func PingHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "pong"})
}

func TestRouter_Single(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		method     string
		path       string
		methodFunc func(r ginopenapi.Router) SingleRouteFunc
	}{
		{http.MethodGet, "/ping", func(r ginopenapi.Router) SingleRouteFunc { return r.GET }},
		{http.MethodPost, "/ping", func(r ginopenapi.Router) SingleRouteFunc { return r.POST }},
		{http.MethodPut, "/ping", func(r ginopenapi.Router) SingleRouteFunc { return r.PUT }},
		{http.MethodDelete, "/ping", func(r ginopenapi.Router) SingleRouteFunc { return r.DELETE }},
		{http.MethodPatch, "/ping", func(r ginopenapi.Router) SingleRouteFunc { return r.PATCH }},
		{http.MethodHead, "/ping", func(r ginopenapi.Router) SingleRouteFunc { return r.HEAD }},
		{http.MethodOptions, "/ping", func(r ginopenapi.Router) SingleRouteFunc { return r.OPTIONS }},
	}
	for _, tt := range tests {
		t.Run(tt.method, func(t *testing.T) {
			app := gin.New()
			r := ginopenapi.NewRouter(app)

			route := tt.methodFunc(r)(tt.path, PingHandler).With(
				option.OperationID("ping"),
				option.Summary("Ping the server"),
				option.Description("Returns a simple pong response"),
				option.Response(200, new(struct {
					Message string `json:"message" example:"pong"`
				})),
			)

			assert.NotNil(t, route, "expected route to be created")

			req, _ := http.NewRequest(tt.method, tt.path, nil)
			rec := httptest.NewRecorder()
			app.ServeHTTP(rec, req)

			assert.Equal(t, http.StatusOK, rec.Code, "expected status code 200 for %s", tt.method)
			assert.Equal(t, `{"message":"pong"}`, rec.Body.String(), "expected response body to be 'pong' for %s", tt.method)

			schema, err := r.GenerateSchema()
			require.NoError(t, err, "failed to generate OpenAPI schema for %s", tt.method)
			assert.Contains(t, string(schema), "operationId: ping", "expected operationId in schema for %s", tt.method)
		})
	}
	t.Run("Static", func(t *testing.T) {
		// Create temp dir
		tmpDir := t.TempDir()

		// Create test file
		fileName := "hello.txt"
		fileContent := []byte("Hello, static!")
		err := os.WriteFile(filepath.Join(tmpDir, fileName), fileContent, 0644)
		assert.NoError(t, err)

		// Setup Gin
		gin.SetMode(gin.TestMode)
		g := gin.New()
		r := ginopenapi.NewRouter(g)
		r.Static("/static", tmpDir)

		// Create test server
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/static/"+fileName, nil)
		g.ServeHTTP(w, req)

		// Assertions
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, string(fileContent), w.Body.String())
	})

	t.Run("StaticFS", func(t *testing.T) {
		// Create temp dir
		tmpDir := t.TempDir()

		// Create a subfolder or file
		fileName := "test.txt"
		content := []byte("Hello from StaticFS!")
		err := os.WriteFile(filepath.Join(tmpDir, fileName), content, 0644)
		assert.NoError(t, err)

		// Setup Gin
		gin.SetMode(gin.TestMode)
		g := gin.New()
		r := ginopenapi.NewRouter(g)

		// Serve the temp dir using StaticFS
		r.StaticFS("/assets", http.Dir(tmpDir))

		// Make request
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/assets/"+fileName, nil)
		g.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, string(content), w.Body.String())
	})
	t.Run("StaticFile", func(t *testing.T) {
		// Create temp file
		tmpFile, err := os.CreateTemp("", "static-file-*.txt")
		require.NoError(t, err)
		defer func() {
			_ = os.Remove(tmpFile.Name())
		}()

		// Write content to temp file
		fileContent := []byte("Hello, static file!")
		_, err = tmpFile.Write(fileContent)
		require.NoError(t, err)

		// Setup Gin
		gin.SetMode(gin.TestMode)
		g := gin.New()
		r := ginopenapi.NewRouter(g)
		r.StaticFile("/static-file", tmpFile.Name())

		// Create test server
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/static-file", nil)
		g.ServeHTTP(w, req)

		// Assertions
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, string(fileContent), w.Body.String())
	})
	t.Run("StaticFileFS", func(t *testing.T) {
		// Setup temp dir and file
		tmpDir := t.TempDir()
		fileName := "foo.txt"
		content := []byte("This is served by StaticFileFS!")

		err := os.WriteFile(filepath.Join(tmpDir, fileName), content, 0644)
		assert.NoError(t, err)

		g := gin.New()
		r := ginopenapi.NewRouter(g)

		// Serve the single file at /myfile
		r.StaticFileFS("/myfile", fileName, http.Dir(tmpDir))

		// Make request
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/myfile", nil)
		g.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, string(content), w.Body.String())
	})
}

func TestRouter_Group(t *testing.T) {
	gin.SetMode(gin.TestMode)

	app := gin.New()
	r := ginopenapi.NewRouter(app)

	// Create a group with a prefix and middleware
	group := r.Group("/api", func(c *gin.Context) {
		c.Next()
	})

	// Add a route to the group
	group.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	}).With(
		option.OperationID("pingHandler"),
		option.Summary("Ping the server"),
		option.Description("Returns a simple pong response"),
		option.Response(200, new(struct {
			Message string `json:"message" example:"pong"`
		})),
	)

	assert.NotNil(t, group, "expected group to be created")

	req, _ := http.NewRequest("GET", "/api/ping", nil)
	rec := httptest.NewRecorder()
	app.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code, "expected status code 200 for /api/ping")
	assert.Equal(t, `{"message":"pong"}`, rec.Body.String(), "expected response body to be 'pong' for /api/ping")

	schema, err := r.GenerateSchema()
	require.NoError(t, err, "failed to generate OpenAPI schema for /api/ping")
	assert.Contains(t, string(schema), "operationId: pingHandler", "expected operationId in schema for /api/ping")
}

func TestRouter_Middleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Use", func(t *testing.T) {
		called := false
		middleware := func(c *gin.Context) {
			called = true
			c.Next()
		}
		app := gin.New()
		r := ginopenapi.NewRouter(app)
		r.Use(middleware)
		r.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "test"})
		}).With(
			option.OperationID("testHandler"),
		)

		req := httptest.NewRequest("GET", "/test", nil)
		rec := httptest.NewRecorder()
		app.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code, "expected status code 200 for /test")
		assert.Equal(t, `{"message":"test"}`, rec.Body.String(), "expected response body to be 'test'")
		assert.True(t, called, "expected middleware to be called")
	})
}

func TestGenerator_Docs(t *testing.T) {
	gin.SetMode(gin.TestMode)

	app := gin.New()
	r := ginopenapi.NewRouter(app)

	// Register a simple route
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	}).With(
		option.OperationID("pingHandler"),
		option.Summary("Ping the server"),
		option.Description("Returns a simple pong response"),
		option.Response(200, new(struct {
			Message string `json:"message" example:"pong"`
		})),
	)

	// Validate the router
	err := r.Validate()
	assert.NoError(t, err, "expected no error when validating router")

	// Test the OpenAPI documentation endpoint
	t.Run("should serve docs", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/docs", nil)
		rec := httptest.NewRecorder()
		app.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code, "expected status code 200 for /docs")
		assert.Contains(t, rec.Body.String(), "Gin OpenAPI", "expected API documentation in response body")
	})
	t.Run("should serve OpenAPI YAML", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/docs/openapi.yaml", nil)
		rec := httptest.NewRecorder()
		app.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code, "expected status code 200 for /docs/openapi.yaml")
		assert.NotEmpty(t, rec.Body.String(), "expected non-empty OpenAPI YAML response")
		assert.Contains(t, rec.Header().Get("Content-Type"), "application/x-yaml", "expected Content-Type to be application/x-yaml")
		assert.Contains(t, rec.Body.String(), "openapi: 3.0.3", "expected OpenAPI version in response body")
	})
}

func TestGenerator_DisableDocs(t *testing.T) {
	gin.SetMode(gin.TestMode)

	app := gin.New()
	r := ginopenapi.NewGenerator(app, option.WithDisableDocs(true))

	// Register a simple route
	r.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "test"})
	})

	t.Run("should not register docs routes", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/docs", nil)
		rec := httptest.NewRecorder()
		app.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusNotFound, rec.Code, "expected status code 404 for /docs when OpenAPI is disabled")
	})
	t.Run("should not register OpenAPI YAML route", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/docs/openapi.yaml", nil)
		rec := httptest.NewRecorder()
		app.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusNotFound, rec.Code, "expected status code 404 for /docs/openapi.yaml when OpenAPI is disabled")
	})
}

func TestGenerator_WriteSchemaTo(t *testing.T) {
	app := gin.Default()
	r := ginopenapi.NewRouter(app)

	// Register a simple route
	r.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "test"})
	})

	err := r.Validate()
	require.NoError(t, err, "failed to validate OpenAPI configuration")

	tempFile, err := os.CreateTemp("", "openapi-schema-*.yaml")
	require.NoError(t, err, "failed to create temporary file for OpenAPI schema")
	defer func() {
		err := os.Remove(tempFile.Name())
		require.NoError(t, err, "failed to remove temporary file")
	}()

	err = r.WriteSchemaTo(tempFile.Name())
	require.NoError(t, err, "failed to write OpenAPI schema to file")

	schema, err := os.ReadFile(tempFile.Name())
	require.NoError(t, err, "failed to read OpenAPI schema from file")
	assert.NotEmpty(t, schema, "expected non-empty OpenAPI schema")
}

func TestGenerator_MarshalYAML(t *testing.T) {
	gin.SetMode(gin.TestMode)

	app := gin.New()
	r := ginopenapi.NewRouter(app)

	// Register a simple route
	r.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "test"})
	})

	err := r.Validate()
	require.NoError(t, err, "failed to validate OpenAPI configuration")

	schema, err := r.MarshalYAML()
	require.NoError(t, err, "failed to marshal OpenAPI schema to YAML")
	assert.NotEmpty(t, schema, "expected non-empty OpenAPI schema in YAML format")
	assert.Contains(t, string(schema), "openapi:", "expected OpenAPI schema to contain 'openapi:' field")
}

func TestGenerator_MarshalJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)

	app := gin.New()
	r := ginopenapi.NewRouter(app)

	// Register a simple route
	r.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "test"})
	})

	err := r.Validate()
	require.NoError(t, err, "failed to validate OpenAPI configuration")

	schema, err := r.MarshalJSON()
	require.NoError(t, err, "failed to marshal OpenAPI schema to JSON")
	assert.NotEmpty(t, schema, "expected non-empty OpenAPI schema in JSON format")
	assert.Contains(t, string(schema), `"openapi":`, "expected OpenAPI schema to contain 'openapi' field")
}

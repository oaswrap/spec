package ginopenapi_test

import (
	"flag"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/oaswrap/spec/adapters/ginopenapi"
	"github.com/oaswrap/spec/option"
	"github.com/oaswrap/spec/pkg/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var update = flag.Bool("update", false, "update golden files")

type Pet struct {
	ID       int64    `json:"id" example:"1"`
	Name     string   `json:"name" binding:"required" example:"doggie"`
	Category Category `json:"category"`
	PhotoURL []string `json:"photoUrls" binding:"required"`
	Tags     []Tag    `json:"tags"`
	Status   string   `json:"status" enums:"available,pending,sold" example:"available"`
}

type Category struct {
	ID   int64  `json:"id" example:"1"`
	Name string `json:"name" example:"Dogs"`
}

type Tag struct {
	ID   int64  `json:"id" example:"1"`
	Name string `json:"name" example:"friendly"`
}

type CreatePetRequest struct {
	Name     string   `json:"name" binding:"required"`
	Category Category `json:"category"`
	PhotoURL []string `json:"photoUrls" binding:"required"`
	Tags     []Tag    `json:"tags"`
	Status   string   `json:"status" enums:"available,pending,sold"`
}

type UpdatePetRequest struct {
	ID       int64    `json:"id" binding:"required"`
	Name     string   `json:"name"`
	Category Category `json:"category"`
	PhotoURL []string `json:"photoUrls"`
	Tags     []Tag    `json:"tags"`
	Status   string   `json:"status" enums:"available,pending,sold"`
}

type UpdatePetFormData struct {
	ID     int64  `path:"petId" binding:"required"`
	Name   string `form:"name" binding:"required"`
	Status string `form:"status" binding:"required" enums:"available,pending,sold"`
}

type FindPetsByStatusRequest struct {
	Status string `form:"status" binding:"required" enums:"available,pending,sold"`
}

type FindPetsByTagsRequest struct {
	Tags []string `form:"tags" binding:"required"`
}

type User struct {
	ID         int64  `json:"id" example:"1"`
	Username   string `json:"username" example:"john_doe"`
	FirstName  string `json:"firstName" example:"John"`
	LastName   string `json:"lastName" example:"Doe"`
	Email      string `json:"email" example:"john_doe@example.com"`
	Password   string `json:"password" example:"password123"`
	Phone      string `json:"phone" example:"123-456-7890"`
	UserStatus int64  `json:"userStatus" example:"1"`
}

type ApiResponse struct {
	Code    int64  `json:"code" example:"200"`
	Type    string `json:"type" example:"success"`
	Message string `json:"message" example:"Pet created successfully"`
}

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
			setup: func(r ginopenapi.Router) {
				pets := r.Group("/pets").With(option.GroupTags("Pets"))
				dummyHandler := func(c *gin.Context) {
					c.JSON(200, gin.H{"message": "dummy handler"})
				}

				pets.GET("", dummyHandler).With(
					option.Summary("Get all pets"),
					option.Description("Returns a list of all pets in the store"),
					option.Response(200, new([]Pet)),
				)
				pets.POST("", dummyHandler).With(
					option.Summary("Create a new pet"),
					option.Description("Creates a new pet in the store"),
					option.Request(new(CreatePetRequest)),
					option.Response(201, new(Pet)),
				)
				pets.PUT("", dummyHandler).With(
					option.Summary("Update an existing pet"),
					option.Description("Updates an existing pet in the store"),
					option.Request(new(UpdatePetRequest)),
					option.Response(200, new(Pet)),
				)
				pets.GET("/findByStatus", dummyHandler).With(
					option.Summary("Find pets by status"),
					option.Description("Returns a list of pets based on their status"),
					option.Request(new(FindPetsByStatusRequest)),
					option.Response(200, new([]Pet)),
				)
				pets.GET("/findByTags", dummyHandler).With(
					option.Summary("Find pets by tags"),
					option.Description("Returns a list of pets based on their tags"),
					option.Request(new(FindPetsByTagsRequest)),
					option.Response(200, new([]Pet)),
				)
				pets.GET("/:petId", dummyHandler).With(
					option.Summary("Get pet by ID"),
					option.Description("Returns a single pet by its ID"),
					option.Request(new(struct {
						ID int64 `path:"petId"`
					})),
					option.Response(200, new(Pet)),
				)
				pets.POST("/:petId", dummyHandler).With(
					option.Summary("Update pet by form data"),
					option.Description("Updates a pet using form data"),
					option.Request(new(UpdatePetFormData)),
					option.Response(200, new(Pet)),
				)
				pets.DELETE("/:petId", dummyHandler).With(
					option.Summary("Delete a pet"),
					option.Description("Deletes a pet from the store"),
					option.Request(new(struct {
						ID int64 `path:"petId"`
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

	versions := map[string]string{
		"3.0.3": "3",
		"3.1.0": "31",
	}

	for _, tt := range tests {
		for version, suffix := range versions {
			app := gin.Default()
			opts := []option.OpenAPIOption{
				option.WithOpenAPIVersion(version),
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
				continue
			}
			err := r.Validate()
			assert.NoError(t, err, "failed to validate OpenAPI configuration")

			// Test the OpenAPI schema generation
			schema, err := r.GenerateSchema()

			require.NoError(t, err, "failed to generate OpenAPI schema")
			golden := filepath.Join("testdata", tt.golden+"_"+suffix+".yaml")

			if *update {
				err = r.WriteSchemaTo(golden)
				require.NoError(t, err, "failed to write golden file")
				t.Logf("Updated golden file: %s", golden)
			}

			want, err := os.ReadFile(golden)
			require.NoError(t, err, "failed to read golden file %s", golden)

			testutil.EqualYAML(t, want, schema)
		}
	}
}

func TestRouter_Gin(t *testing.T) {
	gin.SetMode(gin.TestMode)

	totalCalled := 0
	middleware := func(c *gin.Context) {
		totalCalled++
		c.Next()
	}
	pingHandler := func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	}
	app := gin.Default()
	r := ginopenapi.NewRouter(app)
	r.Use(middleware)
	opts := []option.OperationOption{
		option.OperationID("ping"),
		option.Summary("Ping the server"),
		option.Description("Returns a simple pong response"),
		option.Response(200, new(struct {
			Message string `json:"message" example:"pong"`
		})),
	}
	r.GET("/ping", pingHandler).With(opts...)
	r.POST("/ping", pingHandler).With(opts...)
	r.PUT("/ping", pingHandler).With(opts...)
	r.DELETE("/ping", pingHandler).With(opts...)
	r.PATCH("/ping", pingHandler).With(opts...)
	r.HEAD("/ping", pingHandler).With(opts...)
	r.OPTIONS("/ping", pingHandler).With(opts...)

	methods := []string{"GET", "POST", "PUT", "DELETE", "PATCH", "HEAD", "OPTIONS"}
	for i, method := range methods {
		t.Run(method, func(t *testing.T) {
			req, _ := http.NewRequest(method, "/ping", nil)
			rec := httptest.NewRecorder()
			app.ServeHTTP(rec, req)
			assert.Equal(t, http.StatusOK, rec.Code, "expected status code 200 for %s", method)
			assert.Equal(t, `{"message":"pong"}`, rec.Body.String(), "expected response body to be 'pong' for %s", method)
			assert.Equal(t, i+1, totalCalled, "middleware should be called exactly once for %s", method)
		})
	}

	t.Run("must register docs path", func(t *testing.T) {
		docsPath := "/docs"
		req, _ := http.NewRequest("GET", docsPath, nil)
		rec := httptest.NewRecorder()
		app.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code, "expected status code 200 for docs path")
	})
	t.Run("must register OpenAPI YAML path", func(t *testing.T) {
		openAPIPath := "/docs/openapi.yaml"
		req, _ := http.NewRequest("GET", openAPIPath, nil)
		rec := httptest.NewRecorder()
		app.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code, "expected status code 200 for OpenAPI YAML path")
		assert.Contains(t, rec.Header().Get("Content-Type"), "application/yaml", "expected content type to be application/yaml")
		assert.NotEmpty(t, rec.Body.String(), "expected non-empty response body for OpenAPI YAML path")
	})
	t.Run("static directory serving", func(t *testing.T) {
		// Create temp dir
		tmpDir := t.TempDir()

		// Create test file
		fileName := "hello.txt"
		fileContent := []byte("Hello, static!")
		err := os.WriteFile(filepath.Join(tmpDir, fileName), fileContent, 0644)
		assert.NoError(t, err)

		// Setup Gin
		gin.SetMode(gin.TestMode)
		g := gin.Default()
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
	t.Run("static file serving", func(t *testing.T) {
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
		g := gin.Default()
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
	t.Run("static fs serving", func(t *testing.T) {
		// Create temp dir
		tmpDir := t.TempDir()

		// Create a subfolder or file
		fileName := "test.txt"
		content := []byte("Hello from StaticFS!")
		err := os.WriteFile(filepath.Join(tmpDir, fileName), content, 0644)
		assert.NoError(t, err)

		// Setup Gin
		gin.SetMode(gin.TestMode)
		g := gin.Default()
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
	t.Run("static file fs serving", func(t *testing.T) {
		// Setup temp dir and file
		tmpDir := t.TempDir()
		fileName := "foo.txt"
		content := []byte("This is served by StaticFileFS!")

		err := os.WriteFile(filepath.Join(tmpDir, fileName), content, 0644)
		assert.NoError(t, err)

		g := gin.Default()
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

func TestGenerator_DisableDocs(t *testing.T) {
	gin.SetMode(gin.TestMode)

	app := gin.Default()
	r := ginopenapi.NewGenerator(app, option.WithDisableDocs(true))

	// Register a simple route
	r.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "test"})
	})

	// Validate the router
	err := r.Validate()
	assert.NoError(t, err, "expected no error when validating router with OpenAPI disabled")

	// Test the registered route
	req, _ := http.NewRequest("GET", "/test", nil)
	rec := httptest.NewRecorder()
	app.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code, "expected status code 200 for /test")
	assert.Equal(t, `{"message":"test"}`, rec.Body.String(), "expected response body to be 'test'")

	// Ensure OpenAPI paths are not registered
	req, _ = http.NewRequest("GET", "/docs", nil)
	rec = httptest.NewRecorder()
	app.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusNotFound, rec.Code, "expected status code 404 for /docs when OpenAPI is disabled")

	// Ensure OpenAPI YAML path is not registered
	req, _ = http.NewRequest("GET", "/docs/openapi.yaml", nil)
	rec = httptest.NewRecorder()
	app.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusNotFound, rec.Code, "expected status code 404 for /docs/openapi.yaml when OpenAPI is disabled")
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

	app := gin.Default()
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

	app := gin.Default()
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

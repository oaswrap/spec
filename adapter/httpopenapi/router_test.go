package httpopenapi_test

import (
	"encoding/json"
	"flag"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/oaswrap/spec/adapter/httpopenapi"
	"github.com/oaswrap/spec/openapi"
	"github.com/oaswrap/spec/option"
	"github.com/oaswrap/spec/pkg/dto"
	"github.com/oaswrap/spec/pkg/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var update = flag.Bool("update", false, "update golden files")

func TestRouter_Spec(t *testing.T) {
	tests := []struct {
		name      string
		golden    string
		opts      []option.OpenAPIOption
		setup     func(r httpopenapi.Router)
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
			setup: func(r httpopenapi.Router) {
				pet := r.Group("pet").With(
					option.GroupTags("pet"),
					option.GroupSecurity("petstore_auth", "write:pets", "read:pets"),
				)
				pet.HandleFunc("PUT /", nil).With(
					option.OperationID("updatePet"),
					option.Summary("Update an existing pet"),
					option.Description("Update the details of an existing pet in the store."),
					option.Request(new(dto.Pet)),
					option.Response(200, new(dto.Pet)),
				)
				pet.HandleFunc("POST /", nil).With(
					option.OperationID("addPet"),
					option.Summary("Add a new pet"),
					option.Description("Add a new pet to the store."),
					option.Request(new(dto.Pet)),
					option.Response(201, new(dto.Pet)),
				)
				pet.HandleFunc("GET /findByStatus", nil).With(
					option.OperationID("findPetsByStatus"),
					option.Summary("Find pets by status"),
					option.Description("Finds Pets by status. Multiple status values can be provided with comma separated strings."),
					option.Request(new(struct {
						Status string `query:"status" enum:"available,pending,sold"`
					})),
					option.Response(200, new([]dto.Pet)),
				)
				pet.HandleFunc("GET /findByTags", nil).With(
					option.OperationID("findPetsByTags"),
					option.Summary("Find pets by tags"),
					option.Description("Finds Pets by tags. Multiple tags can be provided with comma separated strings."),
					option.Request(new(struct {
						Tags []string `query:"tags"`
					})),
					option.Response(200, new([]dto.Pet)),
				)
				pet.HandleFunc("POST /{petId}/uploadImage", nil).With(
					option.OperationID("uploadFile"),
					option.Summary("Upload an image for a pet"),
					option.Description("Uploads an image for a pet."),
					option.Request(new(dto.UploadImageRequest)),
					option.Response(200, new(dto.ApiResponse)),
				)
				pet.HandleFunc("GET /{petId}", nil).With(
					option.OperationID("getPetById"),
					option.Summary("Get pet by ID"),
					option.Description("Retrieve a pet by its ID."),
					option.Request(new(struct {
						ID int `path:"petId" required:"true"`
					})),
					option.Response(200, new(dto.Pet)),
				)
				pet.HandleFunc("POST /{petId}", nil).With(
					option.OperationID("updatePetWithForm"),
					option.Summary("Update pet with form"),
					option.Description("Updates a pet in the store with form data."),
					option.Request(new(dto.UpdatePetWithFormRequest)),
					option.Response(200, nil),
				)
				pet.HandleFunc("DELETE /{petId}", nil).With(
					option.OperationID("deletePet"),
					option.Summary("Delete a pet"),
					option.Description("Delete a pet from the store by its ID."),
					option.Request(new(dto.DeletePetRequest)),
					option.Response(204, nil),
				)

				store := r.Group("store").With(
					option.GroupTags("store"),
				)
				store.HandleFunc("POST /order", nil).With(
					option.OperationID("placeOrder"),
					option.Summary("Place an order"),
					option.Description("Place a new order for a pet."),
					option.Request(new(dto.Order)),
					option.Response(201, new(dto.Order)),
				)
				store.HandleFunc("GET /order/{orderId}", nil).With(
					option.OperationID("getOrderById"),
					option.Summary("Get order by ID"),
					option.Description("Retrieve an order by its ID."),
					option.Request(new(struct {
						ID int `path:"orderId" required:"true"`
					})),
					option.Response(200, new(dto.Order)),
					option.Response(404, nil),
				)
				store.HandleFunc("DELETE /order/{orderId}", nil).With(
					option.OperationID("deleteOrder"),
					option.Summary("Delete an order"),
					option.Description("Delete an order by its ID."),
					option.Request(new(struct {
						ID int `path:"orderId" required:"true"`
					})),
					option.Response(204, nil),
				)
				user := r.Group("user").With(
					option.GroupTags("user"),
				)
				user.HandleFunc("POST /createWithList", nil).With(
					option.OperationID("createUsersWithList"),
					option.Summary("Create users with list"),
					option.Description("Create multiple users in the store with a list."),
					option.Request(new([]dto.PetUser)),
					option.Response(201, nil),
				)
				user.HandleFunc("POST /", nil).With(
					option.OperationID("createUser"),
					option.Summary("Create a new user"),
					option.Description("Create a new user in the store."),
					option.Request(new(dto.PetUser)),
					option.Response(201, new(dto.PetUser)),
				)
				user.HandleFunc("GET /{username}", nil).With(
					option.OperationID("getUserByName"),
					option.Summary("Get user by username"),
					option.Description("Retrieve a user by their username."),
					option.Request(new(struct {
						Username string `path:"username" required:"true"`
					})),
					option.Response(200, new(dto.PetUser)),
					option.Response(404, nil),
				)
				user.HandleFunc("PUT /{username}", nil).With(
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
				user.HandleFunc("DELETE /{username}", nil).With(
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
			app := http.NewServeMux()
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
			r := httpopenapi.NewRouter(app, opts...)

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

func pingHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"message": "pong"})
}

func TestRouter_Handle(t *testing.T) {
	t.Run("GET", func(t *testing.T) {
		mux := http.NewServeMux()
		r := httpopenapi.NewRouter(mux)
		r.Handle("GET /ping", http.HandlerFunc(pingHandler)).With(
			option.OperationID("pingHandler"),
		)

		req := httptest.NewRequest("GET", "/ping", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), "pong")

		schema, err := r.GenerateSchema()
		require.NoError(t, err, "failed to generate OpenAPI schema")
		assert.Contains(t, string(schema), "operationId: pingHandler", "expected operationId in schema")
	})
	t.Run("CONNECT", func(t *testing.T) {
		mux := http.NewServeMux()
		r := httpopenapi.NewRouter(mux)
		r.Handle("CONNECT /ping", http.HandlerFunc(pingHandler)).With(
			option.OperationID("pingHandler"),
		)

		req := httptest.NewRequest("CONNECT", "/ping", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), "pong")

		schema, err := r.GenerateSchema()
		require.NoError(t, err, "failed to generate OpenAPI schema")
		assert.NotContains(t, string(schema), "operationId: pingHandler", "expected operationId in schema")
	})
}

func TestRouter_HandleFunc(t *testing.T) {
	t.Run("GET", func(t *testing.T) {
		mux := http.NewServeMux()
		r := httpopenapi.NewRouter(mux)
		r.HandleFunc("GET /ping", pingHandler).With(
			option.OperationID("pingHandler"),
		)

		req := httptest.NewRequest("GET", "/ping", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), "pong")

		schema, err := r.GenerateSchema()
		require.NoError(t, err, "failed to generate OpenAPI schema")
		assert.Contains(t, string(schema), "operationId: pingHandler", "expected operationId in schema")
	})
	t.Run("CONNECT", func(t *testing.T) {
		mux := http.NewServeMux()
		r := httpopenapi.NewRouter(mux)
		r.HandleFunc("CONNECT /ping", pingHandler).With(
			option.OperationID("pingHandler"),
		)

		req := httptest.NewRequest("CONNECT", "/ping", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), "pong")

		schema, err := r.GenerateSchema()
		require.NoError(t, err, "failed to generate OpenAPI schema")
		assert.NotContains(t, string(schema), "operationId: pingHandler", "expected operationId in schema")
	})
}

func TestRouter_Group(t *testing.T) {
	t.Run("Group", func(t *testing.T) {
		mux := http.NewServeMux()
		r := httpopenapi.NewRouter(mux)

		group := r.Group("/api").With(
			option.GroupTags("api"),
		)
		group.HandleFunc("GET /ping", pingHandler).With(
			option.OperationID("pingHandler"),
		)

		req := httptest.NewRequest("GET", "/api/ping", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), "pong")

		schema, err := r.GenerateSchema()
		require.NoError(t, err, "failed to generate OpenAPI schema")
		assert.Contains(t, string(schema), "operationId: pingHandler", "expected operationId in schema")
	})
	t.Run("Group with Middleware", func(t *testing.T) {
		mux := http.NewServeMux()
		r := httpopenapi.NewRouter(mux)
		called := false
		middleware := func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				called = true
				next.ServeHTTP(w, r)
			})
		}
		group := r.Group("/api", middleware).With(
			option.GroupTags("api"),
		)
		group.HandleFunc("GET /ping", pingHandler).With(
			option.OperationID("pingHandler"),
		)

		req := httptest.NewRequest("GET", "/api/ping", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		assert.True(t, called, "middleware should be called")
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), "pong")

		schema, err := r.GenerateSchema()
		require.NoError(t, err, "failed to generate OpenAPI schema")
		assert.Contains(t, string(schema), "operationId: pingHandler", "expected operationId in schema")
	})
	t.Run("Route", func(t *testing.T) {
		mux := http.NewServeMux()
		r := httpopenapi.NewRouter(mux)

		r.Route("/api", func(r httpopenapi.Router) {
			r.HandleFunc("GET /ping", pingHandler).With(
				option.OperationID("pingHandler"),
			)
		})

		req := httptest.NewRequest("GET", "/api/ping", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), "pong")

		schema, err := r.GenerateSchema()
		require.NoError(t, err, "failed to generate OpenAPI schema")
		assert.Contains(t, string(schema), "operationId: pingHandler", "expected operationId in schema")
	})
	t.Run("Route with Middleware", func(t *testing.T) {
		mux := http.NewServeMux()
		r := httpopenapi.NewRouter(mux)
		called := false
		middleware := func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				called = true
				next.ServeHTTP(w, r)
			})
		}
		r.Route("/api", func(r httpopenapi.Router) {
			r.HandleFunc("GET /ping", pingHandler).With(
				option.OperationID("pingHandler"),
			)
		}, middleware)

		req := httptest.NewRequest("GET", "/api/ping", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		assert.True(t, called, "middleware should be called")
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), "pong")

		schema, err := r.GenerateSchema()
		require.NoError(t, err, "failed to generate OpenAPI schema")
		assert.Contains(t, string(schema), "operationId: pingHandler", "expected operationId in schema")
	})
}

func TestGenerator_Docs(t *testing.T) {
	mux := http.NewServeMux()
	r := httpopenapi.NewRouter(mux)

	r.HandleFunc("GET /ping", pingHandler).With(
		option.OperationID("pingHandler"),
	)

	// Test the OpenAPI documentation endpoint
	t.Run("should serve docs", func(t *testing.T) {
		docsReq := httptest.NewRequest("GET", "/docs", nil)
		docsRec := httptest.NewRecorder()
		mux.ServeHTTP(docsRec, docsReq)

		assert.Equal(t, http.StatusOK, docsRec.Code)
		assert.Contains(t, docsRec.Body.String(), "HTTP OpenAPI")
	})
	t.Run("should serve docs file", func(t *testing.T) {
		docsFileReq := httptest.NewRequest("GET", "/docs/openapi.yaml", nil)
		docsFileRec := httptest.NewRecorder()
		mux.ServeHTTP(docsFileRec, docsFileReq)

		assert.Equal(t, http.StatusOK, docsFileRec.Code)
		assert.NotEmpty(t, docsFileRec.Body.String())
		assert.Contains(t, docsFileRec.Header().Get("Content-Type"), "application/x-yaml")
		assert.Contains(t, docsFileRec.Body.String(), "openapi: 3.0.3")
	})
}

func TestGenerator_DisableDocs(t *testing.T) {
	mux := http.NewServeMux()
	r := httpopenapi.NewRouter(mux, option.WithDisableDocs(true))

	r.HandleFunc("GET /ping", pingHandler).With(
		option.OperationID("pingHandler"),
	)

	// Test that docs are not served when disabled
	t.Run("should not serve docs when disabled", func(t *testing.T) {
		docsReq := httptest.NewRequest("GET", "/docs", nil)
		docsRec := httptest.NewRecorder()
		mux.ServeHTTP(docsRec, docsReq)

		assert.Equal(t, http.StatusNotFound, docsRec.Code, "expected 404 when docs are disabled")
	})
	t.Run("should not serve docs file when disabled", func(t *testing.T) {
		docsFileReq := httptest.NewRequest("GET", "/docs/openapi.yaml", nil)
		docsFileRec := httptest.NewRecorder()
		mux.ServeHTTP(docsFileRec, docsFileReq)

		assert.Equal(t, http.StatusNotFound, docsFileRec.Code, "expected 404 when docs are disabled")
	})
}

func TestGenerator_MarshalJSON(t *testing.T) {
	mux := http.NewServeMux()
	r := httpopenapi.NewRouter(mux)

	r.HandleFunc("GET /ping", pingHandler).With(
		option.OperationID("pingHandler"),
	)

	jsonData, err := r.MarshalJSON()
	require.NoError(t, err, "failed to marshal OpenAPI schema to JSON")

	assert.NotEmpty(t, jsonData, "expected non-empty JSON data")
	assert.Contains(t, string(jsonData), `"operationId": "pingHandler"`, "expected operationId in JSON schema")
}

func TestGenerator_MarshalYAML(t *testing.T) {
	mux := http.NewServeMux()
	r := httpopenapi.NewRouter(mux)

	r.HandleFunc("GET /ping", pingHandler).With(
		option.OperationID("pingHandler"),
	)

	yamlData, err := r.MarshalYAML()
	require.NoError(t, err, "failed to marshal OpenAPI schema to YAML")

	assert.NotEmpty(t, yamlData, "expected non-empty YAML data")
	assert.Contains(t, string(yamlData), "operationId: pingHandler", "expected operationId in YAML schema")
}

func TestGenerator_WriteSchemaTo(t *testing.T) {
	mux := http.NewServeMux()
	r := httpopenapi.NewRouter(mux)

	r.HandleFunc("GET /ping", pingHandler).With(
		option.OperationID("pingHandler"),
	)

	tempFile, err := os.CreateTemp("", "openapi-schema-*.yaml")
	require.NoError(t, err, "failed to create temporary file")
	defer func() {
		_ = os.Remove(tempFile.Name())
	}()

	err = r.WriteSchemaTo(tempFile.Name())
	require.NoError(t, err, "failed to write OpenAPI schema to file")

	schemaData, err := os.ReadFile(tempFile.Name())
	require.NoError(t, err, "failed to read OpenAPI schema from file")

	assert.NotEmpty(t, schemaData, "expected non-empty schema data")
	assert.Contains(t, string(schemaData), "operationId: pingHandler", "expected operationId in written schema")
}

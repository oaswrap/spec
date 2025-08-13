package httprouteropenapi_test

import (
	"encoding/json"
	"flag"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/julienschmidt/httprouter"
	"github.com/oaswrap/spec/adapter/httprouteropenapi"
	"github.com/oaswrap/spec/openapi"
	"github.com/oaswrap/spec/option"
	"github.com/oaswrap/spec/pkg/dto"
	"github.com/oaswrap/spec/pkg/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var update = flag.Bool("update", false, "update golden files")

func DummyHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"message": "Hello, World!"})
}

func TestRouter_Spec(t *testing.T) {
	tests := []struct {
		name      string
		golden    string
		opts      []option.OpenAPIOption
		setup     func(r httprouteropenapi.Router)
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
			setup: func(r httprouteropenapi.Router) {
				pet := r.Group("/pet").With(
					option.GroupTags("pet"),
					option.GroupSecurity("petstore_auth", "write:pets", "read:pets"),
				)
				pet.PUT("/", DummyHandler).With(
					option.OperationID("updatePet"),
					option.Summary("Update an existing pet"),
					option.Description("Update the details of an existing pet in the store."),
					option.Request(new(dto.Pet)),
					option.Response(200, new(dto.Pet)),
				)
				pet.POST("/", DummyHandler).With(
					option.OperationID("addPet"),
					option.Summary("Add a new pet"),
					option.Description("Add a new pet to the store."),
					option.Request(new(dto.Pet)),
					option.Response(201, new(dto.Pet)),
				)
				pet.GET("/findByStatus", DummyHandler).With(
					option.OperationID("findPetsByStatus"),
					option.Summary("Find pets by status"),
					option.Description("Finds Pets by status. Multiple status values can be provided with comma separated strings."),
					option.Request(new(struct {
						Status string `query:"status" enum:"available,pending,sold"`
					})),
					option.Response(200, new([]dto.Pet)),
				)
				pet.GET("/findByTags", DummyHandler).With(
					option.OperationID("findPetsByTags"),
					option.Summary("Find pets by tags"),
					option.Description("Finds Pets by tags. Multiple tags can be provided with comma separated strings."),
					option.Request(new(struct {
						Tags []string `query:"tags"`
					})),
					option.Response(200, new([]dto.Pet)),
				)
				pet.POST("/:petId/uploadImage", DummyHandler).With(
					option.OperationID("uploadFile"),
					option.Summary("Upload an image for a pet"),
					option.Description("Uploads an image for a pet."),
					option.Request(new(dto.UploadImageRequest)),
					option.Response(200, new(dto.ApiResponse)),
				)
				pet.GET("/id/:petId", DummyHandler).With(
					option.OperationID("getPetById"),
					option.Summary("Get pet by ID"),
					option.Description("Retrieve a pet by its ID."),
					option.Request(new(struct {
						ID int `path:"petId" required:"true"`
					})),
					option.Response(200, new(dto.Pet)),
				)
				pet.POST("/:petId", DummyHandler).With(
					option.OperationID("updatePetWithForm"),
					option.Summary("Update pet with form"),
					option.Description("Updates a pet in the store with form data."),
					option.Request(new(dto.UpdatePetWithFormRequest)),
					option.Response(200, nil),
				)
				pet.DELETE("/:petId", DummyHandler).With(
					option.OperationID("deletePet"),
					option.Summary("Delete a pet"),
					option.Description("Delete a pet from the store by its ID."),
					option.Request(new(dto.DeletePetRequest)),
					option.Response(204, nil),
				)
				store := r.Group("/store").With(
					option.GroupTags("store"),
				)
				store.POST("/order", DummyHandler).With(
					option.OperationID("placeOrder"),
					option.Summary("Place an order"),
					option.Description("Place a new order for a pet."),
					option.Request(new(dto.Order)),
					option.Response(201, new(dto.Order)),
				)
				store.GET("/order/:orderId", DummyHandler).With(
					option.OperationID("getOrderById"),
					option.Summary("Get order by ID"),
					option.Description("Retrieve an order by its ID."),
					option.Request(new(struct {
						ID int `path:"orderId" required:"true"`
					})),
					option.Response(200, new(dto.Order)),
					option.Response(404, nil),
				)
				store.DELETE("/order/:orderId", DummyHandler).With(
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
				user.POST("/createWithList", DummyHandler).With(
					option.OperationID("createUsersWithList"),
					option.Summary("Create users with list"),
					option.Description("Create multiple users in the store with a list."),
					option.Request(new([]dto.PetUser)),
					option.Response(201, nil),
				)
				user.POST("/", DummyHandler).With(
					option.OperationID("createUser"),
					option.Summary("Create a new user"),
					option.Description("Create a new user in the store."),
					option.Request(new(dto.PetUser)),
					option.Response(201, new(dto.PetUser)),
				)
				user.GET("/:username", DummyHandler).With(
					option.OperationID("getUserByName"),
					option.Summary("Get user by username"),
					option.Description("Retrieve a user by their username."),
					option.Request(new(struct {
						Username string `path:"username" required:"true"`
					})),
					option.Response(200, new(dto.PetUser)),
					option.Response(404, nil),
				)
				user.PUT("/:username", DummyHandler).With(
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
				user.DELETE("/:username", DummyHandler).With(
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
			router := httprouter.New()
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
			r := httprouteropenapi.NewRouter(router, opts...)

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

type SingleRouteFunc func(path string, handle httprouter.Handle) httprouteropenapi.Route

func PingHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"message": "pong"})
}

func TestRouter_SingleRoute(t *testing.T) {
	tests := []struct {
		method     string
		path       string
		methodFunc func(r httprouteropenapi.Router) SingleRouteFunc
	}{
		{
			method: http.MethodGet,
			path:   "/ping",
			methodFunc: func(r httprouteropenapi.Router) SingleRouteFunc {
				return r.GET
			},
		},
		{
			method: http.MethodPost,
			path:   "/ping",
			methodFunc: func(r httprouteropenapi.Router) SingleRouteFunc {
				return r.POST
			},
		},
		{
			method: http.MethodPut,
			path:   "/ping",
			methodFunc: func(r httprouteropenapi.Router) SingleRouteFunc {
				return r.PUT
			},
		},
		{
			method: http.MethodDelete,
			path:   "/ping",
			methodFunc: func(r httprouteropenapi.Router) SingleRouteFunc {
				return r.DELETE
			},
		},
		{
			method: http.MethodPatch,
			path:   "/ping",
			methodFunc: func(r httprouteropenapi.Router) SingleRouteFunc {
				return r.PATCH
			},
		},
		{
			method: http.MethodOptions,
			path:   "/ping",
			methodFunc: func(r httprouteropenapi.Router) SingleRouteFunc {
				return r.OPTIONS
			},
		},
		{
			method: http.MethodHead,
			path:   "/ping",
			methodFunc: func(r httprouteropenapi.Router) SingleRouteFunc {
				return r.HEAD
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.method, func(t *testing.T) {
			router := httprouter.New()
			r := httprouteropenapi.NewRouter(router)

			tt.methodFunc(r)(tt.path, PingHandler).With(
				option.Summary("Ping the server"),
				option.Description("Returns a simple pong response"),
				option.Response(200, new(struct {
					Message string `json:"message" example:"pong"`
				})),
			)

			// Test the route
			req := httptest.NewRequest(tt.method, tt.path, nil)
			rec := httptest.NewRecorder()
			router.ServeHTTP(rec, req)

			assert.Equal(t, http.StatusOK, rec.Code)
			assert.JSONEq(t, `{"message":"pong"}`, rec.Body.String(), "response body should match")

			schema, err := r.GenerateSchema()
			assert.NoError(t, err, "failed to generate OpenAPI schema")
			assert.NotEmpty(t, schema, "OpenAPI schema should not be empty")
			assert.Contains(t, string(schema), "summary: Ping the server", "OpenAPI schema should contain the summary")
		})
	}
}

func TestRouter_Group(t *testing.T) {
	logs := []string{}
	middleware1 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logs = append(logs, "middleware1")
			next.ServeHTTP(w, r)
			logs = append(logs, "middleware1")
		})
	}
	middleware2 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logs = append(logs, "middleware2")
			next.ServeHTTP(w, r)
			logs = append(logs, "middleware2")
		})
	}
	router := httprouter.New()
	r := httprouteropenapi.NewRouter(router)
	api := r.Group("/api/v1", middleware1, middleware2).With(
		option.GroupTags("apiv1"),
	)
	dummyHandler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "pong"})
	}
	api.GET("/ping", PingHandler).With(
		option.Summary("Ping the server"),
		option.Description("Returns a simple pong response"),
		option.Response(200, new(struct {
			Message string `json:"message" example:"pong"`
		})),
	)
	api.HandlerFunc(http.MethodGet, "/dummy", dummyHandler).With(
		option.Summary("Dummy endpoint"),
		option.Description("Returns a simple dummy response"),
		option.Response(200, new(struct {
			Message string `json:"message" example:"pong"`
		})),
	)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/ping", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, `{"message":"pong"}`, rec.Body.String(), "response body should match")

	assert.Equal(t, []string{
		"middleware1",
		"middleware2",
		"middleware2",
		"middleware1",
	}, logs, "middleware logs should match")

	schema, err := r.GenerateSchema()
	assert.NoError(t, err, "failed to generate OpenAPI schema")
	assert.NotEmpty(t, schema, "OpenAPI schema should not be empty")
	assert.Contains(t, string(schema), "apiv1", "OpenAPI schema should contain the group tags")

	req = httptest.NewRequest(http.MethodGet, "/api/v1/dummy", nil)
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, `{"message":"pong"}`, rec.Body.String())
}

func TestGenerator_Docs(t *testing.T) {
	router := httprouter.New()
	r := httprouteropenapi.NewRouter(router)

	r.GET("/ping", PingHandler).With(
		option.Summary("Ping the server"),
		option.Description("Returns a simple pong response"),
		option.Response(200, new(struct {
			Message string `json:"message" example:"pong"`
		})),
	)

	t.Run("should serve /docs", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/docs", nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.NotEmpty(t, rec.Body.String(), "response body should not be empty")
	})
	t.Run("should serve /docs/openapi.yaml", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/docs/openapi.yaml", nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.NotEmpty(t, rec.Body.String(), "response body should not be empty")
	})
}

func TestGenerator_DisableDocs(t *testing.T) {
	router := httprouter.New()
	r := httprouteropenapi.NewRouter(router,
		option.WithDisableDocs(),
	)

	r.GET("/ping", PingHandler).With(
		option.Summary("Ping the server"),
		option.Description("Returns a simple pong response"),
		option.Response(200, new(struct {
			Message string `json:"message" example:"pong"`
		})),
	)

	t.Run("should not serve /docs", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/docs", nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
	})
	t.Run("should not serve /docs/openapi.yaml", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/docs/openapi.yaml", nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
	})
}

func TestGenerator_MarshalYAML(t *testing.T) {
	router := httprouter.New()
	r := httprouteropenapi.NewRouter(router)

	r.GET("/ping", PingHandler).With(
		option.Summary("Ping the server"),
		option.Description("Returns a simple pong response"),
		option.Response(200, new(struct {
			Message string `json:"message" example:"pong"`
		})),
	)

	schema, err := r.MarshalYAML()
	assert.NoError(t, err, "failed to marshal OpenAPI schema to YAML")
	assert.NotEmpty(t, schema, "YAML schema should not be empty")
	assert.Contains(t, string(schema), "summary: Ping the server", "YAML schema should contain the summary")
}

func TestGenerator_MarshalJSON(t *testing.T) {
	router := httprouter.New()
	r := httprouteropenapi.NewRouter(router)

	r.GET("/ping", PingHandler).With(
		option.Summary("Ping the server"),
		option.Description("Returns a simple pong response"),
		option.Response(200, new(struct {
			Message string `json:"message" example:"pong"`
		})),
	)

	schema, err := r.MarshalJSON()
	assert.NoError(t, err, "failed to marshal OpenAPI schema to JSON")
	assert.NotEmpty(t, schema, "JSON schema should not be empty")
	assert.Contains(t, string(schema), `"summary": "Ping the server"`, "JSON schema should contain the summary")
}

func TestGenerator_WriteSchemaTo(t *testing.T) {
	router := httprouter.New()
	r := httprouteropenapi.NewRouter(router)

	r.GET("/ping", PingHandler).With(
		option.Summary("Ping the server"),
		option.Description("Returns a simple pong response"),
		option.Response(200, new(struct {
			Message string `json:"message" example:"pong"`
		})),
	)

	dir := t.TempDir()
	path := filepath.Join(dir, "openapi.yaml")
	err := r.WriteSchemaTo(path)
	assert.NoError(t, err, "failed to write OpenAPI schema to directory")
	assert.FileExists(t, path, "openapi.yaml should be created")
}

package chiopenapi_test

import (
	"encoding/json"
	"flag"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/oaswrap/spec/adapter/chiopenapi"
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
		setup     func(r chiopenapi.Router)
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
			setup: func(r chiopenapi.Router) {
				r.Route("/pet", func(r chiopenapi.Router) {
					r.Put("/", nil).With(
						option.OperationID("updatePet"),
						option.Summary("Update an existing pet"),
						option.Description("Update the details of an existing pet in the store."),
						option.Request(new(dto.Pet)),
						option.Response(200, new(dto.Pet)),
					)
					r.Post("/", nil).With(
						option.OperationID("addPet"),
						option.Summary("Add a new pet"),
						option.Description("Add a new pet to the store."),
						option.Request(new(dto.Pet)),
						option.Response(201, new(dto.Pet)),
					)
					r.Get("/findByStatus", nil).With(
						option.OperationID("findPetsByStatus"),
						option.Summary("Find pets by status"),
						option.Description("Finds Pets by status. Multiple status values can be provided with comma separated strings."),
						option.Request(new(struct {
							Status string `query:"status" enum:"available,pending,sold"`
						})),
						option.Response(200, new([]dto.Pet)),
					)
					r.Get("/findByTags", nil).With(
						option.OperationID("findPetsByTags"),
						option.Summary("Find pets by tags"),
						option.Description("Finds Pets by tags. Multiple tags can be provided with comma separated strings."),
						option.Request(new(struct {
							Tags []string `query:"tags"`
						})),
						option.Response(200, new([]dto.Pet)),
					)
					r.Post("/{petId}/uploadImage", nil).With(
						option.OperationID("uploadFile"),
						option.Summary("Upload an image for a pet"),
						option.Description("Uploads an image for a pet."),
						option.Request(new(dto.UploadImageRequest)),
						option.Response(200, new(dto.ApiResponse)),
					)
					r.Get("/{petId}", nil).With(
						option.OperationID("getPetById"),
						option.Summary("Get pet by ID"),
						option.Description("Retrieve a pet by its ID."),
						option.Request(new(struct {
							ID int `path:"petId" required:"true"`
						})),
						option.Response(200, new(dto.Pet)),
					)
					r.Post("/{petId}", nil).With(
						option.OperationID("updatePetWithForm"),
						option.Summary("Update pet with form"),
						option.Description("Updates a pet in the store with form data."),
						option.Request(new(dto.UpdatePetWithFormRequest)),
						option.Response(200, nil),
					)
					r.Delete("/{petId}", nil).With(
						option.OperationID("deletePet"),
						option.Summary("Delete a pet"),
						option.Description("Delete a pet from the store by its ID."),
						option.Request(new(dto.DeletePetRequest)),
						option.Response(204, nil),
					)
				}, option.GroupTags("pet"),
					option.GroupSecurity("petstore_auth", "write:pets", "read:pets"),
				)

				r.Route("/store", func(r chiopenapi.Router) {
					r.Post("/order", nil).With(
						option.OperationID("placeOrder"),
						option.Summary("Place an order"),
						option.Description("Place a new order for a pet."),
						option.Request(new(dto.Order)),
						option.Response(201, new(dto.Order)),
					)
					r.Get("/order/{orderId}", nil).With(
						option.OperationID("getOrderById"),
						option.Summary("Get order by ID"),
						option.Description("Retrieve an order by its ID."),
						option.Request(new(struct {
							ID int `path:"orderId" required:"true"`
						})),
						option.Response(200, new(dto.Order)),
						option.Response(404, nil),
					)
					r.Delete("/order/{orderId}", nil).With(
						option.OperationID("deleteOrder"),
						option.Summary("Delete an order"),
						option.Description("Delete an order by its ID."),
						option.Request(new(struct {
							ID int `path:"orderId" required:"true"`
						})),
						option.Response(204, nil),
					)
				}, option.GroupTags("store"))

				r.Route("/user", func(r chiopenapi.Router) {
					r.Post("/createWithList", nil).With(
						option.OperationID("createUsersWithList"),
						option.Summary("Create users with list"),
						option.Description("Create multiple users in the store with a list."),
						option.Request(new([]dto.PetUser)),
						option.Response(201, nil),
					)
					r.Post("/", nil).With(
						option.OperationID("createUser"),
						option.Summary("Create a new user"),
						option.Description("Create a new user in the store."),
						option.Request(new(dto.PetUser)),
						option.Response(201, new(dto.PetUser)),
					)
					r.Get("/{username}", nil).With(
						option.OperationID("getUserByName"),
						option.Summary("Get user by username"),
						option.Description("Retrieve a user by their username."),
						option.Request(new(struct {
							Username string `path:"username" required:"true"`
						})),
						option.Response(200, new(dto.PetUser)),
						option.Response(404, nil),
					)
					r.Put("/{username}", nil).With(
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
					r.Delete("/{username}", nil).With(
						option.OperationID("deleteUser"),
						option.Summary("Delete a user"),
						option.Description("Delete a user from the store by their username."),
						option.Request(new(struct {
							Username string `path:"username" required:"true"`
						})),
						option.Response(204, nil),
					)
				}, option.GroupTags("user"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := chi.NewRouter()
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
			r := chiopenapi.NewRouter(app, opts...)

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

type SingleRouteFunc func(path string, handler http.HandlerFunc) chiopenapi.Route

func TestRouter_Single(t *testing.T) {
	tests := []struct {
		method     string
		path       string
		methodFunc func(r chiopenapi.Router) SingleRouteFunc
	}{
		{"GET", "/ping", func(r chiopenapi.Router) SingleRouteFunc { return r.Get }},
		{"POST", "/ping", func(r chiopenapi.Router) SingleRouteFunc { return r.Post }},
		{"PUT", "/ping", func(r chiopenapi.Router) SingleRouteFunc { return r.Put }},
		{"DELETE", "/ping", func(r chiopenapi.Router) SingleRouteFunc { return r.Delete }},
		{"HEAD", "/ping", func(r chiopenapi.Router) SingleRouteFunc { return r.Head }},
		{"OPTIONS", "/ping", func(r chiopenapi.Router) SingleRouteFunc { return r.Options }},
		{"TRACE", "/ping", func(r chiopenapi.Router) SingleRouteFunc { return r.Trace }},
		{"PATCH", "/ping", func(r chiopenapi.Router) SingleRouteFunc { return r.Patch }},
		{"CONNECT", "/ping", func(r chiopenapi.Router) SingleRouteFunc { return r.Connect }},
	}
	for _, tt := range tests {
		t.Run(tt.method, func(t *testing.T) {
			c := chi.NewRouter()
			r := chiopenapi.NewRouter(c)
			tt.methodFunc(r)(tt.path, pingHandler).With(
				option.OperationID(tt.method+"Ping"),
				option.Summary("Ping the server with "+tt.method),
				option.Description("This endpoint is used to check if the server is running with a "+tt.method+" request."),
			)

			err := r.Validate()
			require.NoError(t, err, "failed to validate OpenAPI configuration")

			req := httptest.NewRequest(tt.method, tt.path, nil)
			rr := httptest.NewRecorder()
			c.ServeHTTP(rr, req)

			assert.Equal(t, http.StatusOK, rr.Code, "expected status OK for %s method", tt.method)
			assert.Contains(t, rr.Body.String(), "pong", "expected response body to be 'pong' for %s method", tt.method)

			if tt.method == "CONNECT" {
				return
			}

			schema, err := r.GenerateSchema()
			require.NoError(t, err, "failed to generate OpenAPI schema")

			assert.Contains(t, string(schema), tt.method, "expected OpenAPI schema to contain method %s", tt.method)
		})
	}
	t.Run("Method", func(t *testing.T) {
		c := chi.NewRouter()
		r := chiopenapi.NewRouter(c)
		r.Method("GET", "/ping", http.HandlerFunc(pingHandler)).With(
			option.OperationID("getPing"),
			option.Summary("Ping the server with GET"),
			option.Description("This endpoint is used to check if the server is running with a GET request."),
		)

		err := r.Validate()
		require.NoError(t, err, "failed to validate OpenAPI configuration")

		req := httptest.NewRequest("GET", "/ping", nil)
		rr := httptest.NewRecorder()
		c.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code, "expected status OK for GET method")
		assert.Contains(t, rr.Body.String(), "pong", "expected response body to be 'pong' for GET method")

		schema, err := r.GenerateSchema()
		require.NoError(t, err, "failed to generate OpenAPI schema")

		assert.Contains(t, string(schema), "GET", "expected OpenAPI schema to contain method GET")
	})
	t.Run("Method with Connect", func(t *testing.T) {
		c := chi.NewRouter()
		r := chiopenapi.NewRouter(c)
		r.Method("CONNECT", "/ping", http.HandlerFunc(pingHandler)).With(
			option.OperationID("connectPing"),
			option.Summary("Ping the server with CONNECT"),
			option.Description("This endpoint is used to check if the server is running with a CONNECT request."),
		)

		err := r.Validate()
		require.NoError(t, err, "failed to validate OpenAPI configuration")

		req := httptest.NewRequest("CONNECT", "/ping", nil)
		rr := httptest.NewRecorder()
		c.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code, "expected status OK for CONNECT method")
		assert.Contains(t, rr.Body.String(), "pong", "expected response body to be 'pong' for CONNECT method")

		schema, err := r.GenerateSchema()
		require.NoError(t, err, "failed to generate OpenAPI schema")

		assert.NotContains(t, string(schema), "CONNECT", "expected OpenAPI schema not to contain method CONNECT")
	})
	t.Run("Handle", func(t *testing.T) {
		c := chi.NewRouter()
		r := chiopenapi.NewRouter(c)
		r.Handle("/ping", http.HandlerFunc(pingHandler))

		err := r.Validate()
		require.NoError(t, err, "failed to validate OpenAPI configuration")

		req := httptest.NewRequest("GET", "/ping", nil)
		rr := httptest.NewRecorder()
		c.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code, "expected status OK for Handle method")
		assert.Contains(t, rr.Body.String(), "pong", "expected response body to be 'pong' for Handle method")
	})
	t.Run("HandleFunc", func(t *testing.T) {
		c := chi.NewRouter()
		r := chiopenapi.NewRouter(c)
		r.HandleFunc("/ping", pingHandler)

		err := r.Validate()
		require.NoError(t, err, "failed to validate OpenAPI configuration")

		req := httptest.NewRequest("GET", "/ping", nil)
		rr := httptest.NewRecorder()
		c.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code, "expected status OK for HandleFunc method")
		assert.Contains(t, rr.Body.String(), "pong", "expected response body to be 'pong' for HandleFunc method")
	})
	t.Run("Mount", func(t *testing.T) {
		c := chi.NewRouter()
		r := chiopenapi.NewRouter(c)
		subRouter := chiopenapi.NewRouter(chi.NewRouter())
		subRouter.Get("/ping", pingHandler).With(
			option.OperationID("getPing"),
			option.Summary("Ping the server with Mount"),
			option.Description("This endpoint is used to check if the server is running with a Mount request."),
		)
		r.Mount("/sub", subRouter)

		err := r.Validate()
		require.NoError(t, err, "failed to validate OpenAPI configuration")

		req := httptest.NewRequest("GET", "/sub/ping", nil)
		rr := httptest.NewRecorder()
		c.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code, "expected status OK for Mount method")
		assert.Contains(t, rr.Body.String(), "pong", "expected response body to be 'pong' for Mount method")
	})
}

func TestRouter_Group(t *testing.T) {
	c := chi.NewRouter()
	r := chiopenapi.NewRouter(c)
	r.Group(func(r chiopenapi.Router) {
		r.Get("/ping", pingHandler).With(
			option.OperationID("getPing"),
			option.Summary("Ping the server with Group"),
			option.Description("This endpoint is used to check if the server is running with a Group request."),
		)
	}).UseOptions(option.GroupTags("ping"), option.GroupDeprecated(true))

	err := r.Validate()
	require.NoError(t, err, "failed to validate OpenAPI configuration")

	req := httptest.NewRequest("GET", "/ping", nil)
	rr := httptest.NewRecorder()
	c.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code, "expected status OK for Group method")
	assert.Contains(t, rr.Body.String(), "pong", "expected response body to be 'pong' for Group method")
	schema, err := r.GenerateSchema()
	require.NoError(t, err, "failed to generate OpenAPI schema")
	assert.Contains(t, string(schema), "getPing", "expected OpenAPI schema to contain operation ID getPing")
	assert.Contains(t, string(schema), "Ping the server with Group", "expected OpenAPI schema to contain summary for getPing")
	assert.Contains(t, string(schema), "deprecated", "expected OpenAPI schema to contain deprecated flag for getPing")
}

func TestRouter_Middleware(t *testing.T) {
	t.Run("Use", func(t *testing.T) {
		called := false
		middleware := func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				called = true
				next.ServeHTTP(w, r)
			})
		}
		c := chi.NewRouter()
		r := chiopenapi.NewRouter(c)
		r.Group(func(r chiopenapi.Router) {
			r.Use(middleware)
			r.Get("/ping", pingHandler)
		})

		req := httptest.NewRequest("GET", "/ping", nil)
		rec := httptest.NewRecorder()
		c.ServeHTTP(rec, req)
		assert.True(t, called, "expected middleware to be called")
		assert.Equal(t, http.StatusOK, rec.Code, "expected status OK for Middleware method")
		assert.Contains(t, rec.Body.String(), "pong", "expected response body to be 'pong' for Middleware method")
	})
	t.Run("With", func(t *testing.T) {
		called := false
		middleware := func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				called = true
				next.ServeHTTP(w, r)
			})
		}
		c := chi.NewRouter()
		r := chiopenapi.NewRouter(c)
		r.With(middleware).Get("/ping", pingHandler)

		req := httptest.NewRequest("GET", "/ping", nil)
		rec := httptest.NewRecorder()
		c.ServeHTTP(rec, req)
		assert.True(t, called, "expected middleware to be called")
		assert.Equal(t, http.StatusOK, rec.Code, "expected status OK for With method")
		assert.Contains(t, rec.Body.String(), "pong", "expected response body to be 'pong' for With method")
	})
}

func TestRouter_NotFound(t *testing.T) {
	c := chi.NewRouter()
	r := chiopenapi.NewRouter(c)
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Not Found", http.StatusNotFound)
	})

	err := r.Validate()
	require.NoError(t, err, "failed to validate OpenAPI configuration")

	req := httptest.NewRequest(http.MethodGet, "/not-found", nil)
	rr := httptest.NewRecorder()
	c.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code, "expected status Not Found")
	assert.Equal(t, "Not Found\n", rr.Body.String(), "expected response body to be 'Not Found'")
}

func TestRouter_MethodNotAllowed(t *testing.T) {
	c := chi.NewRouter()
	r := chiopenapi.NewRouter(c)
	r.Get("/ping", pingHandler).With(
		option.OperationID("getPing"),
		option.Summary("Ping the server"),
		option.Description("This endpoint is used to check if the server is running."),
	)
	r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	})

	err := r.Validate()
	require.NoError(t, err, "failed to validate OpenAPI configuration")

	req := httptest.NewRequest(http.MethodPost, "/ping", nil)
	rr := httptest.NewRecorder()
	c.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusMethodNotAllowed, rr.Code, "expected status Method Not Allowed")
	assert.Equal(t, "Method Not Allowed\n", rr.Body.String(), "expected response body to be 'Method Not Allowed'")
}

func TestGenerator_Docs(t *testing.T) {
	c := chi.NewRouter()
	r := chiopenapi.NewRouter(c)
	r.Get("/ping", pingHandler).With(
		option.OperationID("getPing"),
		option.Summary("Ping the server"),
		option.Description("This endpoint is used to check if the server is running."),
	)

	t.Run("should register docs route", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/docs", nil)
		rr := httptest.NewRecorder()
		c.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code, "expected status OK for /docs route")
		assert.Contains(t, rr.Body.String(), "Chi OpenAPI", "expected response body to contain 'Chi OpenAPI'")
	})
	t.Run("should register docs file route", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/docs/openapi.yaml", nil)
		rr := httptest.NewRecorder()
		c.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code, "expected status OK for /docs/openapi.yaml route")
		assert.Contains(t, rr.Body.String(), "openapi: 3.0.3", "expected response body to contain 'openapi: 3.0.3'")
	})
}

func TestGenerator_DisableDocs(t *testing.T) {
	c := chi.NewRouter()
	r := chiopenapi.NewRouter(c, option.WithDisableDocs(true))
	r.Get("/ping", pingHandler).With(
		option.OperationID("getPing"),
		option.Summary("Ping the server"),
		option.Description("This endpoint is used to check if the server is running."),
	)

	t.Run("should not register docs route", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/docs", nil)
		rr := httptest.NewRecorder()
		c.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusNotFound, rr.Code, "expected status Not Found for /docs route")
	})
	t.Run("should not register docs file route", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/docs/openapi.yaml", nil)
		rr := httptest.NewRecorder()
		c.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusNotFound, rr.Code, "expected status Not Found for /docs/openapi.yaml route")
	})
}

func TestGenerator_MarshalJSON(t *testing.T) {
	c := chi.NewRouter()
	r := chiopenapi.NewRouter(c)
	r.Get("/ping", pingHandler).With(
		option.OperationID("getPing"),
		option.Summary("Ping the server"),
		option.Description("This endpoint is used to check if the server is running."),
	)

	err := r.Validate()
	require.NoError(t, err, "failed to validate OpenAPI configuration")

	schema, err := r.MarshalJSON()
	require.NoError(t, err, "failed to marshal OpenAPI schema to JSON")
	assert.NotEmpty(t, schema, "expected non-empty OpenAPI schema JSON")
	assert.Contains(t, string(schema), `"openapi": "3.0.3"`, "expected OpenAPI version in schema JSON")
	assert.Contains(t, string(schema), `"title": "Chi OpenAPI"`, "expected title in schema JSON")
}

func TestGenerator_MarshalYAML(t *testing.T) {
	c := chi.NewRouter()
	r := chiopenapi.NewRouter(c)
	r.Get("/ping", pingHandler).With(
		option.OperationID("getPing"),
		option.Summary("Ping the server"),
		option.Description("This endpoint is used to check if the server is running."),
	)

	err := r.Validate()
	require.NoError(t, err, "failed to validate OpenAPI configuration")

	schema, err := r.MarshalYAML()
	require.NoError(t, err, "failed to marshal OpenAPI schema to YAML")
	assert.NotEmpty(t, schema, "expected non-empty OpenAPI schema YAML")
	assert.Contains(t, string(schema), "openapi: 3.0.3", "expected OpenAPI version in schema YAML")
	assert.Contains(t, string(schema), "title: Chi OpenAPI", "expected title in schema YAML")
}

func TestGenerator_WriteSchemaTo(t *testing.T) {
	c := chi.NewRouter()
	r := chiopenapi.NewRouter(c)
	r.Get("/ping", pingHandler).With(
		option.OperationID("getPing"),
		option.Summary("Ping the server"),
		option.Description("This endpoint is used to check if the server is running."),
	)

	err := r.Validate()
	require.NoError(t, err, "failed to validate OpenAPI configuration")

	testDir := t.TempDir()
	goldenPath := filepath.Join(testDir, "openapi.yaml")
	err = r.WriteSchemaTo(goldenPath)
	require.NoError(t, err, "failed to write OpenAPI schema to file")
	assert.FileExists(t, goldenPath, "expected OpenAPI schema file to be created")
	schema, err := os.ReadFile(goldenPath)
	require.NoError(t, err, "failed to read OpenAPI schema file")
	assert.NotEmpty(t, schema, "expected non-empty OpenAPI schema file")
	assert.Contains(t, string(schema), "openapi: 3.0.3", "expected OpenAPI version in schema file")
	assert.Contains(t, string(schema), "title: Chi OpenAPI", "expected title in schema file")
}

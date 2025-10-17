package muxopenapi_test

import (
	"encoding/json"
	"flag"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/gorilla/mux"
	"github.com/oaswrap/spec/adapter/muxopenapi"
	"github.com/oaswrap/spec/openapi"
	"github.com/oaswrap/spec/option"
	"github.com/oaswrap/spec/pkg/dto"
	"github.com/oaswrap/spec/pkg/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//nolint:gochecknoglobals // test flag for golden file updates
var update = flag.Bool("update", false, "update golden files")

func TestRouter_Spec(t *testing.T) {
	tests := []struct {
		name      string
		golden    string
		opts      []option.OpenAPIOption
		setup     func(r muxopenapi.Router)
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
			setup: func(r muxopenapi.Router) {
				pet := r.PathPrefix("pet").Subrouter().With(
					option.GroupTags("pet"),
					option.GroupSecurity("petstore_auth", "write:pets", "read:pets"),
				)
				pet.HandleFunc("/", nil).Methods("GET").With(
					option.OperationID("updatePet"),
					option.Summary("Update an existing pet"),
					option.Description("Update the details of an existing pet in the store."),
					option.Request(new(dto.Pet)),
					option.Response(200, new(dto.Pet)),
				)
				pet.HandleFunc("/", nil).Methods("POST").With(
					option.OperationID("addPet"),
					option.Summary("Add a new pet"),
					option.Description("Add a new pet to the store."),
					option.Request(new(dto.Pet)),
					option.Response(201, new(dto.Pet)),
				)
				pet.HandleFunc("/findByStatus", nil).Methods("GET").With(
					option.OperationID("findPetsByStatus"),
					option.Summary("Find pets by status"),
					option.Description("Finds Pets by status. Multiple status values can be provided with comma separated strings."),
					option.Request(new(struct {
						Status string `query:"status" enum:"available,pending,sold"`
					})),
					option.Response(200, new([]dto.Pet)),
				)
				pet.HandleFunc("/findByTags", nil).Methods("GET").With(
					option.OperationID("findPetsByTags"),
					option.Summary("Find pets by tags"),
					option.Description("Finds Pets by tags. Multiple tags can be provided with comma separated strings."),
					option.Request(new(struct {
						Tags []string `query:"tags"`
					})),
					option.Response(200, new([]dto.Pet)),
				)
				pet.HandleFunc("/{petId}/uploadImage", nil).Methods("POST").With(
					option.OperationID("uploadFile"),
					option.Summary("Upload an image for a pet"),
					option.Description("Uploads an image for a pet."),
					option.Request(new(dto.UploadImageRequest)),
					option.Response(200, new(dto.APIResponse)),
				)
				pet.HandleFunc("/{petId}", nil).Methods("GET").With(
					option.OperationID("getPetById"),
					option.Summary("Get pet by ID"),
					option.Description("Retrieve a pet by its ID."),
					option.Request(new(struct {
						ID int `path:"petId" required:"true"`
					})),
					option.Response(200, new(dto.Pet)),
				)
				pet.HandleFunc("/{petId}", nil).Methods("POST").With(
					option.OperationID("updatePetWithForm"),
					option.Summary("Update pet with form"),
					option.Description("Updates a pet in the store with form data."),
					option.Request(new(dto.UpdatePetWithFormRequest)),
					option.Response(200, nil),
				)
				pet.HandleFunc("/delete/{petId}", nil).Methods("DELETE").With(
					option.OperationID("deletePet"),
					option.Summary("Delete a pet"),
					option.Description("Delete a pet from the store by its ID."),
					option.Request(new(dto.DeletePetRequest)),
					option.Response(204, nil),
				)

				store := r.PathPrefix("store").Subrouter().With(
					option.GroupTags("store"),
				)
				store.HandleFunc("/order", nil).Methods("POST").With(
					option.OperationID("placeOrder"),
					option.Summary("Place an order"),
					option.Description("Place a new order for a pet."),
					option.Request(new(dto.Order)),
					option.Response(201, new(dto.Order)),
				)
				store.HandleFunc("/order/{orderId}", nil).Methods("GET").With(
					option.OperationID("getOrderById"),
					option.Summary("Get order by ID"),
					option.Description("Retrieve an order by its ID."),
					option.Request(new(struct {
						ID int `path:"orderId" required:"true"`
					})),
					option.Response(200, new(dto.Order)),
					option.Response(404, nil),
				)
				store.HandleFunc("/order/{orderId}", nil).Methods("DELETE").With(
					option.OperationID("deleteOrder"),
					option.Summary("Delete an order"),
					option.Description("Delete an order by its ID."),
					option.Request(new(struct {
						ID int `path:"orderId" required:"true"`
					})),
					option.Response(204, nil),
				)

				user := r.PathPrefix("user").Subrouter().With(
					option.GroupTags("user"),
				)
				user.HandleFunc("/createWithList", nil).Methods("POST").With(
					option.OperationID("createUsersWithList"),
					option.Summary("Create users with list"),
					option.Description("Create multiple users in the store with a list."),
					option.Request(new([]dto.PetUser)),
					option.Response(201, nil),
				)
				user.HandleFunc("/", nil).Methods("POST").With(
					option.OperationID("createUser"),
					option.Summary("Create a new user"),
					option.Description("Create a new user in the store."),
					option.Request(new(dto.PetUser)),
					option.Response(201, new(dto.PetUser)),
				)
				user.HandleFunc("/{username}", nil).Methods("GET").With(
					option.OperationID("getUserByName"),
					option.Summary("Get user by username"),
					option.Description("Retrieve a user by their username."),
					option.Request(new(struct {
						Username string `path:"username" required:"true"`
					})),
					option.Response(200, new(dto.PetUser)),
					option.Response(404, nil),
				)
				user.HandleFunc("/{username}", nil).Methods("PUT").With(
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
				user.HandleFunc("/{username}", nil).Methods("DELETE").With(
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
			mux := mux.NewRouter()
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
			r := muxopenapi.NewRouter(mux, opts...)

			if tt.setup != nil {
				tt.setup(r)
			}

			if tt.shouldErr {
				err := r.Validate()
				require.Error(t, err, "expected error for invalid OpenAPI configuration")
				return
			}
			err := r.Validate()
			require.NoError(t, err, "failed to validate OpenAPI configuration")

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

func PingHandler(w http.ResponseWriter, _ *http.Request) {
	_ = json.NewEncoder(w).Encode(map[string]string{"message": "pong"})
}

func TestRouter_HandleFunc(t *testing.T) {
	methods := []string{"GET", "POST", "PUT", "DELETE", "PATCH", "HEAD", "OPTIONS", "TRACE", "CONNECT"}
	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			mux := mux.NewRouter()
			r := muxopenapi.NewRouter(mux)

			route := r.HandleFunc("/ping", PingHandler).Methods(method).With(
				option.Summary("Ping endpoint"),
			).Name("ping")

			assert.NotNil(t, r.Get("ping"))
			assert.NotNil(t, r.GetRoute("ping"))
			assert.NotNil(t, route.GetHandler())
			assert.Equal(t, "ping", route.GetName())
			require.NoError(t, route.GetError())

			// Test the /ping endpoint
			req := httptest.NewRequest(method, "/ping", nil)
			rec := httptest.NewRecorder()
			mux.ServeHTTP(rec, req)

			assert.Equal(t, http.StatusOK, rec.Code)
			assert.JSONEq(t, `{"message": "pong"}`, rec.Body.String())

			schema, err := r.GenerateSchema()
			require.NoError(t, err)

			assert.NotNil(t, schema)

			if method == "CONNECT" {
				assert.NotContains(t, string(schema), "/ping")
				return
			}

			assert.Contains(t, string(schema), "/ping")
		})
	}
}

func TestRouter_Handle(t *testing.T) {
	methods := []string{"GET", "POST", "PUT", "DELETE", "PATCH", "HEAD", "OPTIONS", "TRACE", "CONNECT"}
	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			mux := mux.NewRouter()
			r := muxopenapi.NewRouter(mux)

			route := r.Handle("/ping", http.HandlerFunc(PingHandler)).Methods(method).With(
				option.Summary("Ping endpoint"),
			).Name("ping")

			assert.NotNil(t, r.Get("ping"))
			assert.NotNil(t, r.GetRoute("ping"))
			assert.NotNil(t, route.GetHandler())
			assert.Equal(t, "ping", route.GetName())
			require.NoError(t, route.GetError())

			// Test the /ping endpoint
			req := httptest.NewRequest(method, "/ping", nil)
			rec := httptest.NewRecorder()
			mux.ServeHTTP(rec, req)

			assert.Equal(t, http.StatusOK, rec.Code)
			assert.JSONEq(t, `{"message": "pong"}`, rec.Body.String())

			schema, err := r.GenerateSchema()
			require.NoError(t, err)

			assert.NotNil(t, schema)

			if method == "CONNECT" {
				assert.NotContains(t, string(schema), "/ping")
				return
			}

			assert.Contains(t, string(schema), "/ping")
		})
	}
}

func TestRouter_Queries(t *testing.T) {
	mux := mux.NewRouter()
	r := muxopenapi.NewRouter(mux)

	r.Queries("foo", "bar").Methods("GET").Path("/ping").HandlerFunc(PingHandler).With(
		option.Summary("Ping endpoint"),
	)
	t.Run("when match should return 200", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/ping?foo=bar", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.JSONEq(t, `{"message": "pong"}`, rec.Body.String())
	})
	t.Run("when not match should return 404", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/ping?foo=baz", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
	})
}

func TestRouter_Headers(t *testing.T) {
	mux := mux.NewRouter()
	r := muxopenapi.NewRouter(mux)

	r.Headers("X-Foo", "bar").Methods("GET").Path("/ping").HandlerFunc(PingHandler).With(
		option.Summary("Ping endpoint"),
	)

	t.Run("when match should return 200", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/ping", nil)
		req.Header.Set("X-Foo", "bar")
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.JSONEq(t, `{"message": "pong"}`, rec.Body.String())
	})
	t.Run("when not match should return 404", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/ping", nil)
		req.Header.Set("X-Foo", "baz")
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
	})
}

func TestRouter_Subrouter(t *testing.T) {
	mux := mux.NewRouter()
	r := muxopenapi.NewRouter(mux)

	api := r.PathPrefix("/api").Subrouter()
	api.HandleFunc("/ping", PingHandler).Methods("GET").With(
		option.Summary("Ping API endpoint"),
	)

	// Test the /api/ping endpoint
	req := httptest.NewRequest(http.MethodGet, "/api/ping", nil)

	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.JSONEq(t, `{"message": "pong"}`, rr.Body.String())

	schema, err := r.GenerateSchema()
	require.NoError(t, err)

	assert.NotNil(t, schema)
	assert.Contains(t, string(schema), "/api/ping")
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
		mux := mux.NewRouter()
		r := muxopenapi.NewRouter(mux)
		r.Use(middleware)
		r.HandleFunc("/ping", PingHandler).Methods("GET").With(
			option.Summary("Ping endpoint"),
		)

		req := httptest.NewRequest(http.MethodGet, "/ping", nil)
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		assert.True(t, called)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.JSONEq(t, `{"message": "pong"}`, rec.Body.String())
	})
}

func TestGenerator_Docs(t *testing.T) {
	mux := mux.NewRouter()
	r := muxopenapi.NewRouter(mux)

	r.HandleFunc("/ping", PingHandler).Methods("GET").With(
		option.Summary("Ping endpoint"),
	)

	t.Run("should serve docs", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/docs", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), "Mux OpenAPI")
	})
	t.Run("should serve docs/openapi.yaml", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/docs/openapi.yaml", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), "/ping")
	})
}

func TestGenerator_DisableDocs(t *testing.T) {
	mux := mux.NewRouter()
	r := muxopenapi.NewRouter(mux, option.WithDisableDocs(true))

	r.HandleFunc("/ping", PingHandler).Methods("GET").With(
		option.Summary("Ping endpoint"),
	)

	t.Run("should not serve docs", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/docs", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("should not serve docs/openapi.yaml", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/docs/openapi.yaml", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
	})
}

func TestGenerator_MarshalJSON(t *testing.T) {
	mux := mux.NewRouter()
	r := muxopenapi.NewRouter(mux)

	r.HandleFunc("/ping", PingHandler).Methods("GET").With(
		option.Summary("Ping endpoint"),
	)

	schema, err := r.GenerateSchema()
	require.NoError(t, err)

	assert.NotNil(t, schema)

	// Test JSON marshaling
	jsonData, err := r.MarshalJSON()
	require.NoError(t, err)

	assert.NotNil(t, jsonData)
	assert.Contains(t, string(jsonData), "/ping")
}

func TestGenerator_MarshalYAML(t *testing.T) {
	mux := mux.NewRouter()
	r := muxopenapi.NewRouter(mux)

	r.HandleFunc("/ping", PingHandler).Methods("GET").With(
		option.Summary("Ping endpoint"),
	)

	schema, err := r.GenerateSchema()
	require.NoError(t, err)

	assert.NotNil(t, schema)

	// Test YAML marshaling
	yamlData, err := r.MarshalYAML()
	require.NoError(t, err)

	assert.NotNil(t, yamlData)
	assert.Contains(t, string(yamlData), "/ping")
}

func TestGenerator_WriteSchemaTo(t *testing.T) {
	mux := mux.NewRouter()
	r := muxopenapi.NewRouter(mux)

	r.HandleFunc("/ping", PingHandler).Methods("GET").With(
		option.Summary("Ping endpoint"),
	)

	schema, err := r.GenerateSchema()
	require.NoError(t, err)

	assert.NotNil(t, schema)

	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "openapi.yaml")

	err = r.WriteSchemaTo(filePath)
	require.NoError(t, err, "failed to write schema to file")
	defer func() {
		_ = os.Remove(filePath)
	}()

	want, err := os.ReadFile(filePath)
	require.NoError(t, err, "failed to read schema from file")

	testutil.EqualYAML(t, want, schema)
}

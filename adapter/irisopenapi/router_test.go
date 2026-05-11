package irisopenapi_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	stoplightemb "github.com/oaswrap/spec-ui/stoplightemb"
	"github.com/oaswrap/spec/internal/testutil"
	"github.com/oaswrap/spec/internal/testutil/dto"
	"github.com/oaswrap/spec/openapi"
	"github.com/oaswrap/spec/option"

	"github.com/oaswrap/spec/adapter/irisopenapi"
)

func TestRouter_Spec(t *testing.T) {
	tests := []struct {
		name      string
		golden    string
		opts      []option.OpenAPIOption
		setup     func(r irisopenapi.Router)
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
						Implicit: &openapi.OAuthFlow{
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
			setup: func(r irisopenapi.Router) {
				pet := r.Party("/pet").With(
					option.GroupTags("pet"),
					option.GroupSecurity("petstore_auth", "write:pets", "read:pets"),
				)
				pet.Put("/", nil).With(
					option.OperationID("updatePet"),
					option.Summary("Update an existing pet"),
					option.Description("Update the details of an existing pet in the store."),
					option.Request(new(dto.Pet)),
					option.Response(200, new(dto.Pet)),
				)
				pet.Post("/", nil).With(
					option.OperationID("addPet"),
					option.Summary("Add a new pet"),
					option.Description("Add a new pet to the store."),
					option.Request(new(dto.Pet)),
					option.Response(201, new(dto.Pet)),
				)
				pet.Get("/findByStatus", nil).With(
					option.OperationID("findPetsByStatus"),
					option.Summary("Find pets by status"),
					option.Description("Finds Pets by status. Multiple status values can be provided with comma separated strings."),
					option.Request(new(struct {
						Status string `query:"status" enum:"available,pending,sold"`
					})),
					option.Response(200, new([]dto.Pet)),
				)
				pet.Get("/findByTags", nil).With(
					option.OperationID("findPetsByTags"),
					option.Summary("Find pets by tags"),
					option.Description("Finds Pets by tags. Multiple tags can be provided with comma separated strings."),
					option.Request(new(struct {
						Tags []string `query:"tags"`
					})),
					option.Response(200, new([]dto.Pet)),
				)
				pet.Post("/{petId}/uploadImage", nil).With(
					option.OperationID("uploadFile"),
					option.Summary("Upload an image for a pet"),
					option.Description("Uploads an image for a pet."),
					option.Request(new(dto.UploadImageRequest)),
					option.Response(200, new(dto.APIResponse)),
				)
				pet.Get("/{petId}", nil).With(
					option.OperationID("getPetById"),
					option.Summary("Get pet by ID"),
					option.Description("Retrieve a pet by its ID."),
					option.Request(new(struct {
						ID int `param:"petId" required:"true"`
					})),
					option.Response(200, new(dto.Pet)),
				)
				pet.Post("/{petId}", nil).With(
					option.OperationID("updatePetWithForm"),
					option.Summary("Update pet with form"),
					option.Description("Updates a pet in the store with form data."),
					option.Request(new(dto.UpdatePetWithFormRequest)),
					option.Response(200, nil),
				)
				pet.Delete("/{petId}", nil).With(
					option.OperationID("deletePet"),
					option.Summary("Delete a pet"),
					option.Description("Delete a pet from the store by its ID."),
					option.Request(new(dto.DeletePetRequest)),
					option.Response(204, nil),
				)
				store := r.Party("/store").With(
					option.GroupTags("store"),
				)
				store.Post("/order", nil).With(
					option.OperationID("placeOrder"),
					option.Summary("Place an order"),
					option.Description("Place a new order for a pet."),
					option.Request(new(dto.Order)),
					option.Response(201, new(dto.Order)),
				)
				store.Get("/order/{orderId}", nil).With(
					option.OperationID("getOrderById"),
					option.Summary("Get order by ID"),
					option.Description("Retrieve an order by its ID."),
					option.Request(new(struct {
						ID int `param:"orderId" required:"true"`
					})),
					option.Response(200, new(dto.Order)),
					option.Response(404, nil),
				)
				store.Delete("/order/{orderId}", nil).With(
					option.OperationID("deleteOrder"),
					option.Summary("Delete an order"),
					option.Description("Delete an order by its ID."),
					option.Request(new(struct {
						ID int `param:"orderId" required:"true"`
					})),
					option.Response(204, nil),
				)

				user := r.Party("/user").With(
					option.GroupTags("user"),
				)
				user.Post("/createWithList", nil).With(
					option.OperationID("createUsersWithList"),
					option.Summary("Create users with list"),
					option.Description("Create multiple users in the store with a list."),
					option.Request(new([]dto.PetUser)),
					option.Response(201, nil),
				)
				user.Post("/", nil).With(
					option.OperationID("createUser"),
					option.Summary("Create a new user"),
					option.Description("Create a new user in the store."),
					option.Request(new(dto.PetUser)),
					option.Response(201, new(dto.PetUser)),
				)
				user.Get("/{username}", nil).With(
					option.OperationID("getUserByName"),
					option.Summary("Get user by username"),
					option.Description("Retrieve a user by their username."),
					option.Request(new(struct {
						Username string `param:"username" required:"true"`
					})),
					option.Response(200, new(dto.PetUser)),
					option.Response(404, nil),
				)
				user.Put("/{username}", nil).With(
					option.OperationID("updateUser"),
					option.Summary("Update an existing user"),
					option.Description("Update the details of an existing user."),
					option.Request(new(struct {
						dto.PetUser

						Username string `param:"username" required:"true"`
					})),
					option.Response(200, new(dto.PetUser)),
					option.Response(404, nil),
				)
				user.Delete("/{username}", nil).With(
					option.OperationID("deleteUser"),
					option.Summary("Delete a user"),
					option.Description("Delete a user from the store by their username."),
					option.Request(new(struct {
						Username string `param:"username" required:"true"`
					})),
					option.Response(204, nil),
				)
			},
		},
		{
			name: "Invalid Open API Version",
			opts: []option.OpenAPIOption{
				option.WithOpenAPIVersion("2.0.0"),
			},
			shouldErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := iris.New()
			opts := []option.OpenAPIOption{
				option.WithTitle("Test API " + tt.name),
				option.WithVersion("1.0.0"),
				option.WithDescription("This is a test API for " + tt.name),
				option.WithReflectorConfig(
					option.StripDefNamePrefix("IrisopenapiTest"),
				),
			}
			if len(tt.opts) > 0 {
				opts = append(opts, tt.opts...)
			}
			r := irisopenapi.NewRouter(app, opts...)

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

			schema, err := r.GenerateSchema()
			require.NoError(t, err, "failed to generate OpenAPI schema")
			testutil.AssertGolden(t, schema, filepath.Join("testdata", tt.golden+".yaml"))
		})
	}
}

type SingleRouteFunc func(path string, handlers ...context.Handler) irisopenapi.Route

func PingHandler(ctx iris.Context) {
	_ = ctx.JSON(iris.Map{"message": "pong"})
}

func TestRouter_Single(t *testing.T) {
	tests := []struct {
		method     string
		path       string
		methodFunc func(r irisopenapi.Router) SingleRouteFunc
	}{
		{http.MethodGet, "/ping", func(r irisopenapi.Router) SingleRouteFunc { return r.Get }},
		{http.MethodPost, "/ping", func(r irisopenapi.Router) SingleRouteFunc { return r.Post }},
		{http.MethodPut, "/ping", func(r irisopenapi.Router) SingleRouteFunc { return r.Put }},
		{http.MethodDelete, "/ping", func(r irisopenapi.Router) SingleRouteFunc { return r.Delete }},
		{http.MethodPatch, "/ping", func(r irisopenapi.Router) SingleRouteFunc { return r.Patch }},
		{http.MethodHead, "/ping", func(r irisopenapi.Router) SingleRouteFunc { return r.Head }},
		{http.MethodOptions, "/ping", func(r irisopenapi.Router) SingleRouteFunc { return r.Options }},
		{http.MethodTrace, "/ping", func(r irisopenapi.Router) SingleRouteFunc { return r.Trace }},
		{http.MethodConnect, "/ping", func(r irisopenapi.Router) SingleRouteFunc { return r.Connect }},
	}

	for _, tt := range tests {
		t.Run(tt.method, func(t *testing.T) {
			app := iris.New()
			r := irisopenapi.NewRouter(app)

			route := tt.methodFunc(r)(tt.path, PingHandler).With(
				option.OperationID(fmt.Sprintf("ping-%s", tt.method)),
			)
			assert.NotNil(t, route, "expected route to be created")
			assert.Equal(t, tt.method, route.Method(), "expected method")
			assert.Equal(t, tt.path, route.Path(), "expected path")

			err := app.Build()
			require.NoError(t, err)

			req := httptest.NewRequest(tt.method, tt.path, nil)
			rec := httptest.NewRecorder()
			app.ServeHTTP(rec, req)
			assert.Equal(t, http.StatusOK, rec.Code, "expected status code 200")

			schema, err := r.GenerateSchema()
			require.NoError(t, err, "failed to generate OpenAPI schema")
			if tt.method == http.MethodConnect {
				assert.NotContains(t, string(schema), "operationId: ping-CONNECT")
				return
			}

			assert.Contains(t, string(schema), fmt.Sprintf("operationId: ping-%s", tt.method))
		})
	}
}

func TestRouter_Group(t *testing.T) {
	app := iris.New()
	r := irisopenapi.NewRouter(app)

	called := false
	group := r.Party("/api", func(ctx iris.Context) {
		called = true
		ctx.Next()
	})

	group.Get("/ping", PingHandler).With(
		option.OperationID("pingHandler"),
	)

	err := app.Build()
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodGet, "/api/ping", nil)
	rec := httptest.NewRecorder()
	app.ServeHTTP(rec, req)

	assert.True(t, called, "expected middleware to be called")
	assert.Equal(t, http.StatusOK, rec.Code)

	schema, err := r.GenerateSchema()
	require.NoError(t, err, "failed to generate OpenAPI schema")
	assert.Contains(t, string(schema), "operationId: pingHandler")
}

func TestRouter_Use(t *testing.T) {
	app := iris.New()
	r := irisopenapi.NewRouter(app)

	called := false
	r.Use(func(ctx iris.Context) {
		called = true
		ctx.Next()
	})
	r.Get("/test", PingHandler).With(option.OperationID("testHandler"))

	err := app.Build()
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()
	app.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.True(t, called, "expected middleware to be called")
}

func TestGenerator_Docs(t *testing.T) {
	app := iris.New()
	r := irisopenapi.NewRouter(app)

	r.Get("/ping", PingHandler).With(
		option.OperationID("pingHandler"),
	)

	err := app.Build()
	require.NoError(t, err)

	t.Run("should serve docs", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/docs", nil)
		rec := httptest.NewRecorder()
		app.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), "Iris OpenAPI")
	})

	t.Run("should serve OpenAPI YAML", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/docs/openapi.yaml", nil)
		rec := httptest.NewRecorder()
		app.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.NotEmpty(t, rec.Body.String())
		assert.Contains(t, rec.Header().Get("Content-Type"), "application/x-yaml")
		assert.Contains(t, rec.Body.String(), "openapi: 3.1.2")
	})
}

func TestGenerator_DisableDocs(t *testing.T) {
	app := iris.New()
	r := irisopenapi.NewGenerator(app, option.WithDisableDocs(true))

	r.Get("/test", PingHandler)

	err := app.Build()
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodGet, "/docs", nil)
	rec := httptest.NewRecorder()
	app.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusNotFound, rec.Code)

	req = httptest.NewRequest(http.MethodGet, "/docs/openapi.yaml", nil)
	rec = httptest.NewRecorder()
	app.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestGenerator_Assets(t *testing.T) {
	app := iris.New()
	r := irisopenapi.NewRouter(app, option.WithUIOption(stoplightemb.WithUI()))

	r.Get("/ping", PingHandler).With(
		option.OperationID("pingHandler"),
	)

	err := app.Build()
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodGet, "/docs/_assets/styles.min.css", nil)
	rec := httptest.NewRecorder()
	app.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestGenerator_WriteSchemaTo(t *testing.T) {
	app := iris.New()
	r := irisopenapi.NewRouter(app)

	r.Get("/test", PingHandler).With(option.OperationID("pingHandler"))

	err := r.Validate()
	require.NoError(t, err, "failed to validate OpenAPI configuration")

	tempFile, err := os.CreateTemp(t.TempDir(), "openapi-schema-*.yaml")
	require.NoError(t, err)
	defer func() {
		_ = os.Remove(tempFile.Name())
	}()

	err = r.WriteSchemaTo(tempFile.Name())
	require.NoError(t, err)

	schema, err := os.ReadFile(tempFile.Name())
	require.NoError(t, err)
	assert.NotEmpty(t, schema)
}

func TestGenerator_MarshalYAML(t *testing.T) {
	app := iris.New()
	r := irisopenapi.NewRouter(app)

	r.Get("/test", PingHandler).With(option.OperationID("pingHandler"))

	err := r.Validate()
	require.NoError(t, err)

	schema, err := r.MarshalYAML()
	require.NoError(t, err)
	assert.NotEmpty(t, schema)
	assert.Contains(t, string(schema), "openapi:")
}

func TestGenerator_MarshalJSON(t *testing.T) {
	app := iris.New()
	r := irisopenapi.NewRouter(app)

	r.Get("/test", PingHandler).With(option.OperationID("pingHandler"))

	err := r.Validate()
	require.NoError(t, err)

	schema, err := r.MarshalJSON()
	require.NoError(t, err)
	assert.NotEmpty(t, schema)
	assert.Contains(t, string(schema), `"openapi":`)
}

func TestGenerator_ValidateReport(t *testing.T) {
	app := iris.New()
	r := irisopenapi.NewRouter(app,
		option.WithContact(openapi.Contact{Name: "Support"}),
		option.WithLicense(openapi.License{Name: "MIT"}),
		option.WithServer("https://example.com"),
	)
	err := r.ValidateReport()
	assert.NoError(t, err)
}

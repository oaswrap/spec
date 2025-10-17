package fiberopenapi_test

import (
	"flag"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/oaswrap/spec/adapter/fiberopenapi"
	"github.com/oaswrap/spec/openapi"
	"github.com/oaswrap/spec/option"
	"github.com/oaswrap/spec/pkg/dto"
	"github.com/oaswrap/spec/pkg/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//nolint:gochecknoglobals // test flag for golden file updates
var update = flag.Bool("update", false, "update golden files")

func PingHandler(c *fiber.Ctx) error {
	return c.SendString("pong")
}

func TestRouter_Spec(t *testing.T) {
	tests := []struct {
		name      string
		golden    string
		options   []option.OpenAPIOption
		setup     func(r fiberopenapi.Router)
		shouldErr bool
	}{
		{
			name:   "Pet Store API",
			golden: "petstore.yaml",
			options: []option.OpenAPIOption{
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
			setup: func(r fiberopenapi.Router) {
				pet := r.Group("/pet").With(
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
				pet.Post("/:petId/uploadImage", nil).With(
					option.OperationID("uploadFile"),
					option.Summary("Upload an image for a pet"),
					option.Description("Uploads an image for a pet."),
					option.Request(new(dto.UploadImageRequest)),
					option.Response(200, new(dto.APIResponse)),
				)
				pet.Get("/:petId", nil).With(
					option.OperationID("getPetById"),
					option.Summary("Get pet by ID"),
					option.Description("Retrieve a pet by its ID."),
					option.Request(new(struct {
						ID int `path:"petId" required:"true"`
					})),
					option.Response(200, new(dto.Pet)),
				)
				pet.Post("/:petId", nil).With(
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
				store := r.Group("/store").With(
					option.GroupTags("store"),
				)
				store.Post("/order", nil).With(
					option.OperationID("placeOrder"),
					option.Summary("Place an order"),
					option.Description("Place a new order for a pet."),
					option.Request(new(dto.Order)),
					option.Response(201, new(dto.Order)),
				)
				store.Get("/order/:orderId", nil).With(
					option.OperationID("getOrderById"),
					option.Summary("Get order by ID"),
					option.Description("Retrieve an order by its ID."),
					option.Request(new(struct {
						ID int `path:"orderId" required:"true"`
					})),
					option.Response(200, new(dto.Order)),
					option.Response(404, nil),
				)
				store.Delete("/order/:orderId", nil).With(
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
				user.Get("/:username", nil).With(
					option.OperationID("getUserByName"),
					option.Summary("Get user by username"),
					option.Description("Retrieve a user by their username."),
					option.Request(new(struct {
						Username string `path:"username" required:"true"`
					})),
					option.Response(200, new(dto.PetUser)),
					option.Response(404, nil),
				)
				user.Put("/:username", nil).With(
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
				user.Delete("/:username", nil).With(
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
			name: "Invalid OpenAPI Version",
			options: []option.OpenAPIOption{
				option.WithTitle("Invalid OpenAPI Version"),
				option.WithOpenAPIVersion("2.0"), // Intentionally invalid for testing
				option.WithDescription("This is a test API with an invalid OpenAPI version"),
			},
			shouldErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := fiber.New()
			opts := []option.OpenAPIOption{
				option.WithTitle("Test API " + tt.name),
				option.WithVersion("1.0.0"),
				option.WithDescription("This is a test API for " + tt.name),
				option.WithReflectorConfig(
					option.RequiredPropByValidateTag(),
				),
			}
			if len(tt.options) > 0 {
				opts = append(opts, tt.options...)
			}
			r := fiberopenapi.NewRouter(app, opts...)

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
			goldenFile := filepath.Join("testdata", tt.golden)

			if *update {
				err = r.WriteSchemaTo(goldenFile)
				require.NoError(t, err, "failed to write golden file")
				t.Logf("Updated golden file: %s", goldenFile)
			}

			want, err := os.ReadFile(goldenFile)
			require.NoError(t, err, "failed to read golden file %s", goldenFile)

			testutil.EqualYAML(t, want, schema)
		})
	}
}

type SingleRouteFunc func(string, ...fiber.Handler) fiberopenapi.Route

func TestRouter_Single(t *testing.T) {
	tests := []struct {
		method     string
		path       string
		methodFunc func(r fiberopenapi.Router) SingleRouteFunc
	}{
		{"GET", "/ping", func(r fiberopenapi.Router) SingleRouteFunc { return r.Get }},
		{"HEAD", "/ping", func(r fiberopenapi.Router) SingleRouteFunc { return r.Head }},
		{"POST", "/ping", func(r fiberopenapi.Router) SingleRouteFunc { return r.Post }},
		{"PUT", "/ping", func(r fiberopenapi.Router) SingleRouteFunc { return r.Put }},
		{"PATCH", "/ping", func(r fiberopenapi.Router) SingleRouteFunc { return r.Patch }},
		{"DELETE", "/ping", func(r fiberopenapi.Router) SingleRouteFunc { return r.Delete }},
		{"CONNECT", "/ping", func(r fiberopenapi.Router) SingleRouteFunc { return r.Connect }},
		{"OPTIONS", "/ping", func(r fiberopenapi.Router) SingleRouteFunc { return r.Options }},
		{"TRACE", "/ping", func(r fiberopenapi.Router) SingleRouteFunc { return r.Trace }},
	}
	for _, tt := range tests {
		t.Run(tt.method, func(t *testing.T) {
			app := fiber.New()
			r := fiberopenapi.NewRouter(app)

			routeFunc := tt.methodFunc(r)
			route := routeFunc(tt.path, func(c *fiber.Ctx) error {
				return c.SendString("pong")
			}).With(
				option.OperationID("ping"),
				option.Summary("Ping Endpoint"),
			).Name("ping")

			assert.NotNil(t, route, "expected route to be created for %s %s", tt.method, tt.path)
			fr := app.GetRoute("ping")
			assert.NotEmpty(t, fr.Name, "expected route name to be set for %s %s", tt.method, tt.path)

			req, _ := http.NewRequest(tt.method, tt.path, nil)
			res, err := app.Test(req, -1)
			require.NoError(t, err, "failed to test %s request", tt.method)
			assert.Equal(t, http.StatusOK, res.StatusCode, "expected status OK for %s request", tt.method)

			if tt.method != "HEAD" {
				var body []byte
				body, err = io.ReadAll(res.Body)
				require.NoError(t, err, "failed to read response body for %s request", tt.method)
				assert.Equal(t, "pong", string(body), "expected response body to be 'pong' for %s request", tt.method)
			}
			if tt.method == "CONNECT" {
				return // CONNECT method is not supported by OpenAPI, so we skip it
			}

			schema, err := r.GenerateSchema()
			require.NoError(t, err, "failed to generate OpenAPI schema for %s request", tt.method)
			assert.NotEmpty(t, schema, "expected non-empty OpenAPI schema for %s request", tt.method)

			// Check if the route is registered in the OpenAPI schema
			assert.Contains(
				t,
				string(schema),
				"operationId: ping",
				"expected operationId 'ping' in OpenAPI schema for %s request",
				tt.method,
			)
		})
	}

	t.Run("Static", func(t *testing.T) {
		app := fiber.New()
		r := fiberopenapi.NewRouter(app)
		r.Static("/static", "./testdata", fiber.Static{})
		req, _ := http.NewRequest(http.MethodGet, "/static/petstore.yaml", nil)
		res, err := app.Test(req, -1)
		require.NoError(t, err, "failed to test static file request")
		assert.Equal(t, http.StatusOK, res.StatusCode, "expected status OK for static file request")
	})
}

func TestRouter_Group(t *testing.T) {
	t.Run("Group", func(t *testing.T) {
		app := fiber.New()
		r := fiberopenapi.NewRouter(app)

		group := r.Group("/api", func(c *fiber.Ctx) error {
			return c.Next()
		})

		group.Get("/ping", func(c *fiber.Ctx) error {
			return c.SendString("pong")
		}).With(
			option.OperationID("ping"),
			option.Summary("Ping Endpoint"),
		)

		req, _ := http.NewRequest(http.MethodGet, "/api/ping", nil)
		res, err := app.Test(req, -1)
		require.NoError(t, err, "failed to test group route")
		assert.Equal(t, http.StatusOK, res.StatusCode, "expected status OK for group route")

		body, err := io.ReadAll(res.Body)
		require.NoError(t, err, "failed to read response body for group route")
		assert.Equal(t, "pong", string(body), "expected response body to be 'pong' for group route")
	})
	t.Run("Route", func(t *testing.T) {
		app := fiber.New()
		r := fiberopenapi.NewRouter(app)

		r.Route("/api", func(r fiberopenapi.Router) {
			r.Get("/ping", func(c *fiber.Ctx) error {
				return c.SendString("pong")
			}).With(
				option.OperationID("ping"),
				option.Summary("Ping Endpoint"),
			)
		})

		req, _ := http.NewRequest(http.MethodGet, "/api/ping", nil)
		res, err := app.Test(req, -1)
		require.NoError(t, err, "failed to test route")
		assert.Equal(t, http.StatusOK, res.StatusCode, "expected status OK for route")

		body, err := io.ReadAll(res.Body)
		require.NoError(t, err, "failed to read response body for route")
		assert.Equal(t, "pong", string(body), "expected response body to be 'pong' for route")
	})
}

func TestRouter_Middleware(t *testing.T) {
	t.Run("Use", func(t *testing.T) {
		called := false
		middleware := func(c *fiber.Ctx) error {
			called = true
			return c.Next()
		}
		app := fiber.New()
		r := fiberopenapi.NewRouter(app)
		r.Use(middleware)
		r.Get("/ping", func(c *fiber.Ctx) error {
			return c.SendString("pong")
		})

		req, _ := http.NewRequest(http.MethodGet, "/ping", nil)
		res, err := app.Test(req, -1)
		require.NoError(t, err, "failed to test middleware route")
		assert.Equal(t, http.StatusOK, res.StatusCode, "expected status OK for middleware route")
		assert.True(t, called, "expected middleware to be called")
	})
}

func TestGenerator_Docs(t *testing.T) {
	// Test that the docs route is registered
	app := fiber.New()
	r := fiberopenapi.NewRouter(app)
	r.Get("/ping", PingHandler).With(
		option.Summary("Ping Endpoint"),
	)

	t.Run("should serve docs", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/docs", nil)
		res, err := app.Test(req, -1)
		require.NoError(t, err, "failed to test docs route")
		assert.Equal(t, http.StatusOK, res.StatusCode, "expected status OK for docs route")

		body, err := io.ReadAll(res.Body)
		require.NoError(t, err, "failed to read response body for docs route")
		assert.Contains(
			t,
			string(body),
			"Fiber OpenAPI",
			"expected response body to contain 'Fiber OpenAPI' for docs route",
		)
	})
	t.Run("should serve OpenAPI YAML", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/docs/openapi.yaml", nil)
		res, err := app.Test(req, -1)
		require.NoError(t, err, "failed to test OpenAPI YAML route")
		assert.Equal(t, http.StatusOK, res.StatusCode, "expected status OK for OpenAPI YAML route")

		body, err := io.ReadAll(res.Body)
		require.NoError(t, err, "failed to read response body for OpenAPI YAML route")
		assert.NotEmpty(t, body, "expected non-empty response body for OpenAPI YAML route")
		assert.Contains(
			t,
			string(body),
			"openapi: 3.0.3",
			"expected OpenAPI version in response body for OpenAPI YAML route",
		)
	})
}

func TestGenerator_DisableDocs(t *testing.T) {
	pingHandler := func(c *fiber.Ctx) error {
		return c.SendString("pong")
	}
	app := fiber.New()
	r := fiberopenapi.NewRouter(app, option.WithDisableDocs())
	r.Get("/ping", pingHandler).With(
		option.Summary("Ping Endpoint"),
		option.Description("Endpoint to test ping functionality"),
	)

	t.Run("should not register docs route", func(t *testing.T) {
		reqDocs, _ := http.NewRequest(http.MethodGet, "/docs", nil)
		resDocs, err := app.Test(reqDocs, -1)
		require.NoError(t, err, "failed to test docs route")
		assert.Equal(t, http.StatusNotFound, resDocs.StatusCode, "expected status Not Found for docs route")
		_ = resDocs.Body.Close()
	})
	t.Run("should not register openapi.yaml route", func(t *testing.T) {
		reqOpenAPI, _ := http.NewRequest(http.MethodGet, "/docs/openapi.yaml", nil)
		resOpenAPI, err := app.Test(reqOpenAPI, -1)
		require.NoError(t, err, "failed to test OpenAPI YAML route")
		assert.Equal(t, http.StatusNotFound, resOpenAPI.StatusCode, "expected status Not Found for OpenAPI YAML route")
		_ = resOpenAPI.Body.Close()
	})
}

func TestGenerator_WriteSchemaTo(t *testing.T) {
	app := fiber.New()
	r := fiberopenapi.NewGenerator(app,
		option.WithTitle("Test API Write Schema"),
		option.WithVersion("1.0.0"),
		option.WithDescription("This is a test API for writing OpenAPI schema to file"),
	)

	r.Get("/ping", PingHandler).With(
		option.Summary("Ping Endpoint"),
		option.Description("Endpoint to test ping functionality"),
	)

	err := r.Validate()
	require.NoError(t, err, "failed to validate OpenAPI configuration")

	tempFile, err := os.CreateTemp(t.TempDir(), "openapi-schema-*.yaml")
	require.NoError(t, err, "failed to create temporary file for OpenAPI schema")
	defer func() {
		err = os.Remove(tempFile.Name())
		require.NoError(t, err, "failed to remove temporary file")
	}()

	err = r.WriteSchemaTo(tempFile.Name())
	require.NoError(t, err, "failed to write OpenAPI schema to file")

	schema, err := os.ReadFile(tempFile.Name())
	require.NoError(t, err, "failed to read OpenAPI schema from file")
	assert.NotEmpty(t, schema, "expected non-empty OpenAPI schema")
}

func TestGenerator_MarshalYAML(t *testing.T) {
	app := fiber.New()
	r := fiberopenapi.NewRouter(app,
		option.WithTitle("Test API Marshall YAML"),
		option.WithVersion("1.0.0"),
		option.WithDescription("This is a test API for marshalling OpenAPI schema to YAML"),
	)

	r.Get("/ping", PingHandler).With(
		option.Summary("Ping Endpoint"),
		option.Description("Endpoint to test ping functionality"),
	)

	err := r.Validate()
	require.NoError(t, err, "failed to validate OpenAPI configuration")

	yamlData, err := r.MarshalYAML()
	require.NoError(t, err, "failed to marshal OpenAPI schema to YAML")
	assert.NotEmpty(t, yamlData, "expected non-empty YAML data")
}

func TestGeneratorMarshalJSON(t *testing.T) {
	app := fiber.New()
	r := fiberopenapi.NewRouter(app,
		option.WithTitle("Test API Marshall JSON"),
		option.WithVersion("1.0.0"),
		option.WithDescription("This is a test API for marshalling OpenAPI schema to JSON"),
	)

	r.Get("/ping", PingHandler).With(
		option.Summary("Ping Endpoint"),
		option.Description("Endpoint to test ping functionality"),
	)

	err := r.Validate()
	require.NoError(t, err, "failed to validate OpenAPI configuration")

	jsonData, err := r.MarshalJSON()
	require.NoError(t, err, "failed to marshal OpenAPI schema to JSON")
	assert.NotEmpty(t, jsonData, "expected non-empty JSON data")
}

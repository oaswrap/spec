package fiberopenapi_test

import (
	"flag"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
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
				pet.Post("/{petId}/uploadImage", nil).With(
					option.OperationID("uploadFile"),
					option.Summary("Upload an image for a pet"),
					option.Description("Uploads an image for a pet."),
					option.Request(new(dto.UploadImageRequest)),
					option.Response(200, new(dto.ApiResponse)),
				)
				pet.Get("/{petId}", nil).With(
					option.OperationID("getPetById"),
					option.Summary("Get pet by ID"),
					option.Description("Retrieve a pet by its ID."),
					option.Request(new(struct {
						ID int `path:"petId" required:"true"`
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
				store.Get("/order/{orderId}", nil).With(
					option.OperationID("getOrderById"),
					option.Summary("Get order by ID"),
					option.Description("Retrieve an order by its ID."),
					option.Request(new(struct {
						ID int `path:"orderId" required:"true"`
					})),
					option.Response(200, new(dto.Order)),
					option.Response(404, nil),
				)
				store.Delete("/order/{orderId}", nil).With(
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
				user.Get("/{username}", nil).With(
					option.OperationID("getUserByName"),
					option.Summary("Get user by username"),
					option.Description("Retrieve a user by their username."),
					option.Request(new(struct {
						Username string `path:"username" required:"true"`
					})),
					option.Response(200, new(dto.PetUser)),
					option.Response(404, nil),
				)
				user.Put("/{username}", nil).With(
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
				user.Delete("/{username}", nil).With(
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
				assert.Error(t, err, "expected error for invalid OpenAPI configuration")
				return
			}
			err := r.Validate()
			assert.NoError(t, err, "failed to validate OpenAPI configuration")

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

func TestRouter_Fiber(t *testing.T) {
	totalCalled := 0
	testMiddleware := func(c *fiber.Ctx) error {
		totalCalled++
		return c.Next()
	}
	pingHandler := func(c *fiber.Ctx) error {
		return c.SendString("pong")
	}

	app := fiber.New()
	r := fiberopenapi.NewRouter(app)
	r.Use(testMiddleware)
	r.Get("/ping", pingHandler).With(
		option.Summary("Ping Endpoint"),
		option.Description("Endpoint to test ping functionality"),
	).Name("ping.get")
	r.Post("/ping", pingHandler).With(
		option.Summary("Ping Endpoint with POST"),
		option.Description("Endpoint to test ping functionality with POST method"),
	).Name("ping.post")
	r.Put("/ping", pingHandler).With(
		option.Summary("Ping Endpoint with PUT"),
		option.Description("Endpoint to test ping functionality with PUT method"),
	).Name("ping.put")
	r.Patch("/ping", pingHandler).With(
		option.Summary("Ping Endpoint with PATCH"),
		option.Description("Endpoint to test ping functionality with PATCH method"),
	).Name("ping.patch")
	r.Delete("/ping", pingHandler).With(
		option.Summary("Ping Endpoint with DELETE"),
		option.Description("Endpoint to test ping functionality with DELETE method"),
	).Name("ping.delete")
	r.Head("/ping", pingHandler).With(
		option.Summary("Ping Endpoint with HEAD"),
		option.Description("Endpoint to test ping functionality with HEAD method"),
	).Name("ping.head")
	r.Options("/ping", pingHandler).With(
		option.Summary("Ping Endpoint with OPTIONS"),
		option.Description("Endpoint to test ping functionality with OPTIONS method"),
	).Name("ping.options")
	r.Connect("/ping", pingHandler).With(
		option.Summary("Ping Endpoint with CONNECT"),
		option.Description("Endpoint to test ping functionality with CONNECT method"),
	).Name("ping.connect")
	r.Trace("/ping", pingHandler).With(
		option.Summary("Ping Endpoint with TRACE"),
		option.Description("Endpoint to test ping functionality with TRACE method"),
	).Name("ping.trace")
	err := r.Validate()
	require.NoError(t, err, "failed to validate OpenAPI configuration")

	methods := []string{
		"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS", "CONNECT", "TRACE",
	}
	totalCalled = 0 // Reset totalCalled for each method test
	for i, method := range methods {
		t.Run(method, func(t *testing.T) {
			r := app.GetRoute("ping." + strings.ToLower(method))
			assert.NotNil(t, r, "expected route to be registered for %s request", method)

			assert.Equal(t, i, totalCalled, "expected middleware to be called %d times, got %d", i, totalCalled)

			req, _ := http.NewRequest(method, "/ping", nil)
			res, err := app.Test(req, -1)
			require.NoError(t, err, "failed to test %s request", method)
			assert.Equal(t, http.StatusOK, res.StatusCode, "expected status OK for %s request", method)

			if method == "HEAD" {
				return // HEAD requests do not have a body
			}
			body, err := io.ReadAll(res.Body)
			require.NoError(t, err, "failed to read response body for %s request", method)
			assert.Equal(t, "pong", string(body), "expected response body to be 'pong' for %s request", method)

			_ = res.Body.Close()
		})
	}

	t.Run("Static File Request", func(t *testing.T) {
		r.Static("/static", "./testdata", fiber.Static{})
		req, _ := http.NewRequest("GET", "/static/petstore.yaml", nil)
		res, err := app.Test(req, -1)
		require.NoError(t, err, "failed to test static file request")
		assert.Equal(t, http.StatusOK, res.StatusCode, "expected status OK for static file request")
	})

	t.Run("must register docs route", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/docs", nil)
		res, err := app.Test(req, -1)
		require.NoError(t, err, "failed to test docs route")
		assert.Equal(t, http.StatusOK, res.StatusCode, "expected status OK for docs route")

		body, err := io.ReadAll(res.Body)
		require.NoError(t, err, "failed to read response body for docs route")
		assert.NotEmpty(t, body, "expected non-empty response body for docs route")
		_ = res.Body.Close()
	})
	t.Run("must register OpenAPI YAML route", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/docs/openapi.yaml", nil)
		res, err := app.Test(req, -1)
		require.NoError(t, err, "failed to test OpenAPI YAML route")
		assert.Equal(t, http.StatusOK, res.StatusCode, "expected status OK for OpenAPI YAML route")

		body, err := io.ReadAll(res.Body)
		require.NoError(t, err, "failed to read response body for OpenAPI YAML route")
		assert.NotEmpty(t, body, "expected non-empty response body for OpenAPI YAML route")
		_ = res.Body.Close()
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
	err := r.Validate()
	require.NoError(t, err, "failed to validate OpenAPI configuration")
	req, _ := http.NewRequest("GET", "/ping", nil)
	res, err := app.Test(req, -1)
	require.NoError(t, err, "failed to test ping route")
	assert.Equal(t, http.StatusOK, res.StatusCode, "expected status OK for ping route")

	_ = res.Body.Close()

	// Ensure OpenAPI routes are not registered
	reqDocs, _ := http.NewRequest("GET", "/docs", nil)
	resDocs, err := app.Test(reqDocs, -1)
	require.NoError(t, err, "failed to test docs route")
	assert.Equal(t, http.StatusNotFound, resDocs.StatusCode, "expected status Not Found for docs route")
	_ = resDocs.Body.Close()

	// Ensure OpenAPI YAML route is not registered
	reqOpenAPI, _ := http.NewRequest("GET", "/docs/openapi.yaml", nil)
	resOpenAPI, err := app.Test(reqOpenAPI, -1)
	require.NoError(t, err, "failed to test OpenAPI YAML route")
	assert.Equal(t, http.StatusNotFound, resOpenAPI.StatusCode, "expected status Not Found for OpenAPI YAML route")
	_ = resOpenAPI.Body.Close()
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

func TestGenerator_MarshallYAML(t *testing.T) {
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

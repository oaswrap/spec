package fiberopenapi_test

import (
	"bytes"
	"flag"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	openapiwrapper "github.com/faizlabs/openapi-wrapper"
	"github.com/faizlabs/openapi-wrapper/fiberopenapi"
	"github.com/faizlabs/openapi-wrapper/internal/dto"
	"github.com/faizlabs/openapi-wrapper/internal/testutil"
	"github.com/faizlabs/openapi-wrapper/option"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var update = flag.Bool("update", false, "update golden files")

func PingHandler(c *fiber.Ctx) error {
	return c.SendString("pong")
}

func TestRouter(t *testing.T) {
	tests := []struct {
		name    string
		golden  string
		method  string
		path    string
		options []option.OpenAPIOption
		setup   func(r fiberopenapi.Router)
	}{
		{
			name:   "basic_data_types",
			golden: "basic_data_types.yaml",
			method: "POST",
			path:   "/data-types",
			setup: func(r fiberopenapi.Router) {
				r.Post("/data-types", PingHandler).With(
					option.Summary("All Basic Data Types"),
					option.Description("Endpoint to test all basic data types"),
					option.Request(new(dto.AllBasicDataTypes)),
					option.Response(200, new(dto.AllBasicDataTypes)),
				)
			},
		},
		{
			name:   "basic_data_types_pointers",
			golden: "basic_data_types_pointers.yaml",
			method: "PUT",
			path:   "/data-types-pointers",
			setup: func(r fiberopenapi.Router) {
				r.Put("/data-types-pointers", PingHandler).With(
					option.Summary("All Basic Data Types Pointers"),
					option.Description("Endpoint to test all basic data types with pointers"),
					option.Request(new(dto.AllBasicDataTypesPointers)),
					option.Response(200, new(dto.AllBasicDataTypesPointers)),
				)
			},
		},
		{
			name:   "auth_login",
			golden: "auth_login.yaml",
			method: "POST",
			path:   "/auth/login",
			setup: func(r fiberopenapi.Router) {
				r.Post("/auth/login", PingHandler).With(
					option.Summary("User Login"),
					option.Description("Endpoint for user login"),
					option.Request(new(dto.LoginRequest)),
					option.Response(200, new(dto.Response[dto.Token])),
					option.Response(400, new(dto.ErrorResponse)),
					option.Response(422, new(dto.ValidationResponse)),
				)
			},
		},
		{
			name:   "auth_profile",
			golden: "auth_profile.yaml",
			method: "GET",
			path:   "/auth/me",
			options: []option.OpenAPIOption{
				option.WithSecurity("bearerAuth", option.SecurityHTTPBearer()),
			},
			setup: func(r fiberopenapi.Router) {
				r.Get("/auth/me", PingHandler).With(
					option.Summary("Get User Profile"),
					option.Description("Endpoint to get the authenticated user's profile"),
					option.Security("bearerAuth"),
					option.Response(200, new(dto.Response[dto.UserProfile])),
					option.Response(401, new(dto.ErrorResponse)),
				)
			},
		},
		{
			name:   "petstore",
			golden: "petstore.yaml",
			options: []option.OpenAPIOption{
				option.WithTitle("Pet Store API - OpenAPI 3.1"),
				option.WithVersion("1.0.0"),
				option.WithDescription("This is a sample Pet Store API using OpenAPI 3.1"),
				option.WithDocsPath("/docs"),
				option.WithServer("https://petstore3.swagger.io", "Pet Store Server"),
				option.WithSecurity("petstore_auth", option.SecurityOAuth2(
					openapiwrapper.OAuthFlows{
						Implicit: &openapiwrapper.OAuthFlowsDefsImplicit{
							AuthorizationURL: "https://petstore3.swagger.io/oauth/authorize",
							Scopes: map[string]string{
								"write:pets": "modify pets in your account",
								"read:pets":  "read your pets",
							},
						},
					},
				)),
			},
			setup: func(r fiberopenapi.Router) {
				api := r.Group("/api")
				v3 := api.Group("/v3")
				v3.Route("/pet", func(r fiberopenapi.Router) {
					r.Get("/findByStatus", nil).With(
						option.Summary("Finds Pets by status."),
						option.Description("Multiple status values can be provided with comma separated strings"),
						option.Request(new(dto.FindPetsByStatusRequest)),
						option.Response(200, new([]dto.Pet)),
					)
					r.Get("/findByTags", nil).With(
						option.Summary("Finds Pets by tags."),
						option.Description("Multiple tags can be provided with comma separated strings. Use tag1, tag2, tag3 for testing."),
						option.Request(new(dto.FindPetsByTagsRequest)),
						option.Response(200, new([]dto.Pet)),
					)
					r.Get("/:petId", nil).With(
						option.Summary("Find a pet by ID."),
						option.Description("Returns a single pet."),
						option.Request(new(dto.FindPetByIdRequest)),
						option.Response(200, new(dto.Pet)),
					)
					r.Post("/:petId", nil).With(
						option.Summary("Updates a pet in the store with form data."),
						option.Description("Update a pet resource based on form data."),
						option.Request(new(dto.UpdatePetFormDataRequest)),
						option.Response(200, new(dto.Pet)),
					)
					r.Delete("/:petId", nil).With(
						option.Summary("Deletes a pet."),
						option.Request(new(dto.DeletePetRequest)),
					)
					r.Post("/:petId/uploadImage", nil).With(
						option.Summary("Uploads an image."),
						option.Description("Uploads image of the pet."),
						option.Request(new(dto.UploadImageRequest)),
						option.Response(200, new(dto.ApiResponse)),
					)
					r.Post("/", nil).With(
						option.Summary("Add a new pet to the store."),
						option.Request(new(dto.Pet)),
						option.Response(200, new(dto.Pet)),
					)
					r.Put("/", nil).With(
						option.Summary("Update an existing pet."),
						option.Description("Update an existing pet by Id."),
						option.Request(new(dto.Pet)),
						option.Response(200, new(dto.Pet)),
					)
				}).With(option.RouteTags("pet"), option.RouteSecurity("petstore_auth", "write:pets", "read:pets"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := fiber.New()
			opts := []option.OpenAPIOption{
				option.WithTitle("Test API " + tt.name),
				option.WithVersion("1.0.0"),
				option.WithDescription("This is a test API for " + tt.name),
			}
			if len(tt.options) > 0 {
				opts = append(opts, tt.options...)
			}
			r := fiberopenapi.NewRouter(app, opts...)

			tt.setup(r)

			// Test the route registration
			if tt.method != "" && tt.path != "" {
				req := httptest.NewRequest(tt.method, tt.path, nil)
				resp, err := app.Test(req)
				assert.NoError(t, err, "failed to test request %s %s", tt.method, tt.path)
				assert.Equal(t, 200, resp.StatusCode, "expected status code 200 for %s %s", tt.method, tt.path)
				var buffer bytes.Buffer
				_, err = buffer.ReadFrom(resp.Body)
				assert.NoError(t, err, "failed to read response body for %s %s", tt.method, tt.path)
				assert.NotEmpty(t, buffer.String(), "expected non-empty response body for %s %s", tt.method, tt.path)
			}

			// Test the OpenAPI schema generation
			schema, err := r.GenerateOpenAPISchema()

			require.NoError(t, err, "failed to generate OpenAPI schema")
			goldenFile := filepath.Join("testdata", tt.golden)

			if *update {
				err = os.WriteFile(goldenFile, schema, 0644)
				require.NoError(t, err, "failed to write golden file")
				t.Logf("Updated golden file: %s", goldenFile)
			}

			want, err := os.ReadFile(goldenFile)
			require.NoError(t, err, "failed to read golden file %s", goldenFile)

			testutil.EqualYAML(t, want, schema)
		})
	}
}

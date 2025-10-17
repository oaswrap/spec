package spec_test

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/oaswrap/spec"
	"github.com/oaswrap/spec/openapi"
	"github.com/oaswrap/spec/option"
	"github.com/oaswrap/spec/pkg/dto"
	"github.com/oaswrap/spec/pkg/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//nolint:gochecknoglobals // test flag for golden file updates
var update = flag.Bool("update", false, "update golden files")

type AllBasicDataTypes struct {
	Int     int     `json:"int"`
	Int8    int8    `json:"int8"`
	Int16   int16   `json:"int16"`
	Int32   int32   `json:"int32"`
	Int64   int64   `json:"int64"`
	Uint    uint    `json:"uint"`
	Uint8   uint8   `json:"uint8"`
	Uint16  uint16  `json:"uint16"`
	Uint32  uint32  `json:"uint32"`
	Uint64  uint64  `json:"uint64"`
	Float32 float32 `json:"float32"`
	Float64 float64 `json:"float64"`
	Byte    byte    `json:"byte"`
	Rune    rune    `json:"rune"`
	String  string  `json:"string"`
	Bool    bool    `json:"bool"`
}

type AllBasicDataTypesPointers struct {
	Int     *int     `json:"int"`
	Int8    *int8    `json:"int8"`
	Int16   *int16   `json:"int16"`
	Int32   *int32   `json:"int32"`
	Int64   *int64   `json:"int64"`
	Uint    *uint    `json:"uint"`
	Uint8   *uint8   `json:"uint8"`
	Uint16  *uint16  `json:"uint16"`
	Uint32  *uint32  `json:"uint32"`
	Uint64  *uint64  `json:"uint64"`
	Float32 *float32 `json:"float32"`
	Float64 *float64 `json:"float64"`
	Byte    *byte    `json:"byte"`
	Rune    *rune    `json:"rune"`
	String  *string  `json:"string"`
	Bool    *bool    `json:"bool"`
}

type LoginRequest struct {
	Username string `json:"username" example:"john_doe"    validate:"required"`
	Password string `json:"password" example:"password123" validate:"required"`
}

type Response[T any] struct {
	Status int `json:"status" example:"200"`
	Data   T   `json:"data"`
}

type Token struct {
	Token string `json:"token" example:"abc123"`
}

type NullString struct {
	String string
	Valid  bool
}
type NullTime struct {
	Time  time.Time
	Valid bool
}

type User struct {
	ID        int        `json:"id"`
	Username  string     `json:"username"`
	Email     NullString `json:"email"`
	Age       *int       `json:"age,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt NullTime   `json:"updated_at"`
}

type CustomParser struct {
	re *regexp.Regexp
}

func NewCustomParser() *CustomParser {
	return &CustomParser{
		re: regexp.MustCompile(`:([a-zA-Z_][a-zA-Z0-9_]*)`),
	}
}

func (p *CustomParser) Parse(path string) (string, error) {
	return p.re.ReplaceAllString(path, "{$1}"), nil
}

type ErrorCustomParser struct{}

func (p *ErrorCustomParser) Parse(path string) (string, error) {
	return "", fmt.Errorf("failed to parse path: %s", path)
}

func TestRouter(t *testing.T) {
	tests := []struct {
		name        string
		golden      string
		opts        []option.OpenAPIOption
		setup       func(r spec.Router)
		shouldError bool
	}{
		{
			name:   "Basic Data Types",
			golden: "basic_data_types",
			setup: func(r spec.Router) {
				r.Post("/basic-data-types",
					option.OperationID("getBasicDataTypes"),
					option.Summary("Get Basic Data Types"),
					option.Description("This operation returns all basic data types."),
					option.Request(new(AllBasicDataTypes)),
					option.Response(200, new(AllBasicDataTypes)),
				)
			},
		},
		{
			name:   "Basic Data Types Pointers",
			golden: "basic_data_types_pointers",
			setup: func(r spec.Router) {
				r.Put("/basic-data-types-pointers",
					option.OperationID("getBasicDataTypesPointers"),
					option.Summary("Get Basic Data Types Pointers"),
					option.Description("This operation returns all basic data types as pointers."),
					option.Request(new(AllBasicDataTypesPointers)),
					option.Response(200, new(AllBasicDataTypesPointers)),
				)
			},
		},
		{
			name:   "All methods",
			golden: "all_methods",
			setup: func(r spec.Router) {
				type UserDetailRequest struct {
					ID int `path:"id" validate:"required"`
				}
				r.Get("/user", option.OperationID("getUser"), option.Summary("Get User"))
				r.Post(
					"/user",
					option.OperationID("createUser"),
					option.Summary("Create User"),
					option.Response(201, new(string), option.ContentType("plain/text")),
				)
				r.Put(
					"/user/{id}",
					option.OperationID("updateUser"),
					option.Summary("Update User"),
					option.Request(new(UserDetailRequest)),
				)
				r.Patch(
					"/user/{id}",
					option.OperationID("patchUser"),
					option.Summary("Patch User"),
					option.Request(new(UserDetailRequest)),
				)
				r.Delete(
					"/user/{id}",
					option.OperationID("deleteUser"),
					option.Summary("Delete User"),
					option.Request(new(UserDetailRequest)),
				)
				r.Head(
					"/user/{id}",
					option.OperationID("headUser"),
					option.Summary("Head User"),
					option.Request(new(UserDetailRequest)),
				)
				r.Options("/user", option.OperationID("optionsUser"), option.Summary("Options User"))
				r.Trace(
					"/user/{id}",
					option.OperationID("traceUser"),
					option.Summary("Trace User"),
					option.Request(new(UserDetailRequest)),
				)
			},
		},
		{
			name:   "Generic Response",
			golden: "generic_response",
			opts: []option.OpenAPIOption{
				option.WithTags(openapi.Tag{
					Name:        "Authentication",
					Description: "Operations related to user authentication",
				}),
			},
			setup: func(r spec.Router) {
				r.Post("/login",
					option.OperationID("login"),
					option.Summary("User Login"),
					option.Description("This operation allows users to log in."),
					option.Tags("Authentication"),
					option.Request(new(LoginRequest)),
					option.Response(200, new(Response[Token])),
				)
			},
		},
		{
			name:   "Custom Type Mapping",
			golden: "custom_type_mapping",
			opts: []option.OpenAPIOption{
				option.WithSecurity("bearerAuth", option.SecurityHTTPBearer("Bearer")),
				option.WithReflectorConfig(
					option.TypeMapping(NullString{}, new(string)),
					option.TypeMapping(NullTime{}, new(time.Time)),
				),
			},
			setup: func(r spec.Router) {
				r.Get("/auth/me",
					option.OperationID("getUserProfile"),
					option.Summary("Get User Profile"),
					option.Description("This operation retrieves the authenticated user's profile."),
					option.Security("bearerAuth"),
					option.Request(new(User)),
					option.Response(200, new(User)),
				)
			},
		},
		{
			name:   "Pet Store",
			golden: "petstore",
			opts: []option.OpenAPIOption{option.WithTitle("Petstore API"),
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
			setup: func(r spec.Router) {
				pet := r.Group("/pet",
					option.GroupTags("pet"),
					option.GroupSecurity("petstore_auth", "write:pets", "read:pets"),
				)
				pet.Put("/",
					option.OperationID("updatePet"),
					option.Summary("Update an existing pet"),
					option.Description("Update the details of an existing pet in the store."),
					option.Request(new(dto.Pet)),
					option.Response(200, new(dto.Pet)),
				)
				pet.Post("/",
					option.OperationID("addPet"),
					option.Summary("Add a new pet"),
					option.Description("Add a new pet to the store."),
					option.Request(new(dto.Pet)),
					option.Response(201, new(dto.Pet)),
				)
				pet.Get(
					"/findByStatus",
					option.OperationID("findPetsByStatus"),
					option.Summary("Find pets by status"),
					option.Description(
						"Finds Pets by status. Multiple status values can be provided with comma separated strings.",
					),
					option.Request(new(struct {
						Status string `query:"status" enum:"available,pending,sold"`
					})),
					option.Response(200, new([]dto.Pet)),
				)
				pet.Get(
					"/findByTags",
					option.OperationID("findPetsByTags"),
					option.Summary("Find pets by tags"),
					option.Description(
						"Finds Pets by tags. Multiple tags can be provided with comma separated strings.",
					),
					option.Request(new(struct {
						Tags []string `query:"tags"`
					})),
					option.Response(200, new([]dto.Pet)),
				)
				pet.Post("/{petId}/uploadImage",
					option.OperationID("uploadFile"),
					option.Summary("Upload an image for a pet"),
					option.Description("Uploads an image for a pet."),
					option.Request(new(dto.UploadImageRequest)),
					option.Response(200, new(dto.APIResponse)),
				)
				pet.Get("/{petId}",
					option.OperationID("getPetById"),
					option.Summary("Get pet by ID"),
					option.Description("Retrieve a pet by its ID."),
					option.Request(new(struct {
						ID int `path:"petId" required:"true"`
					})),
					option.Response(200, new(dto.Pet)),
				)
				pet.Post("/{petId}",
					option.OperationID("updatePetWithForm"),
					option.Summary("Update pet with form"),
					option.Description("Updates a pet in the store with form data."),
					option.Request(new(dto.UpdatePetWithFormRequest)),
					option.Response(200, nil),
				)
				pet.Delete("/{petId}",
					option.OperationID("deletePet"),
					option.Summary("Delete a pet"),
					option.Description("Delete a pet from the store by its ID."),
					option.Request(new(dto.DeletePetRequest)),
					option.Response(204, nil),
				)
				store := r.Group("/store",
					option.GroupTags("store"),
				)
				store.Post("/order",
					option.OperationID("placeOrder"),
					option.Summary("Place an order"),
					option.Description("Place a new order for a pet."),
					option.Request(new(dto.Order)),
					option.Response(201, new(dto.Order)),
				)
				store.Get("/order/{orderId}",
					option.OperationID("getOrderById"),
					option.Summary("Get order by ID"),
					option.Description("Retrieve an order by its ID."),
					option.Request(new(struct {
						ID int `path:"orderId" required:"true"`
					})),
					option.Response(200, new(dto.Order)),
					option.Response(404, nil),
				)
				store.Delete("/order/{orderId}",
					option.OperationID("deleteOrder"),
					option.Summary("Delete an order"),
					option.Description("Delete an order by its ID."),
					option.Request(new(struct {
						ID int `path:"orderId" required:"true"`
					})),
					option.Response(204, nil),
				)

				user := r.Group("/user",
					option.GroupTags("user"),
				)
				user.Post("/createWithList",
					option.OperationID("createUsersWithList"),
					option.Summary("Create users with list"),
					option.Description("Create multiple users in the store with a list."),
					option.Request(new([]dto.PetUser)),
					option.Response(201, nil),
				)
				user.Post("/",
					option.OperationID("createUser"),
					option.Summary("Create a new user"),
					option.Description("Create a new user in the store."),
					option.Request(new(dto.PetUser)),
					option.Response(201, new(dto.PetUser)),
				)
				user.Get("/{username}",
					option.OperationID("getUserByName"),
					option.Summary("Get user by username"),
					option.Description("Retrieve a user by their username."),
					option.Request(new(struct {
						Username string `path:"username" required:"true"`
					})),
					option.Response(200, new(dto.PetUser)),
					option.Response(404, nil),
				)
				user.Put("/{username}",
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
				user.Delete("/{username}",
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
			name:   "All Operation Options",
			golden: "all_operation_options",
			opts: []option.OpenAPIOption{
				option.WithSecurity("apiKey", option.SecurityAPIKey("x-api-key", "header")),
			},
			setup: func(r spec.Router) {
				r.Post("/operation/options",
					option.OperationID("postOperationOptions"),
					option.Summary("Post Operation Options"),
					option.Description("This operation retrieves all operation options."),
					option.Security("apiKey"),
					option.Tags("Operation Options"),
					option.Deprecated(),
					option.Request(new(LoginRequest),
						option.ContentType("application/json"),
						option.ContentDescription("Request body for operation options"),
					),
					option.Response(200, new(Response[User]),
						option.ContentType("application/json"),
						option.ContentDescription("Response body for operation options"),
						option.ContentDefault(true),
					),
				)
			},
		},
		{
			name:   "Hide Operation",
			golden: "hide_operation",
			setup: func(r spec.Router) {
				r.Get("/hidden/operation",
					option.OperationID("hiddenOperation"),
					option.Summary("Hidden Operation"),
					option.Description("This operation is hidden and should not appear in the spec."),
					option.Hidden(),
					option.Request(new(LoginRequest)),
					option.Response(200, new(Response[User])),
				)
			},
		},
		{
			name:   "All Reflector Options",
			golden: "all_reflector_options",
			opts: []option.OpenAPIOption{
				option.WithReflectorConfig(
					option.InlineRefs(),
					option.RootRef(),
					option.RootNullable(),
					option.StripDefNamePrefix("Test", "Mock"),
					option.InterceptDefNameFunc(func(_ reflect.Type, defaultDefName string) string {
						return defaultDefName + "_Custom"
					}),
					option.InterceptPropFunc(func(_ openapi.InterceptPropParams) error {
						return nil
					}),
					option.RequiredPropByValidateTag(),
					option.InterceptSchemaFunc(func(_ openapi.InterceptSchemaParams) (bool, error) {
						return false, nil
					}),
					option.TypeMapping(NullString{}, new(string)),
					option.TypeMapping(NullTime{}, new(time.Time)),
				),
			},
			setup: func(r spec.Router) {
				r.Get("/reflector/options",
					option.OperationID("getReflectorOptions"),
					option.Summary("Get Reflector Options"),
					option.Description("This operation retrieves the OpenAPI reflector options."),
					option.Request(new(LoginRequest)),
					option.Response(200, new(Response[User])),
				)
			},
		},
		{
			name:   "Group Routes",
			golden: "group_routes",
			opts: []option.OpenAPIOption{
				option.WithSecurity("bearerAuth", option.SecurityHTTPBearer("Bearer")),
				option.WithReflectorConfig(
					option.TypeMapping(NullString{}, new(string)),
					option.TypeMapping(NullTime{}, new(time.Time)),
				),
			},
			setup: func(r spec.Router) {
				api := r.Group("/api")
				v1 := api.Group("/v1")
				v1.Route("/auth", func(r spec.Router) {
					r.Post("/login",
						option.Summary("User Login v1"),
						option.Request(new(LoginRequest)),
						option.Response(200, new(Token)),
					)
					r.Get("/me",
						option.Summary("Get Profile v1"),
						option.Response(200, new(User)),
					).With(option.Security("bearerAuth"))
				}, option.GroupDeprecated(), option.GroupTags("Authentication"))
				v1.Route("/profile", func(r spec.Router) {
					r.Get("/",
						option.Summary("Get Profile v1"),
						option.Response(200, new(User)),
					)
				}, option.GroupHidden())
				v2 := api.Group("/v2")
				v2.Route("/auth", func(r spec.Router) {
					r.Post("/login",
						option.Summary("User Login v2"),
						option.Request(new(LoginRequest)),
						option.Response(200, new(Token)),
					)
					auth := r.Group("/", option.GroupSecurity("bearerAuth"))
					auth.Get("/me",
						option.Summary("Get Profile v2"),
						option.Response(200, new(User)),
					).With(option.Tags("Profile"))
				}, option.GroupTags("Authentication"))
				v2.Route("/profile", func(r spec.Router) {
					r.Put("/",
						option.Summary("Update Profile v2"),
						option.Request(new(User)),
						option.Response(200, new(User)),
					)
				}, option.GroupSecurity("bearerAuth")).With(option.GroupTags("Profile"))
			},
		},
		{
			name:   "Custom Path Parser",
			golden: "custom_path_parser",
			opts: []option.OpenAPIOption{
				option.WithPathParser(NewCustomParser()),
			},
			setup: func(r spec.Router) {
				r.Get("/user/:id",
					option.OperationID("getUserById"),
					option.Summary("Get User by ID"),
					option.Description("This operation retrieves a user by ID."),
					option.Request(new(struct {
						ID int `path:"id" validate:"required"`
					})),
					option.Response(200, new(User)),
				)
			},
		},
		{
			name:   "Server Variables",
			golden: "server_variables",
			opts: []option.OpenAPIOption{
				option.WithServer("https://api.example.com/{version}",
					option.ServerDescription("Production Server"),
					option.ServerVariables(map[string]openapi.ServerVariable{
						"version": {
							Default:     "v1",
							Enum:        []string{"v1", "v2"},
							Description: "API version",
						},
					}),
				),
				option.WithServer("https://api.example.dev/{version}",
					option.ServerDescription("Development Server"),
					option.ServerVariables(map[string]openapi.ServerVariable{
						"version": {
							Default:     "v1",
							Enum:        []string{"v1", "v2"},
							Description: "API version",
						},
					}),
				),
			},
		},
		{
			name:   "Spec Information",
			golden: "spec_information",
			opts: []option.OpenAPIOption{
				option.WithContact(openapi.Contact{
					Name:  "Support Team",
					URL:   "https://support.example.com",
					Email: "support@example.com",
				}),
				option.WithLicense(openapi.License{
					Name: "MIT License",
					URL:  "https://opensource.org/licenses/MIT",
				}),
				option.WithExternalDocs("https://docs.example.com", "API Documentation"),
			},
		},
		{
			name:   "Mux Route",
			golden: "mux_route",
			opts: []option.OpenAPIOption{
				option.WithSecurity("bearerAuth", option.SecurityHTTPBearer("Bearer")),
			},
			setup: func(r spec.Router) {
				api := r.Group("api")
				v1 := api.Group("v1")
				v1.NewRoute().Method("POST").Path("/login").With(
					option.Summary("User Login v1"),
					option.Request(new(LoginRequest)),
					option.Response(200, new(Token)),
				)
				auth := v1.Group("/", option.GroupSecurity("bearerAuth"))
				auth.NewRoute().Method("GET").Path("/profile").With(
					option.Summary("Get Profile v1"),
					option.Response(200, new(User)),
				)
			},
		},
		{
			name: "Invalid OpenAPI Version",
			opts: []option.OpenAPIOption{
				option.WithOpenAPIVersion("2.0.0"), // Invalid version for OpenAPI 3.x
			},
			shouldError: true,
		},
		{
			name: "Invalid URL Path Parameter",
			setup: func(r spec.Router) {
				r.Get("/user/{id}",
					option.OperationID("getUserById"),
					option.Summary("Get User by ID"),
					option.Description("This operation retrieves a user by ID."),
					option.Request(new(struct {
						ID int `params:"id" validate:"required"`
					})),
				)
			},
			shouldError: true, // Invalid path parameter without a proper tag
		},
		{
			name: "Error Custom Path Parser",
			opts: []option.OpenAPIOption{
				option.WithPathParser(&ErrorCustomParser{}),
			},
			setup: func(r spec.Router) {
				r.Get("/user/:id",
					option.OperationID("getUserById"),
					option.Summary("Get User by ID"),
					option.Description("This operation retrieves a user by ID."),
					option.Request(new(struct {
						ID int `path:"id" validate:"required"`
					})),
					option.Response(200, new(User)),
				)
			},
			shouldError: true, // Custom parser should fail
		},
	}

	versions := map[string]string{
		"3.0.0": "3",
		"3.1.0": "31",
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for version, alias := range versions {
				opts := []option.OpenAPIOption{
					option.WithTitle("API Doc: " + tt.name),
					option.WithDescription("This is the API documentation for " + tt.name),
					option.WithOpenAPIVersion(version),
					option.WithVersion("1.0.0"),
					option.WithReflectorConfig(option.RequiredPropByValidateTag()),
				}
				if len(tt.opts) > 0 {
					opts = append(opts, tt.opts...)
				}
				r := spec.NewRouter(opts...)

				if tt.setup != nil {
					tt.setup(r)
				}

				if tt.shouldError {
					require.Error(t, r.Validate(), "Expected router to fail validation")
					continue
				}
				require.NoError(t, r.Validate(), "Router validation failed")

				schema, err := r.GenerateSchema("yaml")
				require.NoError(t, err)

				golden := fmt.Sprintf("%s_%s.yaml", tt.golden, alias)

				goldenFile := filepath.Join("testdata", golden)

				if *update {
					err = os.WriteFile(goldenFile, schema, 0644)
					require.NoError(t, err, "failed to write golden file")
					t.Logf("Updated golden file: %s", goldenFile)
				}

				want, err := os.ReadFile(goldenFile)
				require.NoError(t, err, "failed to read golden file %s", goldenFile)

				testutil.EqualYAML(t, want, schema)
			}
		})
	}
}

func TestRouter_GenerateSchema(t *testing.T) {
	tests := []struct {
		name        string
		formats     []string
		expectError bool
		errorMsg    string
	}{
		{
			name:    "Default format (YAML)",
			formats: nil,
		},
		{
			name:    "Explicit YAML format",
			formats: []string{"yaml"},
		},
		{
			name:    "JSON format",
			formats: []string{"json"},
		},
		{
			name:        "Unsupported format",
			formats:     []string{"xml"},
			expectError: true,
			errorMsg:    "unsupported format: xml, expected one of json, yaml, yml",
		},
		{
			name:        "Empty string format",
			formats:     []string{""},
			expectError: true,
			errorMsg:    "unsupported format: , expected one of json, yaml, yml",
		},
		{
			name:        "Invalid format",
			formats:     []string{"invalid"},
			expectError: true,
			errorMsg:    "unsupported format: invalid, expected one of json, yaml, ym",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := spec.NewRouter(
				option.WithOpenAPIVersion("3.1.0"),
				option.WithTitle("Test API"),
				option.WithVersion("1.0.0"),
			)

			// Add a simple operation to ensure we have some content
			r.Add("GET", "/test",
				option.OperationID("test"),
				option.Summary("Test operation"),
				option.Description("This is a test operation."),
			)

			schema, err := r.GenerateSchema(tt.formats...)

			if tt.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
				assert.Nil(t, schema)
				return
			}

			require.NoError(t, err)
			assert.NotNil(t, schema)
			assert.NotEmpty(t, schema)

			// Verify format-specific content
			if len(tt.formats) == 0 || tt.formats[0] == "yaml" {
				// YAML format should contain YAML-specific syntax
				assert.Contains(t, string(schema), "openapi:")
				assert.Contains(t, string(schema), "info:")
			} else if tt.formats[0] == "json" {
				// JSON format should be valid JSON with proper indentation
				assert.True(t, json.Valid(schema), "Generated JSON should be valid")
				assert.Contains(t, string(schema), "{\n  \"openapi\":")
				assert.Contains(t, string(schema), "\"info\":")
			}
		})
	}
}

func TestRouter_WriteSchemaTo(t *testing.T) {
	tests := []struct {
		name        string
		path        string
		expectError bool
		expectJSON  bool
	}{
		{
			name:       "Write YAML file",
			path:       "test_schema.yaml",
			expectJSON: false,
		},
		{
			name:       "Write JSON file",
			path:       "test_schema.json",
			expectJSON: true,
		},
		{
			name:        "Write file without extension (defaults to YAML)",
			path:        "test_schema",
			expectError: true,
		},
		{
			name:       "Write file with .yml extension (YAML)",
			path:       "test_schema.yml",
			expectJSON: false,
		},
		{
			name:        "Write to invalid path",
			path:        "/invalid/path/that/does/not/exist/test.yaml",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create router with test configuration
			r := spec.NewRouter(
				option.WithOpenAPIVersion("3.1.0"),
				option.WithTitle("Test API"),
				option.WithVersion("1.0.0"),
			)

			// Add a simple operation to ensure we have content
			r.Add("GET", "/test",
				option.OperationID("test"),
				option.Summary("Test operation"),
				option.Description("This is a test operation."),
			)

			// Construct full path
			var fullPath string
			if tt.expectError {
				fullPath = tt.path // Use invalid path as-is
			} else {
				fullPath = filepath.Join(t.TempDir(), tt.path)
			}

			// Write schema to file
			err := r.WriteSchemaTo(fullPath)

			if tt.expectError {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)

			// Verify file was created
			assert.FileExists(t, fullPath)

			// Read and verify file content
			content, err := os.ReadFile(fullPath)
			require.NoError(t, err)
			assert.NotEmpty(t, content)

			if tt.expectJSON {
				// Verify JSON format
				assert.True(t, json.Valid(content), "File content should be valid JSON")
				assert.Contains(t, string(content), "{\n  \"openapi\":")
				assert.Contains(t, string(content), "\"info\":")
			} else {
				// Verify YAML format
				assert.Contains(t, string(content), "openapi:")
				assert.Contains(t, string(content), "info:")
				// Ensure it's not JSON format
				assert.False(t, strings.HasPrefix(string(content), "{"))
			}
		})
	}
}

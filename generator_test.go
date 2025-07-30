package spec_test

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/oaswrap/spec"
	"github.com/oaswrap/spec/openapi"
	"github.com/oaswrap/spec/option"
	"github.com/oaswrap/spec/pkg/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
	Username string `json:"username" example:"john_doe" validate:"required"`
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

type UserProfile struct {
	ID        int        `json:"id"`
	Username  string     `json:"username"`
	Email     NullString `json:"email"`
	Age       *int       `json:"age,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt NullTime   `json:"updated_at"`
}

func TestGenerator(t *testing.T) {
	tests := []struct {
		name   string
		golden string
		opts   []option.OpenAPIOption
		setup  func(g *spec.Generator)
	}{
		{
			name:   "Basic Data Types",
			golden: "basic_data_types",
			setup: func(g *spec.Generator) {
				g.Post("/basic-data-types",
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
			setup: func(g *spec.Generator) {
				g.Put("/basic-data-types-pointers",
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
			setup: func(g *spec.Generator) {
				type UserDetailRequest struct {
					ID int `path:"id" validate:"required"`
				}
				g.Get("/user", option.OperationID("getUser"), option.Summary("Get User"))
				g.Post("/user", option.OperationID("createUser"), option.Summary("Create User"), option.Response(201, new(string), option.WithContentType("plain/text")))
				g.Put("/user/{id}", option.OperationID("updateUser"), option.Summary("Update User"), option.Request(new(UserDetailRequest)))
				g.Patch("/user/{id}", option.OperationID("patchUser"), option.Summary("Patch User"), option.Request(new(UserDetailRequest)))
				g.Delete("/user/{id}", option.OperationID("deleteUser"), option.Summary("Delete User"), option.Request(new(UserDetailRequest)))
				g.Head("/user/{id}", option.OperationID("headUser"), option.Summary("Head User"), option.Request(new(UserDetailRequest)))
				g.Options("/user", option.OperationID("optionsUser"), option.Summary("Options User"))
				g.Trace("/user/{id}", option.OperationID("traceUser"), option.Summary("Trace User"), option.Request(new(UserDetailRequest)))
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
			setup: func(g *spec.Generator) {
				g.Post("/login",
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
				option.WithTypeMapping(NullString{}, new(string)),
				option.WithTypeMapping(NullTime{}, new(time.Time)),
			},
			setup: func(g *spec.Generator) {
				g.Get("/auth/me",
					option.OperationID("getUserProfile"),
					option.Summary("Get User Profile"),
					option.Description("This operation retrieves the authenticated user's profile."),
					option.Security("bearerAuth"),
					option.Request(new(UserProfile)),
					option.Response(200, new(UserProfile)),
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
				}
				if len(tt.opts) > 0 {
					opts = append(opts, tt.opts...)
				}
				gen := spec.NewGenerator(opts...)

				if tt.setup != nil {
					tt.setup(gen)
				}

				assert.NoError(t, gen.Validate(), "Generator validation failed")

				schema, err := gen.GenerateSchema("yaml")
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

func TestGenerator_GenerateSchema(t *testing.T) {
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
			errorMsg:    "unsupported format: xml, only 'json' and 'yaml' are supported",
		},
		{
			name:        "Empty string format",
			formats:     []string{""},
			expectError: true,
			errorMsg:    "unsupported format: , only 'json' and 'yaml' are supported",
		},
		{
			name:        "Invalid format",
			formats:     []string{"invalid"},
			expectError: true,
			errorMsg:    "unsupported format: invalid, only 'json' and 'yaml' are supported",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gen := spec.NewGenerator(
				option.WithOpenAPIVersion("3.1.0"),
				option.WithTitle("Test API"),
				option.WithVersion("1.0.0"),
			)

			// Add a simple operation to ensure we have some content
			gen.Add("GET", "/test",
				option.OperationID("test"),
				option.Summary("Test operation"),
				option.Description("This is a test operation."),
			)

			schema, err := gen.GenerateSchema(tt.formats...)

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
				assert.Nil(t, schema)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, schema)
			assert.Greater(t, len(schema), 0, "Schema should not be empty")

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

func TestGenerator_WriteSchemaTo(t *testing.T) {
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
			name:       "Write file without extension (defaults to YAML)",
			path:       "test_schema",
			expectJSON: false,
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
			// Create generator with test configuration
			gen := spec.NewGenerator(
				option.WithOpenAPIVersion("3.1.0"),
				option.WithTitle("Test API"),
				option.WithVersion("1.0.0"),
			)

			// Add a simple operation to ensure we have content
			gen.Add("GET", "/test",
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
			err := gen.WriteSchemaTo(fullPath)

			if tt.expectError {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)

			// Verify file was created
			assert.FileExists(t, fullPath)

			// Read and verify file content
			content, err := os.ReadFile(fullPath)
			require.NoError(t, err)
			assert.Greater(t, len(content), 0, "File should not be empty")

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

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
	"github.com/oaswrap/spec/option"
	"github.com/oaswrap/spec/pkg/testutil"
	"github.com/oaswrap/spec/pkg/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/swaggest/openapi-go"
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
		setup  func(t *testing.T, g *spec.Generator)
	}{
		{
			name:   "Basic Data Types",
			golden: "basic_data_types",
			setup: func(t *testing.T, g *spec.Generator) {
				op, err := g.NewOperationContext("GET", "/basic-data-types")
				assert.NoError(t, err)
				op.SetID("getBasicDataTypes")
				op.SetSummary("Get Basic Data Types")
				op.SetDescription("This operation returns all basic data types.")
				op.AddReqStructure(new(AllBasicDataTypes))
				op.AddRespStructure(new(AllBasicDataTypes), openapi.WithHTTPStatus(200))
				assert.NoError(t, g.AddOperation(op))
			},
		},
		{
			name:   "Basic Data Types Pointers",
			golden: "basic_data_types_pointers",
			setup: func(t *testing.T, g *spec.Generator) {
				op, err := g.NewOperationContext("GET", "/basic-data-types-pointers")
				assert.NoError(t, err)
				op.SetID("getBasicDataTypesPointers")
				op.SetSummary("Get Basic Data Types Pointers")
				op.SetDescription("This operation returns all basic data types as pointers.")
				op.AddReqStructure(new(AllBasicDataTypes))
				op.AddRespStructure(new(AllBasicDataTypes), openapi.WithHTTPStatus(200))
				assert.NoError(t, g.AddOperation(op))
			},
		},
		{
			name:   "Login Request",
			golden: "login_request",
			setup: func(t *testing.T, g *spec.Generator) {
				op, err := g.NewOperationContext("POST", "/login")
				assert.NoError(t, err)
				op.SetID("login")
				op.SetSummary("User Login")
				op.SetDescription("This operation allows users to log in.")
				op.AddReqStructure(new(LoginRequest))
				op.AddRespStructure(new(Response[Token]), openapi.WithHTTPStatus(200))
				assert.NoError(t, g.AddOperation(op))
			},
		},
		{
			name:   "User Profile",
			golden: "user_profile",
			opts: []option.OpenAPIOption{
				option.WithSecurity("bearerAuth", option.SecurityHTTPBearer("Bearer")),
				option.WithTypeMapping(NullString{}, new(string)),
				option.WithTypeMapping(NullTime{}, new(time.Time)),
			},
			setup: func(t *testing.T, g *spec.Generator) {
				op, err := g.NewOperationContext("GET", "/auth/me")
				assert.NoError(t, err)
				op.SetID("getUserProfile")
				op.SetSummary("Get User Profile")
				op.SetDescription("This operation retrieves the authenticated user's profile.")
				op.AddSecurity("bearerAuth")
				op.AddRespStructure(new(UserProfile), openapi.WithHTTPStatus(200))
				assert.NoError(t, g.AddOperation(op))
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
					option.WithOpenAPIVersion(version),
					option.WithTitle("Test API"),
					option.WithVersion("1.0.0"),
				}
				if len(tt.opts) > 0 {
					opts = append(opts, tt.opts...)
				}
				gen, err := spec.NewGenerator(opts...)
				require.NoError(t, err)

				tt.setup(t, gen)

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

func TestGenerator_New(t *testing.T) {
	tests := []struct {
		name        string
		opts        []option.OpenAPIOption
		shouldError bool
		expected    *option.OpenAPI
	}{
		{
			name: "With Default Config",
			opts: []option.OpenAPIOption{},
		},
		{
			name: "With Custom Config",
			opts: []option.OpenAPIOption{
				option.WithOpenAPIVersion("3.1.0"),
				option.WithTitle("Custom API"),
				option.WithVersion("2.0.0"),
				option.WithDescription("This is a custom API documentation."),
				option.WithBaseURL("https://api.example.com"),
			},
		},
		{
			name: "With Servers 3.0.0",
			opts: []option.OpenAPIOption{
				option.WithOpenAPIVersion("3.0.0"),
				option.WithServer("https://api.example.com/v1",
					option.ServerDescription("Production Server"),
					option.ServerVariables(map[string]option.ServerVariable{
						"version": {
							Default:     "v1",
							Enum:        []string{"v1", "v2"},
							Description: util.PtrOf("API version"),
						},
					}),
				),
			},
		},
		{
			name: "With Servers 3.1.0",
			opts: []option.OpenAPIOption{
				option.WithOpenAPIVersion("3.1.0"),
				option.WithServer("https://api.example.com/v1",
					option.ServerDescription("Production Server"),
					option.ServerVariables(map[string]option.ServerVariable{
						"version": {
							Default:     "v1",
							Enum:        []string{"v1", "v2"},
							Description: util.PtrOf("API version"),
						},
					}),
				),
			},
		},
		{
			name: "With Security Schemes ApiKey 3.0.0",
			opts: []option.OpenAPIOption{
				option.WithOpenAPIVersion("3.0.0"),
				option.WithSecurity("apiKey", option.SecurityAPIKey("x-api-key", "header")),
			},
		},
		{
			name: "With Security Schemes ApiKey 3.1.0",
			opts: []option.OpenAPIOption{
				option.WithOpenAPIVersion("3.1.0"),
				option.WithSecurity("apiKey", option.SecurityAPIKey("x-api-key", "header")),
			},
		},
		{
			name: "With Security Schemes Bearer 3.0.0",
			opts: []option.OpenAPIOption{
				option.WithOpenAPIVersion("3.0.0"),
				option.WithSecurity("bearer", option.SecurityHTTPBearer("Bearer")),
			},
		},
		{
			name: "With Security Schemes Bearer 3.1.0",
			opts: []option.OpenAPIOption{
				option.WithOpenAPIVersion("3.1.0"),
				option.WithSecurity("bearer", option.SecurityHTTPBearer("Bearer")),
			},
		},
		{
			name: "With Security Schemes OAuth2 3.0.0",
			opts: []option.OpenAPIOption{
				option.WithOpenAPIVersion("3.0.0"),
				option.WithSecurity("oauth2", option.SecurityOAuth2(option.OAuthFlows{
					Implicit: &option.OAuthFlowsDefsImplicit{
						AuthorizationURL: "https://auth.example.com/oauth/authorize",
						Scopes: map[string]string{
							"read":  "Read access",
							"write": "Write access",
						},
					},
					Password: &option.OAuthFlowsDefsPassword{
						TokenURL: "https://auth.example.com/oauth/token",
						Scopes: map[string]string{
							"read":  "Read access",
							"write": "Write access",
						},
					},
					ClientCredentials: &option.OAuthFlowsDefsClientCredentials{
						TokenURL: "https://auth.example.com/oauth/token",
						Scopes: map[string]string{
							"read":  "Read access",
							"write": "Write access",
						},
					},
					AuthorizationCode: &option.OAuthFlowsDefsAuthorizationCode{
						AuthorizationURL: "https://auth.example.com/oauth/authorize",
						TokenURL:         "https://auth.example.com/oauth/token",
						Scopes: map[string]string{
							"read":  "Read access",
							"write": "Write access",
						},
					},
				})),
			},
		},
		{
			name: "With Security Schemes OAuth2 3.1.0",
			opts: []option.OpenAPIOption{
				option.WithOpenAPIVersion("3.1.0"),
				option.WithSecurity("oauth2", option.SecurityOAuth2(option.OAuthFlows{
					Implicit: &option.OAuthFlowsDefsImplicit{
						AuthorizationURL: "https://auth.example.com/oauth/authorize",
						Scopes: map[string]string{
							"read":  "Read access",
							"write": "Write access",
						},
					},
					Password: &option.OAuthFlowsDefsPassword{
						TokenURL: "https://auth.example.com/oauth/token",
						Scopes: map[string]string{
							"read":  "Read access",
							"write": "Write access",
						},
					},
					ClientCredentials: &option.OAuthFlowsDefsClientCredentials{
						TokenURL: "https://auth.example.com/oauth/token",
						Scopes: map[string]string{
							"read":  "Read access",
							"write": "Write access",
						},
					},
					AuthorizationCode: &option.OAuthFlowsDefsAuthorizationCode{
						AuthorizationURL: "https://auth.example.com/oauth/authorize",
						TokenURL:         "https://auth.example.com/oauth/token",
						Scopes: map[string]string{
							"read":  "Read access",
							"write": "Write access",
						},
					},
				})),
			},
		},
		{
			name: "With Invalid OpenAPI Version",
			opts: []option.OpenAPIOption{
				option.WithOpenAPIVersion("4.0.0"), // Invalid version
			},
			shouldError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := spec.NewGenerator(tt.opts...)
			if tt.shouldError {
				assert.Error(t, err, "expected error but got none")
				return
			}
			assert.NoError(t, err, "unexpected error creating generator")
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
			gen, err := spec.NewGenerator(
				option.WithOpenAPIVersion("3.1.0"),
				option.WithTitle("Test API"),
				option.WithVersion("1.0.0"),
			)
			require.NoError(t, err)

			// Add a simple operation to ensure we have some content
			op, err := gen.NewOperationContext("GET", "/test")
			require.NoError(t, err)
			op.SetID("test")
			op.SetSummary("Test operation")
			require.NoError(t, gen.AddOperation(op))

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
			gen, err := spec.NewGenerator(
				option.WithOpenAPIVersion("3.1.0"),
				option.WithTitle("Test API"),
				option.WithVersion("1.0.0"),
			)
			require.NoError(t, err)

			// Add a simple operation to ensure we have content
			op, err := gen.NewOperationContext("GET", "/test")
			require.NoError(t, err)
			op.SetID("test")
			op.SetSummary("Test operation")
			require.NoError(t, gen.AddOperation(op))

			// Construct full path
			var fullPath string
			if tt.expectError {
				fullPath = tt.path // Use invalid path as-is
			} else {
				fullPath = filepath.Join(t.TempDir(), tt.path)
			}

			// Write schema to file
			err = gen.WriteSchemaTo(fullPath)

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

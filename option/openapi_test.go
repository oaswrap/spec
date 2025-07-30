package option_test

import (
	"log"
	"testing"

	"github.com/oaswrap/spec/option"
	"github.com/oaswrap/spec/pkg/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWithOpenAPIVersion(t *testing.T) {
	config := &option.OpenAPI{}
	opt := option.WithOpenAPIVersion("3.0.0")
	opt(config)

	assert.Equal(t, "3.0.0", config.OpenAPIVersion)
}

func TestWithDisableOpenAPI(t *testing.T) {
	tests := []struct {
		name     string
		disable  []bool
		expected bool
	}{
		{"default true", []bool{}, true},
		{"explicit true", []bool{true}, true},
		{"explicit false", []bool{false}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &option.OpenAPI{}
			opt := option.WithDisableOpenAPI(tt.disable...)
			opt(config)

			assert.Equal(t, tt.expected, config.DisableOpenAPI)
		})
	}
}

func TestWithBaseURL(t *testing.T) {
	config := &option.OpenAPI{}
	opt := option.WithBaseURL("https://api.example.com")
	opt(config)

	assert.Equal(t, "https://api.example.com", config.BaseURL)
}

func TestWithTitle(t *testing.T) {
	config := &option.OpenAPI{}
	opt := option.WithTitle("My API")
	opt(config)

	assert.Equal(t, "My API", config.Title)
}

func TestWithVersion(t *testing.T) {
	config := &option.OpenAPI{}
	opt := option.WithVersion("1.0.0")
	opt(config)

	assert.Equal(t, "1.0.0", config.Version)
}

func TestWithDescription(t *testing.T) {
	config := &option.OpenAPI{}
	opt := option.WithDescription("API description")
	opt(config)

	require.NotNil(t, config.Description)
	assert.Equal(t, "API description", *config.Description)
}

func TestWithServer(t *testing.T) {
	tests := []struct {
		name     string
		url      string
		opts     []option.ServerOption
		expected option.Server
	}{
		{
			name: "without description",
			url:  "https://api.example.com",
			expected: option.Server{
				URL: "https://api.example.com",
			},
		},
		{
			name: "with description",
			url:  "https://api.example.com",
			opts: []option.ServerOption{option.ServerDescription("Production server")},
			expected: option.Server{
				URL:         "https://api.example.com",
				Description: util.PtrOf("Production server"),
			},
		},
		{
			name: "with variables",
			url:  "https://api.example.com",
			opts: []option.ServerOption{
				option.ServerVariables(map[string]option.ServerVariable{
					"version": {
						Default:     "v1",
						Description: util.PtrOf("API version"),
						Enum:        []string{"v1", "v2"},
					},
				}),
			},
			expected: option.Server{
				URL: "https://api.example.com",
				Variables: map[string]option.ServerVariable{
					"version": {
						Default:     "v1",
						Description: util.PtrOf("API version"),
						Enum:        []string{"v1", "v2"},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &option.OpenAPI{}
			opt := option.WithServer(tt.url, tt.opts...)
			opt(config)

			require.Len(t, config.Servers, 1)
			assert.Equal(t, tt.expected.URL, config.Servers[0].URL)
			if tt.expected.Description != nil {
				require.NotNil(t, config.Servers[0].Description)
				assert.Equal(t, *tt.expected.Description, *config.Servers[0].Description)
			} else {
				assert.Nil(t, config.Servers[0].Description)
			}
		})
	}
}

func TestWithDocsPath(t *testing.T) {
	config := &option.OpenAPI{}
	opt := option.WithDocsPath("/docs")
	opt(config)

	assert.Equal(t, "/docs", config.DocsPath)
}

func TestWithSecurity(t *testing.T) {
	tests := []struct {
		name     string
		scheme   string
		opts     []option.SecurityOption
		expected *option.SecurityScheme
	}{
		{
			name:   "API Key Scheme",
			scheme: "apiKey",
			opts: []option.SecurityOption{
				option.SecurityAPIKey("x-api-key", "header"),
				option.SecurityDescription("API key for authentication"),
			},
			expected: &option.SecurityScheme{
				Description: util.PtrOf("API key for authentication"),
				APIKey: &option.SecuritySchemeAPIKey{
					Name: "x-api-key",
					In:   "header",
				},
			},
		},
		{
			name:   "HTTP Bearer Scheme",
			scheme: "bearerAuth",
			opts: []option.SecurityOption{
				option.SecurityHTTPBearer("Bearer"),
				option.SecurityDescription(""),
			},
			expected: &option.SecurityScheme{
				HTTPBearer: &option.SecuritySchemeHTTPBearer{
					Scheme: "Bearer",
				},
			},
		},
		{
			name:   "OAuth2 Scheme",
			scheme: "oauth2",
			opts: []option.SecurityOption{
				option.SecurityOAuth2(option.OAuthFlows{
					Implicit: &option.OAuthFlowsDefsImplicit{
						AuthorizationURL: "https://auth.example.com/authorize",
						Scopes: map[string]string{
							"read":  "Read access",
							"write": "Write access",
						},
					},
				}),
			},
			expected: &option.SecurityScheme{
				OAuth2: &option.SecuritySchemeOAuth2{
					Flows: option.OAuthFlows{
						Implicit: &option.OAuthFlowsDefsImplicit{
							AuthorizationURL: "https://auth.example.com/authorize",
							Scopes: map[string]string{
								"read":  "Read access",
								"write": "Write access",
							},
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &option.OpenAPI{}
			opt := option.WithSecurity(tt.scheme, tt.opts...)
			opt(config)

			require.NotNil(t, config.SecuritySchemes)
			require.Len(t, config.SecuritySchemes, 1)
			assert.Equal(t, tt.expected, config.SecuritySchemes[tt.scheme])
		})
	}
}

func TestWithSwaggerConfig(t *testing.T) {
	tests := []struct {
		name     string
		cfg      []option.SwaggerConfig
		expected *option.SwaggerConfig
	}{
		{
			name:     "no config",
			cfg:      []option.SwaggerConfig{},
			expected: nil,
		},
		{
			name:     "nil config",
			cfg:      []option.SwaggerConfig{},
			expected: nil,
		},
		{
			name: "valid config",
			cfg: []option.SwaggerConfig{
				{
					ShowTopBar: true,
					HideCurl:   false,
				},
			},
			expected: &option.SwaggerConfig{
				ShowTopBar: true,
				HideCurl:   false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &option.OpenAPI{}
			opt := option.WithSwaggerConfig(tt.cfg...)
			opt(config)

			assert.Equal(t, tt.expected, config.SwaggerConfig)
		})
	}
}

func TestWithDebug(t *testing.T) {
	tests := []struct {
		name      string
		debug     []bool
		expectLog bool
	}{
		{"default true", []bool{}, true},
		{"explicit true", []bool{true}, true},
		{"explicit false", []bool{false}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &option.OpenAPI{}
			opt := option.WithDebug(tt.debug...)
			opt(config)

			if tt.expectLog {
				assert.Equal(t, log.Default(), config.Logger)
			} else {
				assert.IsType(t, &option.NoopLogger{}, config.Logger)
			}
		})
	}
}

func TestWithContact(t *testing.T) {
	tests := []struct {
		name     string
		contact  option.Contact
		expected option.Contact
	}{
		{
			name: "full contact info",
			contact: option.Contact{
				Name:  "API Support",
				URL:   "https://example.com/support",
				Email: "support@example.com",
			},
			expected: option.Contact{
				Name:  "API Support",
				URL:   "https://example.com/support",
				Email: "support@example.com",
			},
		},
		{
			name: "minimal contact info",
			contact: option.Contact{
				Name: "Support Team",
			},
			expected: option.Contact{
				Name: "Support Team",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &option.OpenAPI{}
			opt := option.WithContact(tt.contact)
			opt(config)

			require.NotNil(t, config.Contact)
			assert.Equal(t, tt.expected, *config.Contact)
		})
	}
}

func TestWithLicense(t *testing.T) {
	tests := []struct {
		name     string
		license  option.License
		expected option.License
	}{
		{
			name: "license with URL",
			license: option.License{
				Name: "MIT",
				URL:  "https://opensource.org/licenses/MIT",
			},
			expected: option.License{
				Name: "MIT",
				URL:  "https://opensource.org/licenses/MIT",
			},
		},
		{
			name: "license without URL",
			license: option.License{
				Name: "Apache 2.0",
			},
			expected: option.License{
				Name: "Apache 2.0",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &option.OpenAPI{}
			opt := option.WithLicense(tt.license)
			opt(config)

			require.NotNil(t, config.License)
			assert.Equal(t, tt.expected, *config.License)
		})
	}
}

func TestWithExternalDocs(t *testing.T) {
	tests := []struct {
		name     string
		url      string
		desc     []string
		expected *option.ExternalDocumentation
	}{
		{
			name: "with description",
			url:  "https://example.com/docs",
			desc: []string{"External documentation"},
			expected: &option.ExternalDocumentation{
				URL:         "https://example.com/docs",
				Description: "External documentation",
			},
		},
		{
			name: "without description",
			url:  "https://example.com/docs",
			desc: []string{},
			expected: &option.ExternalDocumentation{
				URL: "https://example.com/docs",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &option.OpenAPI{}
			opt := option.WithExternalDocs(tt.url, tt.desc...)
			opt(config)

			require.NotNil(t, config.ExternalDocs)
			assert.Equal(t, tt.expected.URL, config.ExternalDocs.URL)
			assert.Equal(t, tt.expected.Description, config.ExternalDocs.Description)
		})
	}
}

func TestWithTags(t *testing.T) {
	tags := []option.Tag{
		{
			Name:        "users",
			Description: "User management",
		},
		{
			Name:        "orders",
			Description: "Order management",
		},
	}

	config := &option.OpenAPI{}
	opt := option.WithTags(tags...)
	opt(config)

	require.Len(t, config.Tags, 2)
	assert.Equal(t, "users", config.Tags[0].Name)
	assert.Equal(t, "User management", config.Tags[0].Description)
	assert.Equal(t, "orders", config.Tags[1].Name)
	assert.Equal(t, "Order management", config.Tags[1].Description)
}

func TestOpenAPIConfigDefaults(t *testing.T) {
	config := &option.OpenAPI{}

	// Test that default values are properly set
	assert.Empty(t, config.OpenAPIVersion)
	assert.False(t, config.DisableOpenAPI)
	assert.Empty(t, config.BaseURL)
	assert.Empty(t, config.Title)
	assert.Empty(t, config.Version)
	assert.Nil(t, config.Description)
	assert.Empty(t, config.Servers)
	assert.Empty(t, config.DocsPath)
	assert.Nil(t, config.SecuritySchemes)
	assert.Nil(t, config.SwaggerConfig)
	assert.Nil(t, config.Logger)
	assert.Nil(t, config.Contact)
	assert.Nil(t, config.License)
	assert.Empty(t, config.Tags)
}

package option_test

import (
	"testing"

	"github.com/oaswrap/spec/internal/util"

	"github.com/oaswrap/spec/openapi"
	"github.com/oaswrap/spec/option"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWithOpenAPIVersion(t *testing.T) {
	config := &openapi.Config{}
	opt := option.WithOpenAPIVersion("3.0.0")
	opt(config)

	assert.Equal(t, "3.0.0", config.OpenAPIVersion)
}

func TestWithDisableDocs(t *testing.T) {
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
			config := &openapi.Config{}
			opt := option.WithDisableDocs(tt.disable...)
			opt(config)

			assert.Equal(t, tt.expected, config.DisableDocs)
		})
	}
}

func TestWithBaseURL(t *testing.T) {
	config := &openapi.Config{}
	opt := option.WithBaseURL("https://api.example.com")
	opt(config)

	assert.Equal(t, "https://api.example.com", config.BaseURL)
}

func TestWithTitle(t *testing.T) {
	config := &openapi.Config{}
	opt := option.WithTitle("My API")
	opt(config)

	assert.Equal(t, "My API", config.Title)
}

func TestWithVersion(t *testing.T) {
	config := &openapi.Config{}
	opt := option.WithVersion("1.0.0")
	opt(config)

	assert.Equal(t, "1.0.0", config.Version)
}

func TestWithDescription(t *testing.T) {
	config := &openapi.Config{}
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
		expected openapi.Server
	}{
		{
			name: "without description",
			url:  "https://api.example.com",
			expected: openapi.Server{
				URL: "https://api.example.com",
			},
		},
		{
			name: "with description",
			url:  "https://api.example.com",
			opts: []option.ServerOption{option.ServerDescription("Production server")},
			expected: openapi.Server{
				URL:         "https://api.example.com",
				Description: util.PtrOf("Production server"),
			},
		},
		{
			name: "with variables",
			url:  "https://api.example.com",
			opts: []option.ServerOption{
				option.ServerVariables(map[string]openapi.ServerVariable{
					"version": {
						Default:     "v1",
						Description: "API version",
						Enum:        []string{"v1", "v2"},
					},
				}),
			},
			expected: openapi.Server{
				URL: "https://api.example.com",
				Variables: map[string]openapi.ServerVariable{
					"version": {
						Default:     "v1",
						Description: "API version",
						Enum:        []string{"v1", "v2"},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &openapi.Config{}
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
	config := &openapi.Config{}
	opt := option.WithDocsPath("/docs")
	opt(config)

	assert.Equal(t, "/docs", config.DocsPath)
}

func TestWithSecurity(t *testing.T) {
	tests := []struct {
		name     string
		scheme   string
		opts     []option.SecurityOption
		expected *openapi.SecurityScheme
	}{
		{
			name:   "API Key Scheme",
			scheme: "apiKey",
			opts: []option.SecurityOption{
				option.SecurityAPIKey("x-api-key", "header"),
				option.SecurityDescription("API key for authentication"),
			},
			expected: &openapi.SecurityScheme{
				Description: util.PtrOf("API key for authentication"),
				APIKey: &openapi.SecuritySchemeAPIKey{
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
			expected: &openapi.SecurityScheme{
				HTTPBearer: &openapi.SecuritySchemeHTTPBearer{
					Scheme: "Bearer",
				},
			},
		},
		{
			name:   "OAuth2 Scheme",
			scheme: "oauth2",
			opts: []option.SecurityOption{
				option.SecurityOAuth2(openapi.OAuthFlows{
					Implicit: &openapi.OAuthFlowsDefsImplicit{
						AuthorizationURL: "https://auth.example.com/authorize",
						Scopes: map[string]string{
							"read":  "Read access",
							"write": "Write access",
						},
					},
				}),
			},
			expected: &openapi.SecurityScheme{
				OAuth2: &openapi.SecuritySchemeOAuth2{
					Flows: openapi.OAuthFlows{
						Implicit: &openapi.OAuthFlowsDefsImplicit{
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
			config := &openapi.Config{}
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
		cfg      openapi.SwaggerConfig
		expected *openapi.SwaggerConfig
	}{
		{
			name:     "no config",
			cfg:      openapi.SwaggerConfig{},
			expected: &openapi.SwaggerConfig{},
		},
		{
			name: "valid config",
			cfg: openapi.SwaggerConfig{
				ShowTopBar: true,
				HideCurl:   false,
			},
			expected: &openapi.SwaggerConfig{
				ShowTopBar: true,
				HideCurl:   false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &openapi.Config{}
			opt := option.WithSwaggerConfig(tt.cfg)
			opt(config)

			assert.Equal(t, tt.expected, config.SwaggerConfig)
		})
	}
}

func TestWithContact(t *testing.T) {
	tests := []struct {
		name     string
		contact  openapi.Contact
		expected openapi.Contact
	}{
		{
			name: "full contact info",
			contact: openapi.Contact{
				Name:  "API Support",
				URL:   "https://example.com/support",
				Email: "support@example.com",
			},
			expected: openapi.Contact{
				Name:  "API Support",
				URL:   "https://example.com/support",
				Email: "support@example.com",
			},
		},
		{
			name: "minimal contact info",
			contact: openapi.Contact{
				Name: "Support Team",
			},
			expected: openapi.Contact{
				Name: "Support Team",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &openapi.Config{}
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
		license  openapi.License
		expected openapi.License
	}{
		{
			name: "license with URL",
			license: openapi.License{
				Name: "MIT",
				URL:  "https://opensource.org/licenses/MIT",
			},
			expected: openapi.License{
				Name: "MIT",
				URL:  "https://opensource.org/licenses/MIT",
			},
		},
		{
			name: "license without URL",
			license: openapi.License{
				Name: "Apache 2.0",
			},
			expected: openapi.License{
				Name: "Apache 2.0",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &openapi.Config{}
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
		expected *openapi.ExternalDocs
	}{
		{
			name: "with description",
			url:  "https://example.com/docs",
			desc: []string{"External documentation"},
			expected: &openapi.ExternalDocs{
				URL:         "https://example.com/docs",
				Description: "External documentation",
			},
		},
		{
			name: "without description",
			url:  "https://example.com/docs",
			desc: []string{},
			expected: &openapi.ExternalDocs{
				URL: "https://example.com/docs",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &openapi.Config{}
			opt := option.WithExternalDocs(tt.url, tt.desc...)
			opt(config)

			require.NotNil(t, config.ExternalDocs)
			assert.Equal(t, tt.expected.URL, config.ExternalDocs.URL)
			assert.Equal(t, tt.expected.Description, config.ExternalDocs.Description)
		})
	}
}

func TestWithTags(t *testing.T) {
	tags := []openapi.Tag{
		{
			Name:        "users",
			Description: "User management",
		},
		{
			Name:        "orders",
			Description: "Order management",
		},
	}

	config := &openapi.Config{}
	opt := option.WithTags(tags...)
	opt(config)

	require.Len(t, config.Tags, 2)
	assert.Equal(t, "users", config.Tags[0].Name)
	assert.Equal(t, "User management", config.Tags[0].Description)
	assert.Equal(t, "orders", config.Tags[1].Name)
	assert.Equal(t, "Order management", config.Tags[1].Description)
}

func TestOpenAPIConfigDefaults(t *testing.T) {
	config := &openapi.Config{}

	// Test that default values are properly set
	assert.Empty(t, config.OpenAPIVersion)
	assert.False(t, config.DisableDocs)
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
	assert.Nil(t, config.PathParser)
}

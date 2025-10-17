package option_test

import (
	"testing"

	"github.com/oaswrap/spec-ui/config"
	"github.com/oaswrap/spec/openapi"
	"github.com/oaswrap/spec/option"
	"github.com/oaswrap/spec/pkg/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWithOpenAPIConfig(t *testing.T) {
	tests := []struct {
		name     string
		opts     []option.OpenAPIOption
		validate func(t *testing.T, config *openapi.Config)
	}{
		{
			name: "default configuration",
			opts: []option.OpenAPIOption{},
			validate: func(t *testing.T, config *openapi.Config) {
				assert.Equal(t, "3.0.3", config.OpenAPIVersion)
				assert.Equal(t, "API Documentation", config.Title)
				assert.Equal(t, "/docs", config.DocsPath)
				assert.Equal(t, "/docs/openapi.yaml", config.SpecPath)
				assert.Nil(t, config.Description)
				assert.NotNil(t, config.Logger)
				assert.Empty(t, config.Version)
				assert.Nil(t, config.Contact)
				assert.Nil(t, config.License)
				assert.Empty(t, config.TermsOfService)
				assert.Empty(t, config.Tags)
				assert.Empty(t, config.Servers)
				assert.Nil(t, config.ExternalDocs)
				assert.Nil(t, config.SecuritySchemes)
				assert.Nil(t, config.ReflectorConfig)
				assert.False(t, config.DisableDocs)
				assert.Nil(t, config.SwaggerUIConfig)
				assert.Nil(t, config.StoplightElementsConfig)
				assert.Nil(t, config.ReDocConfig)
				assert.Nil(t, config.ScalarConfig)
				assert.Nil(t, config.RapiDocConfig)
				assert.Nil(t, config.PathParser)
			},
		},
		{
			name: "single option",
			opts: []option.OpenAPIOption{
				option.WithTitle("My Custom API"),
			},
			validate: func(t *testing.T, config *openapi.Config) {
				assert.Equal(t, "3.0.3", config.OpenAPIVersion)
				assert.Equal(t, "My Custom API", config.Title)
				assert.Nil(t, config.Description)
				assert.NotNil(t, config.Logger)
			},
		},
		{
			name: "multiple options",
			opts: []option.OpenAPIOption{
				option.WithOpenAPIVersion("3.1.0"),
				option.WithTitle("Advanced API"),
				option.WithVersion("2.0.0"),
				option.WithDescription("A comprehensive API"),
				option.WithCacheAge(86400),
			},
			validate: func(t *testing.T, config *openapi.Config) {
				assert.Equal(t, "3.1.0", config.OpenAPIVersion)
				assert.Equal(t, "Advanced API", config.Title)
				assert.Equal(t, "2.0.0", config.Version)
				require.NotNil(t, config.Description)
				assert.Equal(t, "A comprehensive API", *config.Description)
				assert.Equal(t, 86400, *config.CacheAge)
			},
		},
		{
			name: "complex configuration",
			opts: []option.OpenAPIOption{
				option.WithTitle("Complete API"),
				option.WithVersion("1.0.0"),
				option.WithDescription("Full-featured API"),
				option.WithContact(openapi.Contact{
					Name:  "Support",
					Email: "support@example.com",
				}),
				option.WithLicense(openapi.License{
					Name: "MIT",
					URL:  "https://opensource.org/licenses/MIT",
				}),
				option.WithServer("https://api.example.com"),
				option.WithTags(openapi.Tag{
					Name:        "users",
					Description: "User operations",
				}),
				option.WithTermsOfService("https://example.com/terms"),
				option.WithDebug(true),
				option.WithDisableDocs(false),
			},
			validate: func(t *testing.T, config *openapi.Config) {
				assert.Equal(t, "Complete API", config.Title)
				assert.Equal(t, "1.0.0", config.Version)
				require.NotNil(t, config.Description)
				assert.Equal(t, "Full-featured API", *config.Description)
				require.NotNil(t, config.Contact)
				assert.Equal(t, "Support", config.Contact.Name)
				assert.Equal(t, "support@example.com", config.Contact.Email)
				require.NotNil(t, config.License)
				assert.Equal(t, "MIT", config.License.Name)
				assert.Equal(t, "https://opensource.org/licenses/MIT", config.License.URL)
				require.Len(t, config.Servers, 1)
				assert.Equal(t, "https://api.example.com", config.Servers[0].URL)
				require.Len(t, config.Tags, 1)
				assert.Equal(t, "users", config.Tags[0].Name)
				assert.Equal(t, "User operations", config.Tags[0].Description)
				assert.Equal(t, "https://example.com/terms", *config.TermsOfService)
				assert.NotNil(t, config.Logger)
				assert.False(t, config.DisableDocs)
			},
		},
		{
			name: "overriding defaults",
			opts: []option.OpenAPIOption{
				option.WithOpenAPIVersion("3.1.0"),
				option.WithTitle("Override Title"),
				option.WithDescription("Override Description"),
				option.WithDebug(false),
			},
			validate: func(t *testing.T, config *openapi.Config) {
				assert.Equal(t, "3.1.0", config.OpenAPIVersion)
				assert.Equal(t, "Override Title", config.Title)
				require.NotNil(t, config.Description)
				assert.Equal(t, "Override Description", *config.Description)
				assert.NotNil(t, config.Logger)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := option.WithOpenAPIConfig(tt.opts...)
			require.NotNil(t, config)
			tt.validate(t, config)
		})
	}
}

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
					Implicit: &openapi.OAuthFlowsImplicit{
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
						Implicit: &openapi.OAuthFlowsImplicit{
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

func TestWithSwaggerUI(t *testing.T) {
	tests := []struct {
		name     string
		cfgs     []config.SwaggerUI
		expected *config.SwaggerUI
	}{
		{
			name:     "no config",
			cfgs:     []config.SwaggerUI{},
			expected: &config.SwaggerUI{},
		},
		{
			name:     "empty config",
			cfgs:     []config.SwaggerUI{{}},
			expected: &config.SwaggerUI{},
		},
		{
			name: "valid config",
			cfgs: []config.SwaggerUI{
				{
					HideCurl: false,
				},
			},
			expected: &config.SwaggerUI{
				HideCurl: false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &openapi.Config{}
			opt := option.WithSwaggerUI(tt.cfgs...)
			opt(config)

			assert.Equal(t, tt.expected, config.SwaggerUIConfig)
		})
	}
}

func TestWithStoplightElements(t *testing.T) {
	tests := []struct {
		name     string
		cfgs     []config.StoplightElements
		expected *config.StoplightElements
	}{
		{
			name:     "no config",
			cfgs:     []config.StoplightElements{},
			expected: &config.StoplightElements{},
		},
		{
			name:     "empty config",
			cfgs:     []config.StoplightElements{{}},
			expected: &config.StoplightElements{},
		},
		{
			name: "valid config",
			cfgs: []config.StoplightElements{
				{
					HideExport:  true,
					HideSchemas: true,
					Logo:        "https://example.com/logo.png",
					Layout:      "sidebar",
					Router:      "hash",
				},
			},
			expected: &config.StoplightElements{
				HideExport:  true,
				HideSchemas: true,
				Logo:        "https://example.com/logo.png",
				Layout:      "sidebar",
				Router:      "hash",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &openapi.Config{}
			opt := option.WithStoplightElements(tt.cfgs...)
			opt(config)

			assert.Equal(t, tt.expected, config.StoplightElementsConfig)
		})
	}
}

func TestWithRedoc(t *testing.T) {
	tests := []struct {
		name     string
		cfgs     []config.ReDoc
		expected *config.ReDoc
	}{
		{
			name:     "no config",
			cfgs:     []config.ReDoc{},
			expected: &config.ReDoc{},
		},
		{
			name:     "empty config",
			cfgs:     []config.ReDoc{},
			expected: &config.ReDoc{},
		},
		{
			name: "valid config",
			cfgs: []config.ReDoc{
				{
					DisableSearch:       true,
					HideDownloadButtons: true,
					HideSchemaTitles:    true,
				},
			},
			expected: &config.ReDoc{
				DisableSearch:       true,
				HideDownloadButtons: true,
				HideSchemaTitles:    true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &openapi.Config{}
			opt := option.WithReDoc(tt.cfgs...)
			opt(config)

			assert.Equal(t, tt.expected, config.ReDocConfig)
		})
	}
}

func TestWithScalar(t *testing.T) {
	tests := []struct {
		name     string
		cfgs     []config.Scalar
		expected *config.Scalar
	}{
		{
			name:     "no config",
			cfgs:     []config.Scalar{},
			expected: &config.Scalar{},
		},
		{
			name:     "empty config",
			cfgs:     []config.Scalar{{}},
			expected: &config.Scalar{},
		},
		{
			name: "valid config",
			cfgs: []config.Scalar{
				{
					HideSidebar: true,
					HideModels:  true,
				},
			},
			expected: &config.Scalar{
				HideSidebar: true,
				HideModels:  true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &openapi.Config{}
			opt := option.WithScalar(tt.cfgs...)
			opt(config)

			assert.Equal(t, tt.expected, config.ScalarConfig)
		})
	}
}

func TestWithRapiDoc(t *testing.T) {
	tests := []struct {
		name     string
		cfgs     []config.RapiDoc
		expected *config.RapiDoc
	}{
		{
			name:     "no config",
			cfgs:     []config.RapiDoc{},
			expected: &config.RapiDoc{},
		},
		{
			name:     "empty config",
			cfgs:     []config.RapiDoc{{}},
			expected: &config.RapiDoc{},
		},
		{
			name: "valid config",
			cfgs: []config.RapiDoc{
				{
					Theme:     "dark",
					HideTryIt: true,
				},
			},
			expected: &config.RapiDoc{
				Theme:     "dark",
				HideTryIt: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &openapi.Config{}
			opt := option.WithRapiDoc(tt.cfgs...)
			opt(config)

			assert.Equal(t, tt.expected, config.RapiDocConfig)
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

func TestWithReflectorConfig(t *testing.T) {
	tests := []struct {
		name     string
		opts     []option.ReflectorOption
		validate func(t *testing.T, config *openapi.Config)
	}{
		{
			name: "creates new reflector config when nil",
			opts: []option.ReflectorOption{},
			validate: func(t *testing.T, config *openapi.Config) {
				require.NotNil(t, config.ReflectorConfig)
			},
		},
		{
			name: "applies single option",
			opts: []option.ReflectorOption{
				option.InlineRefs(),
			},
			validate: func(t *testing.T, config *openapi.Config) {
				require.NotNil(t, config.ReflectorConfig)
				assert.True(t, config.ReflectorConfig.InlineRefs)
			},
		},
		{
			name: "applies multiple options",
			opts: []option.ReflectorOption{
				option.RequiredPropByValidateTag("validate", ","),
				option.StripDefNamePrefix("MyPrefix"),
			},
			validate: func(t *testing.T, config *openapi.Config) {
				require.NotNil(t, config.ReflectorConfig)
				assert.NotNil(t, config.ReflectorConfig.InterceptPropFunc)
				assert.NotEmpty(t, config.ReflectorConfig.StripDefNamePrefix)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &openapi.Config{}

			// For the "preserves existing reflector config" test, pre-populate the config
			if tt.name == "preserves existing reflector config" {
				config.ReflectorConfig = &openapi.ReflectorConfig{}
			}

			opt := option.WithReflectorConfig(tt.opts...)
			opt(config)

			tt.validate(t, config)
		})
	}
}

func TestWithDebug(t *testing.T) {
	config := &openapi.Config{}
	opt := option.WithDebug(true)
	opt(config)

	assert.NotNil(t, config.Logger)

	config = &openapi.Config{}
	opt = option.WithDebug(false)
	opt(config)
	assert.NotNil(t, config.Logger)
}

func TestWithPathParser(t *testing.T) {
	// Mock path parser for testing
	mockParser := &mockPathParser{}

	config := &openapi.Config{}
	opt := option.WithPathParser(mockParser)
	opt(config)

	assert.Equal(t, mockParser, config.PathParser)
	assert.NotNil(t, config.PathParser)
}

// mockPathParser is a test implementation of openapi.PathParser.
type mockPathParser struct{}

func (m *mockPathParser) Parse(path string) (string, error) {
	// Simple mock implementation that converts :param to {param}
	return path, nil
}

func TestOpenAPIConfigDefaults(t *testing.T) {
	config := &openapi.Config{}

	// Test that default values are properly set
	assert.Empty(t, config.OpenAPIVersion)
	assert.False(t, config.DisableDocs)
	assert.Empty(t, config.Title)
	assert.Empty(t, config.Version)
	assert.Nil(t, config.Description)
	assert.Empty(t, config.Servers)
	assert.Empty(t, config.DocsPath)
	assert.Empty(t, config.SpecPath)
	assert.Nil(t, config.SecuritySchemes)
	assert.Nil(t, config.SwaggerUIConfig)
	assert.Nil(t, config.StoplightElementsConfig)
	assert.Nil(t, config.ReDocConfig)
	assert.Nil(t, config.ScalarConfig)
	assert.Nil(t, config.RapiDocConfig)
	assert.Nil(t, config.Logger)
	assert.Nil(t, config.Contact)
	assert.Nil(t, config.License)
	assert.Empty(t, config.Tags)
	assert.Empty(t, config.TermsOfService)
	assert.Nil(t, config.PathParser)
}

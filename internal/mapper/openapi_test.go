package mapper_test

import (
	"testing"

	"github.com/oaswrap/spec/internal/mapper"
	"github.com/oaswrap/spec/openapi"
	"github.com/oaswrap/spec/pkg/util"
	"github.com/stretchr/testify/assert"
	"github.com/swaggest/openapi-go/openapi3"
	"github.com/swaggest/openapi-go/openapi31"
)

func TestOASContact(t *testing.T) {
	tests := []struct {
		name       string
		contact    *openapi.Contact
		expected3  *openapi3.Contact
		expected31 *openapi31.Contact
	}{
		{
			name:      "nil contact",
			contact:   nil,
			expected3: nil,
		},
		{
			name:       "empty contact",
			contact:    &openapi.Contact{},
			expected3:  &openapi3.Contact{},
			expected31: &openapi31.Contact{},
		},
		{
			name:       "contact with name only",
			contact:    &openapi.Contact{Name: "John Doe"},
			expected3:  &openapi3.Contact{Name: util.PtrOf("John Doe")},
			expected31: &openapi31.Contact{Name: util.PtrOf("John Doe")},
		},
		{
			name:       "contact with URL only",
			contact:    &openapi.Contact{URL: "https://example.com"},
			expected3:  &openapi3.Contact{URL: util.PtrOf("https://example.com")},
			expected31: &openapi31.Contact{URL: util.PtrOf("https://example.com")},
		},
		{
			name:       "contact with email only",
			contact:    &openapi.Contact{Email: "john@example.com"},
			expected3:  &openapi3.Contact{Email: util.PtrOf("john@example.com")},
			expected31: &openapi31.Contact{Email: util.PtrOf("john@example.com")},
		},
		{
			name: "contact with all fields",
			contact: &openapi.Contact{
				Name:  "John Doe",
				URL:   "https://example.com",
				Email: "john@example.com",
				MapOfAnything: map[string]interface{}{
					"x-custom": "value",
				},
			},
			expected3: &openapi3.Contact{
				Name:  util.PtrOf("John Doe"),
				URL:   util.PtrOf("https://example.com"),
				Email: util.PtrOf("john@example.com"),
				MapOfAnything: map[string]interface{}{
					"x-custom": "value",
				},
			},
			expected31: &openapi31.Contact{
				Name:  util.PtrOf("John Doe"),
				URL:   util.PtrOf("https://example.com"),
				Email: util.PtrOf("john@example.com"),
				MapOfAnything: map[string]interface{}{
					"x-custom": "value",
				},
			},
		},
		{
			name: "contact with MapOfAnything only",
			contact: &openapi.Contact{
				MapOfAnything: map[string]interface{}{
					"x-vendor": "extension",
				},
			},
			expected3: &openapi3.Contact{
				MapOfAnything: map[string]interface{}{
					"x-vendor": "extension",
				},
			},
			expected31: &openapi31.Contact{
				MapOfAnything: map[string]interface{}{
					"x-vendor": "extension",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result3 := mapper.OAS3Contact(tt.contact)
			assert.Equal(t, tt.expected3, result3, "OpenAPI 3 Contact mapping failed")

			result31 := mapper.OAS31Contact(tt.contact)
			assert.Equal(t, tt.expected31, result31, "OpenAPI 3.1 Contact mapping failed")
		})
	}
}

func TestOASLicense(t *testing.T) {
	tests := []struct {
		name       string
		license    *openapi.License
		expected3  *openapi3.License
		expected31 *openapi31.License
	}{
		{
			name:       "nil license",
			license:    nil,
			expected3:  nil,
			expected31: nil,
		},
		{
			name:       "empty license",
			license:    &openapi.License{},
			expected3:  &openapi3.License{},
			expected31: &openapi31.License{},
		},
		{
			name:       "license with name only",
			license:    &openapi.License{Name: "MIT"},
			expected3:  &openapi3.License{Name: "MIT"},
			expected31: &openapi31.License{Name: "MIT"},
		},
		{
			name: "license with all fields",
			license: &openapi.License{
				Name: "Apache 2.0",
				URL:  "https://www.apache.org/licenses/LICENSE-2.0.html",
				MapOfAnything: map[string]interface{}{
					"x-custom": "license-extension",
				},
			},
			expected3: &openapi3.License{
				Name: "Apache 2.0",
				URL:  util.PtrOf("https://www.apache.org/licenses/LICENSE-2.0.html"),
				MapOfAnything: map[string]interface{}{
					"x-custom": "license-extension",
				},
			},
			expected31: &openapi31.License{
				Name: "Apache 2.0",
				URL:  util.PtrOf("https://www.apache.org/licenses/LICENSE-2.0.html"),
				MapOfAnything: map[string]interface{}{
					"x-custom": "license-extension",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result3 := mapper.OAS3License(tt.license)
			assert.Equal(t, tt.expected3, result3, "OpenAPI 3 License mapping failed")
			result31 := mapper.OAS31License(tt.license)
			assert.Equal(t, tt.expected31, result31, "OpenAPI 3.1 License mapping failed")
		})
	}
}

func TestOASExternalDocs(t *testing.T) {
	tests := []struct {
		name       string
		external   *openapi.ExternalDocs
		expected3  *openapi3.ExternalDocumentation
		expected31 *openapi31.ExternalDocumentation
	}{
		{
			name:       "nil external docs",
			external:   nil,
			expected3:  nil,
			expected31: nil,
		},
		{
			name:       "empty external docs",
			external:   &openapi.ExternalDocs{},
			expected3:  &openapi3.ExternalDocumentation{},
			expected31: &openapi31.ExternalDocumentation{},
		},
		{
			name: "external docs with URL only",
			external: &openapi.ExternalDocs{
				URL: "https://example.com/docs",
			},
			expected3: &openapi3.ExternalDocumentation{
				URL: "https://example.com/docs",
			},
			expected31: &openapi31.ExternalDocumentation{
				URL: "https://example.com/docs",
			},
		},
		{
			name: "external docs with description only",
			external: &openapi.ExternalDocs{
				Description: "API documentation",
			},
			expected3: &openapi3.ExternalDocumentation{
				Description: util.PtrOf("API documentation"),
			},
			expected31: &openapi31.ExternalDocumentation{
				Description: util.PtrOf("API documentation"),
			},
		},
		{
			name: "external docs with all fields",
			external: &openapi.ExternalDocs{
				URL:         "https://example.com/docs",
				Description: "API documentation",
				MapOfAnything: map[string]interface{}{
					"x-custom": "external-docs-extension",
				},
			},
			expected3: &openapi3.ExternalDocumentation{
				URL:         "https://example.com/docs",
				Description: util.PtrOf("API documentation"),
				MapOfAnything: map[string]interface{}{
					"x-custom": "external-docs-extension",
				},
			},
			expected31: &openapi31.ExternalDocumentation{
				URL:         "https://example.com/docs",
				Description: util.PtrOf("API documentation"),
				MapOfAnything: map[string]interface{}{
					"x-custom": "external-docs-extension",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result3 := mapper.OAS3ExternalDocs(tt.external)
			assert.Equal(t, tt.expected3, result3, "OpenAPI 3 ExternalDocs mapping failed")
			result31 := mapper.OAS31ExternalDocs(tt.external)
			assert.Equal(t, tt.expected31, result31, "OpenAPI 3.1 ExternalDocs mapping failed")
		})
	}
}

func TestOASTags(t *testing.T) {
	tests := []struct {
		name       string
		tags       []openapi.Tag
		expected3  []openapi3.Tag
		expected31 []openapi31.Tag
	}{
		{
			name: "single tag",
			tags: []openapi.Tag{
				{Name: "example"},
			},
			expected3: []openapi3.Tag{
				{Name: "example"},
			},
			expected31: []openapi31.Tag{
				{Name: "example"},
			},
		},
		{
			name: "multiple tags",
			tags: []openapi.Tag{
				{Name: "tag1"},
				{Name: "tag2", Description: "Description for tag2"},
			},
			expected3: []openapi3.Tag{
				{Name: "tag1"},
				{Name: "tag2", Description: util.PtrOf("Description for tag2")},
			},
			expected31: []openapi31.Tag{
				{Name: "tag1"},
				{Name: "tag2", Description: util.PtrOf("Description for tag2")},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result3 := mapper.OAS3Tags(tt.tags)
			assert.Equal(t, tt.expected3, result3, "OpenAPI 3 Tags mapping failed")

			result31 := mapper.OAS31Tags(tt.tags)
			assert.Equal(t, tt.expected31, result31, "OpenAPI 3.1 Tags mapping failed")
		})
	}
}

func TestOASTag(t *testing.T) {
	tests := []struct {
		name       string
		tag        openapi.Tag
		expected3  openapi3.Tag
		expected31 openapi31.Tag
	}{
		{
			name: "empty tag",
			tag:  openapi.Tag{},
			expected3: openapi3.Tag{
				Name:          "",
				Description:   nil,
				MapOfAnything: nil,
			},
			expected31: openapi31.Tag{
				Name:          "",
				Description:   nil,
				MapOfAnything: nil,
			},
		},
		{
			name: "tag with name only",
			tag:  openapi.Tag{Name: "example"},
			expected3: openapi3.Tag{
				Name:          "example",
				Description:   nil,
				MapOfAnything: nil,
			},
			expected31: openapi31.Tag{
				Name:          "example",
				Description:   nil,
				MapOfAnything: nil,
			},
		},
		{
			name: "tag with description",
			tag:  openapi.Tag{Name: "example", Description: "An example tag"},
			expected3: openapi3.Tag{
				Name:          "example",
				Description:   util.PtrOf("An example tag"),
				MapOfAnything: nil,
			},
			expected31: openapi31.Tag{
				Name:          "example",
				Description:   util.PtrOf("An example tag"),
				MapOfAnything: nil,
			},
		},
		{
			name: "tag with external docs",
			tag: openapi.Tag{
				Name:        "example",
				Description: "An example tag",
				ExternalDocs: &openapi.ExternalDocs{
					URL:         "https://example.com/docs",
					Description: "External documentation for the example tag",
					MapOfAnything: map[string]interface{}{
						"x-custom": "value",
					},
				},
			},
			expected3: openapi3.Tag{
				Name:        "example",
				Description: util.PtrOf("An example tag"),
				ExternalDocs: &openapi3.ExternalDocumentation{
					URL:         "https://example.com/docs",
					Description: util.PtrOf("External documentation for the example tag"),
					MapOfAnything: map[string]interface{}{
						"x-custom": "value",
					},
				},
			},
			expected31: openapi31.Tag{
				Name:        "example",
				Description: util.PtrOf("An example tag"),
				ExternalDocs: &openapi31.ExternalDocumentation{
					URL:         "https://example.com/docs",
					Description: util.PtrOf("External documentation for the example tag"),
					MapOfAnything: map[string]interface{}{
						"x-custom": "value",
					},
				},
			},
		},
		{
			name: "tag with all fields",
			tag: openapi.Tag{
				Name:        "example",
				Description: "An example tag",
				ExternalDocs: &openapi.ExternalDocs{
					URL:         "https://example.com/docs",
					Description: "External documentation for the example tag",
					MapOfAnything: map[string]interface{}{
						"x-custom": "value",
					},
				},
				MapOfAnything: map[string]interface{}{
					"x-vendor": "example-tag-extension",
				},
			},
			expected3: openapi3.Tag{
				Name:        "example",
				Description: util.PtrOf("An example tag"),
				ExternalDocs: &openapi3.ExternalDocumentation{
					URL:         "https://example.com/docs",
					Description: util.PtrOf("External documentation for the example tag"),
					MapOfAnything: map[string]interface{}{
						"x-custom": "value",
					},
				},
				MapOfAnything: map[string]interface{}{
					"x-vendor": "example-tag-extension",
				},
			},
			expected31: openapi31.Tag{
				Name:        "example",
				Description: util.PtrOf("An example tag"),
				ExternalDocs: &openapi31.ExternalDocumentation{
					URL:         "https://example.com/docs",
					Description: util.PtrOf("External documentation for the example tag"),
					MapOfAnything: map[string]interface{}{
						"x-custom": "value",
					},
				},
				MapOfAnything: map[string]interface{}{
					"x-vendor": "example-tag-extension",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result3 := mapper.OAS3Tag(tt.tag)
			assert.Equal(t, tt.expected3, result3, "OpenAPI 3 Tag mapping failed")

			result31 := mapper.OAS31Tag(tt.tag)
			assert.Equal(t, tt.expected31, result31, "OpenAPI 3.1 Tag mapping failed")
		})
	}
}

func TestOASServers(t *testing.T) {
	tests := []struct {
		name       string
		servers    []openapi.Server
		expected3  []openapi3.Server
		expected31 []openapi31.Server
	}{
		{
			name: "single server",
			servers: []openapi.Server{
				{URL: "https://api.example.com"},
			},
			expected3: []openapi3.Server{
				{URL: "https://api.example.com"},
			},
			expected31: []openapi31.Server{
				{URL: "https://api.example.com"},
			},
		},
		{
			name: "multiple servers",
			servers: []openapi.Server{
				{URL: "https://api1.example.com"},
				{URL: "https://api2.example.com", Description: util.PtrOf("Second API server")},
			},
			expected3: []openapi3.Server{
				{URL: "https://api1.example.com"},
				{URL: "https://api2.example.com", Description: util.PtrOf("Second API server")},
			},
			expected31: []openapi31.Server{
				{URL: "https://api1.example.com"},
				{URL: "https://api2.example.com", Description: util.PtrOf("Second API server")},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result3 := mapper.OAS3Servers(tt.servers)
			assert.Equal(t, tt.expected3, result3, "OpenAPI 3 Servers mapping failed")

			result31 := mapper.OAS31Servers(tt.servers)
			assert.Equal(t, tt.expected31, result31, "OpenAPI 3.1 Servers mapping failed")
		})
	}
}

func TestOASServer(t *testing.T) {
	tests := []struct {
		name       string
		server     openapi.Server
		expected3  openapi3.Server
		expected31 openapi31.Server
	}{
		{
			name:   "empty server",
			server: openapi.Server{},
			expected3: openapi3.Server{
				URL:           "",
				Description:   nil,
				MapOfAnything: nil,
				Variables:     nil,
			},
			expected31: openapi31.Server{
				URL:           "",
				Description:   nil,
				MapOfAnything: nil,
				Variables:     nil,
			},
		},
		{
			name:   "server with URL only",
			server: openapi.Server{URL: "https://api.example.com"},
			expected3: openapi3.Server{
				URL:           "https://api.example.com",
				Description:   nil,
				MapOfAnything: nil,
				Variables:     nil,
			},
			expected31: openapi31.Server{
				URL:           "https://api.example.com",
				Description:   nil,
				MapOfAnything: nil,
				Variables:     nil,
			},
		},
		{
			name:   "server with description",
			server: openapi.Server{URL: "https://api.example.com", Description: util.PtrOf("API server")},
			expected3: openapi3.Server{
				URL:           "https://api.example.com",
				Description:   util.PtrOf("API server"),
				MapOfAnything: nil,
				Variables:     nil,
			},
			expected31: openapi31.Server{
				URL:           "https://api.example.com",
				Description:   util.PtrOf("API server"),
				MapOfAnything: nil,
				Variables:     nil,
			},
		},
		{
			name: "server with variables",
			server: openapi.Server{
				URL:         "https://api.example.com",
				Description: util.PtrOf("API server"),
				Variables: map[string]openapi.ServerVariable{
					"port": {
						Enum:    []string{"8080", "8443"},
						Default: "8080",
					},
				},
			},
			expected3: openapi3.Server{
				URL:           "https://api.example.com",
				Description:   util.PtrOf("API server"),
				MapOfAnything: nil,
				Variables: map[string]openapi3.ServerVariable{
					"port": {
						Enum:    []string{"8080", "8443"},
						Default: "8080",
					},
				},
			},
			expected31: openapi31.Server{
				URL:           "https://api.example.com",
				Description:   util.PtrOf("API server"),
				MapOfAnything: nil,
				Variables: map[string]openapi31.ServerVariable{
					"port": {
						Enum:    []string{"8080", "8443"},
						Default: "8080",
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result3 := mapper.OAS3Server(tt.server)
			assert.Equal(t, tt.expected3, result3, "OpenAPI 3 Server mapping failed")

			result31 := mapper.OAS31Server(tt.server)
			assert.Equal(t, tt.expected31, result31, "OpenAPI 3.1 Server mapping failed")
		})
	}
}

func TestOASSecurityScheme(t *testing.T) {
	tests := []struct {
		name       string
		scheme     *openapi.SecurityScheme
		expected3  *openapi3.SecurityScheme
		expected31 *openapi31.SecurityScheme
	}{
		{
			name:       "nil scheme",
			scheme:     nil,
			expected3:  nil,
			expected31: nil,
		},
		{
			name:       "empty scheme",
			scheme:     &openapi.SecurityScheme{},
			expected3:  nil,
			expected31: nil,
		},
		{
			name: "API Key scheme",
			scheme: &openapi.SecurityScheme{
				APIKey: &openapi.SecuritySchemeAPIKey{Name: "api_key", In: "header"},
			},
			expected3: &openapi3.SecurityScheme{
				APIKeySecurityScheme: &openapi3.APIKeySecurityScheme{
					Name: "api_key",
					In:   openapi3.APIKeySecuritySchemeIn("header"),
				},
			},
			expected31: &openapi31.SecurityScheme{
				APIKey: &openapi31.SecuritySchemeAPIKey{
					Name: "api_key",
					In:   openapi31.SecuritySchemeAPIKeyIn("header"),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result3 := mapper.OAS3SecurityScheme(tt.scheme)
			assert.Equal(t, tt.expected3, result3, "OpenAPI 3 Security Scheme mapping failed")

			result31 := mapper.OAS31SecurityScheme(tt.scheme)
			assert.Equal(t, tt.expected31, result31, "OpenAPI 3.1 Security Scheme mapping failed")
		})
	}
}

func TestOASAPIKey(t *testing.T) {
	tests := []struct {
		name       string
		apiKey     *openapi.SecuritySchemeAPIKey
		expected3  *openapi3.APIKeySecurityScheme
		expected31 *openapi31.SecuritySchemeAPIKey
	}{
		{
			name:       "nil apiKey",
			apiKey:     nil,
			expected3:  nil,
			expected31: nil,
		},
		{
			name: "empty apiKey",
			apiKey: &openapi.SecuritySchemeAPIKey{
				Name: "",
				In:   "",
			},
			expected3: &openapi3.APIKeySecurityScheme{
				Name: "",
				In:   openapi3.APIKeySecuritySchemeIn(""),
			},
			expected31: &openapi31.SecuritySchemeAPIKey{
				Name: "",
				In:   openapi31.SecuritySchemeAPIKeyIn(""),
			},
		},
		{
			name: "apiKey with name only",
			apiKey: &openapi.SecuritySchemeAPIKey{
				Name: "api_key",
				In:   "",
			},
			expected3: &openapi3.APIKeySecurityScheme{
				Name: "api_key",
				In:   openapi3.APIKeySecuritySchemeIn(""),
			},
			expected31: &openapi31.SecuritySchemeAPIKey{
				Name: "api_key",
				In:   openapi31.SecuritySchemeAPIKeyIn(""),
			},
		},
		{
			name: "apiKey with in header",
			apiKey: &openapi.SecuritySchemeAPIKey{
				Name: "X-API-Key",
				In:   "header",
			},
			expected3: &openapi3.APIKeySecurityScheme{
				Name: "X-API-Key",
				In:   openapi3.APIKeySecuritySchemeIn("header"),
			},
			expected31: &openapi31.SecuritySchemeAPIKey{
				Name: "X-API-Key",
				In:   openapi31.SecuritySchemeAPIKeyIn("header"),
			},
		},
		{
			name: "apiKey with in query",
			apiKey: &openapi.SecuritySchemeAPIKey{
				Name: "api_key",
				In:   "query",
			},
			expected3: &openapi3.APIKeySecurityScheme{
				Name: "api_key",
				In:   openapi3.APIKeySecuritySchemeIn("query"),
			},
			expected31: &openapi31.SecuritySchemeAPIKey{
				Name: "api_key",
				In:   openapi31.SecuritySchemeAPIKeyIn("query"),
			},
		},
		{
			name: "apiKey with in cookie",
			apiKey: &openapi.SecuritySchemeAPIKey{
				Name: "sessionId",
				In:   "cookie",
			},
			expected3: &openapi3.APIKeySecurityScheme{
				Name: "sessionId",
				In:   openapi3.APIKeySecuritySchemeIn("cookie"),
			},
			expected31: &openapi31.SecuritySchemeAPIKey{
				Name: "sessionId",
				In:   openapi31.SecuritySchemeAPIKeyIn("cookie"),
			},
		},
		{
			name: "apiKey with all fields",
			apiKey: &openapi.SecuritySchemeAPIKey{
				Name: "Authorization",
				In:   "header",
			},
			expected3: &openapi3.APIKeySecurityScheme{
				Name: "Authorization",
				In:   openapi3.APIKeySecuritySchemeIn("header"),
			},
			expected31: &openapi31.SecuritySchemeAPIKey{
				Name: "Authorization",
				In:   openapi31.SecuritySchemeAPIKeyIn("header"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := mapper.OAS31APIKey(tt.apiKey)
			assert.Equal(t, tt.expected31, result)
		})
	}
}

func TestOASHTTPBearer(t *testing.T) {
	tests := []struct {
		name       string
		scheme     *openapi.SecuritySchemeHTTPBearer
		expected3  *openapi3.HTTPSecurityScheme
		expected31 *openapi31.SecuritySchemeHTTPBearer
	}{
		{
			name:       "nil scheme",
			scheme:     nil,
			expected3:  nil,
			expected31: nil,
		},
		{
			name: "empty scheme",
			scheme: &openapi.SecuritySchemeHTTPBearer{
				Scheme: "",
			},
			expected3: &openapi3.HTTPSecurityScheme{
				Scheme: "",
			},
			expected31: &openapi31.SecuritySchemeHTTPBearer{
				Scheme: "",
			},
		},
		{
			name: "scheme with scheme only",
			scheme: &openapi.SecuritySchemeHTTPBearer{
				Scheme: "bearer",
			},
			expected3: &openapi3.HTTPSecurityScheme{
				Scheme: "bearer",
			},
			expected31: &openapi31.SecuritySchemeHTTPBearer{
				Scheme: "bearer",
			},
		},
		{
			name: "scheme with all fields",
			scheme: &openapi.SecuritySchemeHTTPBearer{
				Scheme:       "bearer",
				BearerFormat: util.PtrOf("JWT"),
			},
			expected3: &openapi3.HTTPSecurityScheme{
				Scheme:       "bearer",
				BearerFormat: util.PtrOf("JWT"),
			},
			expected31: &openapi31.SecuritySchemeHTTPBearer{
				Scheme:       "bearer",
				BearerFormat: util.PtrOf("JWT"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := mapper.OAS31HTTPBearer(tt.scheme)
			assert.Equal(t, tt.expected31, result)
		})
	}
}

func TestOASOauth2Flows(t *testing.T) {
	tests := []struct {
		name       string
		flows      openapi.OAuthFlows
		expected3  openapi3.OAuthFlows
		expected31 openapi31.OauthFlows
	}{
		{
			name:  "empty flows",
			flows: openapi.OAuthFlows{},
			expected3: openapi3.OAuthFlows{
				Implicit:          nil,
				Password:          nil,
				ClientCredentials: nil,
				AuthorizationCode: nil,
			},
			expected31: openapi31.OauthFlows{
				Implicit:          nil,
				Password:          nil,
				ClientCredentials: nil,
				AuthorizationCode: nil,
			},
		},
		{
			name: "flows with implicit only",
			flows: openapi.OAuthFlows{
				Implicit: &openapi.OAuthFlowsDefsImplicit{
					AuthorizationURL: "https://example.com/auth",
					RefreshURL:       util.PtrOf("https://example.com/refresh"),
					Scopes: map[string]string{
						"read":  "Read access",
						"write": "Write access",
					},
					MapOfAnything: map[string]interface{}{
						"x-custom": "implicit-flow",
					},
				},
			},
			expected3: openapi3.OAuthFlows{
				Implicit: &openapi3.ImplicitOAuthFlow{
					AuthorizationURL: "https://example.com/auth",
					RefreshURL:       util.PtrOf("https://example.com/refresh"),
					Scopes: map[string]string{
						"read":  "Read access",
						"write": "Write access",
					},
					MapOfAnything: map[string]interface{}{
						"x-custom": "implicit-flow",
					},
				},
			},
			expected31: openapi31.OauthFlows{
				Implicit: &openapi31.OauthFlowsDefsImplicit{
					AuthorizationURL: "https://example.com/auth",
					RefreshURL:       util.PtrOf("https://example.com/refresh"),
					Scopes: map[string]string{
						"read":  "Read access",
						"write": "Write access",
					},
					MapOfAnything: map[string]interface{}{
						"x-custom": "implicit-flow",
					},
				},
			},
		},
		{
			name: "flows with password only",
			flows: openapi.OAuthFlows{
				Password: &openapi.OAuthFlowsDefsPassword{
					TokenURL:   "https://example.com/token",
					RefreshURL: util.PtrOf("https://example.com/refresh"),
					Scopes: map[string]string{
						"admin": "Admin access",
					},
					MapOfAnything: map[string]interface{}{
						"x-vendor": "password-flow",
					},
				},
			},
			expected3: openapi3.OAuthFlows{
				Password: &openapi3.PasswordOAuthFlow{
					TokenURL:   "https://example.com/token",
					RefreshURL: util.PtrOf("https://example.com/refresh"),
					Scopes: map[string]string{
						"admin": "Admin access",
					},
					MapOfAnything: map[string]interface{}{
						"x-vendor": "password-flow",
					},
				},
			},
			expected31: openapi31.OauthFlows{
				Password: &openapi31.OauthFlowsDefsPassword{
					TokenURL:   "https://example.com/token",
					RefreshURL: util.PtrOf("https://example.com/refresh"),
					Scopes: map[string]string{
						"admin": "Admin access",
					},
					MapOfAnything: map[string]interface{}{
						"x-vendor": "password-flow",
					},
				},
			},
		},
		{
			name: "flows with client credentials only",
			flows: openapi.OAuthFlows{
				ClientCredentials: &openapi.OAuthFlowsDefsClientCredentials{
					TokenURL: "https://example.com/token",
					Scopes: map[string]string{
						"service": "Service access",
					},
					MapOfAnything: map[string]interface{}{
						"x-internal": true,
					},
				},
			},
			expected3: openapi3.OAuthFlows{
				ClientCredentials: &openapi3.ClientCredentialsFlow{
					TokenURL: "https://example.com/token",
					Scopes: map[string]string{
						"service": "Service access",
					},
					MapOfAnything: map[string]interface{}{
						"x-internal": true,
					},
				},
			},
			expected31: openapi31.OauthFlows{
				ClientCredentials: &openapi31.OauthFlowsDefsClientCredentials{
					TokenURL: "https://example.com/token",
					Scopes: map[string]string{
						"service": "Service access",
					},
					MapOfAnything: map[string]interface{}{
						"x-internal": true,
					},
				},
			},
		},
		{
			name: "flows with authorization code only",
			flows: openapi.OAuthFlows{
				AuthorizationCode: &openapi.OAuthFlowsDefsAuthorizationCode{
					AuthorizationURL: "https://example.com/auth",
					TokenURL:         "https://example.com/token",
					RefreshURL:       util.PtrOf("https://example.com/refresh"),
					Scopes: map[string]string{
						"profile": "Profile access",
						"email":   "Email access",
					},
					MapOfAnything: map[string]interface{}{
						"x-flow-type": "authorization-code",
					},
				},
			},
			expected3: openapi3.OAuthFlows{
				AuthorizationCode: &openapi3.AuthorizationCodeOAuthFlow{
					AuthorizationURL: "https://example.com/auth",
					TokenURL:         "https://example.com/token",
					RefreshURL:       util.PtrOf("https://example.com/refresh"),
					Scopes: map[string]string{
						"profile": "Profile access",
						"email":   "Email access",
					},
					MapOfAnything: map[string]interface{}{
						"x-flow-type": "authorization-code",
					},
				},
			},
			expected31: openapi31.OauthFlows{
				Implicit:          nil,
				Password:          nil,
				ClientCredentials: nil,
				AuthorizationCode: &openapi31.OauthFlowsDefsAuthorizationCode{
					AuthorizationURL: "https://example.com/auth",
					TokenURL:         "https://example.com/token",
					RefreshURL:       util.PtrOf("https://example.com/refresh"),
					Scopes: map[string]string{
						"profile": "Profile access",
						"email":   "Email access",
					},
					MapOfAnything: map[string]interface{}{
						"x-flow-type": "authorization-code",
					},
				},
			},
		},
		{
			name: "flows with all flow types",
			flows: openapi.OAuthFlows{
				Implicit: &openapi.OAuthFlowsDefsImplicit{
					AuthorizationURL: "https://example.com/implicit/auth",
					Scopes: map[string]string{
						"read": "Read access",
					},
				},
				Password: &openapi.OAuthFlowsDefsPassword{
					TokenURL: "https://example.com/password/token",
					Scopes: map[string]string{
						"write": "Write access",
					},
				},
				ClientCredentials: &openapi.OAuthFlowsDefsClientCredentials{
					TokenURL: "https://example.com/client/token",
					Scopes: map[string]string{
						"admin": "Admin access",
					},
				},
				AuthorizationCode: &openapi.OAuthFlowsDefsAuthorizationCode{
					AuthorizationURL: "https://example.com/code/auth",
					TokenURL:         "https://example.com/code/token",
					Scopes: map[string]string{
						"full": "Full access",
					},
				},
			},
			expected3: openapi3.OAuthFlows{
				Implicit: &openapi3.ImplicitOAuthFlow{
					AuthorizationURL: "https://example.com/implicit/auth",
					Scopes: map[string]string{
						"read": "Read access",
					},
				},
				Password: &openapi3.PasswordOAuthFlow{
					TokenURL: "https://example.com/password/token",
					Scopes: map[string]string{
						"write": "Write access",
					},
				},
				ClientCredentials: &openapi3.ClientCredentialsFlow{
					TokenURL: "https://example.com/client/token",
					Scopes: map[string]string{
						"admin": "Admin access",
					},
				},
				AuthorizationCode: &openapi3.AuthorizationCodeOAuthFlow{
					AuthorizationURL: "https://example.com/code/auth",
					TokenURL:         "https://example.com/code/token",
					Scopes: map[string]string{
						"full": "Full access",
					},
				},
			},
			expected31: openapi31.OauthFlows{
				Implicit: &openapi31.OauthFlowsDefsImplicit{
					AuthorizationURL: "https://example.com/implicit/auth",
					Scopes: map[string]string{
						"read": "Read access",
					},
				},
				Password: &openapi31.OauthFlowsDefsPassword{
					TokenURL: "https://example.com/password/token",
					Scopes: map[string]string{
						"write": "Write access",
					},
				},
				ClientCredentials: &openapi31.OauthFlowsDefsClientCredentials{
					TokenURL: "https://example.com/client/token",
					Scopes: map[string]string{
						"admin": "Admin access",
					},
				},
				AuthorizationCode: &openapi31.OauthFlowsDefsAuthorizationCode{
					AuthorizationURL: "https://example.com/code/auth",
					TokenURL:         "https://example.com/code/token",
					Scopes: map[string]string{
						"full": "Full access",
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result3 := mapper.OAS3Oauth2Flows(tt.flows)
			assert.Equal(t, tt.expected3, result3, "OpenAPI 3 OAuth2 Flows mapping failed")
			result31 := mapper.OAS31Oauth2Flows(tt.flows)
			assert.Equal(t, tt.expected31, result31, "OpenAPI 3.1 OAuth2 Flows mapping failed")
		})
	}
}

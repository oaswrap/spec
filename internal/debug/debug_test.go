package debug

import (
	"fmt"
	"testing"

	"github.com/oaswrap/spec/openapi"
	"github.com/stretchr/testify/assert"
)

func TestNewLogger(t *testing.T) {
	mockLogger := &mockLogger{}
	logger := NewLogger(mockLogger)

	assert.NotNil(t, logger)
	assert.Equal(t, "[spec]", logger.prefix)
	assert.Equal(t, mockLogger, logger.logger)
}

func TestLogger_Printf(t *testing.T) {
	mockLogger := &mockLogger{}
	logger := NewLogger(mockLogger)

	logger.Printf("test message %s", "value")

	assert.Equal(t, "[spec] test message value", mockLogger.lastMessage)
}

func TestLogger_LogOp(t *testing.T) {
	mockLogger := &mockLogger{}
	logger := NewLogger(mockLogger)

	logger.LogOp("GET", "/users", "added", "endpoint")

	assert.Equal(t, "[spec] GET /users → added: endpoint", mockLogger.lastMessage)
}

func TestLogger_LogAction(t *testing.T) {
	mockLogger := &mockLogger{}
	logger := NewLogger(mockLogger)

	logger.LogAction("validate", "schema")

	assert.Equal(t, "[spec] validate: schema", mockLogger.lastMessage)
}

func TestLogger_LogContact(t *testing.T) {
	mockLogger := &mockLogger{}
	logger := NewLogger(mockLogger)

	t.Run("nil contact", func(t *testing.T) {
		logger.LogContact(nil)
		assert.Empty(t, mockLogger.lastMessage)
	})

	t.Run("full contact", func(t *testing.T) {
		contact := &openapi.Contact{
			Name:  "John Doe",
			Email: "john@example.com",
			URL:   "https://example.com",
		}
		logger.LogContact(contact)
		assert.Equal(t, "[spec] set contact: name: John Doe, email: john@example.com, url: https://example.com", mockLogger.lastMessage)
	})

	t.Run("partial contact", func(t *testing.T) {
		contact := &openapi.Contact{
			Name: "Jane Doe",
		}
		logger.LogContact(contact)
		assert.Equal(t, "[spec] set contact: name: Jane Doe, ", mockLogger.lastMessage)
	})
}

func TestLogger_LogLicense(t *testing.T) {
	mockLogger := &mockLogger{}
	logger := NewLogger(mockLogger)

	t.Run("full license", func(t *testing.T) {
		license := &openapi.License{
			Name: "MIT",
			URL:  "https://opensource.org/licenses/MIT",
		}
		logger.LogLicense(license)
		assert.Equal(t, "[spec] set license: name: MIT, url: https://opensource.org/licenses/MIT", mockLogger.lastMessage)
	})

	t.Run("name only", func(t *testing.T) {
		license := &openapi.License{
			Name: "Apache 2.0",
		}
		logger.LogLicense(license)
		assert.Equal(t, "[spec] set license: name: Apache 2.0, ", mockLogger.lastMessage)
	})
}

func TestLogger_LogExternalDocs(t *testing.T) {
	mockLogger := &mockLogger{}
	logger := NewLogger(mockLogger)

	t.Run("full external docs", func(t *testing.T) {
		docs := &openapi.ExternalDocs{
			URL:         "https://docs.example.com",
			Description: "API Documentation",
		}
		logger.LogExternalDocs(docs)
		assert.Equal(t, "[spec] set external docs: url: https://docs.example.com, description: API Documentation", mockLogger.lastMessage)
	})

	t.Run("url only", func(t *testing.T) {
		docs := &openapi.ExternalDocs{
			URL: "https://docs.example.com",
		}
		logger.LogExternalDocs(docs)
		assert.Equal(t, "[spec] set external docs: url: https://docs.example.com", mockLogger.lastMessage)
	})
}

func TestLogger_LogServer(t *testing.T) {
	mockLogger := &mockLogger{}
	logger := NewLogger(mockLogger)

	t.Run("simple server", func(t *testing.T) {
		server := openapi.Server{
			URL: "https://api.example.com",
		}
		logger.LogServer(server)
		assert.Equal(t, "[spec] set server: url: https://api.example.com", mockLogger.lastMessage)
	})

	t.Run("server with description", func(t *testing.T) {
		desc := "Production server"
		server := openapi.Server{
			URL:         "https://api.example.com",
			Description: &desc,
		}
		logger.LogServer(server)
		assert.Equal(t, "[spec] set server: url: https://api.example.com, description: Production server", mockLogger.lastMessage)
	})

	t.Run("server with variables", func(t *testing.T) {
		server := openapi.Server{
			URL: "https://{env}.example.com",
			Variables: map[string]openapi.ServerVariable{
				"env": {Default: "api"},
			},
		}
		logger.LogServer(server)
		assert.Equal(t, "[spec] set server: url: https://{env}.example.com, variables: env: api", mockLogger.lastMessage)
	})
}

func TestLogger_LogTag(t *testing.T) {
	mockLogger := &mockLogger{}
	logger := NewLogger(mockLogger)

	t.Run("simple tag", func(t *testing.T) {
		tag := openapi.Tag{
			Name: "users",
		}
		logger.LogTag(tag)
		assert.Equal(t, "[spec] add tag: name: users", mockLogger.lastMessage)
	})

	t.Run("tag with description", func(t *testing.T) {
		tag := openapi.Tag{
			Name:        "users",
			Description: "User management endpoints",
		}
		logger.LogTag(tag)
		assert.Equal(t, "[spec] add tag: name: users, description: User management endpoints", mockLogger.lastMessage)
	})

	t.Run("tag with external docs", func(t *testing.T) {
		tag := openapi.Tag{
			Name: "users",
			ExternalDocs: &openapi.ExternalDocs{
				URL:         "https://docs.example.com/users",
				Description: "User API docs",
			},
		}
		logger.LogTag(tag)
		assert.Equal(t, "[spec] add tag: name: users, external docs: https://docs.example.com/users (User API docs)", mockLogger.lastMessage)
	})
}

func TestLogger_LogSecurityScheme(t *testing.T) {
	mockLogger := &mockLogger{}
	logger := NewLogger(mockLogger)

	t.Run("APIKey scheme", func(t *testing.T) {
		scheme := &openapi.SecurityScheme{
			APIKey: &openapi.SecuritySchemeAPIKey{},
		}
		logger.LogSecurityScheme("api_key", scheme)
		assert.Equal(t, "[spec] add security scheme: name: api_key, type: APIKey", mockLogger.lastMessage)
	})

	t.Run("HTTPBearer scheme", func(t *testing.T) {
		desc := "Bearer token authentication"
		scheme := &openapi.SecurityScheme{
			HTTPBearer:  &openapi.SecuritySchemeHTTPBearer{},
			Description: &desc,
		}
		logger.LogSecurityScheme("bearer", scheme)
		assert.Equal(t, "[spec] add security scheme: name: bearer, type: HTTPBearer, description: Bearer token authentication", mockLogger.lastMessage)
	})

	t.Run("OAuth2 scheme", func(t *testing.T) {
		scheme := &openapi.SecurityScheme{
			OAuth2: &openapi.SecuritySchemeOAuth2{},
		}
		logger.LogSecurityScheme("oauth2", scheme)
		assert.Equal(t, "[spec] add security scheme: name: oauth2, type: OAuth2", mockLogger.lastMessage)
	})

	t.Run("unknown scheme", func(t *testing.T) {
		scheme := &openapi.SecurityScheme{}
		logger.LogSecurityScheme("unknown", scheme)
		assert.Equal(t, "[spec] add security scheme: name: unknown, type: Unknown", mockLogger.lastMessage)
	})
}

// mockLogger implements openapi.Logger for testing
type mockLogger struct {
	lastMessage string
}

func (m *mockLogger) Printf(format string, v ...any) {
	m.lastMessage = fmt.Sprintf(format, v...)
}

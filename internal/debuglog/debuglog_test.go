package debuglog

import (
	"fmt"
	"testing"

	"github.com/oaswrap/spec/openapi"
	"github.com/stretchr/testify/assert"
)

type mockLogger struct {
	messages []string
}

func (m *mockLogger) Printf(format string, v ...any) {
	m.messages = append(m.messages, fmt.Sprintf(format, v...))
}

func TestNewLogger(t *testing.T) {
	mockLog := &mockLogger{}
	logger := NewLogger("test", mockLog)

	assert.Equal(t, "[test]", logger.prefix)
	assert.Equal(t, mockLog, logger.logger)
}

func TestLogger_Printf(t *testing.T) {
	mockLog := &mockLogger{}
	logger := NewLogger("test", mockLog)

	logger.Printf("Hello %s", "world")

	expected := "[test] Hello world"
	assert.Len(t, mockLog.messages, 1)
	assert.Equal(t, expected, mockLog.messages[0])
}

func TestLogger_LogOp(t *testing.T) {
	mockLog := &mockLogger{}
	logger := NewLogger("api", mockLog)

	logger.LogOp("GET", "/users", "fetch", "all users")

	expected := "[api] GET /users → fetch: all users"
	assert.Len(t, mockLog.messages, 1)
	assert.Equal(t, expected, mockLog.messages[0])
}

func TestLogger_LogAction(t *testing.T) {
	mockLog := &mockLogger{}
	logger := NewLogger("test", mockLog)

	logger.LogAction("validate", "schema")

	expected := "[test] validate: schema"
	assert.Len(t, mockLog.messages, 1)
	assert.Equal(t, expected, mockLog.messages[0])
}

func TestLogger_LogContact(t *testing.T) {
	mockLog := &mockLogger{}
	logger := NewLogger("test", mockLog)

	// Test with nil contact
	logger.LogContact(nil)
	assert.Empty(t, mockLog.messages)

	// Test with full contact
	contact := &openapi.Contact{
		Name:  "John Doe",
		Email: "john@example.com",
		URL:   "https://example.com",
	}
	logger.LogContact(contact)

	expected := "[test] set contact: name: John Doe, email: john@example.com, url: https://example.com"
	assert.Len(t, mockLog.messages, 1)
	assert.Equal(t, expected, mockLog.messages[0])
}

func TestLogger_LogLicense(t *testing.T) {
	mockLog := &mockLogger{}
	logger := NewLogger("test", mockLog)

	license := &openapi.License{
		Name: "MIT",
		URL:  "https://opensource.org/licenses/MIT",
	}
	logger.LogLicense(license)

	expected := "[test] set license: name: MIT, url: https://opensource.org/licenses/MIT"
	assert.Len(t, mockLog.messages, 1)
	assert.Equal(t, expected, mockLog.messages[0])
}

func TestLogger_LogExternalDocs(t *testing.T) {
	mockLog := &mockLogger{}
	logger := NewLogger("test", mockLog)

	docs := &openapi.ExternalDocs{
		URL:         "https://docs.example.com",
		Description: "API Documentation",
	}
	logger.LogExternalDocs(docs)

	expected := "[test] set external docs: url: https://docs.example.com, description: API Documentation"
	assert.Len(t, mockLog.messages, 1)
	assert.Equal(t, expected, mockLog.messages[0])
}

func TestLogger_LogServer(t *testing.T) {
	mockLog := &mockLogger{}
	logger := NewLogger("test", mockLog)

	desc := "Production server"
	server := openapi.Server{
		URL:         "https://api.example.com",
		Description: &desc,
		Variables: map[string]openapi.ServerVariable{
			"version": {Default: "v1"},
			"env":     {Default: "prod"},
		},
	}
	logger.LogServer(server)

	assert.Len(t, mockLog.messages, 1)
	message := mockLog.messages[0]
	assert.Contains(t, message, "[test] set server: url: https://api.example.com")
	assert.Contains(t, message, "description: Production server")
	assert.Contains(t, message, "variables:")
}

func TestLogger_LogTag(t *testing.T) {
	mockLog := &mockLogger{}
	logger := NewLogger("test", mockLog)

	tag := openapi.Tag{
		Name:        "users",
		Description: "User operations",
		ExternalDocs: &openapi.ExternalDocs{
			URL:         "https://docs.example.com/users",
			Description: "User docs",
		},
	}
	logger.LogTag(tag)

	expected := "[test] add tag: name: users, description: User operations, external docs: https://docs.example.com/users (User docs)"
	assert.Len(t, mockLog.messages, 1)
	assert.Equal(t, expected, mockLog.messages[0])
}

func TestLogger_LogSecurityScheme(t *testing.T) {
	mockLog := &mockLogger{}
	logger := NewLogger("test", mockLog)

	desc := "API Key authentication"
	scheme := &openapi.SecurityScheme{
		APIKey:      &openapi.SecuritySchemeAPIKey{},
		Description: &desc,
	}
	logger.LogSecurityScheme("apiKey", scheme)

	expected := "[test] add security scheme: name: apiKey, type: APIKey, description: API Key authentication"
	assert.Len(t, mockLog.messages, 1)
	assert.Equal(t, expected, mockLog.messages[0])

	desc = "HTTP Bearer authentication"
	scheme = &openapi.SecurityScheme{
		HTTPBearer:  &openapi.SecuritySchemeHTTPBearer{},
		Description: &desc,
	}
	logger.LogSecurityScheme("bearer", scheme)
	expected = "[test] add security scheme: name: bearer, type: HTTPBearer, description: HTTP Bearer authentication"
	assert.Len(t, mockLog.messages, 2)
	assert.Equal(t, expected, mockLog.messages[1])

	desc = "OAuth 2.0 authentication"
	scheme = &openapi.SecurityScheme{
		OAuth2:      &openapi.SecuritySchemeOAuth2{},
		Description: &desc,
	}
	logger.LogSecurityScheme("oauth2", scheme)
	expected = "[test] add security scheme: name: oauth2, type: OAuth2, description: OAuth 2.0 authentication"
	assert.Len(t, mockLog.messages, 3)
	assert.Equal(t, expected, mockLog.messages[2])
}

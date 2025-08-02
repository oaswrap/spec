package openapi

// Contact represents contact information for the API.
// Generated from "#/$defs/contact".
type Contact struct {
	Name  string // Contact name.
	URL   string // Contact URL. Format: uri.
	Email string // Contact email. Format: email.

	// MapOfAnything holds vendor extensions. Keys must match `^x-`.
	MapOfAnything map[string]any
}

// License provides license information for the API.
// Generated from "#/$defs/license".
type License struct {
	Name string // License name (required).

	Identifier string // SPDX identifier.
	URL        string // License URL. Format: uri.

	// MapOfAnything holds vendor extensions. Keys must match `^x-`.
	MapOfAnything map[string]any
}

// Tag adds metadata to an API operation.
// Generated from "#/definitions/Tag".
type Tag struct {
	Name          string         // Tag name (required).
	Description   string         // Tag description.
	ExternalDocs  *ExternalDocs  // Additional external documentation.
	MapOfAnything map[string]any // Vendor extensions. Keys must match `^x-`.
}

// ExternalDocs describes external documentation for a tag or operation.
// Generated from "#/$defs/external-documentation".
type ExternalDocs struct {
	Description string // Description of the documentation.

	URL string // Required. Documentation URL. Format: uri.

	MapOfAnything map[string]any // Vendor extensions. Keys must match `^x-`.
}

// Server describes an API server.
// Generated from "#/$defs/server".
type Server struct {
	URL string // Required. Server URL. Format: uri-reference.

	Description *string                   // Optional server description.
	Variables   map[string]ServerVariable // Server variables for URL templates.

	MapOfAnything map[string]any // Vendor extensions. Keys must match `^x-`.
}

// ServerVariable describes a variable for server URL template substitution.
// Generated from "#/$defs/server-variable".
type ServerVariable struct {
	Enum    []string // Allowed values.
	Default string   // Required. Default value.

	Description   string         // Variable description.
	MapOfAnything map[string]any // Vendor extensions. Keys must match `^x-`.
}

// SecurityScheme describes a security scheme that can be used by operations.
// Generated from "#/$defs/security-scheme".
type SecurityScheme struct {
	Description *string                   // Optional description.
	APIKey      *SecuritySchemeAPIKey     // API key authentication scheme.
	HTTPBearer  *SecuritySchemeHTTPBearer // HTTP Bearer authentication scheme.
	OAuth2      *SecuritySchemeOAuth2     // OAuth2 authentication scheme.

	MapOfAnything map[string]any // Vendor extensions. Keys must match `^x-`.
}

// SecuritySchemeAPIKey defines an API key authentication scheme.
// Generated from "#/$defs/security-scheme/$defs/type-apikey".
type SecuritySchemeAPIKey struct {
	Name string                 // Required. Name of the header, query, or cookie parameter.
	In   SecuritySchemeAPIKeyIn // Required. Location of the API key.
}

// SecuritySchemeAPIKeyIn specifies where the API key is passed.
type SecuritySchemeAPIKeyIn string

const (
	// SecuritySchemeAPIKeyInQuery specifies the key is passed in the query string.
	SecuritySchemeAPIKeyInQuery = SecuritySchemeAPIKeyIn("query")
	// SecuritySchemeAPIKeyInHeader specifies the key is passed in a header.
	SecuritySchemeAPIKeyInHeader = SecuritySchemeAPIKeyIn("header")
	// SecuritySchemeAPIKeyInCookie specifies the key is passed in a cookie.
	SecuritySchemeAPIKeyInCookie = SecuritySchemeAPIKeyIn("cookie")
)

// SecuritySchemeHTTPBearer defines HTTP Bearer authentication.
// Generated from "#/$defs/security-scheme/$defs/type-http-bearer".
type SecuritySchemeHTTPBearer struct {
	Scheme       string  // Required. Must match pattern `^[Bb][Ee][Aa][Rr][Ee][Rr]$`.
	BearerFormat *string // Optional bearer format hint.
}

// SecuritySchemeOAuth2 defines OAuth2 flows.
// Generated from "#/$defs/security-scheme/$defs/type-oauth2".
type SecuritySchemeOAuth2 struct {
	Flows OAuthFlows // Required. Supported OAuth2 flows.
}

// OAuthFlows groups supported OAuth2 flows.
// Generated from "#/$defs/oauth-flows".
type OAuthFlows struct {
	Implicit          *OAuthFlowsImplicit
	Password          *OAuthFlowsPassword
	ClientCredentials *OAuthFlowsClientCredentials
	AuthorizationCode *OAuthFlowsAuthorizationCode

	MapOfAnything map[string]any // Vendor extensions. Keys must match `^x-`.
}

// OAuthFlowsImplicit defines an OAuth2 implicit flow.
// Generated from "#/$defs/oauth-flows/$defs/implicit".
type OAuthFlowsImplicit struct {
	AuthorizationURL string            // Required. Format: uri.
	RefreshURL       *string           // Optional refresh URL. Format: uri.
	Scopes           map[string]string // Required. Available scopes.

	MapOfAnything map[string]any // Vendor extensions. Keys must match `^x-`.
}

// OAuthFlowsPassword defines an OAuth2 password flow.
// Generated from "#/$defs/oauth-flows/$defs/password".
type OAuthFlowsPassword struct {
	TokenURL   string            // Required. Token URL. Format: uri.
	RefreshURL *string           // Optional refresh URL. Format: uri.
	Scopes     map[string]string // Required. Available scopes.

	MapOfAnything map[string]any // Vendor extensions. Keys must match `^x-`.
}

// OAuthFlowsClientCredentials defines an OAuth2 client credentials flow.
// Generated from "#/$defs/oauth-flows/$defs/client-credentials".
type OAuthFlowsClientCredentials struct {
	TokenURL   string            // Required. Token URL. Format: uri.
	RefreshURL *string           // Optional refresh URL. Format: uri.
	Scopes     map[string]string // Required. Available scopes.

	MapOfAnything map[string]any // Vendor extensions. Keys must match `^x-`.
}

// OAuthFlowsAuthorizationCode defines an OAuth2 authorization code flow.
// Generated from "#/$defs/oauth-flows/$defs/authorization-code".
type OAuthFlowsAuthorizationCode struct {
	AuthorizationURL string            // Required. Format: uri.
	TokenURL         string            // Required. Token URL. Format: uri.
	RefreshURL       *string           // Optional refresh URL. Format: uri.
	Scopes           map[string]string // Required. Available scopes.

	MapOfAnything map[string]any // Vendor extensions. Keys must match `^x-`.
}

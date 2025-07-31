package openapi

// Contact structure is generated from "#/$defs/contact".
type Contact struct {
	Name          string
	URL           string         // Format: uri.
	Email         string         // Format: email.
	MapOfAnything map[string]any // Key must match pattern: `^x-`.
}

// License structure is generated from "#/$defs/license".
type License struct {
	Name          string // Required.
	Identifier    string
	URL           string         // Format: uri.
	MapOfAnything map[string]any // Key must match pattern: `^x-`.
}

// Tag structure is generated from "#/definitions/Tag".
type Tag struct {
	Name          string // Required.
	Description   string
	ExternalDocs  *ExternalDocs
	MapOfAnything map[string]any // Key must match pattern: `^x-`.
}

// ExternalDocs structure is generated from "#/$defs/external-documentation".
type ExternalDocs struct {
	Description string
	// Format: uri.
	// Required.
	URL           string
	MapOfAnything map[string]any // Key must match pattern: `^x-`.
}

// Server structure is generated from "#/$defs/server".
type Server struct {
	// Format: uri-reference.
	// Required.
	URL           string
	Description   *string
	Variables     map[string]ServerVariable
	MapOfAnything map[string]any // Key must match pattern: `^x-`.
}

// ServerVariable structure is generated from "#/$defs/server-variable".
type ServerVariable struct {
	Enum          []string
	Default       string // Required.
	Description   string
	MapOfAnything map[string]any // Key must match pattern: `^x-`.
}

// SecurityScheme structure is generated from "#/$defs/security-scheme".
type SecurityScheme struct {
	Description   *string
	APIKey        *SecuritySchemeAPIKey
	HTTPBearer    *SecuritySchemeHTTPBearer
	OAuth2        *SecuritySchemeOAuth2
	MapOfAnything map[string]any // Key must match pattern: `^x-`.
}

// SecuritySchemeAPIKey structure is generated from "#/$defs/security-scheme/$defs/type-apikey".
type SecuritySchemeAPIKey struct {
	Name string                 // Required.
	In   SecuritySchemeAPIKeyIn // Required.
}

// SecuritySchemeAPIKeyIn is an enum type.
type SecuritySchemeAPIKeyIn string

// SecuritySchemeAPIKeyIn values enumeration.
const (
	SecuritySchemeAPIKeyInQuery  = SecuritySchemeAPIKeyIn("query")
	SecuritySchemeAPIKeyInHeader = SecuritySchemeAPIKeyIn("header")
	SecuritySchemeAPIKeyInCookie = SecuritySchemeAPIKeyIn("cookie")
)

// SecuritySchemeHTTPBearer structure is generated from "#/$defs/security-scheme/$defs/type-http-bearer".
type SecuritySchemeHTTPBearer struct {
	// Value must match pattern: `^[Bb][Ee][Aa][Rr][Ee][Rr]$`.
	// Required.
	Scheme       string
	BearerFormat *string
}

// SecuritySchemeOAuth2 structure is generated from "#/$defs/security-scheme/$defs/type-oauth2".
type SecuritySchemeOAuth2 struct {
	Flows OAuthFlows // Required.
}

// OAuthFlows structure is generated from "#/$defs/oauth-flows".
type OAuthFlows struct {
	Implicit          *OAuthFlowsDefsImplicit
	Password          *OAuthFlowsDefsPassword
	ClientCredentials *OAuthFlowsDefsClientCredentials
	AuthorizationCode *OAuthFlowsDefsAuthorizationCode
	MapOfAnything     map[string]any // Key must match pattern: `^x-`.
}

// OAuthFlowsDefsImplicit structure is generated from "#/$defs/oauth-flows/$defs/implicit".
type OAuthFlowsDefsImplicit struct {
	// Format: uri.
	// Required.
	AuthorizationURL string
	RefreshURL       *string           // Format: uri.
	Scopes           map[string]string // Required.
	MapOfAnything    map[string]any    // Key must match pattern: `^x-`.
}

// OAuthFlowsDefsPassword structure is generated from "#/$defs/oauth-flows/$defs/password".
type OAuthFlowsDefsPassword struct {
	// Format: uri.
	// Required.
	TokenURL      string
	RefreshURL    *string           // Format: uri.
	Scopes        map[string]string // Required.
	MapOfAnything map[string]any    // Key must match pattern: `^x-`.
}

// OAuthFlowsDefsClientCredentials structure is generated from "#/$defs/oauth-flows/$defs/client-credentials".
type OAuthFlowsDefsClientCredentials struct {
	// Format: uri.
	// Required.
	TokenURL      string
	RefreshURL    *string           // Format: uri.
	Scopes        map[string]string // Required.
	MapOfAnything map[string]any    // Key must match pattern: `^x-`.
}

// OAuthFlowsDefsAuthorizationCode structure is generated from "#/$defs/oauth-flows/$defs/authorization-code".
type OAuthFlowsDefsAuthorizationCode struct {
	// Format: uri.
	// Required.
	AuthorizationURL string
	// Format: uri.
	// Required.
	TokenURL      string
	RefreshURL    *string           // Format: uri.
	Scopes        map[string]string // Required.
	MapOfAnything map[string]any    // Key must match pattern: `^x-`.
}

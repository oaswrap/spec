package option

import "github.com/oaswrap/spec/pkg/util"

// OperationConfig holds the configuration for an OpenAPI operation.
type OperationConfig struct {
	Hide        bool
	OperationID string
	Description string
	Summary     string
	Deprecated  bool
	Tags        []string
	Security    []RouteSecurityConfig

	Requests  []*ContentConfig
	Responses []*ContentConfig
}

// Operation is a function that configures an OpenAPI operation.
type OperationOption func(*OperationConfig)

// Hide marks the operation as hidden in the OpenAPI documentation.
// This is useful for operations that should not be exposed to the public API.
func Hide(hide ...bool) OperationOption {
	return func(cfg *OperationConfig) {
		cfg.Hide = util.Optional(true, hide...)
	}
}

// OperationID sets the operation ID for the OpenAPI operation.
func OperationID(id string) OperationOption {
	return func(cfg *OperationConfig) {
		cfg.OperationID = id
	}
}

// Description sets the description for the OpenAPI operation.
func Description(description string) OperationOption {
	return func(cfg *OperationConfig) {
		cfg.Description = description
	}
}

// Summary sets the summary for the OpenAPI operation.
func Summary(summary string) OperationOption {
	return func(cfg *OperationConfig) {
		cfg.Summary = summary
		if cfg.Description == "" {
			cfg.Description = summary // Use summary as description if none provided
		}
	}
}

// Deprecated marks the operation as deprecated in the OpenAPI documentation.
func Deprecated(deprecated ...bool) OperationOption {
	return func(cfg *OperationConfig) {
		cfg.Deprecated = util.Optional(true, deprecated...)
	}
}

// Tags adds tags to the OpenAPI operation.
func Tags(tags ...string) OperationOption {
	return func(cfg *OperationConfig) {
		cfg.Tags = append(cfg.Tags, tags...)
	}
}

// Security adds security requirements to the OpenAPI operation.
func Security(securityName string, scopes ...string) OperationOption {
	return func(cfg *OperationConfig) {
		cfg.Security = append(cfg.Security, RouteSecurityConfig{
			Name:   securityName,
			Scopes: scopes,
		})
	}
}

// Request adds a request structure to the OpenAPI operation.
func Request(structure any, options ...ContentOption) OperationOption {
	return func(cfg *OperationConfig) {
		contentConfig := &ContentConfig{
			Structure: structure,
		}
		for _, opt := range options {
			opt(contentConfig)
		}
		cfg.Requests = append(cfg.Requests, contentConfig)
	}
}

// Response adds a response structure to the OpenAPI operation.
func Response(httpStatus int, structure any, options ...ContentOption) OperationOption {
	return func(cfg *OperationConfig) {
		contentConfig := &ContentConfig{
			HTTPStatus: httpStatus,
			Structure:  structure,
		}
		for _, opt := range options {
			opt(contentConfig)
		}
		cfg.Responses = append(cfg.Responses, contentConfig)
	}
}

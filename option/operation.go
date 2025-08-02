package option

import "github.com/oaswrap/spec/pkg/util"

// OperationConfig holds configuration for an OpenAPI operation.
type OperationConfig struct {
	Hide        bool
	OperationID string
	Description string
	Summary     string
	Deprecated  bool
	Tags        []string
	Security    []OperationSecurityConfig

	Requests  []*ContentUnit
	Responses []*ContentUnit
}

// OperationSecurityConfig defines a security requirement for an operation.
type OperationSecurityConfig struct {
	Name   string
	Scopes []string
}

// OperationOption applies configuration to an OpenAPI operation.
type OperationOption func(*OperationConfig)

// Hide marks the operation as hidden in the OpenAPI documentation.
//
// This is useful for internal or non-public endpoints.
func Hide(hide ...bool) OperationOption {
	return func(cfg *OperationConfig) {
		cfg.Hide = util.Optional(true, hide...)
	}
}

// OperationID sets the unique operation ID for the OpenAPI operation.
func OperationID(id string) OperationOption {
	return func(cfg *OperationConfig) {
		cfg.OperationID = id
	}
}

// Description sets the detailed description for the OpenAPI operation.
func Description(description string) OperationOption {
	return func(cfg *OperationConfig) {
		cfg.Description = description
	}
}

// Summary sets a short summary for the OpenAPI operation.
//
// If no description is set, the summary is also used as the description.
func Summary(summary string) OperationOption {
	return func(cfg *OperationConfig) {
		cfg.Summary = summary
		if cfg.Description == "" {
			cfg.Description = summary
		}
	}
}

// Deprecated marks the operation as deprecated.
//
// Deprecated operations should not be used by clients.
func Deprecated(deprecated ...bool) OperationOption {
	return func(cfg *OperationConfig) {
		cfg.Deprecated = util.Optional(true, deprecated...)
	}
}

// Tags adds tags to the OpenAPI operation.
//
// Tags help organize operations in the generated documentation.
func Tags(tags ...string) OperationOption {
	return func(cfg *OperationConfig) {
		cfg.Tags = append(cfg.Tags, tags...)
	}
}

// Security adds a security requirement to the OpenAPI operation.
//
// Example:
//
//	r.Get("/me",
//	    option.Security("bearerAuth"),
//	)
func Security(securityName string, scopes ...string) OperationOption {
	return func(cfg *OperationConfig) {
		cfg.Security = append(cfg.Security, OperationSecurityConfig{
			Name:   securityName,
			Scopes: scopes,
		})
	}
}

// Request adds a request body or parameter structure to the OpenAPI operation.
func Request(structure any, options ...ContentOption) OperationOption {
	return func(cfg *OperationConfig) {
		cu := &ContentUnit{
			Structure: structure,
		}
		for _, opt := range options {
			opt(cu)
		}
		cfg.Requests = append(cfg.Requests, cu)
	}
}

// Response adds a response for the OpenAPI operation.
//
// The HTTP status code defines which response is described.
func Response(httpStatus int, structure any, options ...ContentOption) OperationOption {
	return func(cfg *OperationConfig) {
		cu := &ContentUnit{
			HTTPStatus: httpStatus,
			Structure:  structure,
		}
		for _, opt := range options {
			opt(cu)
		}
		cfg.Responses = append(cfg.Responses, cu)
	}
}

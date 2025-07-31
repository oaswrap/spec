package option

import "github.com/oaswrap/spec/openapi"

// ServerOption is a function that applies configuration to an OpenAPI server.
type ServerOption func(*openapi.Server)

// ServerDescription sets the description for an OpenAPI server.
//
// Example usage:
//
//	option.WithServer("https://api.example.com", option.ServerDescription("Production server"))
func ServerDescription(description string) ServerOption {
	return func(s *openapi.Server) {
		s.Description = &description
	}
}

// ServerVariables sets the variables for an OpenAPI server.
//
// Example usage:
//
//	option.WithServer("https://api.example.com/{version}", option.ServerVariables(map[string]openapi.ServerVariable{
//	    "version": {
//	        Default:     "v1",
//	        Description: "API version",
//	        Enum:        []string{"v1", "v2"},
//	    },
//	}))
func ServerVariables(variables map[string]openapi.ServerVariable) ServerOption {
	return func(s *openapi.Server) {
		s.Variables = variables
	}
}

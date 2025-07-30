package option

import "github.com/oaswrap/spec/openapi"

type ServerOption func(*openapi.Server)

func ServerDescription(description string) ServerOption {
	return func(s *openapi.Server) {
		s.Description = &description
	}
}

func ServerVariables(variables map[string]openapi.ServerVariable) ServerOption {
	return func(s *openapi.Server) {
		s.Variables = variables
	}
}

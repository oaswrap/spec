package option

import "github.com/oaswrap/spec"

type ServerOption func(*spec.Server)

func ServerDescription(description string) ServerOption {
	return func(s *spec.Server) {
		s.Description = &description
	}
}

func ServerVariables(variables map[string]spec.ServerVariable) ServerOption {
	return func(s *spec.Server) {
		s.Variables = variables
	}
}

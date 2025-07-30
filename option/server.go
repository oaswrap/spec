package option

type ServerOption func(*Server)

func ServerDescription(description string) ServerOption {
	return func(s *Server) {
		s.Description = &description
	}
}

func ServerVariables(variables map[string]ServerVariable) ServerOption {
	return func(s *Server) {
		s.Variables = variables
	}
}

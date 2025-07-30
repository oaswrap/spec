package option

import "github.com/oaswrap/spec/pkg/util"

// Operation is a function that configures an OpenAPI operation.
type RouteConfig struct {
	Tags     []string
	Security []RouteSecurityConfig
	Hide     bool
}

// RouteSecurityConfig holds the security configuration for a route.
type RouteSecurityConfig struct {
	Name   string
	Scopes []string
}

// RouteOption is a function that applies configuration to a RouteConfig.
type RouteOption func(*RouteConfig)

// RouteTags adds tags to the route.
//
// It will add tags to all routes in the sub-router.
func RouteTags(tags ...string) RouteOption {
	return func(cfg *RouteConfig) {
		cfg.Tags = append(cfg.Tags, tags...)
	}
}

// RouteSecurity adds security schemes to the route.
//
// It will add security schemes to all routes in the sub-router.
func RouteSecurity(securityName string, scopes ...string) RouteOption {
	return func(cfg *RouteConfig) {
		cfg.Security = append(cfg.Security, RouteSecurityConfig{
			Name:   securityName,
			Scopes: scopes,
		})
	}
}

// RouteHide marks the route as hidden in the OpenAPI documentation.
// This is useful for routes that should not be exposed to the public API.
//
// It will hide all routes in the sub-router.
// If you want to hide only specific routes, use the `With` method on the route
func RouteHide(hide ...bool) RouteOption {
	return func(cfg *RouteConfig) {
		cfg.Hide = util.Optional(true, hide...)
	}
}

package option

import "github.com/oaswrap/spec/internal/util"

// GroupConfig holds the configuration for a group of routes in an OpenAPI specification.
type GroupConfig struct {
	Tags     []string
	Security []OperationSecurityConfig
	Hide     bool
}

// GroupOption is a function that applies configuration to a GroupConfig.
type GroupOption func(*GroupConfig)

// GroupTags adds tags to the group.
//
// It will add tags to all routes in the sub-router.
func GroupTags(tags ...string) GroupOption {
	return func(cfg *GroupConfig) {
		cfg.Tags = append(cfg.Tags, tags...)
	}
}

// GroupSecurity adds security schemes to the group.
//
// It will add security schemes to all routes in the sub-router.
func GroupSecurity(securityName string, scopes ...string) GroupOption {
	return func(cfg *GroupConfig) {
		cfg.Security = append(cfg.Security, OperationSecurityConfig{
			Name:   securityName,
			Scopes: scopes,
		})
	}
}

// GroupHide sets the hide option for the group.
//
// If hide is true, the group will not be included in the OpenAPI specification.
func GroupHide(hide ...bool) GroupOption {
	return func(cfg *GroupConfig) {
		cfg.Hide = util.Optional(true, hide...)
	}
}

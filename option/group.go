package option

import "github.com/oaswrap/spec/pkg/util"

// GroupConfig defines configuration options for a group of routes in an OpenAPI specification.
type GroupConfig struct {
	Tags       []string
	Security   []OperationSecurityConfig
	Deprecated bool
	Hide       bool
}

// GroupOption applies a configuration option to a GroupConfig.
type GroupOption func(*GroupConfig)

// GroupTags sets one or more tags for the group.
//
// These tags will be added to all routes in the sub-router.
func GroupTags(tags ...string) GroupOption {
	return func(cfg *GroupConfig) {
		cfg.Tags = append(cfg.Tags, tags...)
	}
}

// GroupSecurity adds a security scheme to the group.
//
// The security scheme will apply to all routes in the sub-router.
func GroupSecurity(securityName string, scopes ...string) GroupOption {
	return func(cfg *GroupConfig) {
		cfg.Security = append(cfg.Security, OperationSecurityConfig{
			Name:   securityName,
			Scopes: scopes,
		})
	}
}

// GroupHidden sets whether the group should be hidden.
//
// If true, the group and its routes will be excluded from the OpenAPI output.
func GroupHidden(hidden ...bool) GroupOption {
	return func(cfg *GroupConfig) {
		cfg.Hide = util.Optional(true, hidden...)
	}
}

// GroupDeprecated sets whether the group is deprecated.
//
// If true, all routes in the group will be marked as deprecated in the OpenAPI output.
func GroupDeprecated(deprecated ...bool) GroupOption {
	return func(cfg *GroupConfig) {
		cfg.Deprecated = util.Optional(true, deprecated...)
	}
}

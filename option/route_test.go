package option_test

import (
	"testing"

	"github.com/oaswrap/spec/option"
	"github.com/stretchr/testify/assert"
)

func TestRouteTags(t *testing.T) {
	t.Run("adds single tag", func(t *testing.T) {
		cfg := &option.RouteConfig{}
		opt := option.RouteTags("auth")
		opt(cfg)

		assert.Equal(t, []string{"auth"}, cfg.Tags)
	})

	t.Run("adds multiple tags", func(t *testing.T) {
		cfg := &option.RouteConfig{}
		opt := option.RouteTags("auth", "user", "admin")
		opt(cfg)

		assert.Equal(t, []string{"auth", "user", "admin"}, cfg.Tags)
	})

	t.Run("appends to existing tags", func(t *testing.T) {
		cfg := &option.RouteConfig{Tags: []string{"existing"}}
		opt := option.RouteTags("new")
		opt(cfg)

		assert.Equal(t, []string{"existing", "new"}, cfg.Tags)
	})
}

func TestRouteSecurity(t *testing.T) {
	t.Run("adds security without scopes", func(t *testing.T) {
		cfg := &option.RouteConfig{}
		opt := option.RouteSecurity("oauth2")
		opt(cfg)

		expected := []option.RouteSecurityConfig{
			{Name: "oauth2"},
		}
		assert.Equal(t, expected, cfg.Security)
	})

	t.Run("adds security with scopes", func(t *testing.T) {
		cfg := &option.RouteConfig{}
		opt := option.RouteSecurity("oauth2", "read", "write")
		opt(cfg)

		expected := []option.RouteSecurityConfig{
			{Name: "oauth2", Scopes: []string{"read", "write"}},
		}
		assert.Equal(t, expected, cfg.Security)
	})

	t.Run("appends to existing security", func(t *testing.T) {
		cfg := &option.RouteConfig{
			Security: []option.RouteSecurityConfig{
				{Name: "existing", Scopes: []string{"scope1"}},
			},
		}
		opt := option.RouteSecurity("oauth2", "read")
		opt(cfg)

		expected := []option.RouteSecurityConfig{
			{Name: "existing", Scopes: []string{"scope1"}},
			{Name: "oauth2", Scopes: []string{"read"}},
		}
		assert.Equal(t, expected, cfg.Security)
	})
}

func TestRouteHide(t *testing.T) {
	t.Run("hides route by default", func(t *testing.T) {
		cfg := &option.RouteConfig{}
		opt := option.RouteHide()
		opt(cfg)

		assert.True(t, cfg.Hide)
	})

	t.Run("hides route when true", func(t *testing.T) {
		cfg := &option.RouteConfig{}
		opt := option.RouteHide(true)
		opt(cfg)

		assert.True(t, cfg.Hide)
	})

	t.Run("shows route when false", func(t *testing.T) {
		cfg := &option.RouteConfig{}
		opt := option.RouteHide(false)
		opt(cfg)

		assert.False(t, cfg.Hide)
	})

	t.Run("uses first value when multiple provided", func(t *testing.T) {
		cfg := &option.RouteConfig{}
		opt := option.RouteHide(false, true, false)
		opt(cfg)

		assert.False(t, cfg.Hide)
	})
}

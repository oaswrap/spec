package option_test

import (
	"testing"

	"github.com/oaswrap/spec/option"
	"github.com/stretchr/testify/assert"
)

func TestGroupTags(t *testing.T) {
	t.Run("adds single tag", func(t *testing.T) {
		cfg := &option.GroupConfig{}
		opt := option.GroupTags("auth")
		opt(cfg)

		assert.Equal(t, []string{"auth"}, cfg.Tags)
	})

	t.Run("adds multiple tags", func(t *testing.T) {
		cfg := &option.GroupConfig{}
		opt := option.GroupTags("auth", "user", "admin")
		opt(cfg)

		assert.Equal(t, []string{"auth", "user", "admin"}, cfg.Tags)
	})

	t.Run("appends to existing tags", func(t *testing.T) {
		cfg := &option.GroupConfig{Tags: []string{"existing"}}
		opt := option.GroupTags("new")
		opt(cfg)

		assert.Equal(t, []string{"existing", "new"}, cfg.Tags)
	})
}

func TestGroupSecurity(t *testing.T) {
	t.Run("adds security without scopes", func(t *testing.T) {
		cfg := &option.GroupConfig{}
		opt := option.GroupSecurity("oauth2")
		opt(cfg)

		expected := []option.OperationSecurityConfig{
			{Name: "oauth2"},
		}
		assert.Equal(t, expected, cfg.Security)
	})

	t.Run("adds security with scopes", func(t *testing.T) {
		cfg := &option.GroupConfig{}
		opt := option.GroupSecurity("oauth2", "read", "write")
		opt(cfg)

		expected := []option.OperationSecurityConfig{
			{Name: "oauth2", Scopes: []string{"read", "write"}},
		}
		assert.Equal(t, expected, cfg.Security)
	})

	t.Run("appends to existing security", func(t *testing.T) {
		cfg := &option.GroupConfig{
			Security: []option.OperationSecurityConfig{
				{Name: "existing", Scopes: []string{"scope1"}},
			},
		}
		opt := option.GroupSecurity("oauth2", "read")
		opt(cfg)

		expected := []option.OperationSecurityConfig{
			{Name: "existing", Scopes: []string{"scope1"}},
			{Name: "oauth2", Scopes: []string{"read"}},
		}
		assert.Equal(t, expected, cfg.Security)
	})
}

func TestRouteHide(t *testing.T) {
	t.Run("hides route by default", func(t *testing.T) {
		cfg := &option.GroupConfig{}
		opt := option.GroupHide()
		opt(cfg)

		assert.True(t, cfg.Hide)
	})

	t.Run("hides route when true", func(t *testing.T) {
		cfg := &option.GroupConfig{}
		opt := option.GroupHide(true)
		opt(cfg)

		assert.True(t, cfg.Hide)
	})

	t.Run("shows route when false", func(t *testing.T) {
		cfg := &option.GroupConfig{}
		opt := option.GroupHide(false)
		opt(cfg)

		assert.False(t, cfg.Hide)
	})

	t.Run("uses first value when multiple provided", func(t *testing.T) {
		cfg := &option.GroupConfig{}
		opt := option.GroupHide(false, true, false)
		opt(cfg)

		assert.False(t, cfg.Hide)
	})
}

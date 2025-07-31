package spec

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	stdpath "path"
	"strings"
	"sync"

	"github.com/oaswrap/spec/internal/util"
	"github.com/oaswrap/spec/openapi"
	"github.com/oaswrap/spec/option"
)

// generator implements the Generator interface for creating OpenAPI specifications.
type generator struct {
	reflector reflector
	spec      spec
	cfg       *openapi.Config

	prefix string
	groups []*generator
	routes []*route
	opts   []option.GroupOption
	once   sync.Once
}

var _ Generator = (*generator)(nil)

// NewRouter creates a new Router instance with the provided OpenAPI options.
//
// It initializes the reflector and sets up the OpenAPI configuration.
func NewRouter(opts ...option.OpenAPIOption) Generator {
	return NewGenerator(opts...)
}

// NewGenerator creates a new Generator instance with the provided OpenAPI options.
//
// It initializes the reflector and sets up the OpenAPI configuration.
func NewGenerator(opts ...option.OpenAPIOption) Generator {
	cfg := option.WithOpenAPIConfig(opts...)

	reflector := newReflector(cfg)

	return &generator{
		reflector: reflector,
		spec:      reflector.Spec(),
		cfg:       cfg,
	}
}

// Config returns the OpenAPI configuration used by the Router.
func (g *generator) Config() *openapi.Config {
	return g.cfg
}

// Get registers a new GET operation with the specified path and options.
func (g *generator) Get(path string, opts ...option.OperationOption) Route {
	return g.Add("GET", path, opts...)
}

// Post registers a new POST operation with the specified path and options.
func (g *generator) Post(path string, opts ...option.OperationOption) Route {
	return g.Add("POST", path, opts...)
}

// Put registers a new PUT operation with the specified path and options.
func (g *generator) Put(path string, opts ...option.OperationOption) Route {
	return g.Add("PUT", path, opts...)
}

// Delete registers a new DELETE operation with the specified path and options.
func (g *generator) Delete(path string, opts ...option.OperationOption) Route {
	return g.Add("DELETE", path, opts...)
}

// Patch registers a new PATCH operation with the specified path and options.
func (g *generator) Patch(path string, opts ...option.OperationOption) Route {
	return g.Add("PATCH", path, opts...)
}

// Options registers a new OPTIONS operation with the specified path and options.
func (g *generator) Options(path string, opts ...option.OperationOption) Route {
	return g.Add("OPTIONS", path, opts...)
}

// Trace registers a new TRACE operation with the specified path and options.
func (g *generator) Trace(path string, opts ...option.OperationOption) Route {
	return g.Add("TRACE", path, opts...)
}

// Head registers a new HEAD operation with the specified path and options.
func (g *generator) Head(path string, opts ...option.OperationOption) Route {
	return g.Add("HEAD", path, opts...)
}

// Add registers a new operation with the specified method and path.
// It applies the provided operation options to the operation context.
func (g *generator) Add(method, path string, opts ...option.OperationOption) Route {
	if g.prefix != "" {
		path = g.cleanPath(path)
	}
	route := &route{
		method: method,
		path:   path,
		opts:   opts,
	}
	g.routes = append(g.routes, route)

	return route
}

// Route registers a new route with the specified pattern and function.
//
// The function receives a Router instance to define sub-routes.
func (g *generator) Route(pattern string, fn func(router Router), opts ...option.GroupOption) Router {
	subGroup := g.Group(pattern, opts...)
	fn(subGroup)
	return subGroup
}

// Group creates a new sub-router with the specified prefix and options.
func (g *generator) Group(pattern string, opts ...option.GroupOption) Router {
	group := &generator{
		prefix:    g.cleanPath(pattern),
		reflector: g.reflector,
		cfg:       g.cfg,
		opts:      opts,
	}
	g.groups = append(g.groups, group)
	return group
}

// Use applies the provided options to the router.
func (g *generator) Use(opts ...option.GroupOption) Router {
	g.opts = append(g.opts, opts...)
	return g
}

// MarshalYAML marshals the OpenAPI specification to YAML format.
func (g *generator) MarshalYAML() ([]byte, error) {
	if err := g.Validate(); err != nil {
		return nil, err
	}
	return g.spec.MarshalYAML()
}

// MarshalJSON marshals the OpenAPI specification to JSON format with indentation.
func (g *generator) MarshalJSON() ([]byte, error) {
	if err := g.Validate(); err != nil {
		return nil, err
	}
	schema, err := g.spec.MarshalJSON()
	if err != nil {
		return nil, err
	}

	var buffer bytes.Buffer
	if err := json.Indent(&buffer, schema, "", "  "); err != nil {
		return nil, fmt.Errorf("failed to indent OpenAPI JSON schema: %w", err)
	}

	return buffer.Bytes(), nil
}

// GenerateSchema generates the OpenAPI schema in the specified format (JSON or YAML).
//
// By default, it generates YAML. If "json" is specified, it generates JSON.
func (g *generator) GenerateSchema(formats ...string) ([]byte, error) {
	format := util.Optional("yaml", formats...)
	if format != "json" && format != "yaml" && format != "yml" {
		return nil, fmt.Errorf("unsupported format: %s, expected 'json', 'yaml', or 'yml'", format)
	}

	if format == "yaml" || format == "yml" {
		return g.MarshalYAML()
	}

	return g.MarshalJSON()
}

// WriteSchemaTo writes the OpenAPI schema to the specified file path.
//
// The file format is determined by the file extension: ".json" for JSON and ".yaml" for YAML.
func (g *generator) WriteSchemaTo(path string) error {
	format := "yaml"
	if strings.HasSuffix(path, ".json") {
		format = "json"
	} else if !strings.HasSuffix(path, ".yaml") && !strings.HasSuffix(path, ".yml") {
		return fmt.Errorf("unsupported file extension: %s, expected '.json' or '.yaml' or '.yml'", path)
	}
	schema, err := g.GenerateSchema(format)
	if err != nil {
		return err
	}
	return os.WriteFile(path, schema, 0644)
}

// Validate checks if the generated OpenAPI specification is valid.
func (g *generator) Validate() error {
	g.buildOnce()

	return g.reflector.Validate()
}

func (g *generator) buildOnce() {
	g.once.Do(func() {
		for _, r := range g.build() {
			g.reflector.Add(r.method, r.path, r.opts...)
		}
	})
}

func (g *generator) build() []*route {
	var routes []*route
	for _, r := range g.routes {
		var opts []option.OperationOption

		if len(g.opts) > 0 {
			cfg := &option.GroupConfig{}
			for _, opt := range g.opts {
				opt(cfg)
			}
			if cfg.Hide {
				continue
			}
			if len(cfg.Tags) > 0 {
				opts = append(opts, option.Tags(cfg.Tags...))
			}
			if len(cfg.Security) > 0 {
				for _, sec := range cfg.Security {
					opts = append(opts, option.Security(sec.Name, sec.Scopes...))
				}
			}
		}
		if len(r.opts) > 0 {
			r.opts = append(r.opts, opts...)
		}
		routes = append(routes, r)
	}

	for _, group := range g.groups {
		group.opts = append(g.opts, group.opts...)
		routes = append(routes, group.build()...)
	}
	return routes
}

func (g *generator) cleanPath(path string) string {
	cleaned := stdpath.Join(g.prefix, path)
	cleaned = stdpath.Clean(cleaned)
	return cleaned
}

type route struct {
	method string
	path   string
	opts   []option.OperationOption
}

var _ Route = (*route)(nil)

func (r *route) With(opts ...option.OperationOption) Route {
	r.opts = append(r.opts, opts...)
	return r
}

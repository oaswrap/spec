package spec

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/oaswrap/spec/openapi"
	"github.com/oaswrap/spec/option"
)

// Router is the main struct for managing OpenAPI operations and specifications.
type Router struct {
	reflector reflector
	spec      spec
	cfg       *openapi.Config
}

// NewRouter creates a new Router instance with the provided OpenAPI options.
func NewRouter(opts ...option.OpenAPIOption) *Router {
	cfg := option.WithOpenAPIConfig(opts...)

	reflector := newReflector(cfg)

	return &Router{
		reflector: reflector,
		spec:      reflector.Spec(),
		cfg:       cfg,
	}
}

// Config returns the OpenAPI configuration used by the Router.
func (g *Router) Config() *openapi.Config {
	return g.cfg
}

// Get registers a new GET operation with the specified path and options.
func (g *Router) Get(path string, opts ...option.OperationOption) {
	g.Add("GET", path, opts...)
}

// Post registers a new POST operation with the specified path and options.
func (g *Router) Post(path string, opts ...option.OperationOption) {
	g.Add("POST", path, opts...)
}

// Put registers a new PUT operation with the specified path and options.
func (g *Router) Put(path string, opts ...option.OperationOption) {
	g.Add("PUT", path, opts...)
}

// Delete registers a new DELETE operation with the specified path and options.
func (g *Router) Delete(path string, opts ...option.OperationOption) {
	g.Add("DELETE", path, opts...)
}

// Patch registers a new PATCH operation with the specified path and options.
func (g *Router) Patch(path string, opts ...option.OperationOption) {
	g.Add("PATCH", path, opts...)
}

// Options registers a new OPTIONS operation with the specified path and options.
func (g *Router) Options(path string, opts ...option.OperationOption) {
	g.Add("OPTIONS", path, opts...)
}

// Trace registers a new TRACE operation with the specified path and options.
func (g *Router) Trace(path string, opts ...option.OperationOption) {
	g.Add("TRACE", path, opts...)
}

// Head registers a new HEAD operation with the specified path and options.
func (g *Router) Head(path string, opts ...option.OperationOption) {
	g.Add("HEAD", path, opts...)
}

// Add registers a new operation with the specified method and path.
// It applies the provided operation options to the operation context.
func (g *Router) Add(method, path string, opts ...option.OperationOption) {
	g.reflector.Add(method, path, opts...)
}

// GenerateSchema generates the OpenAPI schema in the specified format (JSON or YAML).
//
// By default, it generates YAML. If "json" is specified, it generates JSON.
func (g *Router) GenerateSchema(formats ...string) ([]byte, error) {
	format := "yaml"
	if len(formats) > 0 {
		format = formats[0]
	}

	if format != "json" && format != "yaml" {
		return nil, fmt.Errorf("unsupported format: %s, only 'json' and 'yaml' are supported", format)
	}

	if format == "yaml" {
		schema, err := g.spec.MarshalYAML()
		if err != nil {
			return nil, err
		}
		return schema, nil
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

// WriteSchemaTo writes the OpenAPI schema to the specified file path.
//
// The file format is determined by the file extension: ".json" for JSON and ".yaml" for YAML.
func (g *Router) WriteSchemaTo(path string) error {
	format := "yaml"
	if strings.HasSuffix(path, ".json") {
		format = "json"
	}
	schema, err := g.GenerateSchema(format)
	if err != nil {
		return err
	}
	return os.WriteFile(path, schema, 0644)
}

// Validate checks if the generated OpenAPI specification is valid.
func (g *Router) Validate() error {
	return g.reflector.Validate()
}

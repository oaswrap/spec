package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"

	openapiwrapper "github.com/faizlabs/openapi-wrapper"
	"github.com/swaggest/openapi-go"
)

// Generator is responsible for generating OpenAPI documentation.
type Generator struct {
	reflector Reflector
	spec      Spec
}

// NewGenerator creates a new Generator instance with the provided configuration.
func NewGenerator(cfg *openapiwrapper.Config) (*Generator, error) {
	reflector, err := newReflector(cfg)
	if err != nil {
		return nil, err
	}

	return &Generator{reflector: reflector, spec: reflector.Spec()}, nil
}

// NewOperationContext creates a new operation context for the specified method and path.
func (g *Generator) NewOperationContext(method, path string) (openapi.OperationContext, error) {
	operation, err := g.reflector.NewOperationContext(method, path)
	if err != nil {
		return nil, err
	}

	return operation, nil
}

// AddOperation adds an operation to the OpenAPI documentation.
func (g *Generator) AddOperation(ctx openapi.OperationContext) error {
	if err := g.reflector.AddOperation(ctx); err != nil {
		return err
	}

	return nil
}

// GenerateSchema generates the OpenAPI schema in the specified format (JSON or YAML).
func (g *Generator) GenerateSchema(formats ...string) ([]byte, error) {
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
func (g *Generator) WriteSchemaTo(path string) error {
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

var (
	re3  = regexp.MustCompile(`^3\.0\.\d(-.+)?$`)
	re31 = regexp.MustCompile(`^3\.1\.\d+(-.+)?$`)
)

func newReflector(cfg *openapiwrapper.Config) (Reflector, error) {
	if re3.MatchString(cfg.OpenAPIVersion) {
		return newReflector3(cfg), nil
	} else if re31.MatchString(cfg.OpenAPIVersion) {
		return newReflector31(cfg), nil
	}
	return nil, fmt.Errorf("unsupported OpenAPI version: %s", cfg.OpenAPIVersion)
}

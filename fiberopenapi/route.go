package fiberopenapi

import (
	"github.com/faizlabs/openapi-wrapper/option"
	"github.com/gofiber/fiber/v2"
	"github.com/swaggest/openapi-go"
)

// Route represents a single route in the OpenAPI specification.
type Route interface {
	// Name sets the name for the route.
	Name(name string) Route
	// With applies the given options to the route.
	With(opts ...option.OperationOption) Route
}

type route struct {
	rr        *router
	router    fiber.Router
	operation openapi.OperationContext
	tags      []string
	security  []option.RouteSecurityConfig
	hide      bool
}

// Name sets the name for the route.
func (r *route) Name(name string) Route {
	r.router.Name(name)
	return r
}

// With applies the given options to the route.
func (r *route) With(opts ...option.OperationOption) Route {
	if len(opts) == 0 || r.operation == nil || r.hide {
		return r
	}
	operationCfg := &option.OperationConfig{}
	for _, opt := range opts {
		opt(operationCfg)
	}
	if operationCfg.Hide {
		return r
	}
	r.applyOperation(operationCfg)

	return r
}

func (r *route) applyOperation(cfg *option.OperationConfig) {
	path := r.operation.Method() + " " + r.operation.PathPattern()
	if cfg.OperationID != "" {
		r.rr.logger.Printf("Adding operation with ID: %s for path: %s", cfg.OperationID, path)
		r.operation.SetID(cfg.OperationID)
	}
	if cfg.Description != "" {
		r.rr.logger.Printf("Setting operation description: %s for path: %s", cfg.Description, path)
		r.operation.SetDescription(cfg.Description)
	}
	if cfg.Summary != "" {
		r.rr.logger.Printf("Setting operation summary: %s for path: %s", cfg.Summary, path)
		r.operation.SetSummary(cfg.Summary)
	}
	if cfg.Deprecated {
		r.rr.logger.Printf("Marking operation as deprecated for path: %s", path)
		r.operation.SetIsDeprecated(true)
	}

	r.tags = append(r.tags, cfg.Tags...)
	r.security = append(r.security, cfg.Security...)

	for _, req := range cfg.Requests {
		opts := []openapi.ContentOption{}
		if req.ContentType != "" {
			opts = append(opts, openapi.WithContentType(req.ContentType))
		}
		r.rr.logger.Printf("Adding request structure: %T for path: %s", req.Structure, path)
		r.operation.AddReqStructure(req.Structure, opts...)
	}

	for _, resp := range cfg.Responses {
		opts := []openapi.ContentOption{
			openapi.WithHTTPStatus(resp.HTTPStatus),
		}
		if resp.ContentType != "" {
			opts = append(opts, openapi.WithContentType(resp.ContentType))
		}
		r.rr.logger.Printf("Adding response structure: %T for path: %s", resp.Structure, path)
		r.operation.AddRespStructure(resp.Structure, opts...)
	}
}

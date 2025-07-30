package spec

import (
	"github.com/oaswrap/spec/option"
	"github.com/swaggest/openapi-go"
)

var _ OperationContext = (*operationContext)(nil)

type operationContext struct {
	op     openapi.OperationContext
	cfg    *option.OperationConfig
	logger option.Logger
}

func (oc *operationContext) With(opts ...option.OperationOption) OperationContext {
	for _, opt := range opts {
		opt(oc.cfg)
	}
	return oc
}

func (oc *operationContext) Set(opt option.OperationOption) {
	opt(oc.cfg)
}

func (oc *operationContext) build() openapi.OperationContext {
	path := oc.op.Method() + " " + oc.op.PathPattern()

	cfg := oc.cfg
	if cfg == nil {
		return nil
	}
	if cfg.Hide {
		oc.logger.Printf("Skipping operation %s: hidden", path)
		return nil
	}
	if cfg.Deprecated {
		oc.logger.Printf("Marking operation %s as deprecated", path)
		oc.op.SetIsDeprecated(true)
	}
	if cfg.OperationID != "" {
		oc.logger.Printf("Setting operation ID for %s: %s", path, cfg.OperationID)
		oc.op.SetID(cfg.OperationID)
	}
	if cfg.Summary != "" {
		oc.logger.Printf("Setting summary for operation %s: %s", path, cfg.Summary)
		oc.op.SetSummary(cfg.Summary)
	}
	if cfg.Description != "" {
		oc.logger.Printf("Setting description for operation %s: %s", path, cfg.Description)
		oc.op.SetDescription(cfg.Description)
	}
	if len(cfg.Tags) > 0 {
		oc.logger.Printf("Setting tags for operation %s: %v", path, cfg.Tags)
		oc.op.SetTags(cfg.Tags...)
	}
	if len(cfg.Security) > 0 {
		for _, sec := range cfg.Security {
			oc.logger.Printf("Adding security scheme %s to operation %s", sec.Name, path)
			oc.op.AddSecurity(sec.Name, sec.Scopes...)
		}
	}

	for _, req := range cfg.Requests {
		opts := []openapi.ContentOption{}
		if req.ContentType != "" {
			opts = append(opts, openapi.WithContentType(req.ContentType))
		}
		oc.logger.Printf("Adding request structure for operation %s with content type %s", path, req.ContentType)
		oc.op.AddReqStructure(req.Structure, opts...)
	}

	for _, resp := range cfg.Responses {
		opts := []openapi.ContentOption{
			openapi.WithHTTPStatus(resp.HTTPStatus),
		}
		if resp.ContentType != "" {
			opts = append(opts, openapi.WithContentType(resp.ContentType))
		}
		oc.logger.Printf("Adding response structure for operation %s with HTTP status %d and content type %s", path, resp.HTTPStatus, resp.ContentType)
		oc.op.AddRespStructure(resp.Structure, opts...)
	}

	return oc.op
}

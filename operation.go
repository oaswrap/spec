package spec

import (
	"fmt"
	"strings"

	"github.com/oaswrap/spec/option"
	"github.com/oaswrap/spec/pkg/debuglog"
	"github.com/swaggest/openapi-go"
)

var _ operationContext = (*operationContextImpl)(nil)

type operationContextImpl struct {
	op     openapi.OperationContext
	cfg    *option.OperationConfig
	logger *debuglog.Logger
}

func (oc *operationContextImpl) With(opts ...option.OperationOption) operationContext {
	for _, opt := range opts {
		opt(oc.cfg)
	}
	return oc
}

func (oc *operationContextImpl) build() openapi.OperationContext {
	method := strings.ToUpper(oc.op.Method())
	path := oc.op.PathPattern()

	logger := oc.logger

	cfg := oc.cfg
	if cfg == nil {
		return nil
	}
	if cfg.Hide {
		logger.LogAction("skip operation", fmt.Sprintf("%s %s", method, path))
		return nil
	}
	if cfg.Deprecated {
		oc.op.SetIsDeprecated(true)
		logger.LogOp(method, path, "set is deprecated", "true")
	}
	if cfg.OperationID != "" {
		oc.op.SetID(cfg.OperationID)
		logger.LogOp(method, path, "set operation ID", cfg.OperationID)
	}
	if cfg.Summary != "" {
		oc.op.SetSummary(cfg.Summary)
		logger.LogOp(method, path, "set summary", cfg.Summary)
	}
	if cfg.Description != "" {
		oc.op.SetDescription(cfg.Description)
		logger.LogOp(method, path, "set description", cfg.Description)
	}
	if len(cfg.Tags) > 0 {
		oc.op.SetTags(cfg.Tags...)
		logger.LogOp(method, path, "set tags", fmt.Sprintf("%v", cfg.Tags))
	}
	if len(cfg.Security) > 0 {
		for _, sec := range cfg.Security {
			oc.op.AddSecurity(sec.Name, sec.Scopes...)
		}
		logger.LogOp(method, path, "set security", fmt.Sprintf("%v", cfg.Security))
	}

	for _, req := range cfg.Requests {
		value := fmt.Sprintf("%T", req.Structure)
		opts := []openapi.ContentOption{}
		if req.ContentType != "" {
			value += fmt.Sprintf(" (Content-Type: %s)", req.ContentType)
			opts = append(opts, openapi.WithContentType(req.ContentType))
		}
		oc.op.AddReqStructure(req.Structure, opts...)
		logger.LogOp(method, path, "add request", value)
	}

	for _, resp := range cfg.Responses {
		value := fmt.Sprintf("%T (HTTP %d)", resp.Structure, resp.HTTPStatus)
		opts := []openapi.ContentOption{
			openapi.WithHTTPStatus(resp.HTTPStatus),
		}
		if resp.ContentType != "" {
			value += fmt.Sprintf(" (Content-Type: %s)", resp.ContentType)
			opts = append(opts, openapi.WithContentType(resp.ContentType))
		}
		oc.op.AddRespStructure(resp.Structure, opts...)
		logger.LogOp(method, path, "add response", value)
	}

	return oc.op
}

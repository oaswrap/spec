package spec

import (
	"fmt"
	"strings"

	"github.com/oaswrap/spec/internal/debuglog"
	specopenapi "github.com/oaswrap/spec/openapi"
	"github.com/oaswrap/spec/option"
	"github.com/swaggest/openapi-go"
	"github.com/swaggest/openapi-go/openapi3"
	"github.com/swaggest/openapi-go/openapi31"
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
		opts, value := oc.buildRequestOpts(req)
		oc.op.AddReqStructure(req.Structure, opts...)
		logger.LogOp(method, path, "add request", value)
	}

	for _, resp := range cfg.Responses {
		opts, value := oc.buildResponseOpts(resp)
		oc.op.AddRespStructure(resp.Structure, opts...)
		logger.LogOp(method, path, "add response", value)
	}

	return oc.op
}

func stringMapToEncodingMap3(enc map[string]string) map[string]openapi3.Encoding {
	res := map[string]openapi3.Encoding{}
	for k, v := range enc {
		rv := v
		res[k] = openapi3.Encoding{
			ContentType: &rv,
		}
	}
	return res
}

func stringMapToEncodingMap31(enc map[string]string) map[string]openapi31.Encoding {
	res := map[string]openapi31.Encoding{}
	for k, v := range enc {
		rv := v
		res[k] = openapi31.Encoding{
			ContentType: &rv,
		}
	}
	return res
}

func (oc *operationContextImpl) buildRequestOpts(req *specopenapi.ContentUnit) ([]openapi.ContentOption, string) {
	log := fmt.Sprintf("%T", req.Structure)
	var opts []openapi.ContentOption
	if req.Description != "" {
		opts = append(opts, func(cu *openapi.ContentUnit) {
			cu.Description = req.Description
		})
		log += fmt.Sprintf(" (%s)", req.Description)
	}
	if req.ContentType != "" {
		opts = append(opts, openapi.WithContentType(req.ContentType))
		log += fmt.Sprintf(" (Content-Type: %s)", req.ContentType)
	}
	opts = append(opts, func(cu *openapi.ContentUnit) {
		cu.Customize = func(cor openapi.ContentOrReference) {
			switch v := cor.(type) {
			case *openapi3.RequestBodyOrRef:
				content := map[string]openapi3.MediaType{}
				for k, val := range v.RequestBody.Content {
					content[k] = *val.WithEncoding(stringMapToEncodingMap3(req.Encoding))
				}
				v.RequestBody.WithContent(content)
			case *openapi31.RequestBodyOrReference:
				content := map[string]openapi31.MediaType{}
				for k, val := range v.RequestBody.Content {
					content[k] = *val.WithEncoding(stringMapToEncodingMap31(req.Encoding))
				}
				v.RequestBody.WithContent(content)
			}
		}
	})
	return opts, log
}

func (oc *operationContextImpl) buildResponseOpts(resp *specopenapi.ContentUnit) ([]openapi.ContentOption, string) {
	log := fmt.Sprintf("%T", resp.Structure)
	var opts []openapi.ContentOption
	if resp.IsDefault {
		opts = append(opts, func(cu *openapi.ContentUnit) {
			cu.IsDefault = true
		})
		log += " (default)"
	}
	if resp.HTTPStatus != 0 {
		opts = append(opts, openapi.WithHTTPStatus(resp.HTTPStatus))
		log += fmt.Sprintf(" (HTTP %d)", resp.HTTPStatus)
	}
	if resp.Description != "" {
		opts = append(opts, func(cu *openapi.ContentUnit) {
			cu.Description = resp.Description
		})
		log += fmt.Sprintf(" (%s)", resp.Description)
	}
	if resp.ContentType != "" {
		opts = append(opts, openapi.WithContentType(resp.ContentType))
		log += fmt.Sprintf(" (Content-Type: %s)", resp.ContentType)
	}
	return opts, log
}

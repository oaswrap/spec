package spec

import (
	"fmt"
	"regexp"

	"github.com/oaswrap/spec/internal/debuglog"
	"github.com/oaswrap/spec/internal/errs"
	"github.com/oaswrap/spec/openapi"
	"github.com/oaswrap/spec/option"
)

var (
	re3  = regexp.MustCompile(`^3\.0\.\d(-.+)?$`)
	re31 = regexp.MustCompile(`^3\.1\.\d+(-.+)?$`)
)

func newReflector(cfg *openapi.Config) reflector {
	logger := debuglog.NewLogger("spec", cfg.Logger)

	if re3.MatchString(cfg.OpenAPIVersion) {
		return newReflector3(cfg, logger)
	} else if re31.MatchString(cfg.OpenAPIVersion) {
		return newReflector31(cfg, logger)
	}

	logger.Printf("Unsupported OpenAPI version: %s", cfg.OpenAPIVersion)
	return newInvalidReflector(fmt.Errorf("unsupported OpenAPI version: %s", cfg.OpenAPIVersion))
}

type invalidReflector struct {
	spec   *noopSpec
	errors *errs.SpecError
}

func newInvalidReflector(err error) reflector {
	e := &errs.SpecError{}
	e.Add(err)

	return &invalidReflector{
		errors: e,
		spec:   &noopSpec{},
	}
}

var _ reflector = (*invalidReflector)(nil)

func (r *invalidReflector) Spec() spec {
	return r.spec
}

func (r *invalidReflector) Add(_, _ string, _ ...option.OperationOption) {}

func (r *invalidReflector) Validate() error {
	if r.errors.HasErrors() {
		return r.errors
	}
	return nil
}

type noopSpec struct{}

func (s *noopSpec) MarshalYAML() ([]byte, error) {
	return nil, nil
}

func (s *noopSpec) MarshalJSON() ([]byte, error) {
	return nil, nil
}

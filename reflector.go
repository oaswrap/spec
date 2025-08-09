package spec

import (
	"fmt"
	"regexp"

	"github.com/oaswrap/spec/internal/debuglog"
	"github.com/oaswrap/spec/internal/errors"
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
	errors *errors.SpecError
}

func newInvalidReflector(err error) reflector {
	errors := &errors.SpecError{}
	errors.Add(err)

	return &invalidReflector{
		errors: errors,
		spec:   &noopSpec{},
	}
}

var _ reflector = (*invalidReflector)(nil)

func (r *invalidReflector) Spec() spec {
	return r.spec
}

func (r *invalidReflector) Add(method, path string, opts ...option.OperationOption) {}

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

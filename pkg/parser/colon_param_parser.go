package parser

import (
	"regexp"

	"github.com/oaswrap/spec/openapi"
)

// ColonParamParser is a parser that converts paths with colon-prefixed parameters
// (e.g., "/users/:id") to OpenAPI-style parameters (e.g., "/users/{id}").
type ColonParamParser struct {
	re *regexp.Regexp
}

var _ openapi.PathParser = &ColonParamParser{}

// NewColonParamParser creates a new ColonParamParser instance.
func NewColonParamParser() *ColonParamParser {
	return &ColonParamParser{
		re: regexp.MustCompile(`:([a-zA-Z_][a-zA-Z0-9_]*)`),
	}
}

// Parse converts a path with colon-prefixed parameters to OpenAPI-style parameters.
func (p *ColonParamParser) Parse(colonParam string) (string, error) {
	return p.re.ReplaceAllString(colonParam, "{$1}"), nil
}

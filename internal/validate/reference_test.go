package validate_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/oaswrap/spec"
	"github.com/oaswrap/spec/internal/validate"
	"github.com/oaswrap/spec/openapi"
	"github.com/oaswrap/spec/option"
)

func TestValidateOpenAPI312AllowsReferenceDescription(t *testing.T) {
	description := "Security scheme reference"
	r := spec.NewRouter(
		option.WithOpenAPIVersion(openapi.Version312),
		option.WithComponentParameter("BaseTraceID", &openapi.Parameter{
			Name:     "X-Trace-ID",
			In:       "header",
			Required: false,
			Schema:   &openapi.Schema{Type: "string"},
		}),
		option.WithComponentParameter("TraceID", &openapi.Parameter{
			Ref:         "#/components/parameters/BaseTraceID",
			Summary:     "Trace",
			Description: "Trace identifier",
		}),
		option.WithComponentSecurityScheme("BaseAuth", &openapi.SecurityScheme{
			Type:   "http",
			Scheme: "bearer",
		}),
		option.WithComponentSecurityScheme("Auth", &openapi.SecurityScheme{
			Ref:         "#/components/securitySchemes/BaseAuth",
			Summary:     "Auth",
			Description: &description,
		}),
	)
	r.Get("/refs", option.Response(200, nil))

	assert.NoError(t, r.Validate())
}

func TestValidateOpenAPI320AllowsMediaTypeReferenceSummaryDescription(t *testing.T) {
	r := spec.NewRouter(
		option.WithOpenAPIVersion(openapi.Version320),
		option.WithComponentMediaType("JSON", &openapi.MediaType{
			Schema: &openapi.Schema{Type: "object"},
		}),
		option.WithComponentMediaType("JSONRef", &openapi.MediaType{
			Ref:         "#/components/mediaTypes/JSON",
			Summary:     "JSON",
			Description: "JSON media type",
		}),
	)
	r.Get("/refs", option.CustomizeOperation(func(op *openapi.Operation) {
		op.Responses["200"] = &openapi.Response{
			Description: "OK",
			Content: map[string]openapi.MediaType{
				"application/json": {
					Ref:         "#/components/mediaTypes/JSON",
					Summary:     "JSON",
					Description: "JSON media type",
				},
			},
		}
	}))

	assert.NoError(t, r.Validate())
}

func TestValidateRejectsReferenceSummaryOnNonReferenceObject(t *testing.T) {
	r := spec.NewRouter(
		option.WithDocument(func(doc *openapi.Document) {
			doc.Paths["/summary"].Get.Parameters = append(doc.Paths["/summary"].Get.Parameters, &openapi.Parameter{
				Summary: "not a reference",
				Name:    "q",
				In:      "query",
				Schema:  &openapi.Schema{Type: "string"},
			})
		}),
	)
	r.Get("/summary", option.Response(204, nil))

	err := r.Validate()
	assertValidationContains(t, err, "parameters[0].summary is only allowed with $ref")
}

func TestValidateResponseAndReferenceRules(t *testing.T) {
	r := spec.NewRouter(
		option.WithOpenAPIVersion(openapi.Version304),
		option.WithTitle("Invalid Responses"),
		option.WithVersion("1.0.0"),
		option.WithDocument(func(doc *openapi.Document) {
			doc.Components.Responses = map[string]*openapi.Response{
				"BadRef": {Ref: "#/components/responses/OK", Description: "sibling"},
			}
		}),
	)
	r.Get("/bad",
		option.CustomizeOperation(func(op *openapi.Operation) {
			op.Responses["700"] = &openapi.Response{}
		}),
		option.Response(204, nil),
	)

	err := r.Validate()
	assertValidationContains(t, err,
		"responses.700 must be default, a status code, or a status code range",
		"responses.700.description is required",
		"components.responses.BadRef must not define siblings with $ref",
	)
}

func TestValidateRejectsLocalRefWithoutHashPrefix(t *testing.T) {
	r := spec.NewRouter(
		option.WithTitle("Invalid Ref"),
		option.WithVersion("1.0.0"),
		option.WithDocument(func(doc *openapi.Document) {
			doc.Components.Schemas = map[string]*openapi.Schema{
				"User": {Type: "object"},
			}
			doc.Paths["/users"].Get.Responses["200"] = &openapi.Response{
				Ref: "/components/schemas/User",
			}
		}),
	)
	r.Get("/users", option.Response(200, nil))

	err := r.Validate()
	assertValidationContains(t, err, `$ref "/components/schemas/User" must use #/ for local references`)
}

func TestValidateRejectsMissingLocalRefTarget(t *testing.T) {
	r := spec.NewRouter(
		option.WithTitle("Missing Ref Target"),
		option.WithVersion("1.0.0"),
		option.WithDocument(func(doc *openapi.Document) {
			doc.Paths["/users"].Get.Responses["200"] = &openapi.Response{
				Ref: "#/components/schemas/User",
			}
		}),
	)
	r.Get("/users", option.Response(200, nil))

	err := r.Validate()
	assertValidationContains(t, err, `$ref "#/components/schemas/User" points to a missing target`)
}

func TestReferenceSiblingHelpers(t *testing.T) {
	t.Run("BodyRef", func(t *testing.T) {
		assert.True(
			t,
			validate.BodyRefHasInvalidSiblings(&openapi.RequestBody{Description: "legacy"}, openapi.Version304),
		)
		assert.False(t, validate.BodyRefHasInvalidSiblings(&openapi.RequestBody{Description: "ok"}, openapi.Version312))
	})

	t.Run("LinkRef", func(t *testing.T) {
		assert.True(t, validate.LinkRefHasInvalidSiblings(&openapi.Link{Description: "legacy"}, openapi.Version304))
		assert.False(t, validate.LinkRefHasInvalidSiblings(&openapi.Link{Description: "ok"}, openapi.Version312))
		assert.True(t, validate.LinkRefHasInvalidSiblings(&openapi.Link{OperationID: "op"}, openapi.Version312))
	})

	t.Run("CallbackRef", func(t *testing.T) {
		assert.True(t, validate.CallbackRefHasInvalidSiblings(&openapi.Callback{Summary: "legacy"}, openapi.Version304))
		assert.False(t, validate.CallbackRefHasInvalidSiblings(&openapi.Callback{Summary: "ok"}, openapi.Version312))
		assert.True(
			t,
			validate.CallbackRefHasInvalidSiblings(
				&openapi.Callback{Expressions: map[string]*openapi.PathItem{"/": {}}},
				openapi.Version312,
			),
		)
	})

	t.Run("MediaTypeRef", func(t *testing.T) {
		assert.True(
			t,
			validate.MediaTypeRefHasInvalidSiblings(&openapi.MediaType{Summary: "legacy"}, openapi.Version304),
		)
		assert.False(t, validate.MediaTypeRefHasInvalidSiblings(&openapi.MediaType{Summary: "ok"}, openapi.Version312))
	})

	t.Run("SecuritySchemeRef", func(t *testing.T) {
		assert.True(
			t,
			validate.SecuritySchemeRefHasInvalidSiblings(
				&openapi.SecurityScheme{Summary: "legacy"},
				openapi.Version304,
			),
		)
		assert.False(
			t,
			validate.SecuritySchemeRefHasInvalidSiblings(&openapi.SecurityScheme{Summary: "ok"}, openapi.Version312),
		)
	})

	assert.True(t, validate.HeaderRefHasInvalidSiblings(&openapi.Header{Description: "legacy"}, openapi.Version304))
	assert.False(t, validate.HeaderRefHasInvalidSiblings(&openapi.Header{Description: "ok"}, openapi.Version312))
	assert.True(
		t,
		validate.HeaderRefHasInvalidSiblings(
			&openapi.Header{Extra: map[string]any{"note": "invalid"}},
			openapi.Version312,
		),
	)

	assert.True(t, validate.ExampleRefHasInvalidSiblings(&openapi.Example{Summary: "legacy"}, openapi.Version304))
	assert.False(t, validate.ExampleRefHasInvalidSiblings(&openapi.Example{Summary: "ok"}, openapi.Version312))
	assert.True(
		t,
		validate.ExampleRefHasInvalidSiblings(&openapi.Example{SerializedValue: "hello%20world"}, openapi.Version312),
	)
	assert.True(
		t,
		//nolint:staticcheck // SA1019: deliberately exercising the deprecated field.
		validate.ExampleRefHasInvalidSiblings(&openapi.Example{SerializedExample: "legacy"}, openapi.Version320),
	)
	assert.True(
		t,
		validate.ExampleRefHasInvalidSiblings(
			&openapi.Example{Extra: map[string]any{"note": "invalid"}},
			openapi.Version312,
		),
	)
	assert.False(
		t,
		validate.ExampleRefHasInvalidSiblings(
			&openapi.Example{Extra: map[string]any{"summary": "ok"}},
			openapi.Version312,
		),
	)
}

func TestExtraHelpers(t *testing.T) {
	assert.True(t, validate.HasInvalidReferenceExtra(map[string]any{"summary": "legacy"}, openapi.Version304))
	assert.False(t, validate.HasInvalidReferenceExtra(map[string]any{
		"summary":     "ok",
		"description": "ok",
	}, openapi.Version312))
	assert.True(t, validate.HasInvalidReferenceExtra(map[string]any{"invalid": true}, openapi.Version312))
	assert.False(t, validate.HasInvalidReferenceExtra(map[string]any{"x-meta": "ok"}, openapi.Version312))

	assert.True(t, validate.HasNonExtensionExtra(map[string]any{"invalid": true}))
	assert.False(t, validate.HasNonExtensionExtra(map[string]any{"x-valid": true}))
}

func TestValidateReferenceTargets_InvalidURIReference(t *testing.T) {
	doc := &openapi.Document{
		OpenAPI: openapi.Version304,
		Info:    openapi.Info{Title: "T", Version: "1.0.0"},
		Paths: map[string]*openapi.PathItem{
			"/x": {
				Get: &openapi.Operation{
					Responses: map[string]*openapi.Response{
						"200": {Ref: "%zz"},
					},
				},
			},
		},
	}

	errs := validate.ValidateReferenceTargets(doc)
	assert.NotEmpty(t, errs)
	assert.Contains(t, errs[0].Error(), "must be a URI reference")
}

# oaswrap/spec

[![CI](https://github.com/oaswrap/spec/actions/workflows/ci.yml/badge.svg)](https://github.com/oaswrap/spec/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/oaswrap/spec/graph/badge.svg?token=RIEIM9BAIW)](https://codecov.io/gh/oaswrap/spec)
[![Go Reference](https://pkg.go.dev/badge/github.com/oaswrap/spec.svg)](https://pkg.go.dev/github.com/oaswrap/spec)
[![Go Report Card](https://goreportcard.com/badge/github.com/oaswrap/spec)](https://goreportcard.com/report/github.com/oaswrap/spec)
[![Go Version](https://img.shields.io/github/go-mod/go-version/oaswrap/spec)](https://github.com/oaswrap/spec/blob/main/go.mod)
[![License](https://img.shields.io/github/license/oaswrap/spec)](LICENSE)

Code-first, framework-agnostic OpenAPI 3.x spec builder for Go. Generate docs from route registrations and Go structs — no annotations, no vendor lock-in.

---

## Why oaswrap/spec?

- **Native OpenAPI builder** — paths, operations, components, validation, and schema reflection are all implemented in this repository without third-party OpenAPI dependencies.
- **Framework-agnostic core** — use `spec.NewRouter` for static generation, or drop in adapters for Chi, Echo, Gin, Fiber, net/http, and Mux.
- **Code-first route documentation** — register routes and their documentation together using Go functions and typed options.
- **Version-aware output** — defaults to OpenAPI `3.1.2`, with full support for `3.0.x` and `3.2.0` features when selected.
- **Direct model escape hatches** — use typed `openapi` structs, `Extensions` for `x-*` fields, and `Extra` for official or future fields not yet wrapped by a helper option.
- **Deterministic output** — generated documents are stable enough for golden-file snapshot tests and CI documentation checks.

---

## Features

- JSON and YAML generation.
- Framework-agnostic route registration with `NewRouter`, groups, route builders, and HTTP method helpers.
- Webhook registration helpers for OpenAPI `3.1.x` and `3.2.0`.
- OpenAPI `3.2.0` `QUERY` method and additional operation support.
- Request and response schema reflection from Go structs.
- Parameter reflection from `path`, `query`, `header`, `cookie`, and OpenAPI `3.2.0` `querystring` tags.
- JSON Schema/OpenAPI schema tags for examples, enums, constraints, nullability, read/write flags, content metadata, and XML metadata.
- Security helpers for API key, HTTP bearer/basic-style schemes, OAuth2, OpenID Connect, and mutual TLS.
- Duplicate response merging into `oneOf`.
- Low-level typed OpenAPI model for direct field control.
- `spec.OneOf` for explicit one-of schemas.
- `SchemaExposer` and `StaticSchemaExposer` hooks for custom reflected schemas.

---

## Installation

Requirements:

- Go `1.22+`

```bash
go get -u github.com/oaswrap/spec
```

---

## Quick Start

```go
package main

import (
	"log"

	"github.com/oaswrap/spec"
	"github.com/oaswrap/spec/openapi"
	"github.com/oaswrap/spec/option"
)

func main() {
	r := spec.NewRouter(
		option.WithOpenAPIVersion(openapi.Version320),
		option.WithTitle("Users API"),
		option.WithVersion("1.0.0"),
		option.WithServer("https://api.example.com"),
		option.WithSecurity("bearerAuth", option.SecurityHTTPBearer("bearer")),
	)

	v1 := r.Group("/api/v1", option.GroupTags("users"))
	v1.Post("/login",
		option.Summary("Login"),
		option.Request(new(LoginRequest)),
		option.Response(200, new(LoginResponse)),
	)
	v1.Get("/users/{id}",
		option.Summary("Get user"),
		option.Request(new(GetUserRequest)),
		option.Response(200, new(User)),
	)

	if err := r.WriteSchemaTo("openapi.yaml"); err != nil {
		log.Fatal(err)
	}
}

type LoginRequest struct {
	Username string `json:"username" required:"true"`
	Password string `json:"password" required:"true" writeOnly:"true"`
}

type LoginResponse struct {
	Token string `json:"token" required:"true"`
}

type GetUserRequest struct {
	ID string `path:"id" required:"true" description:"User identifier"`
}

type User struct {
	ID   string `json:"id" required:"true"`
	Name string `json:"name"`
}
```

---

## Output

`spec.NewRouter` and `spec.NewGenerator` expose the following output methods:

| Method | Description |
| --- | --- |
| `GenerateSchema()` | Defaults to YAML. |
| `GenerateSchema("yaml")` / `GenerateSchema("yml")` | Returns YAML. |
| `GenerateSchema("json")` | Returns pretty-printed JSON. |
| `MarshalYAML()` | Validates and serializes YAML. |
| `MarshalJSON()` | Validates and serializes pretty-printed JSON. |
| `WriteSchemaTo("openapi.yaml")` | Infers format from file extension (`.yaml`, `.yml`, `.json`). |
| `Document()` | Returns the built `*openapi.Document`. |
| `Validate()` | Builds the document and checks OpenAPI invariants. Returns only `SeverityError` findings. |
| `ValidateReport()` | Builds and validates, returning all findings including warnings and info as `ValidationErrors`. |
| `Config()` | Returns the effective OpenAPI configuration. |

---

## Framework Adapters

Use the core `spec` package for static generation and CI workflows. Use an adapter when you want route registration, spec generation, and docs UI wiring in a single framework-specific router.

| Framework | Package |
| --- | --- |
| Chi | [`chiopenapi`](/adapter/chiopenapi) |
| Echo v4 | [`echoopenapi`](/adapter/echoopenapi) |
| Echo v5 | [`echov5openapi`](/adapter/echov5openapi) |
| Fiber v2 | [`fiberopenapi`](/adapter/fiberopenapi) |
| Fiber v3 | [`fiberv3openapi`](/adapter/fiberv3openapi) |
| Gin | [`ginopenapi`](/adapter/ginopenapi) |
| Iris | [`irisopenapi`](/adapter/irisopenapi) |
| net/http | [`httpopenapi`](/adapter/httpopenapi) |
| Mux | [`muxopenapi`](/adapter/muxopenapi) |

Adapter module import paths:

- `github.com/oaswrap/spec/adapter/chiopenapi`
- `github.com/oaswrap/spec/adapter/echoopenapi`
- `github.com/oaswrap/spec/adapter/echov5openapi`
- `github.com/oaswrap/spec/adapter/fiberopenapi`
- `github.com/oaswrap/spec/adapter/fiberv3openapi`
- `github.com/oaswrap/spec/adapter/ginopenapi`
- `github.com/oaswrap/spec/adapter/irisopenapi`
- `github.com/oaswrap/spec/adapter/httpopenapi`
- `github.com/oaswrap/spec/adapter/muxopenapi`

---

## OpenAPI Options

Root options are passed to `spec.NewRouter` or `spec.NewGenerator`.

```go
r := spec.NewRouter(
	option.WithOpenAPIVersion(openapi.Version312),
	option.WithTitle("Payments API"),
	option.WithInfoSummary("Payments API"),
	option.WithVersion("1.4.0"),
	option.WithDescription("Payment operations."),
	option.WithTermsOfService("https://example.com/terms"),
	option.WithContact(openapi.Contact{Name: "API Team", Email: "api@example.com"}),
	option.WithLicense(openapi.License{Name: "Apache 2.0", URL: "https://www.apache.org/licenses/LICENSE-2.0.html"}),
	option.WithExternalDocs("https://docs.example.com", "Full documentation"),
	option.WithServer("https://api.example.com", option.ServerDescription("Production")),
	option.WithTag("payments", option.TagDescription("Payment operations")),
	option.WithStripTrailingSlash(),
)
```

| Option | Purpose |
| --- | --- |
| `WithOpenAPIConfig(opts...)` | Build an `*openapi.Config` with defaults and apply options. |
| `WithOpenAPIVersion(version)` | Set `openapi`; default is `openapi.Version312`. Constants are available for `3.0.0`–`3.0.4`, `3.1.0`–`3.1.2`, and `3.2.0`. |
| `WithSelf(uri)` | Set OpenAPI `3.2.0` `$self`. |
| `WithJSONSchemaDialect(uri)` | Set root `jsonSchemaDialect`. |
| `WithTitle(title)` | Set `info.title`. |
| `WithInfoSummary(summary)` | Set `info.summary`. |
| `WithVersion(version)` | Set `info.version`. |
| `WithDescription(description)` | Set `info.description`. |
| `WithContact(contact)` | Set `info.contact`. |
| `WithLicense(license)` | Set `info.license`. |
| `WithTermsOfService(url)` | Set `info.termsOfService`. |
| `WithTags(tags...)` | Add root `tags` from `openapi.Tag` values. |
| `WithTag(name, tagOpts...)` | Add one root tag without constructing `openapi.Tag` manually. |
| `WithServer(url, opts...)` | Add a root server. |
| `WithExternalDocs(url, description...)` | Set root `externalDocs`. |
| `WithSecurity(name, opts...)` | Add a reusable security scheme. |
| `WithGlobalSecurity(name, scopes...)` | Add a root security requirement. |
| `WithReflectorConfig(opts...)` | Configure schema reflection. |
| `WithStripTrailingSlash(strip...)` | Trim trailing slashes from operation paths except `/`. |
| `WithPathParser(parser)` | Convert framework-style paths, e.g. `:id` → `{id}`. |
| `WithDocument(fn)` | Mutate the final low-level document before validation and serialization. |
| `WithComponentSchema(name, schema)` | Register a reusable schema component. |
| `WithComponentResponse(name, response)` | Register a reusable response component. |
| `WithComponentParameter(name, parameter)` | Register a reusable parameter component. |
| `WithComponentExample(name, example)` | Register a reusable example component. |
| `WithComponentRequestBody(name, requestBody)` | Register a reusable request body component. |
| `WithComponentHeader(name, header)` | Register a reusable header component. |
| `WithComponentSecurityScheme(name, scheme)` | Register a reusable security scheme component. |
| `WithComponentLink(name, link)` | Register a reusable link component. |
| `WithComponentCallback(name, callback)` | Register a reusable callback component. |
| `WithComponentPathItem(name, pathItem)` | Register a reusable path item component. |
| `WithComponentMediaType(name, mediaType)` | Register a reusable media type component (OpenAPI `3.2.0`). |
| `WithDisableDocs(disable...)` | Disable adapter docs endpoints. |
| `WithDocsPath(path)` | Set adapter docs UI path. |
| `WithSpecPath(path)` | Set adapter OpenAPI spec path. |
| `WithCacheAge(cacheAge)` | Set docs/spec cache age. |
| `WithUIOption(opt)` | Set a low-level `spec-ui` option. |
| `WithSwaggerUI(cfg...)` | Use Swagger UI. |
| `WithStoplightElements(cfg...)` | Use Stoplight Elements. |
| `WithReDoc(cfg...)` | Use ReDoc. |
| `WithScalar(cfg...)` | Use Scalar. |
| `WithRapiDoc(cfg...)` | Use RapiDoc. |

> Adapter-only options: `WithDisableDocs`, `WithDocsPath`, `WithSpecPath`, `WithCacheAge`, `WithUIOption`, `WithSwaggerUI`, `WithStoplightElements`, `WithReDoc`, `WithScalar`, and `WithRapiDoc` affect adapter docs/spec endpoints, not core static generation output.

**Tag options:** `TagSummary`, `TagDescription`, `TagExternalDocs`, `TagParent` (3.2.0), `TagKind` (3.2.0).

**Server options:** `ServerDescription`, `ServerVariables`, `ServerName` (3.2.0).

---

## Routes and Groups

```go
api := r.Group("/api/v1", option.GroupTags("v1"))

api.Get("/users/{id}",
	option.OperationID("getUser"),
	option.Summary("Get user"),
	option.Description("Returns one user."),
	option.Tags("users"),
	option.Security("bearerAuth"),
	option.Request(new(GetUserRequest)),
	option.Response(200, new(User), option.ContentDescription("OK")),
	option.Response(404, nil, option.ContentDescription("Not Found")),
)
```

**Router methods:**

`Get`, `Post`, `Put`, `Delete`, `Patch`, `Options`, `Head`, `Trace`, `Query` (3.2.0), `Add`, `Webhook` (3.1.x+), `AddWebhook` (3.1.x+), `NewRoute(...).Method(...).Path(...).With(...)`, `Group`, `Route`, `With`.

**Group options:**

| Option | Purpose |
| --- | --- |
| `GroupTags(tags...)` | Apply tags to all routes in the group. |
| `GroupSecurity(name, scopes...)` | Apply an operation security requirement to all routes in the group. |
| `GroupDeprecated(deprecated...)` | Mark all routes in the group as deprecated. |
| `GroupHidden(hide...)` | Exclude all routes in the group from the generated document. |

**Operation options:**

| Option | Purpose |
| --- | --- |
| `OperationID(id)` | Set `operationId`. |
| `Summary(summary)` | Set `summary`; also sets `description` when description is empty. |
| `Description(description)` | Set `description`. |
| `ExternalDocs(url, description...)` | Set operation `externalDocs`. |
| `Tags(tags...)` | Add operation tags. |
| `Security(name, scopes...)` | Add an operation security requirement. |
| `Deprecated(deprecated...)` | Mark the operation deprecated. |
| `Hidden(hide...)` | Exclude the operation from the generated document. |
| `Request(structure, contentOpts...)` | Add parameters and/or a request body from a Go value. |
| `Response(status, structure, contentOpts...)` | Add a response. Duplicate status/content pairs are merged into `oneOf`. |
| `CustomizeOperation(fn)` | Mutate the low-level `openapi.Operation` directly. |

**Content options:**

| Option | Purpose |
| --- | --- |
| `ContentType(contentType)` | Set media type; default is `application/json`. |
| `ContentDescription(description)` | Set request/response description. |
| `ContentSummary(summary)` | Set request/response summary (OpenAPI `3.2.0`). |
| `ContentDefault(isDefault...)` | Mark response as `default`. |
| `ContentEncoding(prop, enc)` | Add media type encoding metadata for a property. |
| `ContentExample(value)` | Set media type `example`. |
| `ContentNamedExample(name, value, opts...)` | Add one named media type example. |
| `ContentExamples(examples)` | Set media type `examples`. |
| `ContentRequired(required...)` | Mark request body as required. |
| `ContentFormat(format)` | Override the reflected content schema format. |

**Example options:** `ExampleSummary`, `ExampleDescription`, `ExampleExternalValue`, `ExampleDataValue` (3.2.0), `ExampleSerializedValue` (3.2.0).

---

## Security

```go
r := spec.NewRouter(
	option.WithSecurity("apiKey", option.SecurityAPIKey("X-API-Key", openapi.SecuritySchemeAPIKeyInHeader)),
	option.WithSecurity("bearerAuth", option.SecurityHTTPBearer("bearer")),
	option.WithSecurity("mtls", option.SecurityMutualTLS()),
	option.WithSecurity("oauth2", option.SecurityOAuth2AuthorizationCode(
		"https://auth.example.com/oauth/authorize",
		"https://auth.example.com/oauth/token",
		map[string]string{
			"read":  "Read access",
			"write": "Write access",
		},
	)),
	option.WithSecurity("oidc", option.SecurityOpenIDConnect("https://auth.example.com/.well-known/openid-configuration")),
)
```

**Security helpers:**

- `SecurityAPIKey(name, in)`
- `SecurityHTTPBearer(scheme, bearerFormat...)`
- `SecurityOAuth2(flows)`
- `SecurityOAuth2Implicit(authorizationURL, scopes, flowOpts...)`
- `SecurityOAuth2Password(tokenURL, scopes, flowOpts...)`
- `SecurityOAuth2ClientCredentials(tokenURL, scopes, flowOpts...)`
- `SecurityOAuth2AuthorizationCode(authorizationURL, tokenURL, scopes, flowOpts...)`
- `SecurityOAuth2DeviceAuthorization(deviceAuthorizationURL, tokenURL, scopes, flowOpts...)` — OpenAPI `3.2.0`
- `OAuthRefreshURL(url)`
- `SecurityOpenIDConnect(url)`
- `SecurityMutualTLS()`
- `SecurityDescription(description)`
- `SecurityOAuth2MetadataURL(url)` — OpenAPI `3.2.0`
- `SecurityDeprecated(deprecated...)` — OpenAPI `3.2.0`

---

## Reflection Tags

Request structs are split into parameters and request body:

- Fields with `path`, `query`, `header`, `cookie`, or `querystring` tags become OpenAPI parameters.
- `querystring` is only valid for OpenAPI `3.2.0`.
- Body fields use `json` names by default; for `application/x-www-form-urlencoded` and `multipart/form-data` they use `form` names, falling back to `json`.
- `json:"-"` skips a field.
- Path parameters are always marked required.
- Adapter packages can add framework-specific parameter tags with `ParameterTagMapping` (e.g. Gin uses `uri`, Fiber uses `params`, Echo uses `param`).

```go
type SearchRequest struct {
	ID       string `path:"id" required:"true" description:"Resource ID"`
	Status   string `query:"status" enum:"active,inactive"`
	TraceID  string `header:"X-Trace-ID"`
	Session  string `cookie:"session"`
	Name     string `json:"name" minLength:"2" maxLength:"80"`
	File     []byte `form:"file" format:"binary" description:"Upload file"`
	Internal string `json:"-"`
}
```

**Naming and location tags:**

| Tag | Purpose |
| --- | --- |
| `json:"name"` | Body/schema property name. |
| `path:"name"` | Path parameter. |
| `query:"name"` | Query parameter. |
| `header:"name"` | Header parameter. |
| `cookie:"name"` | Cookie parameter. |
| `querystring:"name"` | OpenAPI `3.2.0` whole-query-string parameter. |
| `mediaType:"..."` | Media type for `querystring` parameter content; defaults to `application/x-www-form-urlencoded`. OpenAPI `3.2.0` only. |
| `form:"name"` | Form body property name for form content types. |

**Schema tags:**

| Tag | Output |
| --- | --- |
| `required:"true"` | Adds property to parent `required`; path parameters are required automatically. |
| `type:"string,null"` | Overrides `type`; comma-separated unions emitted only for OpenAPI `3.1.x`/`3.2.0`. |
| `title:"..."` | `title`. |
| `description:"..."` | `description`. |
| `format:"email"` | `format`. |
| `pattern:"..."` | `pattern`. |
| `default:"..."` | `default`; JSON values are decoded when valid. |
| `example:"..."` | `example`; JSON values are decoded when valid. |
| `examples:"[...]"` | `examples` for OpenAPI `3.1.x`/`3.2.0`; JSON array preferred, comma-separated fallback. |
| `enum:"a,b,c"` | `enum`; JSON array or comma-separated values. |
| `const:"..."` | `const` for OpenAPI `3.1.x`/`3.2.0`. |
| `multipleOf:"2"` | `multipleOf`. |
| `maximum:"100"` | `maximum`. |
| `minimum:"1"` | `minimum`. |
| `exclusiveMaximum:"true"` | OpenAPI `3.0.x`: boolean; OpenAPI `3.1.x`/`3.2.0`: numeric value. |
| `exclusiveMinimum:"true"` | OpenAPI `3.0.x`: boolean; OpenAPI `3.1.x`/`3.2.0`: numeric value. |
| `maxLength:"80"` | `maxLength`. |
| `minLength:"2"` | `minLength`. |
| `maxItems:"10"` | `maxItems`. |
| `minItems:"1"` | `minItems`. |
| `maxProperties:"10"` | `maxProperties`. |
| `minProperties:"1"` | `minProperties`. |
| `uniqueItems:"true"` | `uniqueItems`. |
| `nullable:"true"` | `nullable` for OpenAPI `3.0.x`; `type: [T, null]` for OpenAPI `3.1.x`/`3.2.0` when possible. |
| `deprecated:"true"` | `deprecated`. |
| `readOnly:"true"` | `readOnly`. |
| `writeOnly:"true"` | `writeOnly`. |
| `contentEncoding:"base64"` | `contentEncoding` for OpenAPI `3.1.x`/`3.2.0`. |
| `contentMediaType:"image/png"` | `contentMediaType` for OpenAPI `3.1.x`/`3.2.0`. |
| `xmlName:"name"` | XML object `name`. |
| `xmlNamespace:"uri"` | XML object `namespace`. |
| `xmlPrefix:"p"` | XML object `prefix`. |
| `xmlAttribute:"true"` | XML object `attribute`; skipped when OpenAPI `3.2.0` `xmlNodeType` is set. |
| `xmlWrapped:"true"` | XML object `wrapped`; skipped when OpenAPI `3.2.0` `xmlNodeType` is set. |
| `xmlNodeType:"element"` | OpenAPI `3.2.0` XML object `nodeType`. |

> Reflection intentionally omits keywords that are invalid for the selected OpenAPI version. For example, `const`, `examples`, `contentEncoding`, and `contentMediaType` are not emitted for OpenAPI `3.0.x`.

---

## Reflected Go Types

| Go type | Schema |
| --- | --- |
| `bool` | `type: boolean` |
| Signed integers (except `int64`) | `type: integer`, `format: int32` |
| `int64` | `type: integer`, `format: int64` |
| Unsigned integers (except `uint64`/`uintptr`) | `type: integer`, `format: int32`, `minimum: 0` |
| `uint64`, `uintptr` | `type: integer`, `format: int64`, `minimum: 0` |
| `float32` | `type: number`, `format: float` |
| `float64` | `type: number`, `format: double` |
| `string` | `type: string` |
| `time.Time` | `type: string`, `format: date-time` |
| `[]T`, `[N]T` | `type: array`, `items: T` |
| `[]byte` | `3.0.x`: `type: string`, `format: byte`; `3.1.x`/`3.2.0`: `type: string`, `contentEncoding: base64` |
| `map[string]T` | `type: object`, `additionalProperties: T` |
| Structs | `type: object`, `properties` |
| Named structs (component mode) | `#/components/schemas/{TypeName}` reference |
| Pointers | Nullable schema behavior |

Custom types can expose their own schema when tags are not expressive enough:

```go
type Slug string

func (*Slug) OpenAPISchema(version string) *openapi.Schema {
	if version == openapi.Version304 {
		return &openapi.Schema{Type: "string", Format: "slug"}
	}
	return &openapi.Schema{Type: []string{"string", "null"}, Format: "slug"}
}
```

For static schemas, implement `OpenAPISchema() *openapi.Schema` instead. Field tags are still applied on top of custom schemas.

---

## Reflector Configuration

```go
r := spec.NewRouter(
	option.WithReflectorConfig(
		option.InlineRefs(),
		option.StripDefNamePrefix("DTO"),
		option.TypeMapping(NullString{}, new(string)),
		option.ParameterTagMapping(openapi.ParameterInPath, "params"),
		option.InterceptDefName(func(t reflect.Type, defaultName string) string {
			return defaultName
		}),
	),
)
```

| Option | Purpose |
| --- | --- |
| `InlineRefs(inline...)` | Inline schemas instead of using component references for named structs. |
| `StripDefNamePrefix(prefixes...)` | Strip prefixes from generated component names. |
| `TypeMapping(src, dst)` | Reflect `src` as if it were `dst`. |
| `ParameterTagMapping(in, sourceTag)` | Add a custom tag for a parameter location while keeping the default tag. |
| `InterceptDefName(fn)` | Customize schema component names. |

---

## OpenAPI 3.2

Selecting `openapi.Version320` enables the following additional features:

- `Router.Query(path, opts...)` — new `QUERY` HTTP method.
- Custom HTTP methods via `Add`, emitted as `additionalOperations`.
- `querystring` parameter tags.
- Root `$self` field.
- Server `name` field.
- Response `summary` field.
- Tag `parent` and `kind` fields.
- Security scheme metadata and deprecation fields.
- `components.mediaTypes`.
- Media type and encoding fields: `itemSchema`, `prefixEncoding`, `itemEncoding`.
- Discriminator `defaultMapping`.
- Example `dataValue` and `serializedValue` fields.
- XML `nodeType`.

---

## Low-Level OpenAPI Control

The `openapi` package exposes typed low-level OpenAPI structs. Use them when a feature does not require reflection or doesn't yet have a convenience option.

```go
r := spec.NewRouter(
	option.WithOpenAPIVersion(openapi.Version320),
	option.WithTitle("API"),
	option.WithVersion("1.0.0"),
	option.WithDocument(func(doc *openapi.Document) {
		doc.Webhooks = map[string]*openapi.PathItem{
			"user.created": {
				Post: &openapi.Operation{
					Responses: map[string]*openapi.Response{
						"202": {Description: "Accepted"},
					},
				},
			},
		}
		doc.Extensions = map[string]any{"x-company": "oaswrap"}
		doc.Extra = map[string]any{"futureOfficialField": true}
	}),
)
```

- Use `Extensions` for `x-*` specification extensions.
- Use `Extra` to emit official or future fields not yet typed.
- Use `CustomizeOperation` to mutate an operation directly.
- Use `WithDocument` to mutate the final document before output.

> The library validates many OpenAPI invariants but does not attempt exhaustive semantic validation for every rule in the specification.

---

## Examples

Complete working examples are in the [`examples/`](examples/) directory:

- **[Basic](examples/basic/)** — standalone spec generation.
- **[Petstore](examples/petstore/)** — full Petstore API with routes and models.
- **Adapter examples** — framework-specific examples in [`adapter/*/example`](adapter/).

---

## API Reference

Full documentation is available at [pkg.go.dev/github.com/oaswrap/spec](https://pkg.go.dev/github.com/oaswrap/spec).

| Package | Purpose |
| --- | --- |
| [`spec`](https://pkg.go.dev/github.com/oaswrap/spec) | Core router and spec builder. |
| [`openapi`](https://pkg.go.dev/github.com/oaswrap/spec/openapi) | Owned OpenAPI model. |
| [`option`](https://pkg.go.dev/github.com/oaswrap/spec/option) | Configuration options. |
| [`pkg/parser`](https://pkg.go.dev/github.com/oaswrap/spec/pkg/parser) | Path parsers such as `NewColonParamParser`. |

---

## Contributing

Issues and pull requests are welcome. Please check existing issues and discussions before starting work on new features.

---

## License

[MIT](LICENSE)

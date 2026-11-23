package spec_test

import (
	"mime/multipart"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/oaswrap/spec"
	"github.com/oaswrap/spec/internal/testutil"
	"github.com/oaswrap/spec/openapi"
	"github.com/oaswrap/spec/option"
)

type ComplexRequest struct {
	String   string         `json:"string" required:"true" minLength:"1" maxLength:"10" pattern:"^[a-z]+$"`
	Int      int            `json:"int" minimum:"1" maximum:"100"`
	Number   float64        `json:"number" multipleOf:"0.5"`
	Bool     bool           `json:"bool"`
	Enum     string         `json:"enum" enum:"a,b,c"`
	Array    []string       `json:"array" minItems:"1" maxItems:"5" uniqueItems:"true"`
	Object   *SimpleObject  `json:"object"`
	Map      map[string]int `json:"map"`
	Any      any            `json:"any"`
	Nullable *string        `json:"nullable"`
}

type SimpleObject struct {
	Foo string `json:"foo"`
}

type ComplexResponse struct {
	Data    ComplexRequest `json:"data"`
	Status  string         `json:"status"`
	Message string         `json:"message" deprecated:"true"`
}

type UploadRequest struct {
	File     []byte `form:"file" format:"binary" description:"The file to upload"`
	FileName string `form:"fileName" description:"Optional file name"`
}

type MultipartFileRequest struct {
	EventID     string                  `query:"eventId" required:"true"`
	Title       string                  `form:"title"`
	File        *multipart.FileHeader   `form:"file" required:"true" description:"The image to upload"`
	Attachments []*multipart.FileHeader `form:"attachments"`
}

type OneOfRequest struct {
	Value any `json:"value" oneOf:"string,int"`
}

type AnonymousStructRequest struct {
	Foo struct {
		Bar string `json:"bar"`
	} `json:"foo"`
}

type NestedRequest struct {
	Level1 struct {
		Level2 struct {
			Level3 string `json:"level3"`
		} `json:"level2"`
	} `json:"level1"`
}

type QueryStringRequest struct {
	Payload struct {
		Search string `json:"search"`
		Limit  int    `json:"limit"`
	} `querystring:"payload"`
}

// EmbedRef golden test types.

type EmbedRefBase struct {
	ID   int    `json:"id"`
	Role string `json:"role"`
}

type EmbedRefBaseViaInterface struct {
	Tag string `json:"tag"`
}

func (EmbedRefBaseViaInterface) ReferEmbedded() {}

type EmbedRefRequest struct {
	EmbedRefBase `refer:"true"`

	Name string `json:"name"`
	Age  int    `json:"age"`
}

type EmbedRefMixedRequest struct {
	EmbedRefBaseViaInterface

	Code  string `json:"code"`
	Value int    `json:"value"`
}

type mockPathParser struct{}

func (mockPathParser) Parse(path string) (string, error) {
	return strings.ReplaceAll(path, ":", "{") + "}", nil // simplified mock
}

type BaseResponse[T any] struct {
	Data    T      `json:"data"`
	Message string `json:"message"`
	Success bool   `json:"success"`
}

type ProfileResponse BaseResponse[User]

type AllBasicDataTypes struct {
	Int     int     `json:"int"`
	Int8    int8    `json:"int8"`
	Int16   int16   `json:"int16"`
	Int32   int32   `json:"int32"`
	Int64   int64   `json:"int64"`
	Uint    uint    `json:"uint"`
	Uint8   uint8   `json:"uint8"`
	Uint16  uint16  `json:"uint16"`
	Uint32  uint32  `json:"uint32"`
	Uint64  uint64  `json:"uint64"`
	Uintptr uintptr `json:"uintptr"`
	Float32 float32 `json:"float32"`
	Float64 float64 `json:"float64"`
	Byte    byte    `json:"byte"`
	Rune    rune    `json:"rune"`
	String  string  `json:"string"`
	Bool    bool    `json:"bool"`
}

type AllBasicDataTypesPointers struct {
	Int     *int     `json:"int"`
	Int8    *int8    `json:"int8"`
	Int16   *int16   `json:"int16"`
	Int32   *int32   `json:"int32"`
	Int64   *int64   `json:"int64"`
	Uint    *uint    `json:"uint"`
	Uint8   *uint8   `json:"uint8"`
	Uint16  *uint16  `json:"uint16"`
	Uint32  *uint32  `json:"uint32"`
	Uint64  *uint64  `json:"uint64"`
	Uintptr *uintptr `json:"uintptr"`
	Float32 *float32 `json:"float32"`
	Float64 *float64 `json:"float64"`
	Byte    *byte    `json:"byte"`
	Rune    *rune    `json:"rune"`
	String  *string  `json:"string"`
	Bool    *bool    `json:"bool"`
}

func TestGolden(t *testing.T) {
	allVersions := []string{openapi.Version304, openapi.Version312, openapi.Version320}
	cases := []struct {
		name     string
		versions []string
		opts     []option.OpenAPIOption
		run      func(r spec.Router)
	}{
		{
			name: "generics",
			opts: []option.OpenAPIOption{option.WithTitle("Generics API"), option.WithVersion("1.0.0")},
			run: func(r spec.Router) {
				r.Get("/user", option.Response(200, new(BaseResponse[User])))
				r.Post("/users", option.Response(201, new(BaseResponse[[]User])))
				r.Get("/profile", option.Response(200, new(ProfileResponse)))
			},
		},
		{
			name: "composition",
			opts: []option.OpenAPIOption{option.WithTitle("Composition API"), option.WithVersion("1.0.0")},
			run: func(r spec.Router) {
				r.Post("/oneof",
					option.Request(new(OneOfRequest)),
					option.Response(200, spec.OneOf("string", "int")),
				)
			},
		},
		{
			name: "anonymous_structs",
			run: func(r spec.Router) {
				r.Post("/anonymous", option.Request(new(AnonymousStructRequest)), option.Response(204, nil))
			},
		},
		{
			name: "embed_ref",
			opts: []option.OpenAPIOption{option.WithTitle("EmbedRef API"), option.WithVersion("1.0.0")},
			run: func(r spec.Router) {
				r.Post("/refer-tag",
					option.Request(new(EmbedRefRequest)),
					option.Response(200, new(EmbedRefBase)),
				)
				r.Post("/refer-interface",
					option.Request(new(EmbedRefMixedRequest)),
					option.Response(204, nil),
				)
			},
		},
		{
			name: "nested_structures",
			run: func(r spec.Router) {
				r.Post("/nested", option.Request(new(NestedRequest)), option.Response(204, nil))
			},
		},
		{
			name: "custom_path_parser",
			opts: []option.OpenAPIOption{option.WithPathParser(mockPathParser{})},
			run: func(r spec.Router) {
				r.Get("/users/:id", option.Request(new(GetUserRequest)), option.Response(200, new(User)))
			},
		},
		{
			name: "trailing_slash",
			opts: []option.OpenAPIOption{option.WithStripTrailingSlash(true)},
			run: func(r spec.Router) {
				r.Get("/ping/", option.Response(204, nil))
			},
		},
		{
			name: "multipart_binary",
			opts: []option.OpenAPIOption{option.WithTitle("Binary API"), option.WithVersion("1.0.0")},
			run: func(r spec.Router) {
				r.Post("/upload",
					option.OperationID("uploadFile"),
					option.Request(new(UploadRequest),
						option.ContentType("multipart/form-data"),
						option.ContentEncoding("file", "image/png"),
					),
					option.Response(201, nil, option.ContentDescription("File uploaded")),
				)
				r.Get("/download",
					option.OperationID("downloadFile"),
					option.Response(200, "string",
						option.ContentType("image/png"),
						option.ContentDescription("The image file"),
						option.ContentFormat("binary"),
					),
				)
				r.Post("/upload-raw",
					option.OperationID("uploadRaw"),
					option.Request(nil,
						option.ContentType("image/png"),
						option.ContentFormat("binary"),
					),
					option.Response(204, nil),
				)
			},
		},
		{
			name: "multipart_file",
			opts: []option.OpenAPIOption{option.WithTitle("File API"), option.WithVersion("1.0.0")},
			run: func(r spec.Router) {
				r.Post("/events/image",
					option.OperationID("uploadEventImage"),
					option.Request(new(MultipartFileRequest)),
					option.Response(201, nil, option.ContentDescription("File uploaded")),
				)
			},
		},
		{
			name: "spec_information",
			opts: []option.OpenAPIOption{
				option.WithTitle("Users API"),
				option.WithVersion("1.2.3"),
				option.WithDescription("User account operations."),
				option.WithTermsOfService("https://example.com/terms"),
				option.WithContact(openapi.Contact{Name: "API Team", Email: "api@example.com"}),
				option.WithLicense(
					openapi.License{Name: "Apache 2.0", URL: "https://www.apache.org/licenses/LICENSE-2.0.html"},
				),
				option.WithExternalDocs("https://docs.example.com", "API documentation"),
				option.WithServer("https://api.example.com"),
				option.WithTags(openapi.Tag{Name: "users", Description: "User operations"}),
			},
			run: func(r spec.Router) {
				r.Get("/ping", option.Response(204, nil))
			},
		},
		{
			name: "request_response",
			opts: []option.OpenAPIOption{option.WithTitle("Login API"), option.WithVersion("1.0.0")},
			run: func(r spec.Router) {
				r.Post("/login",
					option.OperationID("login"),
					option.Summary("Login"),
					option.Request(new(LoginRequest)),
					option.Response(200, new(LoginResponse)),
				)
			},
		},
		{
			name: "path_parameters",
			opts: []option.OpenAPIOption{option.WithTitle("Users API"), option.WithVersion("1.0.0")},
			run: func(r spec.Router) {
				r.Get("/users/{id}",
					option.OperationID("getUser"),
					option.Summary("Get user"),
					option.Request(new(GetUserRequest)),
					option.Response(200, new(User)),
				)
			},
		},
		{
			name: "security",
			opts: []option.OpenAPIOption{
				option.WithTitle("Secure API"),
				option.WithVersion("1.0.0"),
				option.WithSecurity("bearerAuth", option.SecurityHTTPBearer("bearer")),
				option.WithSecurity("apiKey", option.SecurityAPIKey("X-API-Key", "header")),
			},
			run: func(r spec.Router) {
				r.Get("/me",
					option.Security("bearerAuth"),
					option.Response(200, new(User)),
				)
			},
		},
		{
			name: "complex_types",
			opts: []option.OpenAPIOption{option.WithTitle("Complex API"), option.WithVersion("1.0.0")},
			run: func(r spec.Router) {
				r.Post("/complex",
					option.OperationID("complex"),
					option.Request(new(ComplexRequest)),
					option.Response(200, new(ComplexResponse)),
				)
			},
		},
		{
			name: "multiple_content_types",
			opts: []option.OpenAPIOption{option.WithTitle("Multi-content API"), option.WithVersion("1.0.0")},
			run: func(r spec.Router) {
				r.Post("/multi",
					option.OperationID("multi"),
					option.Request(new(LoginRequest), option.ContentType("application/json")),
					option.Request(new(LoginRequest), option.ContentType("application/x-www-form-urlencoded")),
					option.Response(200, new(LoginResponse), option.ContentType("application/json")),
					option.Response(200, new(LoginResponse), option.ContentType("application/xml")),
				)
			},
		},
		{
			name: "server_variables",
			opts: []option.OpenAPIOption{
				option.WithTitle("Server Var API"),
				option.WithVersion("1.0.0"),
				option.WithServer("https://{environment}.example.com/v1",
					option.ServerVariables(map[string]openapi.ServerVariable{
						"environment": {
							Default:     "production",
							Description: "API environment",
							Enum:        []string{"production", "staging", "dev"},
						},
					}),
				),
			},
			run: func(r spec.Router) {
				r.Get("/ping", option.Response(204, nil))
			},
		},
		{
			name:     "openapi_320_operations",
			versions: []string{openapi.Version320},
			opts:     []option.OpenAPIOption{option.WithTitle("Search API"), option.WithVersion("1.0.0")},
			run: func(r spec.Router) {
				r.Query("/search", option.Response(200, new([]User)))
				r.Add("PURGE", "/cache", option.Response(204, nil))
			},
		},
		{
			name:     "compatibility_extensions",
			versions: []string{openapi.Version312},
			opts: []option.OpenAPIOption{
				option.WithTitle("Compatibility API"),
				option.WithVersion("1.0.0"),
				option.WithSecurity("mtls", option.SecurityMutualTLS()),
				option.WithDocument(func(doc *openapi.Document) {
					doc.Extensions = map[string]any{"x-root": "ok"}
					if doc.OpenAPI != openapi.Version304 {
						doc.Webhooks = map[string]*openapi.PathItem{
							"user.created": {
								Post: &openapi.Operation{
									Responses: map[string]*openapi.Response{"202": {Description: "Accepted"}},
								},
							},
						}
					}
				}),
			},
			run: func(r spec.Router) {
				r.Get("/users/{id}",
					option.Request(new(GetUserRequest)),
					option.Response(200, new(User)),
					option.CustomizeOperation(func(op *openapi.Operation) {
						op.Extensions = map[string]any{"x-operation": "ok"}
					}),
				)
			},
		},
		{
			name: "basic_data_types",
			opts: []option.OpenAPIOption{option.WithTitle("Basic Types API"), option.WithVersion("1.0.0")},
			run: func(r spec.Router) {
				r.Post("/basic-types",
					option.OperationID("basicTypes"),
					option.Request(new(AllBasicDataTypes)),
					option.Response(200, new(AllBasicDataTypes)),
				)
				r.Post("/basic-types-pointers",
					option.OperationID("basicTypesPointers"),
					option.Request(new(AllBasicDataTypesPointers)),
					option.Response(200, new(AllBasicDataTypesPointers)),
				)
			},
		},
		{
			name:     "webhook_helpers",
			versions: []string{openapi.Version312, openapi.Version320},
			opts:     []option.OpenAPIOption{option.WithTitle("Webhook API"), option.WithVersion("1.0.0")},
			run: func(r spec.Router) {
				r.Webhook("user.created", option.Response(202, nil))
				r.AddWebhook("POST", "cache.invalidate", option.Response(204, nil))
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			versions := tc.versions
			if len(versions) == 0 {
				versions = allVersions
			}
			for _, v := range versions {
				t.Run(v, func(t *testing.T) {
					r := spec.NewRouter(append(tc.opts, option.WithOpenAPIVersion(v))...)
					if tc.run != nil {
						tc.run(r)
					}
					schema, err := r.GenerateSchema("yaml")
					require.NoError(t, err)
					versionSuffix := strings.ReplaceAll(v[:3], ".", "")
					testutil.AssertGolden(t, schema, filepath.Join("testdata", tc.name+".v"+versionSuffix+".yaml"))
				})
			}
		})
	}
}

func TestGoldenPetstore(t *testing.T) {
	versions := []string{openapi.Version304, openapi.Version312, openapi.Version320}
	for _, v := range versions {
		t.Run(v, func(t *testing.T) {
			r := newPetstoreRouter(option.WithOpenAPIVersion(v))
			schema, err := r.GenerateSchema("yaml")
			require.NoError(t, err)
			versionSuffix := strings.ReplaceAll(v[:3], ".", "")
			testutil.AssertGolden(t, schema, filepath.Join("testdata", "petstore.v"+versionSuffix+".yaml"))
		})
	}
}

func TestGoldenOpenAPI320Features(t *testing.T) {
	r := spec.NewRouter(
		option.WithOpenAPIVersion(openapi.Version320),
		option.WithSelf("https://api.example.com/openapi.yaml"),
		option.WithTitle("OpenAPI 3.2 Features API"),
		option.WithVersion("1.0.0"),
		option.WithTags(openapi.Tag{
			Name:        "commerce",
			Summary:     "Commerce",
			Description: "Commerce APIs.",
			Kind:        "nav",
		}),
		option.WithTags(openapi.Tag{
			Name:        "payments",
			Summary:     "Payments",
			Description: "Payment operations.",
			Parent:      "commerce",
			Kind:        "nav",
		}),
		option.WithSecurity("deviceAuth",
			option.SecurityOAuth2(openapi.OAuthFlows{
				DeviceAuthorization: &openapi.OAuthFlow{
					DeviceAuthorizationURL: "https://auth.example.com/device",
					TokenURL:               "https://auth.example.com/token",
					Scopes: map[string]string{
						"payments:read": "Read payments",
					},
				},
			}),
			option.SecurityOAuth2MetadataURL("https://auth.example.com/.well-known/oauth-authorization-server"),
			option.SecurityDeprecated(),
		),
		option.WithGlobalSecurity("deviceAuth", "payments:read"),
	)
	r.Get("/payments/{id}",
		option.Tags("payments"),
		option.Request(new(GetUserRequest)),
		option.Response(200, new(User)),
		option.CustomizeOperation(func(op *openapi.Operation) {
			resp := op.Responses["200"]
			resp.Summary = "Payment found"
			mt := resp.Content["application/json"]
			mt.Examples = map[string]*openapi.Example{
				"encoded-id": {
					Summary:         "Encoded identifier",
					DataValue:       map[string]any{"id": "pay_123"},
					SerializedValue: `{"id":"pay_123"}`,
				},
			}
			resp.Content["application/json"] = mt
		}),
	)

	r.Get("/search",
		option.Tags("commerce"),
		option.Request(new(QueryStringRequest)),
		option.Response(200, new([]User)),
	)

	schema, err := r.GenerateSchema("yaml")
	require.NoError(t, err)
	testutil.AssertGolden(t, schema, filepath.Join("testdata", "openapi_320_features.v32.yaml"))
}

func TestGoldenOpenAPI312ReferenceDescriptions(t *testing.T) {
	r := spec.NewRouter(
		option.WithOpenAPIVersion(openapi.Version312),
		option.WithTitle("Reference Descriptions API"),
		option.WithVersion("1.0.0"),
		option.WithDocument(func(doc *openapi.Document) {
			doc.Components.Schemas = map[string]*openapi.Schema{
				"User": {
					Type:     "object",
					Required: []string{"id"},
					Properties: map[string]*openapi.Schema{
						"id":   {Type: "string"},
						"name": {Type: "string"},
					},
				},
			}
			doc.Components.Parameters = map[string]*openapi.Parameter{
				"TraceID": {
					Name:     "traceId",
					In:       "header",
					Required: true,
					Schema:   &openapi.Schema{Type: "string"},
				},
			}
			doc.Components.RequestBodies = map[string]*openapi.RequestBody{
				"UserBody": {
					Content: map[string]openapi.MediaType{
						"application/json": {Schema: &openapi.Schema{Ref: "#/components/schemas/User"}},
					},
				},
			}
			doc.Components.Responses = map[string]*openapi.Response{
				"UserResponse": {
					Description: "User response",
					Content: map[string]openapi.MediaType{
						"application/json": {Schema: &openapi.Schema{Ref: "#/components/schemas/User"}},
					},
				},
			}
			doc.Components.Examples = map[string]*openapi.Example{
				"UserExample": {
					Value: map[string]string{"id": "user-123", "name": "Ada"},
				},
			}
			doc.Components.Links = map[string]*openapi.Link{
				"FindUser": {OperationID: "findUser"},
			}
		}),
	)
	r.Get("/users/{id}",
		option.OperationID("findUser"),
		option.Request(new(GetUserRequest)),
		option.CustomizeOperation(func(op *openapi.Operation) {
			op.Parameters = append(op.Parameters, &openapi.Parameter{
				Ref:         "#/components/parameters/TraceID",
				Description: "Request trace identifier.",
			})
			op.RequestBody = &openapi.RequestBody{
				Ref:         "#/components/requestBodies/UserBody",
				Description: "Reusable user payload.",
			}
			op.Responses = map[string]*openapi.Response{
				"200": {
					Ref:         "#/components/responses/UserResponse",
					Description: "Reusable user response.",
				},
				"204": {
					Description: "No content",
					Links: map[string]*openapi.Link{
						"findUser": {
							Ref:         "#/components/links/FindUser",
							Description: "Reusable follow-up link.",
						},
					},
				},
			}
		}),
	)

	schema, err := r.GenerateSchema("yaml")
	require.NoError(t, err)
	testutil.AssertGolden(t, schema, filepath.Join("testdata", "openapi_312_reference_descriptions.v31.yaml"))
}

func TestOutputKeepsOpenAPIObjectOrder(t *testing.T) {
	r := spec.NewRouter(option.WithTitle("Ordered API"), option.WithVersion("1.0.0"))
	r.Get("/users/{id}", option.Request(new(GetUserRequest)), option.Response(200, new(User)))

	yamlRaw, err := r.GenerateSchema("yaml")
	require.NoError(t, err)
	assertContainsInOrder(t, string(yamlRaw), "\n", "openapi:", "info:", "paths:", "components:")

	jsonRaw, err := r.GenerateSchema("json")
	require.NoError(t, err)
	assertContainsInOrder(t, string(jsonRaw), "", `"openapi":`, `"info":`, `"paths":`, `"components":`)
}

func assertContainsInOrder(t *testing.T, value, separator string, parts ...string) {
	t.Helper()
	offset := 0
	for _, part := range parts {
		index := strings.Index(value[offset:], part)
		if !assert.GreaterOrEqual(t, index, 0, "expected %q after offset %d in:\n%s", part, offset, value) {
			t.FailNow()
		}
		offset += index + len(part)
		if separator != "" {
			offset += len(separator)
		}
	}
}

package option

import (
	"errors"
	"log/slog"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/oaswrap/spec-ui/config"
	"github.com/oaswrap/spec/openapi"
)

func TestWithOpenAPIConfig(t *testing.T) {
	cfg := WithOpenAPIConfig(
		WithOpenAPIVersion(openapi.Version312),
		WithSelf("https://api.test.com/openapi.yaml"),
		WithJSONSchemaDialect("https://json-schema.org/draft/2020-12/schema"),
		WithTitle("Test API"),
		WithInfoSummary("Short API summary"),
		WithVersion("2.0.0"),
		WithDescription("Test description"),
		WithTermsOfService("https://example.com/terms"),
		WithContact(openapi.Contact{Name: "Test Contact"}),
		WithLicense(openapi.License{Name: "MIT"}),
		WithTags(openapi.Tag{Name: "test"}),
		WithServer("https://api.test.com", ServerDescription("Test server"), ServerName("prod")),
		WithExternalDocs("https://docs.test.com", "Test docs"),
		WithSecurity("apiKey", SecurityAPIKey("X-API-Key", "header")),
		WithGlobalSecurity("apiKey", "read", "write"),
		WithStripTrailingSlash(true),
	)

	assert.Equal(t, openapi.Version312, cfg.OpenAPIVersion)
	assert.Equal(t, "https://api.test.com/openapi.yaml", cfg.Self)
	assert.Equal(t, "https://json-schema.org/draft/2020-12/schema", cfg.JSONSchemaDialect)
	assert.Equal(t, "Test API", cfg.Title)
	assert.Equal(t, "Short API summary", cfg.InfoSummary)
	assert.Equal(t, "2.0.0", cfg.Version)
	assert.Equal(t, "Test description", *cfg.Description)
	assert.Equal(t, "https://example.com/terms", *cfg.TermsOfService)
	assert.Equal(t, "Test Contact", cfg.Contact.Name)
	assert.Equal(t, "MIT", cfg.License.Name)

	if assert.Len(t, cfg.Tags, 1) {
		assert.Equal(t, "test", cfg.Tags[0].Name)
	}
	if assert.Len(t, cfg.Servers, 1) {
		assert.Equal(t, "https://api.test.com", cfg.Servers[0].URL)
		assert.Equal(t, "Test server", *cfg.Servers[0].Description)
		assert.Equal(t, "prod", cfg.Servers[0].Name)
	}
	assert.Equal(t, "https://docs.test.com", cfg.ExternalDocs.URL)
	assert.Equal(t, "Test docs", cfg.ExternalDocs.Description)
	assert.Equal(t, "apiKey", cfg.SecuritySchemes["apiKey"].Type)
	if assert.Len(t, cfg.Security, 1) {
		assert.Equal(t, "read", cfg.Security[0]["apiKey"][0])
	}
	assert.True(t, cfg.StripTrailingSlash)

	t.Run("ReflectorConfig", func(t *testing.T) {
		c := WithOpenAPIConfig(WithReflectorConfig(InlineRefs(true)))
		assert.True(t, c.ReflectorConfig.InlineRefs)
	})

	t.Run("PathParser", func(t *testing.T) {
		c := WithOpenAPIConfig(WithPathParser(nil))
		assert.Nil(t, c.PathParser)
	})

	t.Run("GlobalSecurityWithoutScopes", func(t *testing.T) {
		c := WithOpenAPIConfig(WithGlobalSecurity("apiKey"))
		if assert.Len(t, c.Security, 1) {
			assert.Empty(t, c.Security[0]["apiKey"])
		}
	})

	t.Run("WithDocument", func(t *testing.T) {
		called := false
		c := WithOpenAPIConfig(WithDocument(func(_ *openapi.Document) { called = true }))
		if assert.Len(t, c.DocumentCustomizers, 1) {
			c.DocumentCustomizers[0](nil)
			assert.True(t, called)
		}
	})
}

func TestTagOptions(t *testing.T) {
	cfg := WithOpenAPIConfig(
		WithTag("payments",
			TagSummary("Payments"),
			TagDescription("Payment operations"),
			TagExternalDocs("https://docs.test/payments", "Payment docs"),
			TagParent("commerce"),
			TagKind("nav"),
		),
	)

	if assert.Len(t, cfg.Tags, 1) {
		tag := cfg.Tags[0]
		assert.Equal(t, "payments", tag.Name)
		assert.Equal(t, "Payments", tag.Summary)
		assert.Equal(t, "Payment operations", tag.Description)
		assert.Equal(t, "https://docs.test/payments", tag.ExternalDocs.URL)
		assert.Equal(t, "Payment docs", tag.ExternalDocs.Description)
		assert.Equal(t, "commerce", tag.Parent)
		assert.Equal(t, "nav", tag.Kind)
	}
}

func TestComponentOptions(t *testing.T) {
	doc := &openapi.Document{}
	cfg := WithOpenAPIConfig(
		WithComponentSchema("User", &openapi.Schema{Type: "object"}),
		WithComponentResponse("OK", &openapi.Response{Description: "OK"}),
		WithComponentParameter("TraceID", &openapi.Parameter{Name: "X-Trace-ID", In: "header"}),
		WithComponentExample("UserExample", &openapi.Example{Value: map[string]any{"id": "123"}}),
		WithComponentRequestBody("CreateUser", &openapi.RequestBody{Content: map[string]openapi.MediaType{}}),
		WithComponentHeader("RateLimit", &openapi.Header{Description: "Rate limit"}),
		WithComponentSecurityScheme("ApiKey", &openapi.SecurityScheme{Type: "apiKey", Name: "X-API-Key", In: "header"}),
		WithComponentLink("UserLink", &openapi.Link{OperationID: "getUser"}),
		WithComponentCallback("Event", &openapi.Callback{}),
		WithComponentPathItem("Users", &openapi.PathItem{}),
		WithComponentMediaType("JSON", &openapi.MediaType{Schema: &openapi.Schema{Type: "object"}}),
	)

	for _, customize := range cfg.DocumentCustomizers {
		customize(doc)
	}

	assert.Equal(t, "object", doc.Components.Schemas["User"].Type)
	assert.Equal(t, "OK", doc.Components.Responses["OK"].Description)
	assert.Equal(t, "X-Trace-ID", doc.Components.Parameters["TraceID"].Name)
	assert.NotNil(t, doc.Components.Examples["UserExample"])
	assert.NotNil(t, doc.Components.RequestBodies["CreateUser"])
	assert.Equal(t, "Rate limit", doc.Components.Headers["RateLimit"].Description)
	assert.Equal(t, "apiKey", doc.Components.SecuritySchemes["ApiKey"].Type)
	assert.Equal(t, "getUser", doc.Components.Links["UserLink"].OperationID)
	assert.NotNil(t, doc.Components.Callbacks["Event"])
	assert.NotNil(t, doc.Components.PathItems["Users"])
	assert.NotNil(t, doc.Components.MediaTypes["JSON"])
}

func TestOperationOptions(t *testing.T) {
	cfg := &OperationConfig{}
	opts := []OperationOption{
		Hidden(true),
		OperationID("testOp"),
		Summary("Test Summary"),
		Description("Test Description"),
		ExternalDocs("https://docs.test/operation", "Operation docs"),
		Deprecated(true),
		Tags("tag1", "tag2"),
		Security("auth", "scope1"),
		Request(struct{}{}, ContentType("application/json")),
		Response(200, struct{}{}, ContentDescription("OK")),
		CustomizeOperation(func(op *openapi.Operation) { op.Summary = "Custom" }),
	}

	for _, opt := range opts {
		opt(cfg)
	}

	assert.True(t, cfg.Hide)
	assert.Equal(t, "testOp", cfg.OperationID)
	assert.Equal(t, "Test Summary", cfg.Summary)
	assert.Equal(t, "Test Description", cfg.Description)
	assert.Equal(t, "https://docs.test/operation", cfg.ExternalDocs.URL)
	assert.Equal(t, "Operation docs", cfg.ExternalDocs.Description)
	assert.True(t, cfg.Deprecated)
	assert.Equal(t, []string{"tag1", "tag2"}, cfg.Tags)
	if assert.Len(t, cfg.Security, 1) {
		assert.Equal(t, "auth", cfg.Security[0].Name)
		assert.Equal(t, "scope1", cfg.Security[0].Scopes[0])
	}
	if assert.Len(t, cfg.Requests, 1) {
		assert.Equal(t, "application/json", cfg.Requests[0].ContentType)
	}
	if assert.Len(t, cfg.Responses, 1) {
		assert.Equal(t, "OK", cfg.Responses[0].Description)
	}
	assert.Len(t, cfg.Customizers, 1)
}

func TestGroupOptions(t *testing.T) {
	cfg := &GroupConfig{}
	opts := []GroupOption{
		GroupHidden(true),
		GroupDeprecated(true),
		GroupTags("g1", "g2"),
		GroupSecurity("gauth", "gscope"),
	}

	for _, opt := range opts {
		opt(cfg)
	}

	assert.True(t, cfg.Hide)
	assert.True(t, cfg.Deprecated)
	assert.Equal(t, []string{"g1", "g2"}, cfg.Tags)
	if assert.Len(t, cfg.Security, 1) {
		assert.Equal(t, "gauth", cfg.Security[0].Name)
		assert.Equal(t, "gscope", cfg.Security[0].Scopes[0])
	}
}

func TestReflectorOptions(t *testing.T) {
	cfg := &openapi.ReflectorConfig{}
	opts := []ReflectorOption{
		InlineRefs(true),
		StripDefNamePrefix("Pre"),
		InterceptDefName(func(_ reflect.Type, _ string) string { return "Intercepted" }),
		TypeMapping(1, "one"),
		ParameterTagMapping(openapi.ParameterInQuery, "q"),
		InterceptProp(func(_ openapi.InterceptPropParams) error { return nil }),
		InterceptSchema(func(_ openapi.InterceptSchemaParams) (bool, error) { return false, nil }),
	}

	for _, opt := range opts {
		opt(cfg)
	}

	assert.True(t, cfg.InlineRefs)
	if assert.Len(t, cfg.StripDefNamePrefix, 1) {
		assert.Equal(t, "Pre", cfg.StripDefNamePrefix[0])
	}
	assert.Equal(t, "Intercepted", cfg.InterceptDefName(nil, ""))
	if assert.Len(t, cfg.TypeMappings, 1) {
		assert.Equal(t, 1, cfg.TypeMappings[0].Src)
		assert.Equal(t, "one", cfg.TypeMappings[0].Dst)
	}
	assert.Equal(t, "q", cfg.ParameterTagMapping[openapi.ParameterInQuery])
	assert.NotNil(t, cfg.InterceptProp)
	assert.NotNil(t, cfg.InterceptSchema)

	t.Run("RequiredPropByValidateTagSetsInterceptProp", func(t *testing.T) {
		c := &openapi.ReflectorConfig{}
		RequiredPropByValidateTag()(c)
		assert.NotNil(t, c.InterceptProp)
	})

	t.Run("InterceptSchemaChainsAndShortCircuits", func(t *testing.T) {
		t.Run("calls next when previous continues", func(t *testing.T) {
			c := &openapi.ReflectorConfig{}
			calls := []string{}
			InterceptSchema(func(_ openapi.InterceptSchemaParams) (bool, error) {
				calls = append(calls, "first")
				return false, nil
			})(c)
			InterceptSchema(func(_ openapi.InterceptSchemaParams) (bool, error) {
				calls = append(calls, "second")
				return true, nil
			})(c)

			stop, err := c.InterceptSchema(openapi.InterceptSchemaParams{})
			require.NoError(t, err)
			assert.True(t, stop)
			assert.Equal(t, []string{"first", "second"}, calls)
		})

		t.Run("skips next when previous stops", func(t *testing.T) {
			c := &openapi.ReflectorConfig{}
			calls := []string{}
			InterceptSchema(func(_ openapi.InterceptSchemaParams) (bool, error) {
				calls = append(calls, "first")
				return true, nil
			})(c)
			InterceptSchema(func(_ openapi.InterceptSchemaParams) (bool, error) {
				calls = append(calls, "second")
				return false, nil
			})(c)

			stop, err := c.InterceptSchema(openapi.InterceptSchemaParams{})
			require.NoError(t, err)
			assert.True(t, stop)
			assert.Equal(t, []string{"first"}, calls)
		})

		t.Run("skips next when previous errors", func(t *testing.T) {
			boom := errors.New("boom")
			c := &openapi.ReflectorConfig{}
			calls := []string{}
			InterceptSchema(func(_ openapi.InterceptSchemaParams) (bool, error) {
				calls = append(calls, "first")
				return false, boom
			})(c)
			InterceptSchema(func(_ openapi.InterceptSchemaParams) (bool, error) {
				calls = append(calls, "second")
				return false, nil
			})(c)

			stop, err := c.InterceptSchema(openapi.InterceptSchemaParams{})
			assert.False(t, stop)
			require.ErrorIs(t, err, boom)
			assert.Equal(t, []string{"first"}, calls)
		})
	})

	t.Run("InterceptPropChainsAndShortCircuits", func(t *testing.T) {
		t.Run("calls next when previous succeeds", func(t *testing.T) {
			c := &openapi.ReflectorConfig{}
			calls := []string{}
			InterceptProp(func(_ openapi.InterceptPropParams) error {
				calls = append(calls, "first")
				return nil
			})(c)
			InterceptProp(func(_ openapi.InterceptPropParams) error {
				calls = append(calls, "second")
				return nil
			})(c)

			err := c.InterceptProp(openapi.InterceptPropParams{})
			require.NoError(t, err)
			assert.Equal(t, []string{"first", "second"}, calls)
		})

		t.Run("skips next when previous errors", func(t *testing.T) {
			boom := errors.New("boom")
			c := &openapi.ReflectorConfig{}
			calls := []string{}
			InterceptProp(func(_ openapi.InterceptPropParams) error {
				calls = append(calls, "first")
				return boom
			})(c)
			InterceptProp(func(_ openapi.InterceptPropParams) error {
				calls = append(calls, "second")
				return nil
			})(c)

			err := c.InterceptProp(openapi.InterceptPropParams{})
			require.ErrorIs(t, err, boom)
			assert.Equal(t, []string{"first"}, calls)
		})
	})

	t.Run("RequiredPropByValidateTag", func(t *testing.T) {
		type reqDefault struct {
			ID string `validate:"required,min=3"`
		}
		type reqCustom struct {
			ID string `rules:"min=3|required"`
		}
		type notRequired struct {
			Name string `validate:"min=1"`
		}

		t.Run("adds required on processed field with default tag", func(t *testing.T) {
			field, ok := reflect.TypeFor[reqDefault]().FieldByName("ID")
			assert.True(t, ok)
			parent := &openapi.Schema{}
			c := &openapi.ReflectorConfig{}
			RequiredPropByValidateTag()(c)

			err := c.InterceptProp(openapi.InterceptPropParams{
				Name:         "id",
				Field:        field,
				ParentSchema: parent,
				Processed:    true,
			})
			require.NoError(t, err)
			assert.Equal(t, []string{"id"}, parent.Required)
		})

		t.Run("does not add required before processed", func(t *testing.T) {
			field, ok := reflect.TypeFor[reqDefault]().FieldByName("ID")
			assert.True(t, ok)
			parent := &openapi.Schema{}
			c := &openapi.ReflectorConfig{}
			RequiredPropByValidateTag()(c)

			err := c.InterceptProp(openapi.InterceptPropParams{
				Name:         "id",
				Field:        field,
				ParentSchema: parent,
				Processed:    false,
			})
			require.NoError(t, err)
			assert.Empty(t, parent.Required)
		})

		t.Run("supports custom tag and separator", func(t *testing.T) {
			field, ok := reflect.TypeFor[reqCustom]().FieldByName("ID")
			assert.True(t, ok)
			parent := &openapi.Schema{}
			c := &openapi.ReflectorConfig{}
			RequiredPropByValidateTag("rules", "|")(c)

			err := c.InterceptProp(openapi.InterceptPropParams{
				Name:         "id",
				Field:        field,
				ParentSchema: parent,
				Processed:    true,
			})
			require.NoError(t, err)
			assert.Equal(t, []string{"id"}, parent.Required)
		})

		t.Run("ignores fields without required marker", func(t *testing.T) {
			field, ok := reflect.TypeFor[notRequired]().FieldByName("Name")
			assert.True(t, ok)
			parent := &openapi.Schema{}
			c := &openapi.ReflectorConfig{}
			RequiredPropByValidateTag()(c)

			err := c.InterceptProp(openapi.InterceptPropParams{
				Name:         "name",
				Field:        field,
				ParentSchema: parent,
				Processed:    true,
			})
			require.NoError(t, err)
			assert.Empty(t, parent.Required)
		})
	})
}

func TestSecurityOptions(t *testing.T) {
	t.Run("APIKey", func(t *testing.T) {
		cfg := &securityConfig{}
		SecurityAPIKey("api_key", "query")(cfg)
		assert.Equal(t, "apiKey", cfg.scheme.Type)
		assert.Equal(t, "api_key", cfg.scheme.Name)
		assert.Equal(t, openapi.SecuritySchemeAPIKeyInQuery, cfg.scheme.In)
	})

	t.Run("HTTPBearer", func(t *testing.T) {
		cfg := &securityConfig{}
		SecurityHTTPBearer("bearer", "JWT")(cfg)
		assert.Equal(t, "http", cfg.scheme.Type)
		assert.Equal(t, "bearer", cfg.scheme.Scheme)
		assert.Equal(t, "JWT", *cfg.scheme.BearerFormat)
	})

	t.Run("OAuth2", func(t *testing.T) {
		cfg := &securityConfig{}
		flows := openapi.OAuthFlows{Implicit: &openapi.OAuthFlow{AuthorizationURL: "https://auth.com"}}
		SecurityOAuth2(flows)(cfg)
		assert.Equal(t, "oauth2", cfg.scheme.Type)
		assert.Equal(t, "https://auth.com", cfg.scheme.Flows.Implicit.AuthorizationURL)
	})

	t.Run("OAuth2FlowHelpers", func(t *testing.T) {
		scopes := map[string]string{"read": "Read access"}

		cfg := &securityConfig{}
		SecurityOAuth2Implicit("https://auth.com/authorize", scopes, OAuthRefreshURL("https://auth.com/refresh"))(cfg)
		assert.Equal(t, "https://auth.com/authorize", cfg.scheme.Flows.Implicit.AuthorizationURL)
		assert.Equal(t, "https://auth.com/refresh", *cfg.scheme.Flows.Implicit.RefreshURL)

		cfg = &securityConfig{}
		SecurityOAuth2Password("https://auth.com/token", nil)(cfg)
		assert.Equal(t, "https://auth.com/token", cfg.scheme.Flows.Password.TokenURL)
		assert.Empty(t, cfg.scheme.Flows.Password.Scopes)

		cfg = &securityConfig{}
		SecurityOAuth2ClientCredentials("https://auth.com/token", scopes)(cfg)
		assert.Equal(t, "https://auth.com/token", cfg.scheme.Flows.ClientCredentials.TokenURL)

		cfg = &securityConfig{}
		SecurityOAuth2AuthorizationCode("https://auth.com/authorize", "https://auth.com/token", scopes)(cfg)
		assert.Equal(t, "https://auth.com/authorize", cfg.scheme.Flows.AuthorizationCode.AuthorizationURL)
		assert.Equal(t, "https://auth.com/token", cfg.scheme.Flows.AuthorizationCode.TokenURL)

		cfg = &securityConfig{}
		SecurityOAuth2DeviceAuthorization("https://auth.com/device", "https://auth.com/token", scopes)(cfg)
		assert.Equal(t, "https://auth.com/device", cfg.scheme.Flows.DeviceAuthorization.DeviceAuthorizationURL)
		assert.Equal(t, "https://auth.com/token", cfg.scheme.Flows.DeviceAuthorization.TokenURL)
	})

	t.Run("MutualTLS", func(t *testing.T) {
		cfg := &securityConfig{}
		SecurityMutualTLS()(cfg)
		assert.Equal(t, "mutualTLS", cfg.scheme.Type)
	})

	t.Run("OpenIDConnect", func(t *testing.T) {
		cfg := &securityConfig{}
		SecurityOpenIDConnect("https://oidc.com")(cfg)
		assert.Equal(t, "openIdConnect", cfg.scheme.Type)
		assert.Equal(t, "https://oidc.com", cfg.scheme.OpenIDConnectURL)
	})

	t.Run("Description", func(t *testing.T) {
		cfg := &securityConfig{}
		SecurityDescription("desc")(cfg)
		SecurityAPIKey("k", "h")(cfg)
		assert.Equal(t, "desc", *cfg.scheme.Description)
	})

	t.Run("OAuth2MetadataURL", func(t *testing.T) {
		cfg := &securityConfig{}
		SecurityOAuth2MetadataURL("https://issuer.test/.well-known/oauth-authorization-server")(cfg)
		SecurityHTTPBearer("bearer")(cfg)
		assert.Equal(t, "https://issuer.test/.well-known/oauth-authorization-server", cfg.scheme.OAuth2MetadataURL)
	})

	t.Run("Deprecated", func(t *testing.T) {
		cfg := &securityConfig{}
		SecurityDeprecated()(cfg)
		assert.True(t, cfg.scheme.Deprecated)
		SecurityDeprecated(false)(cfg)
		assert.False(t, cfg.scheme.Deprecated)
	})
}

func TestContentOptions(t *testing.T) {
	cu := &openapi.ContentUnit{}
	opts := []ContentOption{
		ContentType("text/plain"),
		ContentDescription("Text"),
		ContentSummary("Summary"),
		ContentDefault(true),
		ContentEncoding("prop", "enc"),
		ContentExample(map[string]any{"id": "123"}),
		ContentNamedExample("named", map[string]any{"id": "456"}, ExampleSummary("Named")),
		ContentRequired(true),
		ContentFormat("binary"),
	}

	for _, opt := range opts {
		opt(cu)
	}

	assert.Equal(t, "text/plain", cu.ContentType)
	assert.Equal(t, "Text", cu.Description)
	assert.Equal(t, "Summary", cu.Summary)
	assert.True(t, cu.IsDefault)
	assert.Equal(t, "enc", cu.Encoding["prop"])
	assert.Equal(t, map[string]any{"id": "123"}, cu.Example)
	assert.Equal(t, "Named", cu.Examples["named"].Summary)
	assert.Equal(t, map[string]any{"id": "456"}, cu.Examples["named"].Value)
	assert.True(t, cu.Required)
	assert.Equal(t, "binary", cu.Format)

	t.Run("ContentExamples", func(t *testing.T) {
		content := &openapi.ContentUnit{}
		examples := map[string]*openapi.Example{"external": {}}
		ContentExamples(examples)(content)
		assert.Same(t, examples["external"], content.Examples["external"])
	})

	t.Run("ExampleExternalValue", func(t *testing.T) {
		content := &openapi.ContentUnit{}
		ContentNamedExample("external", "ignored", ExampleExternalValue("https://example.com/user.json"))(content)
		assert.Nil(t, content.Examples["external"].Value)
		assert.Equal(t, "https://example.com/user.json", content.Examples["external"].ExternalValue)
	})

	t.Run("ExampleSerializedValue", func(t *testing.T) {
		content := &openapi.ContentUnit{}
		ContentNamedExample("serialized", "ignored", ExampleSerializedValue(`{"id":"123"}`))(content)
		assert.Nil(t, content.Examples["serialized"].Value)
		assert.JSONEq(t, `{"id":"123"}`, content.Examples["serialized"].SerializedValue)
	})

	t.Run("ExampleDescriptionAndDataValue", func(t *testing.T) {
		example := &openapi.Example{Value: "keep"}
		ExampleDescription("Example desc")(example)
		ExampleDataValue(map[string]any{"id": "123"})(example)
		assert.Equal(t, "Example desc", example.Description)
		assert.Nil(t, example.Value)
		assert.Equal(t, map[string]any{"id": "123"}, example.DataValue)
		assert.Empty(t, example.ExternalValue)
	})
}

func TestServerOptions(t *testing.T) {
	s := &openapi.Server{}
	opts := []ServerOption{
		ServerDescription("Desc"),
		ServerVariables(map[string]openapi.ServerVariable{"v": {Default: "d"}}),
	}

	for _, opt := range opts {
		opt(s)
	}

	assert.Equal(t, "Desc", *s.Description)
	assert.Equal(t, "d", s.Variables["v"].Default)
}

func TestConfigAndUIOptions(t *testing.T) {
	t.Run("Basic config setters", func(t *testing.T) {
		cfg := &openapi.Config{}

		WithDisableDocs()(cfg)
		assert.True(t, cfg.DisableDocs)

		WithDisableDocs(false)(cfg)
		assert.False(t, cfg.DisableDocs)

		WithDocsPath("/custom/docs")(cfg)
		assert.Equal(t, "/custom/docs", cfg.DocsPath)

		WithSpecPath("/custom/openapi.yaml")(cfg)
		assert.Equal(t, "/custom/openapi.yaml", cfg.SpecPath)

		WithCacheAge(42)(cfg)
		if assert.NotNil(t, cfg.CacheAge) {
			assert.Equal(t, 42, *cfg.CacheAge)
		}
	})

	t.Run("WithUIOption", func(t *testing.T) {
		cfg := &openapi.Config{}
		called := false
		opt := func(*config.SpecUI) { called = true }

		WithUIOption(opt)(cfg)
		if assert.NotNil(t, cfg.UIOption) {
			cfg.UIOption(&config.SpecUI{})
		}
		assert.True(t, called)
	})

	t.Run("Provider helpers", func(t *testing.T) {
		t.Run("SwaggerUI", func(t *testing.T) {
			cfg := &openapi.Config{}
			WithSwaggerUI(config.SwaggerUI{HideCurl: true})(cfg)
			assert.Equal(t, config.ProviderSwaggerUI, cfg.UIProvider)
			if assert.NotNil(t, cfg.SwaggerUIConfig) {
				assert.True(t, cfg.SwaggerUIConfig.HideCurl)
			}
			assert.NotNil(t, cfg.UIOption)
		})

		t.Run("StoplightElements", func(t *testing.T) {
			cfg := &openapi.Config{}
			WithStoplightElements(config.StoplightElements{
				HideTryIt:              true,
				TryItCredentialsPolicy: config.ElementTryItCredentialsPolicyInclude,
			})(cfg)
			assert.Equal(t, config.ProviderStoplightElements, cfg.UIProvider)
			if assert.NotNil(t, cfg.StoplightElementsConfig) {
				assert.True(t, cfg.StoplightElementsConfig.HideTryIt)
				assert.Equal(
					t,
					config.ElementTryItCredentialsPolicyInclude,
					cfg.StoplightElementsConfig.TryItCredentialsPolicy,
				)
			}
			assert.NotNil(t, cfg.UIOption)
		})

		t.Run("ReDoc", func(t *testing.T) {
			cfg := &openapi.Config{}
			WithReDoc(config.ReDoc{HideDownloadButtons: true})(cfg)
			assert.Equal(t, config.ProviderReDoc, cfg.UIProvider)
			if assert.NotNil(t, cfg.ReDocConfig) {
				assert.True(t, cfg.ReDocConfig.HideDownloadButtons)
			}
			assert.NotNil(t, cfg.UIOption)
		})

		t.Run("Scalar", func(t *testing.T) {
			cfg := &openapi.Config{}
			WithScalar(config.Scalar{DarkMode: true})(cfg)
			assert.Equal(t, config.ProviderScalar, cfg.UIProvider)
			if assert.NotNil(t, cfg.ScalarConfig) {
				assert.True(t, cfg.ScalarConfig.DarkMode)
			}
			assert.NotNil(t, cfg.UIOption)
		})

		t.Run("RapiDoc", func(t *testing.T) {
			cfg := &openapi.Config{}
			WithRapiDoc(config.RapiDoc{HideHeader: true})(cfg)
			assert.Equal(t, config.ProviderRapiDoc, cfg.UIProvider)
			if assert.NotNil(t, cfg.RapiDocConfig) {
				assert.True(t, cfg.RapiDocConfig.HideHeader)
			}
			assert.NotNil(t, cfg.UIOption)
		})
	})
}

func TestOptional(t *testing.T) {
	assert.Equal(t, 1, optional(1))
	assert.Equal(t, 2, optional(1, 2))
	assert.Equal(t, "b", optional("a", "b"))
	assert.False(t, optional(true, false))
}

func TestWithLogger(t *testing.T) {
	t.Run("default logger is not nil", func(t *testing.T) {
		cfg := WithOpenAPIConfig()
		assert.NotNil(t, cfg.Logger)
	})

	t.Run("custom logger is stored", func(t *testing.T) {
		custom := slog.Default()
		cfg := WithOpenAPIConfig(WithLogger(custom))
		assert.Same(t, custom, cfg.Logger)
	})

	t.Run("nil logger falls back to slog.Default", func(t *testing.T) {
		cfg := WithOpenAPIConfig(WithLogger(nil))
		assert.Same(t, slog.Default(), cfg.Logger)
	})
}

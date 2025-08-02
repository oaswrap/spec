package main

import (
	"github.com/oaswrap/spec"
	"github.com/oaswrap/spec/openapi"
	"github.com/oaswrap/spec/option"
)

func createRouter() spec.Generator {
	r := spec.NewRouter(
		option.WithOpenAPIVersion("3.0.3"),
		option.WithTitle("Petstore API"),
		option.WithDescription("This is a sample Petstore server."),
		option.WithVersion("1.0.0"),
		option.WithTermsOfService("https://swagger.io/terms/"),
		option.WithContact(openapi.Contact{
			Email: "apiteam@swagger.io",
		}),
		option.WithLicense(openapi.License{
			Name: "Apache 2.0",
			URL:  "https://www.apache.org/licenses/LICENSE-2.0.html",
		}),
		option.WithExternalDocs("https://swagger.io", "Find more info here about swagger"),
		option.WithServer("https://petstore3.swagger.io/api/v3"),
		option.WithTags(
			openapi.Tag{
				Name:        "pet",
				Description: "Everything about your Pets",
				ExternalDocs: &openapi.ExternalDocs{
					Description: "Find out more about our Pets",
					URL:         "https://swagger.io",
				},
			},
			openapi.Tag{
				Name:        "store",
				Description: "Access to Petstore orders",
				ExternalDocs: &openapi.ExternalDocs{
					Description: "Find out more about our Store",
					URL:         "https://swagger.io",
				},
			},
			openapi.Tag{
				Name:        "user",
				Description: "Operations about user",
			},
		),
		option.WithSecurity("petstore_auth", option.SecurityOAuth2(
			openapi.OAuthFlows{
				Implicit: &openapi.OAuthFlowsImplicit{
					AuthorizationURL: "https://petstore3.swagger.io/oauth/authorize",
					Scopes: map[string]string{
						"write:pets": "modify pets in your account",
						"read:pets":  "read your pets",
					},
				},
			}),
		),
		option.WithSecurity("apiKey", option.SecurityAPIKey("api_key", openapi.SecuritySchemeAPIKeyInHeader)),
	)

	return r
}

package chiopenapi_test

import (
	"flag"
	"os"
	"path/filepath"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/oaswrap/spec/adapter/chiopenapi"
	"github.com/oaswrap/spec/openapi"
	"github.com/oaswrap/spec/option"
	"github.com/oaswrap/spec/pkg/dto"
	"github.com/oaswrap/spec/pkg/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var update = flag.Bool("update", false, "update golden files")

func TestRouter_Spec(t *testing.T) {
	tests := []struct {
		name      string
		golden    string
		opts      []option.OpenAPIOption
		setup     func(r chiopenapi.Router)
		shouldErr bool
	}{
		{
			name:   "Pet Store API",
			golden: "petstore",
			opts: []option.OpenAPIOption{
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
			},
			setup: func(r chiopenapi.Router) {
				r.Route("/pet", func(r chiopenapi.Router) {
					r.Put("/", nil).With(
						option.OperationID("updatePet"),
						option.Summary("Update an existing pet"),
						option.Description("Update the details of an existing pet in the store."),
						option.Request(new(dto.Pet)),
						option.Response(200, new(dto.Pet)),
					)
					r.Post("/", nil).With(
						option.OperationID("addPet"),
						option.Summary("Add a new pet"),
						option.Description("Add a new pet to the store."),
						option.Request(new(dto.Pet)),
						option.Response(201, new(dto.Pet)),
					)
					r.Get("/findByStatus", nil).With(
						option.OperationID("findPetsByStatus"),
						option.Summary("Find pets by status"),
						option.Description("Finds Pets by status. Multiple status values can be provided with comma separated strings."),
						option.Request(new(struct {
							Status string `query:"status" enum:"available,pending,sold"`
						})),
						option.Response(200, new([]dto.Pet)),
					)
					r.Get("/findByTags", nil).With(
						option.OperationID("findPetsByTags"),
						option.Summary("Find pets by tags"),
						option.Description("Finds Pets by tags. Multiple tags can be provided with comma separated strings."),
						option.Request(new(struct {
							Tags []string `query:"tags"`
						})),
						option.Response(200, new([]dto.Pet)),
					)
					r.Post("/{petId}/uploadImage", nil).With(
						option.OperationID("uploadFile"),
						option.Summary("Upload an image for a pet"),
						option.Description("Uploads an image for a pet."),
						option.Request(new(dto.UploadImageRequest)),
						option.Response(200, new(dto.ApiResponse)),
					)
					r.Get("/{petId}", nil).With(
						option.OperationID("getPetById"),
						option.Summary("Get pet by ID"),
						option.Description("Retrieve a pet by its ID."),
						option.Request(new(struct {
							ID int `path:"petId" required:"true"`
						})),
						option.Response(200, new(dto.Pet)),
					)
					r.Post("/{petId}", nil).With(
						option.OperationID("updatePetWithForm"),
						option.Summary("Update pet with form"),
						option.Description("Updates a pet in the store with form data."),
						option.Request(new(dto.UpdatePetWithFormRequest)),
						option.Response(200, nil),
					)
					r.Delete("/{petId}", nil).With(
						option.OperationID("deletePet"),
						option.Summary("Delete a pet"),
						option.Description("Delete a pet from the store by its ID."),
						option.Request(new(dto.DeletePetRequest)),
						option.Response(204, nil),
					)
				}, option.GroupTags("pet"),
					option.GroupSecurity("petstore_auth", "write:pets", "read:pets"),
				)

				r.Route("/store", func(r chiopenapi.Router) {
					r.Post("/order", nil).With(
						option.OperationID("placeOrder"),
						option.Summary("Place an order"),
						option.Description("Place a new order for a pet."),
						option.Request(new(dto.Order)),
						option.Response(201, new(dto.Order)),
					)
					r.Get("/order/{orderId}", nil).With(
						option.OperationID("getOrderById"),
						option.Summary("Get order by ID"),
						option.Description("Retrieve an order by its ID."),
						option.Request(new(struct {
							ID int `path:"orderId" required:"true"`
						})),
						option.Response(200, new(dto.Order)),
						option.Response(404, nil),
					)
					r.Delete("/order/{orderId}", nil).With(
						option.OperationID("deleteOrder"),
						option.Summary("Delete an order"),
						option.Description("Delete an order by its ID."),
						option.Request(new(struct {
							ID int `path:"orderId" required:"true"`
						})),
						option.Response(204, nil),
					)
				}, option.GroupTags("store"))

				r.Route("/user", func(r chiopenapi.Router) {
					r.Post("/createWithList", nil).With(
						option.OperationID("createUsersWithList"),
						option.Summary("Create users with list"),
						option.Description("Create multiple users in the store with a list."),
						option.Request(new([]dto.PetUser)),
						option.Response(201, nil),
					)
					r.Post("/", nil).With(
						option.OperationID("createUser"),
						option.Summary("Create a new user"),
						option.Description("Create a new user in the store."),
						option.Request(new(dto.PetUser)),
						option.Response(201, new(dto.PetUser)),
					)
					r.Get("/{username}", nil).With(
						option.OperationID("getUserByName"),
						option.Summary("Get user by username"),
						option.Description("Retrieve a user by their username."),
						option.Request(new(struct {
							Username string `path:"username" required:"true"`
						})),
						option.Response(200, new(dto.PetUser)),
						option.Response(404, nil),
					)
					r.Put("/{username}", nil).With(
						option.OperationID("updateUser"),
						option.Summary("Update an existing user"),
						option.Description("Update the details of an existing user."),
						option.Request(new(struct {
							Username string `path:"username" required:"true"`
							dto.PetUser
						})),
						option.Response(200, new(dto.PetUser)),
						option.Response(404, nil),
					)
					r.Delete("/{username}", nil).With(
						option.OperationID("deleteUser"),
						option.Summary("Delete a user"),
						option.Description("Delete a user from the store by their username."),
						option.Request(new(struct {
							Username string `path:"username" required:"true"`
						})),
						option.Response(204, nil),
					)
				}, option.GroupTags("user"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := chi.NewRouter()
			opts := []option.OpenAPIOption{
				option.WithOpenAPIVersion("3.0.3"),
				option.WithTitle("Test API " + tt.name),
				option.WithVersion("1.0.0"),
				option.WithDescription("This is a test API for " + tt.name),
				option.WithReflectorConfig(
					option.RequiredPropByValidateTag(),
					option.StripDefNamePrefix("GinopenapiTest"),
				),
			}
			if len(tt.opts) > 0 {
				opts = append(opts, tt.opts...)
			}
			r := chiopenapi.NewRouter(app, opts...)

			if tt.setup != nil {
				tt.setup(r)
			}

			if tt.shouldErr {
				err := r.Validate()
				assert.Error(t, err, "expected error for invalid OpenAPI configuration")
				return
			}
			err := r.Validate()
			assert.NoError(t, err, "failed to validate OpenAPI configuration")

			// Test the OpenAPI schema generation
			schema, err := r.GenerateSchema()

			require.NoError(t, err, "failed to generate OpenAPI schema")
			golden := filepath.Join("testdata", tt.golden+".yaml")

			if *update {
				err = r.WriteSchemaTo(golden)
				require.NoError(t, err, "failed to write golden file")
				t.Logf("Updated golden file: %s", golden)
			}

			want, err := os.ReadFile(golden)
			require.NoError(t, err, "failed to read golden file %s", golden)

			testutil.EqualYAML(t, want, schema)
		})
	}
}

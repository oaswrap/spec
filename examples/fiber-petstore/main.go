package main

import (
	"fmt"
	"log"

	openapiwrapper "github.com/faizlabs/openapi-wrapper"
	"github.com/faizlabs/openapi-wrapper/fiberopenapi"
	"github.com/faizlabs/openapi-wrapper/option"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	r := fiberopenapi.NewRouter(app,
		option.WithTitle("Pet Store API - OpenAPI 3.1"),
		option.WithVersion("1.0.0"),
		option.WithDescription("This is a sample Pet Store API using OpenAPI 3.1"),
		option.WithDocsPath("/docs"),
		option.WithServer("https://petstore3.swagger.io", "Pet Store Server"),
		option.WithSecurity("petstore_auth", option.SecurityOAuth2(
			openapiwrapper.OAuthFlows{
				Implicit: &openapiwrapper.OAuthFlowsDefsImplicit{
					AuthorizationURL: "https://petstore3.swagger.io/oauth/authorize",
					Scopes: map[string]string{
						"write:pets": "modify pets in your account",
						"read:pets":  "read your pets",
					},
				},
			},
		)),
	)

	api := r.Group("/api")
	v3 := api.Group("/v3")
	v3.Route("/pet", func(r fiberopenapi.Router) {
		r.Get("/findByStatus", dummyHandler).With(
			option.Summary("Finds Pets by status."),
			option.Description("Multiple status values can be provided with comma separated strings"),
			option.Request(new(FindPetsByStatusRequest)),
			option.Response(200, new([]Pet)),
		)
		r.Get("/findByTags", dummyHandler).With(
			option.Summary("Finds Pets by tags."),
			option.Description("Multiple tags can be provided with comma separated strings. Use tag1, tag2, tag3 for testing."),
			option.Request(new(FindPetsByTagsRequest)),
			option.Response(200, new([]Pet)),
		)
		r.Get("/:petId", dummyHandler).With(
			option.Summary("Find a pet by ID."),
			option.Description("Returns a single pet."),
			option.Request(new(FindPetByIdRequest)),
			option.Response(200, new(Pet)),
		)
		r.Post("/:petId", dummyHandler).With(
			option.Summary("Updates a pet in the store with form data."),
			option.Description("Update a pet resource based on form data."),
			option.Request(new(UpdatePetFormDataRequest)),
			option.Response(200, new(Pet)),
		)
		r.Delete("/:petId", dummyHandler).With(
			option.Summary("Deletes a pet."),
			option.Request(new(DeletePetRequest)),
		)
		r.Post("/:petId/uploadImage", dummyHandler).With(
			option.Summary("Uploads an image."),
			option.Description("Uploads image of the pet."),
			option.Request(new(UploadImageRequest)),
			option.Response(200, new(ApiResponse)),
		)
		r.Post("/", dummyHandler).With(
			option.Summary("Add a new pet to the store."),
			option.Request(new(Pet)),
			option.Response(200, new(Pet)),
		)
		r.Put("/", dummyHandler).With(
			option.Summary("Update an existing pet."),
			option.Description("Update an existing pet by Id."),
			option.Request(new(Pet)),
			option.Response(200, new(Pet)),
		)
	}).With(option.RouteTags("pet"), option.RouteSecurity("petstore_auth", "write:pets", "read:pets"))

	// Validate the OpenAPI configuration
	if err := r.Validate(); err != nil {
		log.Fatalf("OpenAPI validation failed: %v", err)
	}

	// Write the OpenAPI schema to a file (Optional)
	if err := r.WriteSchemaTo("openapi.yaml"); err != nil {
		log.Fatalf("Failed to write OpenAPI schema: %v", err)
	}
	if err := r.WriteSchemaTo("openapi.json"); err != nil {
		log.Fatalf("Failed to write OpenAPI schema: %v", err)
	}

	fmt.Println("Open http://localhost:3000/docs to view the OpenAPI documentation")

	app.Listen(":3000")
}

func dummyHandler(c *fiber.Ctx) error {
	// Dummy handler for demonstration purposes
	return c.SendString("This is a dummy handler")
}

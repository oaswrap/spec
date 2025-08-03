package main

import (
	"log"

	"github.com/oaswrap/spec/option"
)

func main() {
	// Create a new OpenAPI router
	r := createRouter()

	pet := r.Group("/pet",
		option.GroupTags("pet"),
		option.GroupSecurity("petstore_auth", "write:pets", "read:pets"),
	)
	pet.Put("/",
		option.OperationID("updatePet"),
		option.Summary("Update an existing pet"),
		option.Description("Update the details of an existing pet in the store."),
		option.Request(new(Pet)),
		option.Response(200, new(Pet)),
	)
	pet.Post("/",
		option.OperationID("addPet"),
		option.Summary("Add a new pet"),
		option.Description("Add a new pet to the store."),
		option.Request(new(Pet)),
		option.Response(201, new(Pet)),
	)
	pet.Get("/findByStatus",
		option.OperationID("findPetsByStatus"),
		option.Summary("Find pets by status"),
		option.Description("Finds Pets by status. Multiple status values can be provided with comma separated strings."),
		option.Request(new(struct {
			Status string `query:"status" enum:"available,pending,sold"` // Enum values for pet status
		})),
		option.Response(200, new([]Pet)),
	)
	pet.Get("/findByTags",
		option.OperationID("findPetsByTags"),
		option.Summary("Find pets by tags"),
		option.Description("Finds Pets by tags. Multiple tags can be provided with comma separated strings."),
		option.Request(new(struct {
			Tags []string `query:"tags"` // Tags to filter pets
		})),
		option.Response(200, new([]Pet)),
	)
	pet.Post("/{petId}/uploadImage",
		option.OperationID("uploadFile"),
		option.Summary("Upload an image for a pet"),
		option.Description("Uploads an image for a pet."),
		option.Request(new(UploadImageRequest)),
		option.Response(200, new(ApiResponse)),
	)
	pet.Get("/{petId}",
		option.OperationID("getPetById"),
		option.Summary("Get pet by ID"),
		option.Description("Retrieve a pet by its ID."),
		option.Request(new(struct {
			ID int `path:"petId" required:"true"`
		})),
		option.Response(200, new(Pet)),
	)
	pet.Post("/{petId}",
		option.OperationID("updatePetWithForm"),
		option.Summary("Update pet with form"),
		option.Description("Updates a pet in the store with form data."),
		option.Request(new(UpdatePetWithFormRequest)),
		option.Response(200, nil),
	)
	pet.Delete("/{petId}",
		option.OperationID("deletePet"),
		option.Summary("Delete a pet"),
		option.Description("Delete a pet from the store by its ID."),
		option.Request(new(DeletePetRequest)),
		option.Response(204, nil),
	)

	store := r.Group("/store",
		option.GroupTags("store"),
	)
	store.Post("/order",
		option.OperationID("placeOrder"),
		option.Summary("Place an order"),
		option.Description("Place a new order for a pet."),
		option.Request(new(Order)),
		option.Response(201, new(Order)),
	)
	store.Get("/order/{orderId}",
		option.OperationID("getOrderById"),
		option.Summary("Get order by ID"),
		option.Description("Retrieve an order by its ID."),
		option.Request(new(struct {
			ID int `path:"orderId" required:"true"`
		})),
		option.Response(200, new(Order)),
		option.Response(404, nil),
	)
	store.Delete("/order/{orderId}",
		option.OperationID("deleteOrder"),
		option.Summary("Delete an order"),
		option.Description("Delete an order by its ID."),
		option.Request(new(struct {
			ID int `path:"orderId" required:"true"`
		})),
		option.Response(204, nil),
	)

	user := r.Group("/user",
		option.GroupTags("user"),
	)
	user.Post("/createWithList",
		option.OperationID("createUsersWithList"),
		option.Summary("Create users with list"),
		option.Description("Create multiple users in the store with a list."),
		option.Request(new([]User)),
		option.Response(201, nil),
	)
	user.Post("/",
		option.OperationID("createUser"),
		option.Summary("Create a new user"),
		option.Description("Create a new user in the store."),
		option.Request(new(User)),
		option.Response(201, new(User)),
	)
	user.Get("/{username}",
		option.OperationID("getUserByName"),
		option.Summary("Get user by username"),
		option.Description("Retrieve a user by their username."),
		option.Request(new(struct {
			Username string `path:"username" required:"true"`
		})),
		option.Response(200, new(User)),
		option.Response(404, nil),
	)
	user.Put("/{username}",
		option.OperationID("updateUser"),
		option.Summary("Update an existing user"),
		option.Description("Update the details of an existing user."),
		option.Request(new(struct {
			Username string `path:"username" required:"true"`
			User
		})),
		option.Response(200, new(User)),
		option.Response(404, nil),
	)
	user.Delete("/{username}",
		option.OperationID("deleteUser"),
		option.Summary("Delete a user"),
		option.Description("Delete a user from the store by their username."),
		option.Request(new(struct {
			Username string `path:"username" required:"true"`
		})),
		option.Response(204, nil),
	)

	// Generate OpenAPI spec
	if err := r.WriteSchemaTo("openapi.yaml"); err != nil {
		log.Fatal(err)
	}

	log.Println("✅ OpenAPI spec generated at openapi.yaml")
}

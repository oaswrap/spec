module github.com/oaswrap/spec/adapter/httprouteropenapi/examples/basic

go 1.21

require (
	github.com/julienschmidt/httprouter v1.3.0
	github.com/oaswrap/spec v0.3.4
	github.com/oaswrap/spec/adapter/httprouteropenapi v0.0.0
)

require (
	github.com/oaswrap/spec-ui v0.1.4 // indirect
	github.com/swaggest/jsonschema-go v0.3.78 // indirect
	github.com/swaggest/openapi-go v0.2.60 // indirect
	github.com/swaggest/refl v1.4.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)

replace github.com/oaswrap/spec/adapter/httprouteropenapi => ../..

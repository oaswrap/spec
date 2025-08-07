module github.com/oaswrap/spec/adapter/httpopenapi/examples/basic

go 1.22

require (
	github.com/oaswrap/spec v0.2.4
	github.com/oaswrap/spec/adapter/httpopenapi v0.0.0
)

require (
	github.com/swaggest/jsonschema-go v0.3.78 // indirect
	github.com/swaggest/openapi-go v0.2.59 // indirect
	github.com/swaggest/refl v1.4.0 // indirect
	github.com/swaggest/swgui v1.8.4 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)

replace github.com/oaswrap/spec/adapter/httpopenapi => ../..

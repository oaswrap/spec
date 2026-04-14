module github.com/oaswrap/spec/adapter/echov5openapi/example

go 1.25.0

require (
	github.com/labstack/echo/v5 v5.1.0
	github.com/oaswrap/spec v0.4.1
	github.com/oaswrap/spec/adapter/echov5openapi v0.0.0
)

require (
	github.com/oaswrap/spec-ui v0.2.0 // indirect
	github.com/swaggest/jsonschema-go v0.3.79 // indirect
	github.com/swaggest/openapi-go v0.2.61 // indirect
	github.com/swaggest/refl v1.4.0 // indirect
	golang.org/x/text v0.34.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)

replace github.com/oaswrap/spec/adapter/echov5openapi => ..

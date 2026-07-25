module github.com/oaswrap/spec/adapter/echov5openapi/example

go 1.25.0

require (
	github.com/labstack/echo/v5 v5.3.1
	github.com/oaswrap/spec v0.5.1
	github.com/oaswrap/spec/adapter/echov5openapi v0.0.0
)

require (
	github.com/goccy/go-yaml v1.19.2 // indirect
	github.com/oaswrap/spec-ui v0.2.1 // indirect
)

replace github.com/oaswrap/spec => ../../..

replace github.com/oaswrap/spec/adapter/echov5openapi => ..

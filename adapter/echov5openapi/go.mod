module github.com/oaswrap/spec/adapter/echov5openapi

go 1.25.0

require (
	github.com/labstack/echo/v5 v5.3.1
	github.com/oaswrap/spec v0.5.1
	github.com/oaswrap/spec-ui v0.2.0
	github.com/stretchr/testify v1.11.1
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/goccy/go-yaml v1.19.2 // indirect
	github.com/google/go-cmp v0.7.0 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/oaswrap/spec => ../..

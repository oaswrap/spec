module github.com/oaswrap/spec/adapter/muxopenapi

go 1.22

require (
	github.com/google/go-cmp v0.7.0
	github.com/gorilla/mux v1.8.1
	github.com/oaswrap/spec v0.5.2
	github.com/oaswrap/spec-ui v0.2.1
	github.com/stretchr/testify v1.11.1
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/goccy/go-yaml v1.19.2 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/oaswrap/spec => ../..

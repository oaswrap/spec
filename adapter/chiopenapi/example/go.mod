module github.com/oaswrap/spec/adapter/chiopenapi/example

go 1.24

require (
	github.com/go-chi/chi/v5 v5.5.0
	github.com/oaswrap/spec v0.5.2
	github.com/oaswrap/spec/adapter/chiopenapi v0.0.0
)

require (
	github.com/goccy/go-yaml v1.19.2 // indirect
	github.com/oaswrap/spec-ui v0.2.2 // indirect
)

replace github.com/oaswrap/spec => ../../..

replace github.com/oaswrap/spec/adapter/chiopenapi => ..

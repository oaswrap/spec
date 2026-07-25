module github.com/oaswrap/spec/adapter/muxopenapi/example

go 1.23

require (
	github.com/gorilla/mux v1.8.1
	github.com/oaswrap/spec v0.5.2
	github.com/oaswrap/spec/adapter/muxopenapi v0.3.1
)

require (
	github.com/goccy/go-yaml v1.19.2 // indirect
	github.com/oaswrap/spec-ui v0.2.1 // indirect
)

replace github.com/oaswrap/spec => ../../..

replace github.com/oaswrap/spec/adapter/muxopenapi => ..

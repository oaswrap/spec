module github.com/oaswrap/spec/adapter/httpopenapi/example

go 1.23

require (
	github.com/oaswrap/spec v0.5.1
	github.com/oaswrap/spec/adapter/httpopenapi v0.0.0
)

require (
	github.com/goccy/go-yaml v1.19.2 // indirect
	github.com/oaswrap/spec-ui v0.2.0 // indirect
)

replace github.com/oaswrap/spec => ../../..

replace github.com/oaswrap/spec/adapter/httpopenapi => ..

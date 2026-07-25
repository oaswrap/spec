module github.com/oaswrap/spec/adapter/echoopenapi/example

go 1.25.0

require (
	github.com/labstack/echo/v4 v4.15.4
	github.com/oaswrap/spec v0.5.2
	github.com/oaswrap/spec/adapter/echoopenapi v0.0.0
)

require (
	github.com/goccy/go-yaml v1.19.2 // indirect
	github.com/labstack/gommon v0.5.0 // indirect
	github.com/mattn/go-colorable v0.1.15 // indirect
	github.com/mattn/go-isatty v0.0.24 // indirect
	github.com/oaswrap/spec-ui v0.2.1 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasttemplate v1.2.2 // indirect
	golang.org/x/crypto v0.54.0 // indirect
	golang.org/x/net v0.57.0 // indirect
	golang.org/x/sys v0.47.0 // indirect
	golang.org/x/text v0.40.0 // indirect
)

replace github.com/oaswrap/spec => ../../..

replace github.com/oaswrap/spec/adapter/echoopenapi => ..

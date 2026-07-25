module github.com/oaswrap/spec/adapter/fiberopenapi/example

go 1.25.0

require (
	github.com/gofiber/fiber/v2 v2.52.14
	github.com/oaswrap/spec v0.5.1
	github.com/oaswrap/spec/adapter/fiberopenapi v0.0.0
)

require (
	github.com/andybalholm/brotli v1.2.2 // indirect
	github.com/clipperhouse/uax29/v2 v2.7.0 // indirect
	github.com/goccy/go-yaml v1.19.2 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/klauspost/compress v1.19.1 // indirect
	github.com/mattn/go-colorable v0.1.15 // indirect
	github.com/mattn/go-isatty v0.0.24 // indirect
	github.com/mattn/go-runewidth v0.0.27 // indirect
	github.com/oaswrap/spec-ui v0.2.1 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasthttp v1.72.0 // indirect
	golang.org/x/sys v0.47.0 // indirect
)

replace github.com/oaswrap/spec => ../../..

replace github.com/oaswrap/spec/adapter/fiberopenapi => ..

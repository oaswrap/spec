# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Changed
- Updated `spec-ui` to `v0.2.2` (Swagger UI 5.34.1, Stoplight Elements 9.0.29, Scalar 1.76.4, ReDoc 3.0.1, RapiDoc 12.0.0).

### Added
- `config.StoplightElements.TryItCredentialsPolicy` (`omit`, `include`, `same-origin`) is available through `option.WithStoplightElements`.

### Removed
- **Breaking (UI config):** `config.ReDoc.HideSearch` (ReDoc 3 always enables search) and `config.RapiDoc.HeaderColor` (no equivalent in RapiDoc 11+) were removed upstream; delete them from `option.WithReDoc` / `option.WithRapiDoc` calls.

### Fixed
- `multipart.FileHeader`, `multipart.File`, and slices of them are reflected as `type: string, format: binary` instead of an object schema (or a `$ref` component), and are never nullable.
- Petstore example `uploadFile` operation now emits the `application/octet-stream` binary request body; the unsupported `_ *multipart.File \`contentType:"..."\`` struct field never affected the generated spec. Declare raw binary bodies with `option.Request(nil, option.ContentType(...), option.ContentFormat("binary"))`.

### Documentation
- Added a "File uploads" section to the README covering raw binary and `multipart/form-data` bodies.

## [0.5.2] - 2026-07-25

### Fixed
- Path parameters with a regex constraint (e.g. `/files/{name:[a-z]+}`) are parsed correctly, including escaped and nested curly braces in the constraint; the parameter name is taken from the part before the first `:`.

### Changed
- Updated `spec-ui` to `v0.2.1`.
- Updated adapter framework dependencies.

### Documentation
- Removed the Go Report Card badge from the READMEs.

## [0.5.1] - 2026-05-12

### Added
- Logger debug feature for tracing spec generation behavior.
- `InterceptSchema` and `InterceptProp` reflector hooks for customizing schema and property generation.
- `RequiredPropByValidateTag` hook for marking fields required based on `validate` struct tags.
- `ParentType` field on `InterceptPropParams` for richer hook context.
- Automatic content-type detection from struct type when not explicitly set.
- `encoding.TextMarshaler`/`TextUnmarshaler` support: types implementing both interfaces (without `json.Marshaler`) are reflected as `type: string`.
- `EmbedReferencer` interface and `refer:"true"` struct tag: embedded structs opt into `allOf $ref` instead of field inlining.
- `ParameterTagMapping` now accepts `openapi.ParameterInBody` and `openapi.ParameterInForm` to override the struct tag used for JSON and form request body field names (defaults: `json` and `form`).

### Changed
- Schema component names now always include the Go package name as a prefix (e.g., `models.User` → `ModelsUser`). This eliminates cross-package naming collisions without requiring caller-package detection. Use `InterceptDefName` or `StripDefNamePrefix` to remove the prefix if desired.
- `DefNameCallerPkg` field removed from `ReflectorConfig`; the caller-package detection mechanism in `NewGenerator` is removed.
- Package name sanitization handles multi-segment names (e.g., `spec_test` → `SpecTest`); unexported type names are title-cased when a package prefix is prepended.

### Fixed
- `uint8`/`uint16` reflected as `int32` format; `uint`/`uint32`/`uint64`/`uintptr` reflected as `int64` format.
- `InterceptSchema`/`InterceptProp` hook error propagation and correctness.
- `RefSchema` pre-hook: assign `StructSchema` fields onto existing pointer so pre-hook customizations (extensions, description) survive to the post-hook.
- `ApplyNullable` (OAS 3.1+): merge `"null"` into an existing `[]string` type slice instead of silently skipping.
- `Required` field deduplication to prevent duplicates when both `required` struct tag and `RequiredPropByValidateTag` are used simultaneously.
- `parentSnapshot` restore on `ErrSkipProperty` in post-hooks to roll back parent schema mutations (`AllOf`, `AnyOf`, `OneOf`, extensions).
- Lint violations and YAML tag reading in `MarshalYAML`.

## [0.5.0] - 2026-05-11

### Added
- OpenAPI 3.2.0 support and expanded version-specific validation checks.
- Validation severity levels (`Error`, `Warning`, `Info`) with best-practice reporting.
- `ValidateReport` support across the core router and all adapters.
- Dirty-state tracking for incremental build behavior.
- Automatic injection of missing path parameters.
- Support for `application/x-www-form-urlencoded` querystring parameter content.
- Media type tag support in schema reflection.
- New Iris adapter: `irisopenapi`.

### Changed
- Rebuilt core generator, reflector, and validator internals.
- Aligned all adapters with the updated core behavior.
- Upgraded OpenAPI handling to improve version consistency and validation behavior (including 3.1.2 alignment).
- Prefixed package name in default schema definition names.
- Centralized operation method validation logic.

### Fixed
- Hardened path and HTTP method generation behavior.

### Removed
- Removed `httprouter` from the supported adapter list.
- Removed dependency on external OpenAPI generator packages (including `github.com/swaggest/openapi-go`) as part of the core generator rebuild.

## [0.4.2] - 2026-04-22

### Fixed
- Shortened generic type definition names in generated JSON schema.
- Resolved Echo v5 adapter test breakage caused by upstream library changes.

### Changed
- Updated `swaggest/openapi-go` to `v0.2.61`.

## [0.4.1] - 2026-04-12

### Added
- Config option to strip trailing slashes from paths.

### Fixed
- URL paths with trailing slashes now generate correctly.

### Changed
- Updated Fiber v2 and Echo v5 dependencies (core and adapters).

## [0.4.0] - 2026-04-05

### Added
- Echo v5 adapter (`echov5openapi`).
- Fiber v3 adapter (`fiberv3openapi`).
- Custom `spec-ui` option support, including `spec-ui` upgrade to `v0.2.0`.

### Changed
- Merged duplicate status-code responses into a single `oneOf` response schema.

## [0.3.6] - 2025-11-05

### Added
- Custom parameter tag configuration support.

### Fixed
- Echo route syntax mapping now correctly supports `:id` style parameters.

## [0.3.5] - 2025-10-25

### Added
- Ability to configure the OpenAPI `encoding` object for request bodies.
- Added `RELEASE.md` release workflow documentation.

### Changed
- Updated adapter dependency versions.

[Unreleased]: https://github.com/oaswrap/spec/compare/v0.5.2...HEAD
[0.5.2]: https://github.com/oaswrap/spec/compare/v0.5.1...v0.5.2
[0.5.1]: https://github.com/oaswrap/spec/compare/v0.5.0...v0.5.1
[0.5.0]: https://github.com/oaswrap/spec/compare/v0.4.2...v0.5.0
[0.4.2]: https://github.com/oaswrap/spec/compare/v0.4.1...v0.4.2
[0.4.1]: https://github.com/oaswrap/spec/compare/v0.4.0...v0.4.1
[0.4.0]: https://github.com/oaswrap/spec/compare/v0.3.6...v0.4.0
[0.3.6]: https://github.com/oaswrap/spec/compare/v0.3.5...v0.3.6
[0.3.5]: https://github.com/oaswrap/spec/compare/v0.3.4...v0.3.5

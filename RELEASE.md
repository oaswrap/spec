# ✅ Release Process

## 🗂️  1. Sync adapter dependencies

Run this to bump each adapter’s `go.mod` to the new version:

```
make sync-adapter-deps VERSION=v0.3.0 NO_TIDY=1
```

> 🔍 The `NO_TIDY=1` skips `go mod tidy` for speed.  
> You will tidy after pushing tags.

## ✅ 2. Commit changes

Stage and commit the updated adapter files:

```
git add adapters/*/go.mod adapters/*/go.sum
git commit -m "chore(adapters): bump spec version to v0.3.0"
```

## 🚀 3. Run the release

Create & push the tags:

```
make release VERSION=v0.3.0
```

## 🚀 4. For dev release

Same steps but with a dev version:

```
make sync-adapter-deps VERSION=v0.3.0-dev.1 NO_TIDY=1
git add adapters/*/go.mod adapters/*/go.sum
git commit -m "chore(adapters): bump spec version to v0.3.0-dev.1"
make release-dev VERSION=v0.3.0-dev.1
```

## 🔑 Notes

- `sync-adapter-deps` **does NOT commit** — you must commit before tagging.
- After tags are pushed, `make tidy` runs to clean up.

# 📦 Releasing & Version Management

This project uses a **multi-module mono-repo** with a main module (`spec`) and multiple adapters.

## 🔑 How versioning works

- **Main module (`spec`)** has its own semantic version tag: `v0.x.y`  
- **Adapters** depend on the main module by exact version (`github.com/oaswrap/spec v0.x.y`).
- To keep everything aligned, adapters must update their `go.mod` when `spec` is bumped.

## 🚀 Release workflow

### ✅ 1️⃣ Bump to next dev version (before releasing)

If you are preparing a **new development version**, bump all adapters first:

```bash
make bump-dev NEXT=v0.3.0-dev.1 NO_TIDY=1
git commit -am "chore: bump dev version"
```

> `NO_TIDY=1` skips `go mod tidy` because the tag doesn’t exist yet — tidy will run after pushing.

### ✅ 2️⃣ Create and push dev release

```bash
make release-dev VERSION=v0.3.0-dev.1
```

This:
- Runs final checks
- Creates git tag for `spec` and all adapters
- Pushes all tags
- Runs `go mod tidy` to update `go.sum`

### ✅ 3️⃣ Create and push stable release

When ready for production:

```bash
make release VERSION=v0.3.0
```

Same steps:
- Final checks
- Tags `spec` and all adapters
- Pushes tags
- Runs `tidy` to finalize `go.sum`

### ✅ 4️⃣ Sync adapters to the released version

After pushing the stable tag, you may **re-sync** all adapters:

```bash
make sync-adapter-deps VERSION=v0.3.0
git commit -am "chore: sync adapters to v0.3.0"
```

## ⚠️ Good practice

- Always run `go mod tidy` **after** pushing new tags.
- CI will fail if `go.sum` or `replace` directives are stale.
- Use `NO_TIDY=1` only when bumping to a **version that doesn’t exist yet** — tidy will run after the release push.

## ✅ Commands recap

| Command                     | Use case                                |
|-----------------------------|-----------------------------------------|
| `make bump-dev NEXT=...`    | Prepare adapters for next dev version   |
| `make release-dev VERSION=...` | Tag & push dev version, tidy after push |
| `make release VERSION=...`  | Tag & push stable version, tidy after push |
| `make sync-adapter-deps VERSION=...` | Sync adapters to a released version |

📌 **Keeping all adapters aligned = no broken builds.**  
Use this flow → keep your mono-repo healthy. 🔒✅
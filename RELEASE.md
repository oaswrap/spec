
# 📦 Releasing & Version Management

This project uses a **multi-module mono-repo** with a main module (`spec`) and multiple adapters.

## 🔑 How versioning works

- **Main module (`spec`)** has its own semantic version tag: `v0.x.y`  
- **Adapters** depend on the main module by exact version (`github.com/oaswrap/spec v0.x.y`).
- When `spec` is bumped, adapters must update their `go.mod` to match.

---

## ✅ How release tagging works

Your `Makefile` does:
- **`make release`** → creates a version tag for `spec` **and** tags each adapter **as-is**.
- It does **not** rewrite adapter `go.mod` — that’s up to you.

**So:** After you tag, you should run a sync to update adapters to the new `spec` version.

This means:
- **The tag freezes the code state at that point**.
- **Syncing adapters after** makes the next patch or minor bump correct.

This is normal and matches how major Go mono-repos (e.g., Kubernetes) handle internal module deps.

---

## ✅ Recommended release flow

### 🟢 1️⃣ Bump to next dev version

When starting new work:

```bash
make bump-dev NEXT=v0.3.0-dev.1 NO_TIDY=1
git commit -am "chore: bump dev version"
```

**NO_TIDY=1** skips `go mod tidy` (the new tag doesn’t exist yet).

---

### 🟢 2️⃣ Create and push dev release

```bash
make release-dev VERSION=v0.3.0-dev.1
```

- Tags `spec` and all adapters
- Pushes the tags
- Runs `tidy` after the tags exist

---

### 🟢 3️⃣ Develop, test, merge as usual

Keep merging PRs on the `-dev` version.

---

### 🟢 4️⃣ When ready, release stable

```bash
make release VERSION=v0.3.0
```

- Tags `spec` and all adapters
- Pushes the tags
- Runs `tidy` after tagging

**Adapters will be tagged as they are — so they may still point to the old `spec` version in `go.mod`.**

---

### 🟢 5️⃣ Immediately sync adapters (best practice)

After stable tag is pushed:

```bash
make sync-adapter-deps VERSION=v0.3.0
git commit -am "chore: sync adapters to v0.3.0"
git push
```

This updates each adapter’s `go.mod` to match the new stable version.  
This keeps your next patch or minor version aligned.

---

## ✅ Key rule

**Never `sync-adapter-deps` to a version that does not exist yet.**  
Always tag first → then sync.

---

## ⚡ Final checklist

| Command | Use for |
|----------------------------|-------------------------|
| `make bump-dev NEXT=...` | Prepare next dev version |
| `make release-dev VERSION=...` | Tag & push dev version |
| `make release VERSION=...` | Tag & push stable version |
| `make sync-adapter-deps VERSION=...` | Update adapters to use the stable version |
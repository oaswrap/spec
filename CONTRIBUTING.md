# Contributing to oaswrap/spec

🎉 Thanks for your interest in contributing! We welcome pull requests and feedback to improve the project.

## 📦 **Project Structure**

This project uses a **modular Go workspace**:

* `core`: Main package (`./...`)
* `adapters/`: Pluggable framework adapters (`fiberopenapi`, `ginopenapi`, `echoopenapi`)

## 🛠️ **Development Requirements**

Before you start:

* Go **1.23+** (`go version`)
* [gotestsum](https://github.com/gotestyourself/gotestsum) and [golangci-lint](https://golangci-lint.run/) (`make install-tools`)

## 🧩 **Development Workflow**

1️⃣ **Fork & Clone**

```bash
git clone https://github.com/YOUR-USERNAME/spec.git
cd spec
```

2️⃣ **Install tools**

```bash
make install-tools
```

3️⃣ **Sync workspace**

```bash
make sync tidy
```

4️⃣ **Run tests**

```bash
make test          # Core + Adapters
make test-parallel # Run in parallel
```

5️⃣ **Lint**

```bash
make lint
```

## 🧹 **Best Practices**

✅ Use `go work` for local development.

✅ If you modify adapters, keep `go.mod` clean — **don’t commit local `replace`**.
Run:

```bash
make fix-replace
```

✅ Before pushing:

```bash
make check
```

✅ Before creating a release:

```bash
make check-release
```

## 🚀 **Creating a Dev Version**

If you need to publish a dev version for testing:

```bash
make bump-dev NEXT=vX.Y.Z-dev.N
make release-dev VERSION=vX.Y.Z-dev.N
```

This:

* Updates adapters’ `go.mod` to the new version.
* Removes local `replace`.
* Tags & pushes.

## ✅ **Creating a Release**

1️⃣ Ensure your branch is up to date:

```bash
git pull origin main
```

2️⃣ Run checks:

```bash
make release-check
```

3️⃣ Tag and push:

```bash
make release VERSION=vX.Y.Z
```

GitHub Actions will:

* Validate your tag.
* Run tests + lint + replace cleanup.
* Publish a GitHub Release with changelog.

## 📃 **Commit Style**

* Use clear commit messages (`feat:`, `fix:`, `chore:`)
* Keep PRs focused & small if possible.

## 🙏 **Thank You!**

Your contributions make this project better.
Questions? Open an issue or discussion!

Happy coding! 🚀
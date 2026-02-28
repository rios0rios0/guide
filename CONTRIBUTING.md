# Contributing

Contributions are welcome. By participating, you agree to maintain a respectful and constructive environment.

For coding standards, testing patterns, architecture guidelines, commit conventions, and all
development practices, refer to the **[Development Guide](https://github.com/rios0rios0/guide/wiki)**.

## Prerequisites

- [Git](https://git-scm.com/downloads) 2.30+
- [Go](https://go.dev/dl/) 1.26+ (for building and testing the wiki update tool)
- A Markdown editor or IDE (e.g., [VS Code](https://code.visualstudio.com/) with the Markdown preview extension)

## Development Workflow

1. Fork and clone the repository
2. Create a branch: `git checkout -b feat/my-change`
3. Edit the Markdown files (`.md`) following the existing wiki structure and naming conventions
4. Build the wiki update tool:
   ```bash
   cd .github/workflows/update-wiki
   go build -o update-wiki ./...
   ```
5. Run the wiki update tool tests:
   ```bash
   cd .github/workflows/update-wiki
   go test ./...
   ```
6. Preview your Markdown changes locally (e.g., using VS Code Markdown preview or a local Markdown renderer)
7. Commit following the [commit conventions](https://github.com/rios0rios0/guide/wiki/Life-Cycle/Git-Flow)
8. Open a pull request against `main`

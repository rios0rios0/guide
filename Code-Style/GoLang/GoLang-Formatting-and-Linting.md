# Go Formatting and Linting

> **TL;DR:** Use **gofmt** (built-in) for code formatting, **goimports** for import ordering, and **golangci-lint** as the linter aggregator. These tools are non-negotiable and must be integrated into every project's CI pipeline.

## Overview

Go's toolchain includes a built-in formatter (`gofmt`) that eliminates all debates about code style. Combined with `goimports` for import management and `golangci-lint` for static analysis, this toolchain ensures consistent, high-quality code across all projects.

## Formatter: gofmt

`gofmt` is Go's official code formatter. It ships with the Go toolchain and produces a single canonical formatting for any Go source file. There are no configuration options -- this is by design.

```bash
# Format a single file
gofmt -w main.go

# Format all files in the current module
gofmt -w .
```

**Do not use alternative formatters.** `gofmt` is the universal standard in the Go ecosystem, and all Go code must be formatted with it.

## Import Ordering: goimports

[goimports](https://pkg.go.dev/golang.org/x/tools/cmd/goimports) extends `gofmt` by automatically managing import statements -- adding missing imports, removing unused ones, and grouping them into sections:

1. Standard library
2. Third-party packages
3. Application packages

```bash
# Install
go install golang.org/x/tools/cmd/goimports@latest

# Run
goimports -w .
```

Example output:

```go
import (
	"context"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	"myapp/domain/commands"
	"myapp/infrastructure/controllers"
)
```

## Linter: golangci-lint

[golangci-lint](https://golangci-lint.run/) aggregates dozens of Go linters into a single tool. It is fast, configurable, and must be used in all projects.

### Installation

```bash
# Binary installation (recommended for CI)
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin

# Or via go install
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

### Usage

```bash
# Run all enabled linters
golangci-lint run

# Run on specific directories
golangci-lint run ./domain/... ./infrastructure/...
```

### Configuration

Place a `.golangci.yml` file in the project root. A recommended baseline configuration:

```yaml
run:
  timeout: 5m

linters:
  enable:
    - errcheck
    - govet
    - staticcheck
    - unused
    - gosimple
    - ineffassign
    - typecheck
    - misspell
    - gocyclo
    - revive
    - gocritic
    - nakedret
    - prealloc

linters-settings:
  gocyclo:
    min-complexity: 15
  revive:
    rules:
      - name: unexported-return
        disabled: true

issues:
  exclude-use-default: false
```

## Editor Configuration

### Visual Studio Code

Install the [Go extension](https://marketplace.visualstudio.com/items?itemName=golang.go) and add the following settings:

```json
{
    "go.formatTool": "goimports",
    "go.lintTool": "golangci-lint",
    "go.lintOnSave": "package",
    "[go]": {
        "editor.formatOnSave": true,
        "editor.codeActionsOnSave": {
            "source.organizeImports": "explicit"
        }
    }
}
```

### NeoVim

Use [nvim-lspconfig](https://github.com/neovim/nvim-lspconfig) with `gopls` and configure format-on-save:

```lua
require("lspconfig").gopls.setup({
    settings = {
        gopls = {
            gofumpt = true,
            analyses = {
                unusedparams = true,
                shadow = true,
            },
            staticcheck = true,
        },
    },
})
```

## References

- [gofmt Documentation](https://pkg.go.dev/cmd/gofmt)
- [goimports](https://pkg.go.dev/golang.org/x/tools/cmd/goimports)
- [golangci-lint Documentation](https://golangci-lint.run/)
- [Go Editor Plugins](https://go.dev/wiki/IDEsAndTextEditorPlugins)

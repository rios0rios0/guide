# Personal engineering standards

Before project changes, read documentation and git-flow; before code changes also read architecture, code-style, security, testing, and the relevant language. Read ci-cd before validation or Git operations. Resolve the paths below relative to this AGENTS.md file.

New Go projects use manual constructor injection in container.go. Existing Wire migrations belong in dedicated service PRs; encourage Dig migrations. Preserve project-specific loggers. Unit tests are untagged; mocking libraries require the narrow documented external-abstraction exception in testing.md.

Invoke Codex skills with $skill-name. Claude slash commands and CLAUDE.md are tool-specific; Codex uses AGENTS.md. Project instructions refine these defaults. Preserve existing authorization safeguards.

| Applies to | Full standard |
| --- | --- |
| **/*.go | [golang](instructions/golang.md) |
| **/*.py | [python](instructions/python.md) |
| **/*.java | [java](instructions/java.md) |
| **/*.{js,jsx,ts,tsx} | [javascript](instructions/javascript.md) |
| **/*.{yml,yaml} | [yaml](instructions/yaml.md) |
| code-style | [code-style](instructions/code-style.md) |
| git-flow | [git-flow](instructions/git-flow.md) |
| testing | [testing](instructions/testing.md) |
| architecture | [architecture](instructions/architecture.md) |
| security | [security](instructions/security.md) |
| ci-cd | [ci-cd](instructions/ci-cd.md) |
| documentation | [documentation](instructions/documentation.md) |
| **/*.md | [markdown-formatting](instructions/markdown-formatting.md) |
| design-patterns | [design-patterns](instructions/design-patterns.md) |
| bulk-operations | [bulk-operations](instructions/bulk-operations.md) |

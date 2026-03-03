# Cursorrules Template

> **TL;DR:** Copy the template below into a `.cursorrules` file at the project root. This file gives [Cursor](https://cursor.sh/) the full team development standards so it can write code, review PRs, and enforce conventions automatically -- no manual intervention needed.

## Overview

Cursor reads `.cursorrules` from the project root to understand project conventions. This template embeds the entire [Development Guide](https://github.com/rios0rios0/guide/wiki) so that Cursor already knows every convention, architecture pattern, and language-specific rule the team follows. Install it via `install-ai-rules.sh` or copy the content below.

**Note:** Cursor also supports `.cursor/rules/*.mdc` for modular rules. The `.cursorrules` file is the simplest and most compatible approach.

## Template

````text
# Commands

- `make lint` -- fix all linter issues
- `make test` -- run the full test suite
- `make sast` -- run the complete SAST security suite

Always run lint, test, and sast before considering any task complete.

# Architecture

Follow Hexagonal Architecture (Ports & Adapters) with Domain-Driven Design (DDD) and CQRS.

## Layer Separation

- `domain/` contains contracts: entities, repository interfaces, service interfaces, and commands.
- `infrastructure/` contains implementations: controllers, repository implementations, service implementations, and mappers.
- Dependencies always point inward -- infrastructure depends on domain, never the reverse.
- Entities must be framework-agnostic. No ORM annotations, no framework tags, no external dependencies inside entity definitions.

## Principal Actors

- Entities: core business objects with domain rules. Framework-agnostic. No persistence annotations.
- Controllers: bridge between view (client) and business layer. Delegate to commands.
- Commands: implement business/feature logic. Apply domain rules. Define Listeners for callback outcomes.
- Services: handle application-level concerns (parsing, conversions). Domain defines contract, infrastructure implements.
- Repositories: abstract all data access behind domain-defined interfaces.

## Mappers

Use mappers to isolate layers. Never let framework types leak across layer boundaries:
- Repository mappers: convert between infrastructure models (prefixed with tool name, e.g., PgxUser, JpaItem) and domain entities.
- Controller mappers: convert between request/response DTOs and domain entities. Separate request mappers from response mappers.

## File Structure Tiers

Minimal: domain/commands, domain/entities, domain/repositories, infrastructure/repositories.
Intermediate: add domain/services and infrastructure/services.
Complete: add infrastructure/controllers (with requests/, responses/, mappers/), infrastructure/repositories/mappers and models/.

## Frontend Architecture

5-layer lightweight Clean Architecture: Domain, Service, Infrastructure, Presentation, Main.
Dependencies always point toward Domain. Main handles DI wiring.

# Operations Vocabulary

Use these standard prefixes for naming files, classes, and methods:

- List: retrieve multiple records
- Get: retrieve a single record
- Insert: create a single record
- Update: modify a single record
- Delete: remove a single record
- BatchInsert: create multiple records
- BatchUpdate: modify multiple records
- BatchDelete: remove multiple records
- DeleteAll: remove all records

# Git Conventions

## Commit Messages

Format: `type(SCOPE): message`

- SCOPE is the task ID (e.g., TICKET-123) or a short descriptive scope.
- Use simple past tense: added, changed, fixed, removed.
- No period at the end. Do not capitalize the first letter.
- Wrap code references in backticks.
- Types: feat, fix, refactor, chore, test, docs.
- Long messages: body with bullet points starting with lowercase verb in past tense.

## Branches

Format: `type/TICKET-ID` or `type/short-description` (e.g., feat/TICKET-000, fix/input-mask).

## Branch Synchronization

Always use `git rebase` to synchronize with main. Never merge main into feature branches.

## Breaking Changes

Flag in three places: commit footer (BREAKING CHANGE: description), CHANGELOG.md entry, PR description.

## Semantic Versioning

- MAJOR (X.0.0): breaking changes, new native modules, drastic structural changes.
- MINOR (0.X.0): incremental features, no breaking changes.
- PATCH (0.0.X): bug fixes in production only.

# Testing

## BDD Pattern (mandatory)

All tests must use clearly separated comment blocks:
- Go / JavaScript / TypeScript / Java: `// given`, `// when`, `// then`
- Python: `# given`, `# when`, `# then`

## Test Descriptions

- Commands: "should call <LISTENER> when ..."
- Controllers: "should respond <HTTP_STATUS_CODE> (<HTTP_STATUS>) when ..."
- Services/Repositories: "should ... when ..." -- at least one success and one failure test per public method.

## Test Doubles (preference order)

1. Stubs -- return pre-configured (canned) answers. No in-memory logic.
2. Dummies -- fill required parameters that are never used. Return minimal values.
3. In-memory -- implement logic in memory without external modules.
4. Fakers -- generate realistic fake data via an external library.
5. Mocks -- record and verify method calls. Avoid when possible.

## Builders

Use Builder Design Pattern for constructing complex test objects:
UserBuilder.new().withName("John Doe").withEmail("john@example.com").build()

# Documentation

Every code change must include corresponding documentation updates:

1. CHANGELOG.md -- always. Follow Keep a Changelog format. Add entries under [Unreleased]. Categories: Added, Changed, Deprecated, Removed, Fixed, Security. Start with lowercase verb in simple past tense.
2. README.md -- when behavior, configuration, CLI, or setup changes.
3. CONTRIBUTING.md -- when prerequisites, workflow, or project structure changes.
4. AI instruction files (CLAUDE.md, .cursorrules, AGENTS.md, .github/copilot-instructions.md) -- when architecture, commands, or development workflow changes.

Documentation and code ship as one unit in the same commit.

# YAML Conventions

- Always use `.yaml` extension (never `.yml`). Exception only when tool enforces `.yml`.
- Always single-quote strings: `name: 'my-service'`.
- Double quotes only for interpolation or escape sequences: `url: "${API_HOST}/v1"`.
- Never quote booleans or numbers: `enabled: true`, `replicas: 3`.
- Applies to all YAML: CI/CD pipelines, Kubernetes manifests, Docker Compose, Helm, app configs, YAML in Markdown.

# Security

- Never hard-code secrets, API keys, tokens, or passwords. Use environment variables or secret managers.
- Follow OWASP Top 10: validate/sanitize all user input, parameterized queries, enforce MFA, RBAC, HTTPS/TLS, encrypt at rest.
- Run `make sast` before every push. Suite: CodeQL, Semgrep, Trivy, Hadolint, Gitleaks.
- Rotate leaked secrets immediately. Never commit `.env` files.

# Language-Specific Conventions

## Go

- File names: snake_case (e.g., list_users_command.go).
- Method receivers: one/two letter abbreviation of type (c for Command, r for Repository). Never self/this/me.
- Entities: framework-agnostic. No json/gorm/other tags inside entity structs.
- DI: Uber Dig. Each layer has container.go for provider registration.
- Commands: <operation>_<entity>_command.go / <Operation><Entity>Command / method Execute.
- Controllers: <operation>_<entity>_controller.go / <Operation><Entity>Controller / method Execute.
- Repositories (contract): <entity>_repository.go / <Entity>Repository (interface).
- Repositories (impl): <library>_<entity>_repository.go / <Library><Entity>Repository.
- Mappers: <entity>_mapper.go / <Entity>Mapper.
- Models: prefixed with external tool name (PgxUser, AwsFile).
- Linting: golangci-lint. Logging: Logrus. Testing: testify.
- Services layer: not used in Go projects.
- Repository methods: FindByField, FindAllByField, HasCondition, Save, SaveAll, DeleteByField.
- Structure: cmd/<app>/ (main.go + dig.go), internal/domain/, internal/infrastructure/, test/ (builders, doubles, helpers).

## JavaScript / TypeScript

- File names: snake_case (e.g., my_class.js).
- HTML attributes: camelCase for id and data-test-subj values.
- Formatting: Prettier. Linting: ESLint.
- Prefer arrow functions, template strings, spread operator, optional chaining (?.), nullish coalescing (??).
- Never use `any` -- use `unknown` with type narrowing.
- Immutability: do not reassign variables, modify object properties, or push to arrays. Create new copies.
- Use object and array destructuring.

## Java

- File names: PascalCase (matches class name).
- Framework: Spring Boot. Build: Gradle (always use gradlew).
- DI: constructor injection via Lombok @RequiredArgsConstructor. Never field injection (@Autowired on fields). Mark fields final.
- Annotations: @Component for commands/listeners, @RestController for controllers, @Repository for repos, @Service for services.
- Mapping: MapStruct with INSTANCE pattern.
- Entities: no @Entity/@Table/@Column/JPA annotations. Use records for DTOs.
- Models: prefixed Jpa (JpaItem). Annotated @Entity/@Getter/@Setter/@SuperBuilder/@NoArgsConstructor.
- Commands: <Operation><Entity>Command.java / execute / inner Listeners record.
- Controllers: <Operation><Entity>Controller.java / execute.
- Services (contract): <Operation><Entity>Service.java (interface). Impl: Jpa<Operation><Entity>Service.java.
- Repositories (contract): <Entity>Repository.java. Impl: Jpa<Entity>Repository.java.
- Formatting: Google Java Format via Spotless. Linting: Checkstyle + PMD.
- Logging: SLF4J with Logback. Testing: JUnit 5 + Mockito.
- Database: Liquibase (.yaml, snake_case, single quotes, rollback + pre-conditions, constraints before column name).
- Principles: program to interfaces, composition over inheritance, immutability (final), fail fast.

## Python

- File names: snake_case. Package names cannot contain dashes.
- Classes: PascalCase.
- Formatting: Black. Import sorting: isort. Linting: Flake8.
- Type hints: required on all function signatures.
- Logging: Loguru. Testing: pytest. Package manager: PDM (PEP 621, pyproject.toml).
- Entities: framework-agnostic, no ORM base classes or decorators.
- Commands: <operation>_<entity>_command.py / <Operation><Entity>Command / execute.
- Controllers: <operation>_<entity>_controller.py / <Operation><Entity>Controller / execute.
- Repositories (contract): <entity>_repository.py / <Entity>Repository (ABC with @abstractmethod).
- Repositories (impl): <library>_<entity>_repository.py / <Library><Entity>Repository.
- Use descriptive variable names. Write comments that explain why, not just what.

# Code Review Guardrails

When reviewing code, enforce:

1. Architecture violations: entities must not import from infrastructure. Dependencies point inward.
2. Framework leakage: no ORM/framework annotations inside domain entities.
3. Naming consistency: operations vocabulary (List, Get, Insert, Update, Delete). Language-specific naming patterns.
4. Test quality: given/when/then blocks. Prefer stubs over mocks. Success + failure tests per public method.
5. Documentation: changelog entry mandatory. README update when behavior changes.
6. Commit format: type(SCOPE): message in simple past tense.
7. Security: no hardcoded secrets, parameterized queries, input validation.
8. YAML: .yaml extension, single-quoted strings, unquoted booleans/numbers.
9. Immutability: final fields in Java, no mutation in JS/TS, constructor injection everywhere.
10. Layer isolation: mappers between layers, no direct coupling between controllers and repositories.
````

## References

- [Documentation & Change Control](../Documentation-&-Change-Control.md)
- [Cursor Rules Documentation](https://docs.cursor.com/context/rules-for-ai)
- [Development Guide Wiki](https://github.com/rios0rios0/guide/wiki)

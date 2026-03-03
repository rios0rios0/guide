# CLAUDE.md Template

> **TL;DR:** Copy the template below into a `CLAUDE.md` file at the project root. This file gives [Claude Code](https://docs.anthropic.com/en/docs/claude-code) the full team development standards so it can write code, review PRs, and enforce conventions automatically -- no manual intervention needed.

## Overview

Claude Code reads `CLAUDE.md` from the project root (and parent directories) to understand project conventions. This template embeds the entire [Development Guide](https://github.com/rios0rios0/guide/wiki) so that Claude Code already knows every convention, architecture pattern, and language-specific rule the team follows. Install it via `install-ai-rules.sh` or copy the content below.

## Template

````markdown
# CLAUDE.md

## Commands

```bash
make lint    # fix all linter issues
make test    # run the full test suite
make sast    # run the complete SAST security suite
```

Always run lint, test, and sast before considering any task complete.

## Architecture

Follow **Hexagonal Architecture (Ports & Adapters)** with **Domain-Driven Design (DDD)** and **CQRS**.

### Layer Separation

- `domain/` contains contracts: entities, repository interfaces, service interfaces, and commands.
- `infrastructure/` contains implementations: controllers, repository implementations, service implementations, and mappers.
- Dependencies always point inward -- infrastructure depends on domain, never the reverse.
- Entities must be **framework-agnostic**. No ORM annotations, no framework tags, no external dependencies inside entity definitions.

### Principal Actors

| Actor            | Responsibility                                                                                                    |
|------------------|-------------------------------------------------------------------------------------------------------------------|
| **Entities**     | Core business objects. Contain domain rules about what the modeled objects are and how they behave. Framework-agnostic. |
| **Controllers**  | Bridge between the view (client) and the business layer. Receive requests and delegate to commands.               |
| **Commands**     | Implement business/feature logic and apply domain rules. Define a `Listeners` record for callback outcomes.       |
| **Services**     | Handle application-level concerns: parsing, conversions, transformations. Domain layer defines the contract (interface), infrastructure provides the implementation. |
| **Repositories** | Abstract all data access. Changing the database technology requires modifying only the repository implementation. |

### Mappers

Use mappers to isolate layers. Never let framework types leak across layer boundaries:
- **Repository mappers:** convert between infrastructure models (prefixed with tool name, e.g., `PgxUser`, `JpaItem`) and domain entities.
- **Controller mappers:** convert between request/response DTOs and domain entities. Separate request mappers from response mappers.

### File Structure Tiers

Scale structure with complexity:

**Minimal (3 layers):**
```
domain/commands/
domain/entities/
domain/repositories/
infrastructure/repositories/
```

**Intermediate (with Services):**
```
domain/commands/
domain/entities/
domain/repositories/
domain/services/
infrastructure/repositories/
infrastructure/services/
```

**Complete (with API layer):**
```
domain/commands/
domain/entities/
domain/repositories/
domain/services/
infrastructure/controllers/
infrastructure/controllers/requests/
infrastructure/controllers/responses/
infrastructure/controllers/mappers/
infrastructure/repositories/
infrastructure/repositories/mappers/
infrastructure/repositories/models/
infrastructure/services/
infrastructure/services/mappers/
```

### Frontend Architecture

5-layer lightweight Clean Architecture:
- **Domain** -- most abstract layer, defines entities and contracts.
- **Service** -- communication with APIs and external services.
- **Infrastructure** -- technology-specific implementations.
- **Presentation** -- view-rendering code (components, pages, hooks).
- **Main** -- wiring layer, instantiates providers and services, handles DI.

Dependencies always point toward Domain.

## Operations Vocabulary

Use these standard prefixes consistently for naming files, classes, and methods:

| Operation     | Description               |
|---------------|---------------------------|
| `List`        | Retrieve multiple records |
| `Get`         | Retrieve a single record  |
| `Insert`      | Create a single record    |
| `Update`      | Modify a single record    |
| `Delete`      | Remove a single record    |
| `BatchInsert` | Create multiple records   |
| `BatchUpdate` | Modify multiple records   |
| `BatchDelete` | Remove multiple records   |
| `DeleteAll`   | Remove all records        |

## Git Conventions

### Commit Messages

Format: `type(SCOPE): message`

- **SCOPE** is the task ID (e.g., `TICKET-123`) or a short descriptive scope.
- Use **simple past tense**: `added`, `changed`, `fixed`, `removed`.
- **No period** at the end. **Do not capitalize** the first letter.
- Wrap code references in backticks.

Types: `feat`, `fix`, `refactor`, `chore`, `test`, `docs`.

Long commit messages use a body with bullet points starting with a lowercase verb in past tense.

### Branches

Format: `type/TICKET-ID` or `type/short-description` (e.g., `feat/TICKET-000`, `fix/input-mask`).

### Branch Synchronization

Always use `git rebase` to synchronize with `main`. Never merge `main` into feature branches.

### Breaking Changes

Flag in three places:
1. Commit footer: `**BREAKING CHANGE:** description`
2. `CHANGELOG.md` entry: `- **BREAKING CHANGE:** description`
3. PR description: `**BREAKING CHANGE:** description`

### Semantic Versioning

- **MAJOR** (`X.0.0`): breaking changes, new native modules, drastic structural changes.
- **MINOR** (`0.X.0`): incremental features, no breaking changes.
- **PATCH** (`0.0.X`): bug fixes in production only.

## Testing

### BDD Pattern (mandatory)

All tests must use clearly separated comment blocks:

| Language                | Given      | When      | Then      |
|-------------------------|------------|-----------|-----------||
| Go                      | `// given` | `// when` | `// then` |
| JavaScript / TypeScript | `// given` | `// when` | `// then` |
| Java                    | `// given` | `// when` | `// then` |
| Python                  | `# given`  | `# when`  | `# then`  |

### Test Descriptions

- Commands: `"should call <LISTENER> when ..."`
- Controllers: `"should respond <HTTP_STATUS_CODE> (<HTTP_STATUS>) when ..."`
- Services/Repositories: `"should ... when ..."` -- at least one success and one failure test per public method.

### Test Doubles (preference order)

1. **Stubs** -- return pre-configured (canned) answers. No in-memory logic.
2. **Dummies** -- fill required parameters that are never used. Return minimal values.
3. **In-memory** -- implement logic in memory without external modules.
4. **Fakers** -- generate realistic fake data via an external library.
5. **Mocks** -- record and verify method calls. **Avoid when possible.**

### Builders

Use the Builder Design Pattern for constructing complex test objects:
```
UserBuilder.new().withName("John Doe").withEmail("john@example.com").build()
```

## Documentation

Every code change must include corresponding documentation updates:

1. **`CHANGELOG.md`** -- always. Follow [Keep a Changelog](https://keepachangelog.com/en/1.1.0/) format. Add entries under `[Unreleased]` with categories: Added, Changed, Deprecated, Removed, Fixed, Security. Start each entry with a lowercase verb in simple past tense. Be specific.
2. **`README.md`** -- when behavior, configuration, CLI, or setup changes.
3. **`CONTRIBUTING.md`** -- when prerequisites, workflow, or project structure changes.
4. **AI instruction files** (`CLAUDE.md`, `.cursorrules`, `AGENTS.md`, `.github/copilot-instructions.md`) -- when architecture, commands, or development workflow changes.

Documentation and code ship as one unit in the same commit. Never merge a PR that introduces user-facing or architectural changes without the corresponding documentation update.

## YAML Conventions

- Always use `.yaml` extension (never `.yml`). Exception only when the tool enforces `.yml` (e.g., Azure DevOps `azure-pipelines.yml`).
- Always single-quote strings: `name: 'my-service'`.
- Double quotes only for interpolation or escape sequences: `url: "${API_HOST}/v1"`, `message: "line1\nline2"`.
- Never quote booleans or numbers: `enabled: true`, `replicas: 3`.
- These rules apply to all YAML: CI/CD pipelines, Kubernetes manifests, Docker Compose, Helm values, application configs, and YAML code blocks inside Markdown.

## Security

- Never hard-code secrets, API keys, tokens, or passwords. Use environment variables or secret managers (1Password, Vault, AWS Secrets Manager).
- Follow OWASP Top 10 guidelines: validate and sanitize all user input, use parameterized queries, enforce MFA, implement RBAC, use HTTPS/TLS, encrypt data at rest.
- Run `make sast` before every push. The SAST suite includes CodeQL, Semgrep, Trivy, Hadolint, and Gitleaks.
- If Gitleaks detects a leaked secret, rotate the credential immediately.
- Use `.env` files for local development but never commit them.

## Language-Specific Conventions

### Go

- **File names:** `snake_case` (e.g., `list_users_command.go`).
- **Method receivers:** one or two letter abbreviation of the type (`c` for `Command`, `r` for `Repository`). Never use `self`, `this`, or `me`.
- **Entities:** framework-agnostic. No `json`, `gorm`, or other tags inside entity structs.
- **DI:** use [Uber Dig](https://github.com/uber-go/dig). Each architectural layer has a `container.go` file for provider registration.
- **Naming patterns:**
  - Commands: `<operation>_<entity>_command.go` / `<Operation><Entity>Command` / method `Execute`.
  - Controllers: `<operation>_<entity>_controller.go` / `<Operation><Entity>Controller` / method `Execute`.
  - Repositories (contract): `<entity>_repository.go` / `<Entity>Repository` (interface).
  - Repositories (impl): `<library>_<entity>_repository.go` / `<Library><Entity>Repository`.
  - Mappers: `<entity>_mapper.go` / `<Entity>Mapper`.
  - Models: prefixed with external tool name (e.g., `PgxUser`, `AwsFile`).
- **Linting:** `golangci-lint`.
- **Logging:** `Logrus`.
- **Testing:** `testify` (assert + suite + mock).
- **Services layer:** not used in Go projects.
- **Repository methods:** `FindByField`, `FindAllByField`, `HasCondition`, `Save`, `SaveAll`, `DeleteByField`.
- **Project structure:** `cmd/<app>/` (entry point + `dig.go`), `internal/domain/`, `internal/infrastructure/`, `test/` (builders, doubles, helpers).

### JavaScript / TypeScript

- **File names:** `snake_case` (e.g., `my_class.js`).
- **HTML attributes:** `camelCase` for `id` and `data-test-subj` values.
- **Formatting:** Prettier. **Linting:** ESLint.
- Prefer arrow functions, template strings, spread operator, optional chaining (`?.`), nullish coalescing (`??`).
- **Never use `any`** -- use `unknown` with type narrowing instead.
- **Immutability:** do not reassign variables, modify object properties, or push to arrays. Create new copies.
- **Destructuring:** use object and array destructuring instead of index access.

### Java

- **File names:** `PascalCase` (matches class name, e.g., `ListItemsCommand.java`).
- **Framework:** Spring Boot.
- **DI:** constructor injection via Lombok `@RequiredArgsConstructor`. **Never** use field injection (`@Autowired` on fields). Mark injected fields as `final`.
- **Annotations:** `@Component` for commands/listeners, `@RestController` for controllers, `@Repository` for repositories, `@Service` for services.
- **Mapping:** MapStruct with `@Mapper(unmappedTargetPolicy = ReportingPolicy.IGNORE)` and `INSTANCE` pattern.
- **Entities:** no `@Entity`, `@Table`, `@Column`, or JPA annotations. Use Java records for DTOs and value objects.
- **Models:** prefixed with `Jpa` (e.g., `JpaItem`). Annotated with `@Entity`, `@Getter`, `@Setter`, `@SuperBuilder`, `@NoArgsConstructor`.
- **Naming patterns:**
  - Commands: `<Operation><Entity>Command.java` / method `execute` / inner `Listeners` record.
  - Controllers: `<Operation><Entity>Controller.java` / method `execute`.
  - Services (contract): `<Operation><Entity>Service.java` (interface).
  - Services (impl): `Jpa<Operation><Entity>Service.java`.
  - Repositories (contract): `<Entity>Repository.java` (interface).
  - Repositories (impl): `Jpa<Entity>Repository.java` (extends `JpaRepository`).
- **Formatting:** Google Java Format via Spotless. **Linting:** Checkstyle + PMD.
- **Logging:** SLF4J with Logback.
- **Testing:** JUnit 5 + Mockito.
- **Build:** Gradle (always use `gradlew` wrapper).
- **Database:** Liquibase with `.yaml` format, `snake_case` names, single quotes, rollback + pre-conditions on every changeset, constraints before column name.
- **Principles:** program to interfaces, favor composition over inheritance, immutability by default (`final`), fail fast, convention over configuration.

### Python

- **File names:** `snake_case` (e.g., `list_users_command.py`). Package names cannot contain dashes.
- **Classes:** `PascalCase`.
- **Formatting:** Black. **Import sorting:** isort. **Linting:** Flake8.
- **Type hints:** required on all function signatures.
- **Logging:** Loguru.
- **Testing:** pytest.
- **Package manager:** PDM with PEP 621 metadata in `pyproject.toml`.
- **Entities:** framework-agnostic, no ORM base classes or decorators.
- **Naming patterns:**
  - Commands: `<operation>_<entity>_command.py` / `<Operation><Entity>Command` / method `execute`.
  - Controllers: `<operation>_<entity>_controller.py` / `<Operation><Entity>Controller` / method `execute`.
  - Repositories (contract): `<entity>_repository.py` / `<Entity>Repository` (ABC with `@abstractmethod`).
  - Repositories (impl): `<library>_<entity>_repository.py` / `<Library><Entity>Repository`.
- **DI:** constructor injection. Use framework DI (e.g., FastAPI) or wire manually in composition root.
- Use descriptive variable names. Write comments that explain **why**, not just **what**.

## Code Review Guardrails

When reviewing code, enforce these rules:

1. **Architecture violations:** entities must not import from infrastructure. Dependencies must point inward.
2. **Framework leakage:** no ORM annotations, framework tags, or external tool dependencies inside domain entities.
3. **Naming consistency:** operations vocabulary (`List`, `Get`, `Insert`, `Update`, `Delete`) must be used. File and class naming must follow language-specific patterns.
4. **Test quality:** every test must have `given`/`when`/`then` blocks. Prefer stubs over mocks. At least one success + one failure test per public method.
5. **Documentation:** changelog entry is mandatory. README update is mandatory when behavior changes.
6. **Commit format:** `type(SCOPE): message` in simple past tense.
7. **Security:** no hardcoded secrets, parameterized queries, input validation present.
8. **YAML:** `.yaml` extension, single-quoted strings, unquoted booleans/numbers.
9. **Immutability:** `final` fields in Java, no mutation in JS/TS, constructor injection everywhere.
10. **Layer isolation:** mappers exist between layers, no direct coupling between controllers and repositories.
````

## References

- [Documentation & Change Control](../Documentation-&-Change-Control.md)
- [Claude Code Documentation](https://docs.anthropic.com/en/docs/claude-code)
- [Development Guide Wiki](https://github.com/rios0rios0/guide/wiki)

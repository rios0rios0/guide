---
paths:
  - "**/*.md"
---

# CHANGELOG Formatting

> **TL;DR:** Proper nouns are capitalized, code identifiers use backticks, acronyms and protocol names use their official casing, versions and library names use backticks.

## Rules

### 1. Proper Nouns -- Capitalize Correctly

Product names, company names, and personal names must use their official capitalization.

| Wrong      | Correct    |
|------------|------------|
| github     | GitHub     |
| kubernetes | Kubernetes |
| docker     | Docker     |
| terraform  | Terraform  |
| john doe   | John Doe   |
| elasticsearch | Elasticsearch |
| postgresql | PostgreSQL |
| javascript | JavaScript |
| typescript | TypeScript |

### 2. Code Identifiers -- Use Backticks

Class names, function names, variable names, file names, and any code reference must be wrapped in backticks.

| Wrong                                      | Correct                                        |
|--------------------------------------------|------------------------------------------------|
| added CreateUser command                   | added `CreateUser` command                     |
| fixed bug in handleRequest function        | fixed bug in `handleRequest` function          |
| renamed user_id to userId                  | renamed `user_id` to `userId`                  |
| updated Dockerfile                         | updated `Dockerfile`                           |
| changed settings.json configuration        | changed `settings.json` configuration          |

### 3. Acronyms and Protocol Names -- Use Official Casing

Technical acronyms must be uppercased. Branded protocol and technology names must use their official casing.

| Wrong  | Correct |
|--------|---------|
| http   | HTTP    |
| api    | API     |
| sql    | SQL     |
| css    | CSS     |
| html   | HTML    |
| jwt    | JWT     |
| grpc   | gRPC    |
| graphql | GraphQL |
| rest   | REST    |
| oauth  | OAuth   |

### 4. Versions -- Use Backticks

Version numbers (with or without the `v` prefix) must be wrapped in backticks.

| Wrong                          | Correct                            |
|--------------------------------|------------------------------------|
| bumped to v1.2.0               | bumped to `v1.2.0`                 |
| upgraded from 2.0.0 to 3.0.0  | upgraded from `2.0.0` to `3.0.0`  |

### 5. Library and Package Names -- Use Backticks

Dependency names, package names, and module names must be wrapped in backticks.

| Wrong                              | Correct                                |
|------------------------------------|----------------------------------------|
| upgraded lodash to 4.17.21         | upgraded `lodash` to `4.17.21`         |
| added gin-gonic as HTTP framework  | added `gin-gonic` as HTTP framework    |
| replaced moment with dayjs         | replaced `moment` with `dayjs`         |

## Examples

**Bad:**

```markdown
- added createUser endpoint using express framework
- upgraded golang to v1.23
- fixed sql injection in handleLogin
- integrated with github actions for ci/cd
```

**Good:**

```markdown
- added `CreateUser` endpoint using `express` framework
- upgraded Go to `v1.23`
- fixed SQL injection in `handleLogin`
- integrated with GitHub Actions for CI/CD
```

# YAML Conventions

> **TL;DR:** Always use the `.yaml` extension (not `.yml`). Always quote strings with **single quotes**. Use **double quotes** only when the string contains variable interpolation or escape sequences. Never leave string values unquoted. Do not quote booleans or numbers. These rules apply to all YAML files: pipeline configurations, Kubernetes manifests, infrastructure-as-code, and YAML code blocks inside Markdown.

## Overview

[YAML](https://yaml.org/) (YAML Ain't Markup Language) is a human-readable data serialization format used extensively for configuration files, CI/CD pipelines, Kubernetes manifests, and infrastructure-as-code. The current specification is [YAML 1.2.2](https://yaml.org/spec/1.2.2/), published on October 1, 2021.

This document defines the team's conventions for writing consistent, unambiguous YAML across all projects and documentation.

## File Extension

**Always use `.yaml` as the file extension.**

The [official YAML FAQ](https://yaml.org/faq.html) recommends `.yaml` as the preferred extension. The `.yml` variant exists for historical reasons (legacy Windows three-character extension limits) and must be avoided.

```
# Correct
docker-compose.yaml
config.yaml
deployment.yaml
.golangci.yaml

# Wrong
docker-compose.yml
config.yml
deployment.yml
.golangci.yml
```

### Exceptions

Some tools enforce a specific filename that uses `.yml` and do not accept alternatives. In these rare cases, the tool's requirement takes precedence:

| Tool | Required Filename | Reason |
|------|-------------------|--------|
| Azure DevOps | `azure-pipelines.yml` | Only recognizes this exact filename |
| Docker Compose (legacy) | `docker-compose.yml` | Older versions required this name (modern versions accept `.yaml`) |

If a tool accepts both extensions, always use `.yaml`.

## String Quoting

### Rule: Always Quote Strings with Single Quotes

All string values must be explicitly quoted with **single quotes** (`'...'`). Unquoted strings are technically valid YAML, but they introduce ambiguity -- values like `yes`, `no`, `on`, `off`, `null`, or numeric-looking strings can be silently misinterpreted as booleans, nulls, or numbers by the YAML parser.

```yaml
# Correct
name: 'my-service'
image: 'nginx:1.25-alpine'
region: 'us-east-1'
environment: 'production'

# Wrong -- unquoted strings
name: my-service
image: nginx:1.25-alpine
region: us-east-1
environment: production
```

### Rule: Use Double Quotes Only for Interpolation or Escape Sequences

Use **double quotes** (`"..."`) only when the string contains variable interpolation (e.g., environment variable expansion) or escape sequences (e.g., `\n`, `\t`):

```yaml
# Correct -- double quotes for interpolation
connection_string: "${DATABASE_HOST}:${DATABASE_PORT}"
greeting: "Hello,\tWorld\n"

# Correct -- single quotes for everything else
host: 'localhost'
port_label: '8080'
```

### Rule: Do Not Quote Booleans or Numbers

Booleans and numbers are native YAML types and must **not** be quoted:

```yaml
# Correct
enabled: true
replicas: 3
timeout: 30.5
debug: false

# Wrong -- quoted booleans and numbers
enabled: 'true'
replicas: '3'
timeout: '30.5'
```

### Summary Table

| Value Type | Quoting Style | Example |
|-----------|--------------|---------|
| Plain string | Single quotes | `name: 'my-app'` |
| String with variables | Double quotes | `url: "${API_HOST}/v1"` |
| String with escape sequences | Double quotes | `message: "line1\nline2"` |
| Boolean | No quotes | `enabled: true` |
| Number (integer) | No quotes | `replicas: 3` |
| Number (float) | No quotes | `ratio: 0.75` |
| Null | No quotes | `value: null` |

## YAML in Markdown Code Blocks

When embedding YAML snippets inside Markdown files (READMEs, guides, documentation), the same conventions apply. Code examples must be consistent with production YAML:

````markdown
```yaml
# Correct -- follows all conventions
apiVersion: 'apps/v1'
kind: 'Deployment'
metadata:
  name: 'my-service'
  labels:
    app: 'my-service'
    environment: 'production'
spec:
  replicas: 3
  selector:
    matchLabels:
      app: 'my-service'
```
````

This includes YAML examples in:

- `README.md` files
- Wiki pages and guides
- Pull request descriptions
- Code review comments
- Any other documentation containing YAML snippets

## Scope of Application

These conventions apply to **all** YAML files across the organization:

| Domain | Examples |
|--------|---------|
| **CI/CD pipelines** | GitHub Actions workflows, GitLab CI, Azure DevOps pipelines |
| **Kubernetes** | Deployments, Services, ConfigMaps, Ingresses, Helm values |
| **Infrastructure-as-Code** | Docker Compose, Ansible playbooks, CloudFormation templates |
| **Application configuration** | Spring Boot `application.yaml`, service configs |
| **Tooling** | `.golangci.yaml`, `.hadolint.yaml`, `autoupdate.yaml` |
| **Documentation** | YAML code blocks inside Markdown files |

## References

- [YAML Specification 1.2.2](https://yaml.org/spec/1.2.2/)
- [YAML Official FAQ](https://yaml.org/faq.html)
- [YAML Quoting Guide](https://www.yaml.info/learn/quote.html)
- [Do I Need Quotes for Strings in YAML?](https://stackoverflow.com/questions/19109912/do-i-need-quotes-for-strings-in-yaml)

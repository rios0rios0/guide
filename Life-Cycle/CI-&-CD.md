# CI/CD Pipeline

> **TL;DR:** A complete CI/CD pipeline consists of 10 stages, from source control to continuous feedback. All projects use a shared Makefile with `make lint`, `make test`, and `make sast` targets. The SAST suite runs CodeQL, Semgrep, Trivy, Hadolint, and Gitleaks. Always run these locally before pushing code.

## Overview

This document defines the standards for building Continuous Integration and Continuous Deployment (CI/CD) pipelines. The pipeline described below represents the ideal end-to-end flow; individual projects may implement a subset based on their requirements.

## Pipeline Stages

| Stage                            | Description                                                                                                   |
|----------------------------------|---------------------------------------------------------------------------------------------------------------|
| **1. Source Control Management** | Code is stored in a VCS (Git) and committed to the repository by developers.                                  |
| **2. Build**                     | Code is compiled and packaged into a deployable artifact by the build system (e.g., Jenkins, GitHub Actions). |
| **3. Unit Testing**              | Automated unit tests validate that individual components function correctly, catching bugs early.             |
| **4. Integration Testing**       | The code is integrated with other system components and tested for correct inter-module behavior.             |
| **5. Static Analysis (SAST)**    | Code quality, security, and secret detection checks are performed using the SAST toolchain (see below).       |
| **6. Deployment (Staging)**      | The artifact is deployed to a staging environment for further testing and validation.                         |
| **7. Acceptance Testing**        | Users and stakeholders verify that the system meets business requirements in the staging environment.         |
| **8. Release (Production)**      | After passing all acceptance tests, the artifact is released to production.                                   |
| **9. Monitoring**                | The production system is monitored for correct operation and early detection of issues.                       |
| **10. Continuous Feedback**      | Feedback from users and stakeholders drives improvements to the system and the pipeline itself.               |

## Local Quality Gates

All projects use a shared `Makefile` that imports targets from the [pipelines repository](https://github.com/rios0rios0/pipelines). Developers must run these targets locally before pushing code:

```bash
make lint    # Fix all linter issues
make test    # Run the full test suite
make sast    # Run the complete SAST security suite
```

**Never invoke linters, test runners, or security tools directly.** Always use the Makefile targets to ensure consistency with CI.

## SAST Toolchain

All SAST, SCA, and quality tools are provided by the [Pipelines repository](https://github.com/rios0rios0/pipelines). For the full list of available tools, configurations, and usage instructions, see the [Available Tools & Scripts](https://github.com/rios0rios0/pipelines?tab=readme-ov-file#-available-tools--scripts) section of the pipelines README.

The `make sast` target orchestrates five security tools that run both locally and in CI pipelines:

| Tool                                                 | Purpose                                                                                                     | Output              |
|------------------------------------------------------|-------------------------------------------------------------------------------------------------------------|---------------------|
| **[CodeQL](https://codeql.github.com/)**             | Static Application Security Testing -- detects SQL injection, XSS, path traversal, insecure deserialization | SARIF report        |
| **[Semgrep](https://semgrep.dev/)**                  | Pattern-based static analysis with OWASP Top 10, secrets, and best-practice rule sets                       | JSON report         |
| **[Trivy](https://trivy.dev/)**                      | Infrastructure-as-Code misconfiguration scanning and dependency vulnerability scanning                      | SARIF / JSON report |
| **[Hadolint](https://github.com/hadolint/hadolint)** | Dockerfile linting and best practice enforcement                                                            | SARIF report        |
| **[Gitleaks](https://gitleaks.io/)**                 | Secret and credential detection across the entire Git history                                               | JSON report         |

All reports are generated in the `build/reports/` directory.

### False Positive Management

Each tool supports project-level configuration for suppressing false positives:

| Tool | Configuration File |
|------|-------------------|
| CodeQL | `.codeql-false-positives` |
| Semgrep | `.semgrepignore`, `.semgrepexcluderules` |
| Trivy | `.trivyignore` |
| Hadolint | `.hadolint.yaml` |
| Gitleaks | `.gitleaks.toml` |

## Pre-Commit Hooks

While no formal pre-commit hook standard has been defined yet, developers are strongly encouraged to run `make lint` and `make sast` before every commit. These targets serve as the equivalent of a pre-commit gate, catching issues before they enter the repository.

A future standard may formalize this using tools like [pre-commit](https://pre-commit.com/) to automatically run the following checks on every commit:

- **Linting** (language-specific formatters and linters)
- **Secret detection** (Gitleaks)
- **Static analysis** (Semgrep with targeted rule sets)

Until a formal standard is adopted, treat `make lint && make sast` as a mandatory manual pre-push step.

## Key Principles

- **Automate everything.** Every stage that can be automated should be automated.
- **Fail fast.** Run the fastest and cheapest checks (linting, unit tests) first.
- **Shift left on security.** Run SAST tools locally and in CI, not just in production.
- **Immutable artifacts.** Build once, deploy the same artifact to every environment.
- **Rollback capability.** Every deployment must support quick rollback to the previous version.

## References

- [Pipelines Repository -- Available Tools & Scripts](https://github.com/rios0rios0/pipelines?tab=readme-ov-file#-available-tools--scripts)
- [Continuous Integration -- Martin Fowler](https://martinfowler.com/articles/continuousIntegration.html)
- [Continuous Delivery -- Jez Humble](https://continuousdelivery.com/)
- [The Twelve-Factor App](https://12factor.net/)
- [OWASP DevSecOps Guideline](https://owasp.org/www-project-devsecops-guideline/)
- [pre-commit Framework](https://pre-commit.com/)

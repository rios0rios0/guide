# Security Practices

> **TL;DR:** Follow OWASP Top 10 and MITRE ATT&CK guidelines. Run `make sast` before every push to detect secrets (Gitleaks), code vulnerabilities (CodeQL, Semgrep), IaC misconfigurations (Trivy), and Dockerfile issues (Hadolint). Never hard-code secrets -- use environment variables or secret managers.

## Overview

This document defines the security best practices that must be considered when designing, developing, and deploying applications. These practices are grounded in industry standards from OWASP and MITRE, and enforced through the team's SAST toolchain.

## Security Checklist

### 1. Input Validation and Sanitization (OWASP A1)

Ensure all user input is validated and sanitized to prevent injection attacks (SQL injection, XSS, command injection):

- Validate data types and input length.
- Filter or encode special characters.
- Use parameterized queries for database access.

### 2. Authentication and Authorization (OWASP A2)

Implement robust identity and access management:

- Enforce **multi-factor authentication (MFA)**.
- Implement **role-based access controls (RBAC)**.
- Secure session management with proper token expiration and rotation.

### 3. Encryption (OWASP A3)

Protect sensitive data both in transit and at rest:

- Use **HTTPS/TLS** for all data in transit.
- Encrypt sensitive data stored on disk or in databases.
- Manage encryption keys securely (never hard-code secrets).

### 4. Access Controls (OWASP A4)

Apply the **principle of least privilege** and **need-to-know** across all systems:

- Implement fine-grained access controls.
- Use access control lists (ACLs) where appropriate.
- Regularly review and revoke unnecessary privileges.

### 5. Security Testing (OWASP A5)

Integrate security testing into the development pipeline:

- Perform **penetration testing** regularly.
- Use **Dynamic Application Security Testing (DAST)** in CI/CD pipelines.
- Conduct vulnerability assessments on dependencies (e.g., `npm audit`, `snyk`).

### 6. Incident Response (MITRE ATT&CK)

Maintain a documented incident response plan:

- Define scope identification procedures.
- Establish damage containment protocols.
- Document recovery and post-mortem processes.

### 7. Logging and Monitoring (MITRE ATT&CK)

Implement comprehensive observability:

- Use **centralized logging** and event management systems.
- Deploy **intrusion detection and prevention systems (IDS/IPS)**.
- Set up alerting for anomalous behavior.

### 8. Network Security (MITRE ATT&CK)

Protect the network perimeter and internal communications:

- Implement **firewalls** and **network segmentation**.
- Use **VPNs** for remote access.
- Apply **network access controls** to limit lateral movement.

## SAST Pipeline (Static Application Security Testing)

All projects enforce security through a standardized SAST pipeline provided by the [Pipelines repository](https://github.com/rios0rios0/pipelines). The full list of available security tools, their configuration, and usage is documented in the [Available Tools & Scripts](https://github.com/rios0rios0/pipelines?tab=readme-ov-file#-available-tools--scripts) section of the pipelines README.

The pipeline is invoked locally via `make sast`, which runs the following tools in sequence:

| Tool                                                 | Category  | What It Detects                                                                                     |
|------------------------------------------------------|-----------|-----------------------------------------------------------------------------------------------------|
| **[CodeQL](https://codeql.github.com/)**             | SAST      | SQL injection, XSS, path traversal, insecure deserialization, and other vulnerability patterns      |
| **[Semgrep](https://semgrep.dev/)**                  | SAST      | OWASP Top 10 patterns, hardcoded secrets, language-specific anti-patterns, best practice violations |
| **[Trivy](https://trivy.dev/)**                      | IaC / SCA | Infrastructure-as-Code misconfigurations and dependency vulnerabilities                             |
| **[Hadolint](https://github.com/hadolint/hadolint)** | Linting   | Dockerfile best practice violations (unpinned base images, running as root, missing health checks)  |
| **[Gitleaks](https://gitleaks.io/)**                 | Secrets   | API keys, tokens, passwords, private keys, and other secrets committed to Git history               |

Beyond SAST, the pipelines repository also provides **SCA (Software Composition Analysis)** tools such as `govulncheck` (Go), `Safety` (Python), `OWASP Dependency-Check` (Java), and `yarn npm audit` (JavaScript). See the [full tool reference](https://github.com/rios0rios0/pipelines?tab=readme-ov-file#-available-tools--scripts) for details.

### Running Locally

```bash
make setup   # Clone/update the pipelines repository (first time only)
make sast    # Run the full SAST suite
```

All reports are generated in `build/reports/`. If any tool reports findings, fix them before pushing code.

### Secret Hygiene

Leaked secrets are one of the most common and dangerous security issues. Follow these rules:

- **Never** hard-code secrets, API keys, tokens, or passwords in source code.
- Use **environment variables** or **secret managers** (1Password, HashiCorp Vault, AWS Secrets Manager).
- If Gitleaks detects a leaked secret, **rotate the credential immediately** and remove it from the codebase.
- Use `.env` files for local development, but **never commit them** (add `.env` to `.gitignore`).

## Pre-Commit Hooks

While a formal pre-commit hook standard has not yet been adopted, the team strongly recommends running `make lint && make sast` before every commit. This serves as a manual pre-commit gate.

A future standard may use [pre-commit](https://pre-commit.com/) to automate the following hooks:

| Hook             | Tool                                                   | Purpose                                      |
|------------------|--------------------------------------------------------|----------------------------------------------|
| Secret detection | Gitleaks                                               | Prevent secrets from entering the repository |
| Linting          | Language-specific (golangci-lint, Black, ESLint, etc.) | Enforce code style                           |
| Static analysis  | Semgrep (targeted rules)                               | Catch common security anti-patterns          |
| Commit message   | Custom script                                          | Enforce `type(SCOPE): message` format        |

Until formalized, treat `make lint && make sast` as a **mandatory pre-push step**.

## References

- [Pipelines Repository -- Available Tools & Scripts](https://github.com/rios0rios0/pipelines?tab=readme-ov-file#-available-tools--scripts)
- [OWASP Top 10](https://owasp.org/www-project-top-ten/)
- [MITRE ATT&CK Framework](https://attack.mitre.org/)
- [OWASP Proactive Controls](https://owasp.org/www-project-proactive-controls/)
- [OWASP Cheat Sheet Series](https://cheatsheetseries.owasp.org/)
- [OWASP DevSecOps Guideline](https://owasp.org/www-project-devsecops-guideline/)
- [Gitleaks Documentation](https://github.com/gitleaks/gitleaks)
- [pre-commit Framework](https://pre-commit.com/)

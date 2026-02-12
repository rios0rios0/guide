# Forking Technique

> **TL;DR:** Reserve `main` for the upstream community version and use `custom` as the working branch. Rebase `custom` onto updated `main` when upstream releases new versions. Tag fork releases with an incremental fourth digit (e.g., `1.0.0.1`).

## Overview

When forking an open-source project, it is common for upstream maintainers to take significant time reviewing and merging community contributions. This standard defines a forking strategy that maintains compatibility with upstream releases while preserving custom modifications.

## Strategy

### Branch Convention

- The **`main`** branch mirrors the upstream community version. It is updated by rebasing on upstream releases.
- The **`custom`** branch serves as the team's working branch (equivalent to `main` for internal purposes).

### Synchronization

To stay current with the upstream project:

1. Update `main` by rebasing on the latest upstream release.
2. Rebase `custom` onto the updated `main`, placing custom modifications on top of the newest version.

### Versioning

Fork versions use an **incremental fourth digit** appended to the upstream version:

| Scenario                                  | Version   |
|-------------------------------------------|-----------|
| Upstream at `1.0.0`, fork synced          | `1.0.0.0` |
| New fork release                          | `1.0.0.1` |
| Upstream updates to `1.0.1`, fork rebased | `1.0.1.0` |
| Next fork release                         | `1.0.1.1` |

The fourth digit resets to `0` each time the fork is rebased on a new upstream version.

## Caveats

If your CI/CD or release tooling does not support four-segment version numbers (`X.Y.Z.N`), use a dash separator instead: `X.Y.Z-N`.

## References

- [Semantic Versioning](https://semver.org/)
- [Atlassian - Forking Workflow](https://www.atlassian.com/git/tutorials/comparing-workflows/forking-workflow)

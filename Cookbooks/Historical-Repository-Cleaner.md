# Historical Repository Cleaner

> **TL;DR:** Use [BFG Repo-Cleaner](https://rtyley.github.io/bfg-repo-cleaner/) to remove secrets and sensitive data from Git history. Clone with `--mirror`, run BFG with a `replacements.txt` file, then garbage-collect and force-push.

## Overview

When secrets or credentials are accidentally committed to a Git repository, simply deleting them in a new commit does not remove them from the history. BFG Repo-Cleaner provides a fast, efficient way to rewrite Git history and replace sensitive data across all commits.

## Steps

### 1. Install BFG

```bash
wget https://repo1.maven.org/maven2/com/github/git-tools/bfg/1.14.0/bfg-1.14.0.jar -O bfg.jar
```

### 2. Clone the Repository (Mirror)

```bash
git clone --mirror git://example.com/<your-repo>.git
```

### 3. Prepare `replacements.txt`

Create a `replacements.txt` file containing the secrets to replace. Each line defines a pattern that BFG will substitute across the entire history. Identify secrets using a security scanning tool of your choice.

### 4. Run BFG

Ensure both the repository and `replacements.txt` are in the same directory:

```bash
java -jar bfg.jar --replace-text replacements.txt <your-repo>.git
```

### 5. Clean Up Git History

```bash
cd <your-repo>.git
git reflog expire --expire=now --all
git gc --prune=now --aggressive
```

### 6. Push the Cleaned History

```bash
git push --mirror
```

**Warning:** This rewrites the entire repository history. All collaborators must re-clone after this operation.

## References

- [BFG Repo-Cleaner](https://rtyley.github.io/bfg-repo-cleaner/)
- [GitHub - Removing Sensitive Data](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/removing-sensitive-data-from-a-repository)

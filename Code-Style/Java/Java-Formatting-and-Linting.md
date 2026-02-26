# Java Formatting and Linting

> **TL;DR:** Use **[Google Java Format](https://github.com/google/google-java-format)** (via the [Spotless](https://github.com/diffplug/spotless) Gradle plugin) for code formatting, **[Checkstyle](https://checkstyle.sourceforge.io/)** for style enforcement, and **[PMD](https://pmd.github.io/)** for static code analysis. These tools are non-negotiable and must be integrated into every project's CI pipeline.

## Overview

Java's ecosystem does not include a built-in formatter like Go, so the team standardizes on Google Java Format. Combined with Checkstyle for style rules and PMD for bug detection, this toolchain ensures consistent, high-quality code across all projects.

## Formatter: Google Java Format (via Spotless)

[Google Java Format](https://github.com/google/google-java-format) produces a single canonical formatting for any Java source file. It is applied via the [Spotless](https://github.com/diffplug/spotless) Gradle plugin, which also handles import ordering and license header enforcement.

### Gradle Configuration

```groovy
plugins {
    id 'com.diffplug.spotless' version '6.25.0'
}

spotless {
    java {
        googleJavaFormat('1.22.0')
        removeUnusedImports()
        trimTrailingWhitespace()
        endWithNewline()
    }
}
```

### Usage

```bash
# Auto-format all Java files
gradle spotlessApply

# Check formatting without modifying files (CI mode)
gradle spotlessCheck
```

**Do not use alternative formatters.** Google Java Format is the team standard and all Java code must be formatted with it.

## Import Ordering

Spotless handles import ordering automatically when using Google Java Format. The standard grouping is:

1. `com.*` -- third-party packages
2. `java.*` -- standard library
3. `javax.*` / `jakarta.*` -- extension libraries
4. Static imports (last)

Each group is separated by a blank line. Unused imports are removed automatically.

## Linter: Checkstyle

[Checkstyle](https://checkstyle.sourceforge.io/) enforces coding standards and style rules. The project uses the Google ruleset as a baseline with project-specific customizations.

### Gradle Configuration

```groovy
plugins {
    id 'checkstyle'
}

checkstyle {
    toolVersion = '10.12.1'
    configFile = file('src/main/resources/app/quality/checkstyle-google-ruleset.xml')
    maxWarnings = 0
    maxErrors = 0
}
```

### Usage

```bash
# Run Checkstyle on main sources
gradle checkstyleMain

# Run Checkstyle on test sources
gradle checkstyleTest
```

### Key Rules Enforced

| Rule | Description |
|------|-------------|
| `IndentationCheck` | Enforces consistent indentation (2 spaces) |
| `LineLength` | Maximum line length (100 characters) |
| `NeedBraces` | All control structures must use braces |
| `JavadocMethod` | Public methods require Javadoc (configurable) |
| `UnusedImports` | No unused imports allowed |

## Static Analysis: PMD

[PMD](https://pmd.github.io/) detects potential bugs, dead code, suboptimal code, and overly complicated expressions.

### Gradle Configuration

```groovy
plugins {
    id 'pmd'
}

pmd {
    toolVersion = '7.1.0'
    ruleSetFiles = files('src/main/resources/app/quality/pmd-custom-ruleset.xml')
    consoleOutput = true
}
```

### Usage

```bash
# Run PMD analysis
gradle pmdMain
gradle pmdTest
```

## Security Analysis: SpotBugs

[SpotBugs](https://spotbugs.github.io/) detects potential security vulnerabilities and bug patterns in compiled bytecode.

### Gradle Configuration

```groovy
plugins {
    id 'com.github.spotbugs' version '6.0.22'
}

spotbugs {
    toolVersion = '4.8.4'
    excludeFilter = file('src/main/resources/app/security/spotbugs-security-exclude.xml')
}
```

### Usage

```bash
# Run SpotBugs on main sources
gradle spotbugsMain

# Run SpotBugs on test sources
gradle spotbugsTest
```

## Editor Configuration

### IntelliJ IDEA

1. Install the [google-java-format plugin](https://plugins.jetbrains.com/plugin/8527-google-java-format).
2. Enable it in **Settings > google-java-format Settings > Enable google-java-format**.
3. Enable **Settings > Editor > General > Auto Import > Optimize imports on the fly**.
4. Install the [Checkstyle-IDEA plugin](https://plugins.jetbrains.com/plugin/1065-checkstyle-idea) and point it to the project's ruleset.

### Visual Studio Code

Install the [Language Support for Java](https://marketplace.visualstudio.com/items?itemName=redhat.java) extension and add the following settings:

```json
{
    "java.format.settings.url": "https://raw.githubusercontent.com/google/styleguide/gh-pages/eclipse-java-google-style.xml",
    "java.format.enabled": true,
    "editor.formatOnSave": true,
    "[java]": {
        "editor.defaultFormatter": "redhat.java"
    }
}
```

## References

- [Google Java Format](https://github.com/google/google-java-format)
- [Spotless Gradle Plugin](https://github.com/diffplug/spotless/tree/main/plugin-gradle)
- [Checkstyle Documentation](https://checkstyle.sourceforge.io/)
- [PMD Documentation](https://pmd.github.io/)
- [SpotBugs Documentation](https://spotbugs.github.io/)
- [Google Java Style Guide](https://google.github.io/styleguide/javaguide.html)

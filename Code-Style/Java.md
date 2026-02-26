# Java

> **TL;DR:** Use [Spring Boot](https://spring.io/projects/spring-boot) as the application framework, [Google Java Format](https://github.com/google/google-java-format) (via [Spotless](https://github.com/diffplug/spotless)) for formatting, [Checkstyle](https://checkstyle.sourceforge.io/) + [PMD](https://pmd.github.io/) for linting, [MapStruct](https://mapstruct.org/) for object mapping, [SLF4J](https://www.slf4j.org/) with [Logback](https://logback.qos.ch/) for logging, [JUnit 5](https://junit.org/junit5/) + [Mockito](https://site.mockito.org/) for testing, and [Gradle](https://gradle.org/) for builds.

## Overview

This series of pages outlines Java-specific conventions, with detailed explanations of recommended approaches and the reasoning behind them. For the general baseline, refer to the [Code Style](../Code-Style.md) guide. See the sub-pages for specific topics:

- [Conventions](Java/Java-Conventions.md)
- [Formatting and Linting](Java/Java-Formatting-and-Linting.md)
- [Type System](Java/Java-Type-System.md)
- [Logging](Java/Java-Logging.md)
- [Testing](Java/Java-Testing.md)
- [Project Structure](Java/Java-Project-Structure.md)

## Liquibase

All database versioning is managed through [Liquibase](https://www.liquibase.org/). Follow these mandatory conventions:

| Rule                 | Description                                                                |
|----------------------|----------------------------------------------------------------------------|
| File format          | Always use `.yaml` (not `.yml`)                                            |
| Quoting              | Always use single quotes                                                   |
| Naming               | Always use `snake_case`                                                    |
| Rollbacks            | Every changeset must include a rollback statement                          |
| Pre-conditions       | Every changeset must include pre-conditions to validate the database state |
| Boolean/Time columns | Use descriptive names: `created_at`, `enabled`, `updated_at`, `activated`  |

### Column Definition Order

Constraints must appear **before** the column name in column definitions:

```yaml
# Correct
- column:
    constraints:
        nullable: false
    name: 'email'
    type: 'VARCHAR(255)'

# Wrong
- column:
    name: 'email'
    type: 'VARCHAR(255)'
    constraints:
        nullable: false
```

### Constraint Naming

| Type                   | Pattern                                   | Example                               |
|------------------------|-------------------------------------------|---------------------------------------|
| Unique Index (UIX)     | `<table>_<column>_uix`                    | `user_name_uix`                       |
| Foreign Key (FK)       | `<origin_table>_<destination_table>_fkey` | `user_user_contract_information_fkey` |
| Primary Key (PK)       | `<table>_pkey`                            | `user_service_pkey`                   |
| Materialized View (MV) | `<view>_uix`                              | `rich_user_ix`                        |

### Database Versioning

Follow [Semantic Versioning](https://semver.org/) for database changelogs:

| Change Type                | Version Update |
|----------------------------|----------------|
| New table                  | `X.Y+1.0`      |
| Alter existing table       | `X.Y.Z+1`      |
| Breaking structural change | `X+1.0.0`      |

Packages published for testing must use the `-SNAPSHOT` suffix (e.g., `3.0.0-SNAPSHOT`). Snapshot packages must **never** be used in production.

### Complete Liquibase Example

```yaml
databaseChangeLog:
  - changeSet:
      id: '00-create-table-example'
      author: 'author.name'
      context: 'prod, test'
      preConditions:
        - not:
            tableExists:
              tableName: 'example'
      changes:
        - createTable:
            tableName: 'example'
            columns:
              - column:
                  autoIncrement: true
                  constraints:
                    nullable: false
                    primaryKey: true
                    primaryKeyName: 'example_pkey'
                  name: 'id'
                  type: 'BIGINT'
              - column:
                  constraints:
                    nullable: false
                    unique: true
                    uniqueConstraintName: 'example_name_uix'
                  name: 'name'
                  type: 'VARCHAR(50)'
              - column:
                  name: 'active'
                  type: 'BOOLEAN'
              - column:
                  constraints:
                    nullable: false
                  name: 'created_at'
                  type: 'TIMESTAMP'
      rollback:
        - dropTable:
            tableName: 'example'
```

## Java Design Principles

The following principles guide Java development across all projects:

- **Program to interfaces, not implementations.** Depend on abstractions (interfaces) rather than concrete classes. This enables testability and loose coupling.
- **Favor composition over inheritance.** Build complex behavior by composing objects rather than through deep class hierarchies.
- **Immutability by default.** Use `final` fields, records, and immutable collections wherever possible. Mutable state should be the exception, not the rule.
- **Fail fast.** Validate inputs early and throw meaningful exceptions rather than propagating invalid state through the system.
- **Convention over configuration.** Leverage Spring Boot's auto-configuration and opinionated defaults rather than writing boilerplate.

## References

- [Spring Boot Reference](https://docs.spring.io/spring-boot/reference/)
- [Effective Java -- Joshua Bloch](https://www.oreilly.com/library/view/effective-java/9780134686097/)
- [Liquibase Documentation](https://docs.liquibase.com/)
- [Semantic Versioning](https://semver.org/)
- [Java Language Specification](https://docs.oracle.com/javase/specs/)

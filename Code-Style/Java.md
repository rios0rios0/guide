# Java Conventions

> **TL;DR:** Follow the `<Operation><Entity>` naming pattern for Controllers and Commands with an `execute` method. Use `listAllSuccessfully` / `listAllError` style for test method names. Use Liquibase with YAML, `snake_case`, mandatory rollbacks, and standardized constraint naming.

## Overview

This document defines Java-specific coding conventions. For the general baseline, refer to the [Code Style](../Code-Style.md) guide.

## Naming Conventions

### Controllers

Class names follow the pattern `<Operation><Entity>Controller.java`. The primary method must be named `execute`.

```java
public class ListUsersController {
    public void execute() {}
}
```

### Commands

Class names follow the pattern `<Operation><Entity>Command.java`. The primary method must be named `execute`. The command must define callbacks for all possible outcomes.

### Testing

#### Integration Tests

Test method names must mirror the scenario being tested. Use the tested method name as a prefix, followed by a descriptor:

| Scenario             | Method Name                    |
|----------------------|--------------------------------|
| Successful listing   | `listAllSuccessfully`          |
| Error during listing | `listAllError`                 |
| Additional edge case | `listAllAdditionalInformation` |

### Seeds

Seed files are used for populating blank tables with test data:

- The seed file name must match the **table name** (e.g., `user.sql`).
- For multiple seeds targeting the same table, use numbered suffixes: `user_01.sql`, `user_02.sql`.
- Constants must be named according to the file name.

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
# ✅ Correct
- column:
    constraints:
        nullable: false
    name: 'email'
    type: 'VARCHAR(255)'

# ❌ Wrong
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

## References

- [Liquibase Documentation](https://docs.liquibase.com/)
- [Semantic Versioning](https://semver.org/)
- [Java Code Conventions](https://www.oracle.com/java/technologies/javase/codeconventions-introduction.html)

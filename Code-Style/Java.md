# Java

## Convention Names

### Controllers
When you're creating a new controller, the class name needs to follow the pattern: `<OPERATION><ENTITY>Controller.java` for the class name.
For example, if you need to create a new controller to list all users, the class name will be `ListUsersController`.
The method inside the class should be called `execute`. Inside the command class, it needs to set the callbacks for the possible results.

**Example:**
```java
public class ListUsersController {
    public void execute() {}
}
```

### Commands

When you're creating a new command, the class name needs to follow the pattern: `<OPERATION><ENTITY>Command.java` for the class name.
For example, if you need to create a new command to list all users, the class name will be `ListUsersCommand`.
The method inside the class should be called `execute`. Inside the command class, it needs to set the callbacks for the possible results.

### Testing

#### Integration Tests

You should use the same name as the method that you are testing.
If you are testing a service called `ListUsersService` and the method you're testing is `listAll`, your test methods should have this structure:
- `listAllSuccessfully`
- `listAllError`
- `listAllAdditionalInformation`

### Seeds
When you are creating a seed, respect this kind of standard: the name of the seed must be the name of the table.
Remember: the seed is not a script. The seed is a data population to be inserted in a blank table. If you have more than one seed, use names like `user_01.sql`, `user_02.sql`, etc.
The constants you are using need to be named according to the file name.

### Liquibase
We are using Liquibase to manage our database versioning. It leads us to follow some patterns to use the Liquibase engine:

- Always use `YAML` file with `a` instead of `YML`.
- Always use single quotes instead of double quotes, except in cases where you need escaping.
- Always use `snake_case` (lowercase separated by an underscore `_`).
- Always write a rollback statement. Never commit a changelog without a rollback.
- Always write a pre-conditions statement to check the existing state of the database before applying changes.
- Boolean and Time columns must have meaningful names like `created_at`, `enabled`, `updated_at`, `activated` instead of `created`, `enable`, `updated`, `activated`.

When you are creating a new table with columns, make sure that the constraints entry appears before the column name.

**Use this:**
```yaml
- column:
    constraints:
        nullable: false
    name: 'email'
    type: 'VARCHAR(255)'
```

**Instead of this:**
```yaml
- column:
    name: 'email'
    type: 'VARCHAR(255)'
    constraints:
        nullable: false
```

### Constraint Names

#### Unique (UIX, Unique Index)

Follow this standard: `<table_name>_<column_name>_uix`
**Example:** `user_name_uix`

#### Foreign Key

Follow this standard: `<origin_table_name>_<destination_table_name>_fkey`
**Example:** `user_user_contract_information_fkey`

#### Primary Key

Follow this standard: `<table_name>_pkey`
**Example:** `user_service_pkey`

#### Materialized View

For unique indexes, follow this standard: `<view_name>_uix`
**Example:** `rich_user_ix`

### Versioning
Remember: you need ALWAYS to follow Semantic versioning. This leads us to adopt some standards:
- If you are creating a NEW TABLE, do:
  - `X = X`
  - `Y = Y+1`
  - `Z = 0`
- If you are altering an EXISTING TABLE, do:
  - `X = X`
  - `Y = Y`
  - `Z = Z+1`
- If you are breaking the structure (like column names and relationships between tables), which requires application breaking changes, do:
  - `X = X + 1`
  - `Y = 0`
  - `Z = 0`

Packages published for testing purposes should have the `-SNAPSHOT` suffix (e.g., `3.0.0-SNAPSHOT`). Packages with the `-SNAPSHOT` suffix SHOULD NOT be used for purposes other than tests.

**Example of a complete Liquibase file:**
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

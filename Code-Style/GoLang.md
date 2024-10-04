## Context
The theory base is described in the page [Base](../Code-Style.md), but here it's very specific for GoLang.

## File Structure
```
|__ main
  |__ domain                (contracts)
    |__ commands              only one exception that doesn't have contract and it's implementation
    |__ entities
    |__ repositories
    |__ app.go
    |__ main.go
    |__ wire.go
    |__ wire_gen.go
  |__ infrastructure        (implementations)
    |__ controllers
      |__ mappers
      |__ requests
      |__ responses
    |__ repositories        (orm) begin with tool prefix and return models from the database or the external tool
      |__ mappers
      |__ models
|__ test
  |__ domain
  |__ infrastructure
```

## General Considerations
1. For this language we're using [Snake Case](https://www.alura.com.br/artigos/convencoes-nomenclatura-camel-pascal-kebab-snake-case) as a convention.
2. We usually use the word `self` to name the reference to the structure attached to the method in question. It words like `this` in other languages.
3. We don't attach the method to the struct when it does not need to change the structure state.
4. What is DTO pattern? Read [here](https://www.baeldung.com/java-dto-pattern) please.

## Entities
They are the core of the application. The logic regarding the properties and fields, should be inside it.
They need to be free of any framework and external tools. DON'T use any tag inside it.

## Commands
* **File Name:** `<operation>_<entity>_command`. Example: `list_users_command.go`
* **Struct Name:** `<Operation><Entity>Command`. Example: `ListUsersCommand`
* **Method Name:** `Execute`. Example: `func (self ListUsersCommand) Execute(listeners ListUsersCommandListeners)`

1. **Note:** when the operation is related to multiple entities, the name is going to be in the plural.
2. **Note:** the operation is the kind of "thing" you're doing. Like: `List`, `Get`, `Delete`, `Insert`, `Update`, `BatchDelete`, `BatchInsert`, `BatchUpdate`, `DeleteAll` and so on.
3. **Note:** be careful with the **listeners**, they need to reflect all possible controllers answers.

## Controllers
* **File Name:** `<operation>_<entity>_controller`. Example: `list_users_controller.go`
* **Struct Name:** `<Operation><Entity>Controller`. Example: `ListUsersController`
* **Method Name:** `Execute`. Example: `func (self ListUsersController) Execute()`

1. **Note:** when the operation is related to multiple entities, the name is going to be in the plural.
2. **Note:** the operation is the kind of "thing" you're doing. Like: `List`, `Get`, `Delete`, `Insert`, `Update`, `BatchDelete`, `BatchInsert`, `BatchUpdate`, `DeleteAll` and so on.

## Services
We don't use this layer for this language.

## Repositories
* **File Name (contract):** `<entity>_repository`. Example: `users_repository.go`
* **Struct Name (contract):** `<Entity>Repository`. Example: `UsersRepository`
* **File Name (implementation):** `<library>_<entity>_repository`. Example: `pgx_users_repository.go`
* **Struct Name (implementation):** `<Library><Entity>Repository`. Example: `PgxUsersRepository`

- **Method Names (contract):**

  These methods are following a logical sequence: "find all", "find one", "filters", "save one", "save all", "delete" and "delete all".
  ```go
  // filter just 1 entity by some field
  func (self UsersRepository) FindByTargetField(targetField any) entities.User
  // filter many entities by some field
  func (self UsersRepository) FindAllByTargetField(targetField any) []entities.User

  // filter with a boolean type (if an entity exists)
  func (self UsersRepository) HasBooleanVerification(targetField any) bool

  // save just 1 entity per time
  func (self UsersRepository) Save(user entities.User)
  // save many entities at the same time
  func (self UsersRepository) SaveAll(users []entities.User)

  // delete just 1 entity by some field
  func (self UsersRepository) DeleteByTargetField(targetField any)
  ```
- **Method Names (implementation):**

  The same signatures before are applied, but we need to change the attached struct to be the implementation, so the change will be like this:
  ```go
  // this is the contract
  func (self UsersRepository)

  // this is the implementation
  func (self PgxUsersRepository)
  ```

1. **Note:** the `TargetField` could be: `Id`, `Name` and others.
2. **Note:** the `BooleanVerification` could be: `UserInGroup`, `UserPermission` and others.

## Mappers
* **File Name:** `<entity>_mapper`. Example: `user_mapper.go`
* **Struct Name:** `<Entity>Mapper`. Example: `UserMapper`

- Inside `repositories` we can have:
  * **File:** `user_mapper.go`
  * **Struct:** `UserMapper`
  * **Method Names:**
    ```go
    // mapping from some infrastructure DTO to an entity
    func (self UserMapper) MapToEntity(infra any) entities.User
    // mapping from many infrastructure DTOs to many entities
    func (self UserMapper) MapToEntities(infra []any) []entities.User

    // mapping from an entity to an external DTO
    func (self UserMapper) MapToExternal(user entities.User) models.External
    // mapping from many entities to many external DTOs
    func (self UserMapper) MapToExternals(users []entities.User) []models.External
    ```

1. **Note:** the `infra` field is usually an external DTO structure, like a database model or an API response.
2. **Note:** the `MapToEntity` and `MapToEntities` method names are static and must be followed as a standard.
3. **Note:** the `External` and `any` aren't static, and they need to be changed according to the library, API or technology used as the external source.

- Inside `controllers` we'll have:

  In this case, the entity vary between the folders `request` and `responses`. For example effect, we can have both sides like the below.
  * **File:** `insert_user_request_mapper.go` or `insert_user_response_mapper.go`
  * **Struct:** `InsertUserRequestMapper` or `InsertUserResponseMapper`
  * **Method Names:** respectively each case above
  ```go
    // mapping from the infrastructure request to an entity
    func (self InsertUserRequestMapper) MapToEntity(request InsertUserRequest) entities.User
    // mapping from many infrastructure requests to many entities
    func (self InsertUserRequestMapper) MapToEntities(requests []InsertUserRequest) []entities.User
    // THERE'S NO inverse case in this type of mapping
  ```
  ```go
    // mapping from an entity to an external response
    func (self InsertUserResponseMapper) MapToResponse(user entities.User) responses.InsertUserResponse
    // mapping from many entities to many external responses
    func (self InsertUserResponseMapper) MapToResponses(users []entities.User) []responses.InsertUserResponse
    // THERE'S NO inverse case in this type of mapping
    ```

1. **Note:** the request and response are DTOs inside the infrastructure using the proper field names and external coupling code.
2. **Note:** the `MapToResponse` and `MapToResponses` method names are static and must be followed as a standard.
3. **Note:** DON'T use `json` tags outside the infrastructure layer, because they are strict to requests and responses.

## Models
They are always inside the infrastructure layer, and they act like a DTO for the answers coming from the external world.
They can be: database models, APIs answers, queue parameters, payloads and so on. They are like entities, but they aren't entities.

To make sure that the models are different from the entities, we adopt the prefix naming standard.
Each model is prefixed in the name by the external tool that is used to communicate with.
Consider the following model struct names as examples: `AwsFile`, `ApiDocument`, `ApiPayload`, `PgxUser`.

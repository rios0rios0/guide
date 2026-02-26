# Java Project Structure

> **TL;DR:** Follow the domain/infrastructure layer separation within each business module. Use [Gradle](https://gradle.org/) for builds and dependency management. Place test files under `src/test/java/` mirroring the source structure. Use module-based organization for multi-domain applications.

## Overview

This page defines the standard directory layout and dependency management practices for all Java projects. The architecture follows the [Backend Design](../../Life-Cycle/Architecture/Backend-Design.md) specification, separating code into `domain` (contracts) and `infrastructure` (implementations) layers within each business module.

## Directory Structure

```
src/main/java/.../
  Startup.java                           application entry point (@SpringBootApplication)
  global/                                cross-cutting concerns
    errors/                                global error handlers
    helpers/                               shared utilities
  items/                                 business module (one per domain concept)
    domain/
      commands/
        InsertItemCommand.java             business logic (write operation)
        ListItemsCommand.java              business logic (read operation)
      entities/
        Item.java                          pure domain entity (no annotations)
        ItemEvent.java                     domain event (record)
        ItemCode.java                      domain enum
      services/
        SendItemService.java               service contract (interface)
      helpers/
        ItemHelper.java                    domain utilities
    infrastructure/
      controllers/
        InsertItemController.java          HTTP endpoint (write)
        ListItemsController.java           HTTP endpoint (read)
        requests/
          InsertItemRequest.java           request DTO (record)
        responses/
          ItemResponse.java                response DTO (record)
        mappers/
          InsertItemRequestMapper.java     request-to-entity mapper
          ItemResponseMapper.java          entity-to-response mapper
      listeners/
        ProcessItemListener.java           Kafka listener
        messages/
          ItemMessage.java                 message DTO (record)
        mappers/
          ItemMessageMapper.java           message-to-entity mapper
      repositories/
        contracts/
          ItemsRepository.java             repository interface
          QueryDslItemsRepository.java     QueryDSL extension
        JpaItemsRepository.java            Spring Data JPA implementation
        models/
          JpaItem.java                     JPA entity model
        helpers/
          ItemQueryBuilder.java            QueryDSL query construction
      services/
        JpaSendItemService.java            service implementation
        mappers/
          ItemMapper.java                  entity-to-model mapper

src/main/resources/
  application.yaml                       main configuration
  application-prod.yaml                  production profile
  application-test.yaml                  test profile
  messages.properties                    i18n (English)
  messages_pt_BR.properties              i18n (Portuguese)
  db/
    changelog/                           Liquibase migration scripts (YAML)
    seeds/                               test data seeders (SQL)
  app/
    quality/
      checkstyle-google-ruleset.xml      Checkstyle configuration
      pmd-custom-ruleset.xml             PMD configuration
    security/
      spotbugs-security-exclude.xml      SpotBugs exclusions
      dependency-check-suppress.xml      OWASP dependency suppression

src/test/java/.../
  items/                                 mirrors source module structure
    domain/
      builders/
        ItemEventBuilder.java            test data builder
      commands/
        InsertItemCommandTest.java       command unit tests
    infrastructure/
      builders/
        JpaItemBuilder.java              JPA entity builder
      controllers/
        InsertItemControllerTest.java    controller unit tests
      repositories/
        JpaItemsRepositoryTest.java      repository integration tests
      services/
        JpaSendItemServiceTest.java      service unit tests
      listeners/
        ProcessItemListenerTest.java     listener tests
      doubles/
        SendItemServiceStub.java         stubs
        DummyItemService.java            dummies
        InMemoryItemsRepository.java     in-memory implementations

src/test/resources/
  db/                                    test-specific migrations
  http/                                  HTTP test request data
  queue/                                 Kafka test message data
```

### Key Directories

| Directory                               | Purpose                                                 |
|-----------------------------------------|---------------------------------------------------------|
| `global/`                               | Cross-cutting concerns (error handling, shared helpers) |
| `<module>/domain/commands/`             | Business logic implementations (CQRS write side)        |
| `<module>/domain/entities/`             | Framework-agnostic domain entities and events           |
| `<module>/domain/services/`             | Service contracts (interfaces)                          |
| `<module>/infrastructure/controllers/`  | HTTP controllers (CQRS read side + write endpoints)     |
| `<module>/infrastructure/listeners/`    | Kafka/queue message listeners                           |
| `<module>/infrastructure/repositories/` | Repository implementations with JPA/QueryDSL            |
| `<module>/infrastructure/services/`     | Service implementations with infrastructure tools       |
| `src/main/resources/db/changelog/`      | Liquibase database migration scripts                    |
| `src/main/resources/app/quality/`       | Code quality tool configurations                        |

## Package Manager: Gradle

All Java projects use [Gradle](https://gradle.org/) for build automation and dependency management.

### Project Setup

```groovy
// settings.gradle
rootProject.name = 'my-application'

// build.gradle
plugins {
    id 'java'
    id 'org.springframework.boot' version '3.3.4'
    id 'io.spring.dependency-management' version '1.1.6'
}

group = 'com.example'
version = '1.0.0'

java {
    sourceCompatibility = JavaVersion.VERSION_21
    targetCompatibility = JavaVersion.VERSION_21
}
```

### Managing Dependencies

```groovy
dependencies {
    // Spring Boot starters
    implementation 'org.springframework.boot:spring-boot-starter-web'
    implementation 'org.springframework.boot:spring-boot-starter-data-jpa'
    implementation 'org.springframework.boot:spring-boot-starter-validation'

    // Database
    runtimeOnly 'org.postgresql:postgresql'
    implementation 'org.liquibase:liquibase-core'

    // Messaging
    implementation 'org.springframework.kafka:spring-kafka'

    // Mapping
    implementation 'org.mapstruct:mapstruct:1.6.2'
    annotationProcessor 'org.mapstruct:mapstruct-processor:1.6.2'

    // Utilities
    compileOnly 'org.projectlombok:lombok'
    annotationProcessor 'org.projectlombok:lombok'

    // Testing
    testImplementation 'org.springframework.boot:spring-boot-starter-test'
    testImplementation 'com.github.javafaker:javafaker:1.0.2'
    testImplementation 'org.awaitility:awaitility:4.2.2'
    testImplementation 'org.testcontainers:junit-jupiter'
}
```

### Gradle Wrapper

Always use the Gradle Wrapper (`gradlew`) to ensure consistent build tool versions across environments:

```bash
# Build the project
./gradlew build

# Run tests
./gradlew test

# Format code
./gradlew spotlessApply

# Run quality checks
./gradlew checkstyleMain pmdMain spotbugsMain
```

The `gradle/wrapper/` directory and `gradlew` / `gradlew.bat` scripts must be committed to version control.

## Build & Distribution

### Building

```bash
# Build the JAR
./gradlew bootJar

# Build without running tests
./gradlew bootJar -x test
```

### Running

```bash
# Run via Gradle
./gradlew bootRun

# Run the compiled JAR
java -jar build/libs/my-application-1.0.0.jar
```

### Docker

Use multi-stage builds to produce minimal container images:

```dockerfile
FROM eclipse-temurin:21-jdk-alpine AS builder
WORKDIR /app
COPY gradlew build.gradle settings.gradle ./
COPY gradle/ gradle/
RUN ./gradlew dependencies --no-daemon
COPY src/ src/
RUN ./gradlew bootJar --no-daemon -x test

FROM eclipse-temurin:21-jre-alpine
COPY --from=builder /app/build/libs/*.jar /app/app.jar
ENTRYPOINT ["java", "-jar", "/app/app.jar"]
```

## Key Configuration Files

| File                | Purpose                                                            |
|---------------------|--------------------------------------------------------------------|
| `build.gradle`      | Build configuration, plugins, dependencies, and quality tool setup |
| `settings.gradle`   | Project name and multi-module settings                             |
| `gradle.properties` | Gradle build properties (JVM args, versions)                       |
| `application.yaml`  | Spring Boot configuration (profiles, datasource, kafka, logging)   |
| `lombok.config`     | Lombok behavior configuration                                      |
| `compose.yaml`      | Docker Compose for production                                      |
| `compose.dev.yaml`  | Docker Compose for local development                               |
| `.editorconfig`     | Editor standardization                                             |

## References

- [Gradle User Guide](https://docs.gradle.org/current/userguide/userguide.html)
- [Spring Boot Gradle Plugin](https://docs.spring.io/spring-boot/gradle-plugin/)
- [Docker Multi-stage Builds](https://docs.docker.com/build/building/multi-stage/)
- [Testcontainers](https://testcontainers.com/)

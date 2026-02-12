# Go Testing Conventions

> **TL;DR:** Use build flags (`//go:build unit` or `//go:build integration`) on every test file. Place test files next to production code with the `_test.go` suffix. Use `stretchr/testify` for suites and assertions. Test packages must be **external** to the production package. All tests must follow the BDD pattern with `// given`, `// when`, `// then` comment blocks.

## Overview

Go discovers test files automatically by scanning for files ending in `_test.go`. This document defines the conventions for organizing and writing tests across all Go projects.

## File Structure

```
main/
  domain/
  infrastructure/
    repositories/
      pgx_users_repository.go
      pgx_users_repository_test.go        <-- placed next to production file
test/
  domain/
    builders/                              <-- test data builders
    doubles/
      repositories/                        <-- stubs, dummies, fakes
    helpers/
  infrastructure/
    doubles/
      repositories/
```

## General Conventions

1. **Build flags are mandatory.** Every test file must start with a build flag specifying the test type:
   ```go
   //go:build unit
   ```
   ```go
   //go:build integration
   ```
2. **External test packages.** The test package must be outside the production code package. For example, if the production code is in `package commands`, the test file must use `package commands_test`.
3. **Testing framework.** Use [`stretchr/testify`](https://github.com/stretchr/testify) for test suites and assertions.
4. **File naming.** Test files use the `_test` suffix (e.g., `pgx_users_repository_test.go`).
5. **File placement.** Test files are placed next to the corresponding production file.
6. **BDD structure.** Every test must use `// given`, `// when`, `// then` comment blocks to separate preconditions, actions, and assertions.

## Command Tests

```go
package myapp_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"myapp/cmd"
	"myapp/mocks"
)

type CommandSuite struct {
	suite.Suite
	repository *mocks.MockRepository
	command    *cmd.CreateUserCommand
}

func (s *CommandSuite) SetupTest() {
	s.repository = &mocks.MockRepository{}
	s.command = &cmd.CreateUserCommand{
		Repo: s.repository,
	}
}

func (s *CommandSuite) TestCreateUserCommand() {
	// given
	s.repository.On("Save", mock.Anything).Return(nil)
	user := &cmd.User{
		Name:  "John Doe",
		Email: "johndoe@example.com",
	}

	// when
	err := s.command.Run(user)

	// then
	s.repository.AssertExpectations(s.T())
	assert.Nil(s.T(), err)
}

func TestCommandSuite(t *testing.T) {
	suite.Run(t, new(CommandSuite))
}
```

The `CommandSuite` struct defines a test suite for the `CreateUserCommand`. The `SetupTest` method initializes the command and its mock repository before each test. The `TestCreateUserCommand` method validates that the `Run` method correctly delegates to the repository and returns no error.

## Controller Tests

```go
package myapp_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"myapp/ctrl"
	"myapp/mocks"
)

type ControllerSuite struct {
	suite.Suite
	command *mocks.MockCommand
	ctrl    *ctrl.UserController
}

func (s *ControllerSuite) SetupTest() {
	s.command = &mocks.MockCommand{}
	s.ctrl = &ctrl.UserController{
		Repo: s.command,
	}
}

func (s *ControllerSuite) TestCreateUser() {
	// given
	s.command.On("Save", mock.Anything).Return(nil)

	// when
	req, _ := http.NewRequest("POST", "/users", nil)
	w := httptest.NewRecorder()
	s.ctrl.Create(w, req)

	// then
	s.command.AssertExpectations(s.T())
	assert.Equal(s.T(), http.StatusCreated, w.Code)
}

func (s *ControllerSuite) TestGetUsers() {
	// given
	s.command.On("Execute").Return([]*ctrl.User{}, nil)

	// when
	req, _ := http.NewRequest("GET", "/users", nil)
	w := httptest.NewRecorder()
	s.ctrl.GetAll(w, req)

	// then
	s.command.AssertExpectations(s.T())
	assert.Equal(s.T(), http.StatusOK, w.Code)
}

func TestControllerSuite(t *testing.T) {
	suite.Run(t, new(ControllerSuite))
}
```

Controller tests use `httptest.NewRecorder()` and `httptest.NewRequest()` to simulate HTTP requests and validate response status codes.

## Service Tests

This layer is **not used** in Go projects.

## Repository Tests (Integration)

```go
package myapp_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"myapp/repo"
)

type RepositorySuite struct {
	suite.Suite
	repo *repo.PgxUsersRepository
}

func (s *RepositorySuite) SetupTest() {
	s.repo = repo.NewPgxUsersRepository()
}

func (s *RepositorySuite) TestSaveUser() {
	// given
	user := &repo.User{
		Name:  "John Doe",
		Email: "johndoe@example.com",
	}

	// when
	err := s.repo.Save(user)

	// then
	assert.NoError(s.T(), err)
	retrievedUser, err := s.repo.FindByEmail(user.Email)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), user, retrievedUser)
}

func (s *RepositorySuite) TestFindAll() {
	// given
	users := []*repo.User{
		{Name: "John Doe", Email: "johndoe@example.com"},
		{Name: "Jane Smith", Email: "janesmith@example.com"},
	}
	for _, user := range users {
		s.repo.Save(user)
	}

	// when
	retrievedUsers, err := s.repo.FindAll()

	// then
	assert.NoError(s.T(), err)
	assert.ElementsMatch(s.T(), users, retrievedUsers)
}

func TestRepositorySuite(t *testing.T) {
	suite.Run(t, new(RepositorySuite))
}
```

Repository tests are **integration tests** that connect to a real test database. Use `SetupTest` and `TearDownTest` to manage database state between test runs.

## References

- [stretchr/testify](https://github.com/stretchr/testify)
- [Go Testing Package](https://pkg.go.dev/testing)
- [Go Build Constraints](https://pkg.go.dev/cmd/go#hdr-Build_constraints)
- [Given-When-Then -- Martin Fowler](https://martinfowler.com/bliki/GivenWhenThen.html)

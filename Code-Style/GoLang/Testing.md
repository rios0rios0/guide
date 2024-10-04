## Context
When you run `go test` command in the root of your application,
It will look for all the files that end with `_test.go` and run the test cases written in them.

This structure makes it easy to keep your tests organized and separated from your main application code.

## File Structure
```
|__ main
  |__ domain
  |__ infrastructure
    |__ repositories
      |__ pgx_users_repository.go
      |__ pgx_users_repository_test.go
|__ test
  |__ domain
    |__ builders
    |__ doubles
      |__ repositories
    |__ helpers
  |__ infrastructure
    |__ doubles
      |__ repositories
```

## General Considerations
1. Always start the test with a build flag specifying what's the test type, like: `//go:build unit`, `//go:build integration` and so on.
2. The test package MUST be used OUTSIDE the production code package. For example, if your production code is `package commands` your test will be in `package commands_test`.
3. We use primarily the `stretchr/testify` to make the test suites and the assertions inside the automated code.
4. The file MUST be named with a final prefix `_test`. Like: `pgx_users_repository_test`.
5. The file in the point 4th is the file which is going to be executed, and should be placed near to the production file.

## Commands
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
	command  *cmd.CreateUserCommand
}

func (s *CommandSuite) SetupTest() {
	s.repository = &mocks.MockRepository{}
	s.command = &cmd.CreateUserCommand{
		Repo: s.repository,
	}
}

func (s *CommandSuite) TestCreateUserCommand() {
	s.repository.On("Save", mock.Anything).Return(nil)

	user := &cmd.User{
		Name:  "John Doe",
		Email: "johndoe@example.com",
	}

	err := s.command.Run(user)

	s.repository.AssertExpectations(s.T())
	assert.Nil(s.T(), err)
}

func TestCommandSuite(t *testing.T) {
	suite.Run(t, new(CommandSuite))
}
```

In this example, the `CommandSuite` struct is used to define a test suite for the `CreateUserCommand` struct, which is part of the application's command layer.
The `repository` field is used to create a mock implementation of the repository that the command uses, and the `SetupTest` method is used to initialize the command and the mock repository before each test.
The `TestCreateUserCommand` method is used to test the `Run` method of the `CreateUserCommand` struct.

The mock repository's `Save` method is set to return `nil` and the `Run` method is called with a user struct as an argument.
Then, the `AssertExpectations` method is called on the mock repository to ensure that the `Save` method was called as expected, and the `Nil` method is used to assert that the `Run` method returned a `nil` error.
Finally, the `TestCommandSuite` function is used to run the test suite using the `testing` package's `Run` method.

## Controllers
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
	ctrl     *ctrl.UserController
}

func (s *ControllerSuite) SetupTest() {
	s.command = &mocks.MockCommand{}
	s.ctrl = &ctrl.UserController{
		Repo: s.command,
	}
}

func (s *ControllerSuite) TestCreateUser() {
	s.command.On("Save", mock.Anything).Return(nil)

	req, _ := http.NewRequest("POST", "/users", nil)
	w := httptest.NewRecorder()

	s.ctrl.Create(w, req)

	s.command.AssertExpectations(s.T())
	assert.Equal(s.T(), http.StatusCreated, w.Code)
}

func (s *ControllerSuite) TestGetUsers() {
	s.command.On("Execute").Return([]*ctrl.User{}, nil)

	req, _ := http.NewRequest("GET", "/users", nil)
	w := httptest.NewRecorder()

	s.ctrl.GetAll(w, req)

	s.command.AssertExpectations(s.T())
	assert.Equal(s.T(), http.StatusOK, w.Code)
}

func TestControllerSuite(t *testing.T) {
	suite.Run(t, new(ControllerSuite))
}
```

This is the same from the previous example, but using the Controller layer.

## Services
We don't use this layer for this language.

## Repositories
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
	user := &repo.User{
		Name:  "John Doe",
		Email: "johndoe@example.com",
	}

	err := s.repo.Save(user)
	assert.NoError(s.T(), err)

	retrievedUser, err := s.repo.FindByEmail(user.Email)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), user, retrievedUser)
}

func (s *RepositorySuite) TestFindByEmail() {
	email := "johndoe@example.com"
	user := &repo.User{
		Name:  "John Doe",
		Email: email,
	}
	s.repo.Save(user)

	retrievedUser, err := s.repo.FindByEmail(email)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), user, retrievedUser)
}

func (s *RepositorySuite) TestFindAll() {
	users := []*repo.User{
		{Name: "John Doe", Email: "johndoe@example.com"},
		{Name: "Jane Smith", Email: "janesmith@example.com"},
	}

	for _, user := range users {
		s.repo.Save(user)
	}

	retrievedUsers, err := s.repo.FindAll()
	assert.NoError(s.T(), err)
	assert.ElementsMatch(s.T(), users, retrievedUsers)
}

func TestRepositorySuite(t *testing.T) {
	suite.Run(t, new(RepositorySuite))
}
```

The same from the two previous examples, but here using an integration test as an example.

## References

* https://github.com/stretchr/testify

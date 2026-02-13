# Go Testing Conventions

> **TL;DR:** Use build flags (`//go:build unit` or `//go:build integration`) on every test file. Place test files next to production code with the `_test.go` suffix. Use `stretchr/testify` for suites and assertions. Test packages must be **external** to the production package. All tests must follow the BDD pattern with `// given`, `// when`, `// then` comment blocks. Unit tests must run in **parallel** using `t.Parallel()` + `t.Run()`. Integration tests use **suites** with setup/teardown and are NOT parallel.

## Overview

Go discovers test files automatically by scanning for files ending in `_test.go`. This document defines the conventions for organizing and writing tests across all Go projects.

## File Structure

```
main/
  domain/
  infrastructure/
    repositories/
      sqlx_items_repository.go
      sqlx_items_repository_test.go       <-- placed next to production file
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
4. **File naming.** Test files use the `_test` suffix (e.g., `sqlx_items_repository_test.go`).
5. **File placement.** Test files are placed next to the corresponding production file.
6. **BDD structure.** Every test must use `// given`, `// when`, `// then` comment blocks to separate preconditions, actions, and assertions.
7. **Parallel unit tests.** All unit tests must call `t.Parallel()` and use `t.Run()` sub-tests. All sub-tests within a function execute in parallel.
8. **Sequential integration tests.** Integration tests use `suite.Suite` with setup/teardown and must NOT be parallel, as they share database state.

## Unit Tests (Parallel with `t.Run`)

Unit tests must be structured for **parallel execution**. Each top-level test function calls `t.Parallel()`, and each scenario is a `t.Run()` sub-test. All sub-tests within the same function run concurrently.

### Command Tests

```go
package commands_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"myapp/domain/commands"
	"myapp/test/domain/doubles"
	domainErrors "myapp/domain/errors"
)

func TestDeleteItemCommand(t *testing.T) {
	t.Parallel()

	t.Run("should call OnSuccess when the item is deleted", func(t *testing.T) {
		// given
		repository := doubles.NewItemRepositoryStub()
		command := commands.NewDeleteItemCommand(repository)
		onSuccessCalled := false

		// when
		listeners := commands.DeleteItemListeners{
			OnSuccess: func() {
				onSuccessCalled = true
			},
		}
		command.Execute(context.TODO(), 1, listeners)

		// then
		assert.True(t, onSuccessCalled, "OnSuccess should have been called")
	})

	t.Run("should call OnNotFound when the item is not found", func(t *testing.T) {
		// given
		repository := doubles.NewItemRepositoryStub().WithOnError(domainErrors.ErrRecordNotFound)
		command := commands.NewDeleteItemCommand(repository)
		onNotFoundCalled := false

		// when
		listeners := commands.DeleteItemListeners{
			OnNotFound: func() {
				onNotFoundCalled = true
			},
		}
		command.Execute(context.TODO(), 1, listeners)

		// then
		require.True(t, onNotFoundCalled, "OnNotFound should have been called")
	})

	t.Run("should call OnError when there is an error processing the delete", func(t *testing.T) {
		// given
		dbProcessErr := errors.New("test error")
		repository := doubles.NewItemRepositoryStub().WithOnError(dbProcessErr)
		command := commands.NewDeleteItemCommand(repository)

		// when & then
		listeners := commands.DeleteItemListeners{
			OnError: func(err error) {
				assert.Error(t, err)
			},
		}
		command.Execute(context.TODO(), 1, listeners)
	})
}
```

**Key points:**
- `t.Parallel()` is called at the top of `TestDeleteItemCommand`, enabling all `t.Run()` sub-tests to execute concurrently.
- Each sub-test is self-contained -- it creates its own doubles, command, and listeners.
- Listeners pattern reflects all possible controller responses: `OnSuccess`, `OnNotFound`, `OnError`.

### Controller Tests

```go
package controllers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"myapp/infrastructure/controllers"
	"myapp/test/domain/doubles"
)

func TestListItemsController(t *testing.T) {
	t.Parallel()

	t.Run("should respond 200 (OK) when items are listed successfully", func(t *testing.T) {
		// given
		command := doubles.NewListItemsCommandStub()
		ctrl := controllers.NewListItemsController(command)

		// when
		req, _ := http.NewRequest("GET", "/items", nil)
		w := httptest.NewRecorder()
		ctrl.Execute(w, req)

		// then
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("should respond 500 (Internal Server Error) when command fails", func(t *testing.T) {
		// given
		command := doubles.NewListItemsCommandStub().WithOnError()
		ctrl := controllers.NewListItemsController(command)

		// when
		req, _ := http.NewRequest("GET", "/items", nil)
		w := httptest.NewRecorder()
		ctrl.Execute(w, req)

		// then
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}
```

## Service Tests

This layer is **not used** in Go projects.

## Integration Tests (Suite with Setup/Teardown)

Integration tests use `suite.Suite` because they require shared infrastructure (databases, external services). They are **NOT parallel** due to shared mutable state. Use `SetupSuite`, `SetupTest`, and `TearDownTest` to manage the lifecycle.

Group sub-tests by outcome or feature using `suite.Run()`.

### Repository Tests

```go
package repositories_test

import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/suite"

	"myapp/domain"
	"myapp/domain/entities"
	domainErrors "myapp/domain/errors"
	"myapp/infrastructure/repositories"
	"myapp/test/domain/builders"
	dbTest "myapp/test/infrastructure/helpers"
)

type ItemsRepositorySuite struct {
	dbTest.DatabaseSuite
}

func (s *ItemsRepositorySuite) getSeeders() []dbTest.TableSeed {
	return []dbTest.TableSeed{
		{
			MigrationFileRelativePathFromRoot: "db/seeders/test/items/item.sql",
			TableName:                         "items",
		},
	}
}

func (s *ItemsRepositorySuite) SetupSuite() {
	s.SetupDatabase()
}

func (s *ItemsRepositorySuite) SetupTest() {
	s.RunSeeders(s.DB, s.getSeeders())
}

func (s *ItemsRepositorySuite) TearDownTest() {
	sqlResp := s.DB.MustExec("DELETE FROM items")
	rows, err := sqlResp.RowsAffected()
	s.Require().NoError(err)
	s.Require().Positive(rows)
}

func (s *ItemsRepositorySuite) createItem() entities.Item {
	return builders.NewItemBuilder().
		WithID(0).
		WithExternalID(gofakeit.UUID()).
		WithPayload().
		Build()
}

func (s *ItemsRepositorySuite) getRandomID() int64 {
	var randomID int64
	query := "SELECT id FROM items LIMIT 1"
	err := s.DB.Get(&randomID, query)
	s.Require().NoError(err, "failed to get random row")
	return randomID
}

func (s *ItemsRepositorySuite) newRepository() *repositories.SQLXItemsRepository {
	return repositories.NewSQLXItemsRepository(s.DB)
}

func (s *ItemsRepositorySuite) newContext() context.Context {
	return context.Background()
}

func (s *ItemsRepositorySuite) newPagination(page, size int) domain.Pagination {
	return builders.NewPaginationBuilder().
		WithPage(page).
		WithSize(size).
		Build()
}

func (s *ItemsRepositorySuite) defaultPagination() domain.Pagination {
	return s.newPagination(1, 100)
}

func (s *ItemsRepositorySuite) newListParameters(
	pagination domain.Pagination,
	filters map[string]any,
	sorting []domain.SortField,
) domain.ListParameters {
	if filters == nil {
		filters = make(map[string]any)
	}
	if sorting == nil {
		sorting = []domain.SortField{}
	}
	return domain.ListParameters{
		Pagination: pagination,
		Filters:    filters,
		Sorting:    sorting,
	}
}

// --- Success Cases ---

func (s *ItemsRepositorySuite) TestItemsSuccess() {
	s.Run("should save a new item successfully", func() {
		// given
		item := s.createItem()
		repository := s.newRepository()
		ctx := s.newContext()

		// when
		created, err := repository.Save(ctx, &item)

		// then
		s.Require().NoError(err)
		s.Require().Positive(created.ID)
	})

	s.Run("should list all items successfully", func() {
		// given
		params := s.newListParameters(s.defaultPagination(), nil, nil)
		repository := s.newRepository()
		ctx := s.newContext()

		// when
		list, err := repository.ListAll(ctx, params)

		// then
		s.Greater(len(list.GetContent()), 1)
		s.NoError(err)
	})

	s.Run("should delete an item successfully", func() {
		// given
		itemID := s.getRandomID()
		repository := s.newRepository()
		ctx := s.newContext()

		// when
		err := repository.Delete(ctx, itemID)

		// then
		s.NoError(err)
	})

	s.Run("should get the target item by ID successfully", func() {
		// given
		itemID := s.getRandomID()
		repository := s.newRepository()
		ctx := s.newContext()

		// when
		item, err := repository.GetByID(ctx, itemID)

		// then
		s.Require().NoError(err)
		s.Require().NotNil(item)
		s.Require().Equal(item.ID, itemID)
	})

	s.Run("should update an item successfully", func() {
		// given
		itemID := s.getRandomID()
		repository := s.newRepository()
		ctx := s.newContext()

		existingItem, err := repository.GetByID(ctx, itemID)
		s.Require().NoError(err)

		originalTitle := existingItem.Title
		existingItem.Title = gofakeit.Word() + "_updated"

		// when
		updated, err := repository.Update(ctx, existingItem)

		// then
		s.Require().NoError(err)
		s.Equal(existingItem.ID, updated.ID)
		s.NotEqual(originalTitle, updated.Title)
	})
}

// --- Error Cases ---

func (s *ItemsRepositorySuite) TestItemsError() {
	s.Run("should return an error when saving with invalid data", func() {
		// given
		item := builders.NewItemBuilder().WithPayload().Build()
		item.ExternalID = "" // invalid
		repository := s.newRepository()

		// when
		_, err := repository.Save(s.newContext(), &item)

		// then
		s.Error(err)
	})

	s.Run("should return an error when listing with invalid pagination", func() {
		// given
		pagination := s.newPagination(-1, -1)
		params := s.newListParameters(pagination, nil, nil)
		repository := s.newRepository()
		ctx := s.newContext()

		// when
		list, err := repository.ListAll(ctx, params)

		// then
		s.Require().Error(err)
		s.Require().Nil(list)
	})

	s.Run("should return an error when deleting a non-existent item", func() {
		// given
		repository := s.newRepository()

		// when
		err := repository.Delete(s.newContext(), -9)

		// then
		s.Error(err)
	})

	s.Run("should return an error when getting an item that does not exist", func() {
		// given
		const itemID int64 = 99999
		repository := s.newRepository()
		ctx := s.newContext()

		// when
		item, err := repository.GetByID(ctx, itemID)

		// then
		s.Require().Error(err)
		s.Require().Nil(item)
	})

	s.Run("should return an error when updating a non-existent item", func() {
		// given
		const nonExistentID int64 = 99999
		item := s.createItem()
		item.ID = nonExistentID
		repository := s.newRepository()
		ctx := s.newContext()

		// when
		_, err := repository.Update(ctx, &item)

		// then
		s.Require().Error(err)
		s.Require().ErrorIs(err, domainErrors.ErrRecordNotFound)
	})
}

// --- Filter Cases ---

func (s *ItemsRepositorySuite) TestListAllWithFilters() {
	pagination := s.defaultPagination()
	repository := s.newRepository()
	ctx := s.newContext()

	s.Run("should filter items by category ID", func() {
		// given
		params := s.newListParameters(pagination, map[string]any{
			"category_id": 1,
		}, nil)

		// when
		list, _ := repository.ListAll(ctx, params)

		// then
		s.Require().NotEmpty(list.GetContent())
		for _, item := range list.GetContent() {
			s.Equal(uint(1), item.Category.ID)
		}
	})

	s.Run("should filter items by severity", func() {
		// given
		params := s.newListParameters(pagination, map[string]any{
			"severity": "high",
		}, nil)

		// when
		list, _ := repository.ListAll(ctx, params)

		// then
		s.Require().NotEmpty(list.GetContent())
		for _, item := range list.GetContent() {
			s.Equal("high", item.Level)
		}
	})

	s.Run("should return empty list when filter matches no items", func() {
		// given
		params := s.newListParameters(pagination, map[string]any{
			"category_id": 999,
		}, nil)

		// when
		list, err := repository.ListAll(ctx, params)

		// then
		s.Require().NoError(err)
		s.Require().Empty(list.GetContent())
	})
}

// --- Sorting Cases ---

func (s *ItemsRepositorySuite) TestListAllWithSorting() {
	pagination := s.defaultPagination()
	repository := s.newRepository()
	ctx := s.newContext()

	s.Run("should sort items by severity ascending", func() {
		// given
		params := s.newListParameters(pagination, nil, []domain.SortField{
			{Field: "severity", Direction: "asc"},
		})

		// when
		list, _ := repository.ListAll(ctx, params)

		// then
		s.Greater(len(list.GetContent()), 1)
		content := list.GetContent()
		for i := 1; i < len(content); i++ {
			s.GreaterOrEqual(content[i].Level, content[i-1].Level)
		}
	})

	s.Run("should sort items by severity descending", func() {
		// given
		params := s.newListParameters(pagination, nil, []domain.SortField{
			{Field: "severity", Direction: "desc"},
		})

		// when
		list, err := repository.ListAll(ctx, params)

		// then
		s.Require().NoError(err)
		content := list.GetContent()
		for i := 1; i < len(content); i++ {
			s.LessOrEqual(content[i].Level, content[i-1].Level)
		}
	})
}

// --- Pagination Cases ---

func (s *ItemsRepositorySuite) TestListAllPagination() {
	s.Run("should paginate items correctly", func() {
		// given
		pagination := s.newPagination(1, 2)
		params := s.newListParameters(pagination, nil, nil)
		repository := s.newRepository()
		ctx := s.newContext()

		// when
		list, _ := repository.ListAll(ctx, params)

		// then
		s.LessOrEqual(len(list.GetContent()), 2)
	})
}

// --- Entry Point ---

func TestSQLXItemsRepository(t *testing.T) {
	suite.Run(t, new(ItemsRepositorySuite))
}
```

**Key points:**
- **No `t.Parallel()` in suites.** Integration tests share a database through `SetupTest`/`TearDownTest`, so they must run sequentially.
- **`suite.Run()`** groups related sub-tests within a single test method (e.g., `TestItemsSuccess`, `TestItemsError`).
- **Helper methods** on the suite (e.g., `createItem()`, `getRandomID()`, `newRepository()`) reduce repetition and keep tests focused on the scenario.
- **Builders** construct test entities with fluent APIs (e.g., `NewItemBuilder().WithID(0).Build()`).
- **Seeders** populate the database before each test, and `TearDownTest` cleans up after each test.
- **Test methods are grouped by concern:** success, error, filters, sorting, pagination.

## References

- [stretchr/testify](https://github.com/stretchr/testify)
- [Go Testing Package](https://pkg.go.dev/testing)
- [Go Build Constraints](https://pkg.go.dev/cmd/go#hdr-Build_constraints)
- [Go Subtests and Sub-benchmarks](https://go.dev/blog/subtests)
- [Given-When-Then -- Martin Fowler](https://martinfowler.com/bliki/GivenWhenThen.html)

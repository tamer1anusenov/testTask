package repository

import (
	"context"
	"database/sql"
	"testing"
	"time"
	"todo-app/app/models"
	"todo-app/internal/testutils"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestPostgresTaskRepository_Create(t *testing.T) {
	// Создаем mock БД
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create sqlmock: %v", err)
	}
	defer db.Close()

	repo := NewPostgresTaskRepository(db)
	ctx := context.Background()

	// Тестовая задача
	task := &models.Task{
		Title:       "Test Task",
		Description: "Test Description",
		Status:      models.TaskStatusActive,
		Priority:    models.PriorityMedium,
	}

	// Ожидаемые данные
	expectedID := 1
	expectedTime := time.Now()

	// Настраиваем mock
	mock.ExpectQuery(`INSERT INTO tasks`).
		WithArgs(task.Title, task.Description, task.Status, task.Priority, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).
			AddRow(expectedID, expectedTime, expectedTime))

	// Выполняем тест
	result, err := repo.Create(ctx, task)

	// Проверяем результат
	testutils.AssertNoError(t, err, "Create should not return error")
	testutils.AssertNotNil(t, result, "Result should not be nil")
	testutils.AssertEqual(t, expectedID, result.ID, "ID should match")
	testutils.AssertEqual(t, task.Title, result.Title, "Title should match")
	testutils.AssertEqual(t, task.Description, result.Description, "Description should match")
	testutils.AssertEqual(t, task.Status, result.Status, "Status should match")
	testutils.AssertEqual(t, task.Priority, result.Priority, "Priority should match")

	// Проверяем, что все expectations выполнены
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestPostgresTaskRepository_Create_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create sqlmock: %v", err)
	}
	defer db.Close()

	repo := NewPostgresTaskRepository(db)
	ctx := context.Background()

	task := &models.Task{
		Title:       "Test Task",
		Description: "Test Description",
		Status:      models.TaskStatusActive,
		Priority:    models.PriorityMedium,
	}

	// Настраиваем mock для возврата ошибки
	mock.ExpectQuery(`INSERT INTO tasks`).
		WithArgs(task.Title, task.Description, task.Status, task.Priority, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnError(sql.ErrConnDone)

	// Выполняем тест
	result, err := repo.Create(ctx, task)

	// Проверяем результат
	testutils.AssertError(t, err, "Create should return error")
	if result != nil {
		t.Errorf("Result should be nil when error occurs")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestPostgresTaskRepository_GetByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create sqlmock: %v", err)
	}
	defer db.Close()

	repo := NewPostgresTaskRepository(db)
	ctx := context.Background()

	// Тестовые данные
	expectedID := 1
	expectedTime := time.Now()
	expectedTask := &models.Task{
		ID:          expectedID,
		Title:       "Test Task",
		Description: "Test Description",
		Status:      models.TaskStatusActive,
		Priority:    models.PriorityHigh,
		CreatedAt:   expectedTime,
		UpdatedAt:   expectedTime,
	}

	// Настраиваем mock
	mock.ExpectQuery(`SELECT (.+) FROM tasks WHERE id = \$1`).
		WithArgs(expectedID).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "title", "description", "status", "priority", "due_date", "created_at", "updated_at", "completed_at",
		}).AddRow(
			expectedTask.ID, expectedTask.Title, expectedTask.Description, expectedTask.Status,
			expectedTask.Priority, nil, expectedTask.CreatedAt, expectedTask.UpdatedAt, nil,
		))

	// Выполняем тест
	result, err := repo.GetByID(ctx, expectedID)

	// Проверяем результат
	testutils.AssertNoError(t, err, "GetByID should not return error")
	testutils.AssertNotNil(t, result, "Result should not be nil")
	testutils.AssertEqual(t, expectedTask.ID, result.ID, "ID should match")
	testutils.AssertEqual(t, expectedTask.Title, result.Title, "Title should match")
	testutils.AssertEqual(t, expectedTask.Description, result.Description, "Description should match")
	testutils.AssertEqual(t, expectedTask.Status, result.Status, "Status should match")
	testutils.AssertEqual(t, expectedTask.Priority, result.Priority, "Priority should match")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestPostgresTaskRepository_GetByID_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create sqlmock: %v", err)
	}
	defer db.Close()

	repo := NewPostgresTaskRepository(db)
	ctx := context.Background()

	expectedID := 999

	// Настраиваем mock для возврата пустого результата
	mock.ExpectQuery(`SELECT (.+) FROM tasks WHERE id = \$1`).
		WithArgs(expectedID).
		WillReturnError(sql.ErrNoRows)

	// Выполняем тест
	result, err := repo.GetByID(ctx, expectedID)

	// Проверяем результат
	testutils.AssertError(t, err, "GetByID should return error for non-existent ID")
	if result != nil {
		t.Errorf("Result should be nil when task not found")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestPostgresTaskRepository_Update(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create sqlmock: %v", err)
	}
	defer db.Close()

	repo := NewPostgresTaskRepository(db)
	ctx := context.Background()

	// Тестовая задача для обновления
	task := &models.Task{
		ID:          1,
		Title:       "Updated Task",
		Description: "Updated Description",
		Status:      models.TaskStatusCompleted,
		Priority:    models.PriorityLow,
	}

	expectedTime := time.Now()

	// Настраиваем mock
	mock.ExpectQuery(`UPDATE tasks SET`).
		WithArgs(task.Title, task.Description, task.Status, task.Priority, sqlmock.AnyArg(), sqlmock.AnyArg(), task.ID).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "title", "description", "status", "priority", "due_date", "created_at", "updated_at", "completed_at",
		}).AddRow(
			task.ID, task.Title, task.Description, task.Status, task.Priority,
			nil, expectedTime, expectedTime, &expectedTime,
		))

	// Выполняем тест
	result, err := repo.Update(ctx, task)

	// Проверяем результат
	testutils.AssertNoError(t, err, "Update should not return error")
	testutils.AssertNotNil(t, result, "Result should not be nil")
	testutils.AssertEqual(t, task.ID, result.ID, "ID should match")
	testutils.AssertEqual(t, task.Title, result.Title, "Title should match")
	testutils.AssertEqual(t, task.Description, result.Description, "Description should match")
	testutils.AssertEqual(t, task.Status, result.Status, "Status should match")
	testutils.AssertEqual(t, task.Priority, result.Priority, "Priority should match")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestPostgresTaskRepository_Delete(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create sqlmock: %v", err)
	}
	defer db.Close()

	repo := NewPostgresTaskRepository(db)
	ctx := context.Background()

	taskID := 1

	// Настраиваем mock
	mock.ExpectExec(`DELETE FROM tasks WHERE id = \$1`).
		WithArgs(taskID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	// Выполняем тест
	err = repo.Delete(ctx, taskID)

	// Проверяем результат
	testutils.AssertNoError(t, err, "Delete should not return error")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestPostgresTaskRepository_Delete_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create sqlmock: %v", err)
	}
	defer db.Close()

	repo := NewPostgresTaskRepository(db)
	ctx := context.Background()

	taskID := 999

	// Настраиваем mock для случая, когда задача не найдена
	mock.ExpectExec(`DELETE FROM tasks WHERE id = \$1`).
		WithArgs(taskID).
		WillReturnResult(sqlmock.NewResult(0, 0))

	// Выполняем тест
	err = repo.Delete(ctx, taskID)

	// Проверяем результат
	testutils.AssertError(t, err, "Delete should return error when task not found")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

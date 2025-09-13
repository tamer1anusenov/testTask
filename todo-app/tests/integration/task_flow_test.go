package integration

import (
	"context"
	"testing"
	"todo-app/app/models"
	"todo-app/internal/testutils"
	"todo-app/tests/internal"
)

func TestTaskFlow_CreateUpdateDeleteFlow(t *testing.T) {
	// Настраиваем тестовый контейнер
	container := internal.SetupTestContainer(t)
	defer container.TeardownTestContainer(t)

	// Очищаем данные
	container.ClearTestData(t)

	// Получаем тестовое приложение
	app := container.GetTestApp()
	ctx := context.Background()

	// === Тест 1: Создание задачи ===
	createReq := models.CreateTaskRequest{
		Title:       "Integration Test Task",
		Description: "This is a test task for integration testing",
		Priority:    models.PriorityHigh,
	}

	createdTask, err := app.TaskUseCase.CreateTask(ctx, createReq)
	testutils.AssertNoError(t, err, "Create task should not return error")
	testutils.AssertNotNil(t, createdTask, "Created task should not be nil")
	testutils.AssertEqual(t, createReq.Title, createdTask.Title, "Title should match")
	testutils.AssertEqual(t, createReq.Description, createdTask.Description, "Description should match")
	testutils.AssertEqual(t, createReq.Priority, createdTask.Priority, "Priority should match")
	testutils.AssertEqual(t, models.TaskStatusActive, createdTask.Status, "Status should be active")
	testutils.AssertNotEqual(t, 0, createdTask.ID, "ID should be set")

	// === Тест 2: Получение задачи по ID ===
	retrievedTask, err := app.TaskUseCase.GetTaskByID(ctx, createdTask.ID)
	testutils.AssertNoError(t, err, "Get task by ID should not return error")
	testutils.AssertNotNil(t, retrievedTask, "Retrieved task should not be nil")
	testutils.AssertEqual(t, createdTask.ID, retrievedTask.ID, "ID should match")
	testutils.AssertEqual(t, createdTask.Title, retrievedTask.Title, "Title should match")

	// === Тест 3: Обновление задачи ===
	updateReq := models.UpdateTaskRequest{
		ID:          createdTask.ID,
		Title:       "Updated Integration Test Task",
		Description: "This task has been updated",
		Priority:    models.PriorityLow,
	}

	updatedTask, err := app.TaskUseCase.UpdateTask(ctx, updateReq)
	testutils.AssertNoError(t, err, "Update task should not return error")
	testutils.AssertNotNil(t, updatedTask, "Updated task should not be nil")
	testutils.AssertEqual(t, updateReq.Title, updatedTask.Title, "Title should be updated")
	testutils.AssertEqual(t, updateReq.Description, updatedTask.Description, "Description should be updated")
	testutils.AssertEqual(t, updateReq.Priority, updatedTask.Priority, "Priority should be updated")

	// === Тест 4: Переключение статуса ===
	toggledTask, err := app.TaskUseCase.ToggleTaskStatus(ctx, createdTask.ID)
	testutils.AssertNoError(t, err, "Toggle task status should not return error")
	testutils.AssertNotNil(t, toggledTask, "Toggled task should not be nil")
	testutils.AssertEqual(t, models.TaskStatusCompleted, toggledTask.Status, "Status should be completed")
	testutils.AssertNotNil(t, toggledTask.CompletedAt, "CompletedAt should be set")

	// Переключаем обратно
	toggledTask, err = app.TaskUseCase.ToggleTaskStatus(ctx, createdTask.ID)
	testutils.AssertNoError(t, err, "Toggle task status back should not return error")
	testutils.AssertEqual(t, models.TaskStatusActive, toggledTask.Status, "Status should be active again")
	if toggledTask.CompletedAt != nil {
		t.Errorf("CompletedAt should be nil for active task")
	}

	// === Тест 5: Удаление задачи ===
	err = app.TaskUseCase.DeleteTask(ctx, createdTask.ID)
	testutils.AssertNoError(t, err, "Delete task should not return error")

	// Проверяем, что задача удалена
	_, err = app.TaskUseCase.GetTaskByID(ctx, createdTask.ID)
	testutils.AssertError(t, err, "Get deleted task should return error")
}

func TestTaskFlow_MultipleTasksFlow(t *testing.T) {
	// Настраиваем тестовый контейнер
	container := internal.SetupTestContainer(t)
	defer container.TeardownTestContainer(t)

	// Очищаем данные
	container.ClearTestData(t)

	// Получаем тестовое приложение
	app := container.GetTestApp()
	ctx := context.Background()

	// Создаем несколько задач
	tasks := []models.CreateTaskRequest{
		{
			Title:       "High Priority Task",
			Description: "Important task",
			Priority:    models.PriorityHigh,
		},
		{
			Title:       "Medium Priority Task",
			Description: "Regular task",
			Priority:    models.PriorityMedium,
		},
		{
			Title:       "Low Priority Task",
			Description: "Can wait",
			Priority:    models.PriorityLow,
		},
	}

	var createdTasks []*models.Task
	for _, taskReq := range tasks {
		task, err := app.TaskUseCase.CreateTask(ctx, taskReq)
		testutils.AssertNoError(t, err, "Create task should not return error")
		createdTasks = append(createdTasks, task)
	}

	testutils.AssertEqual(t, len(tasks), len(createdTasks), "All tasks should be created")

	// Завершаем некоторые задачи
	for i, task := range createdTasks {
		if i%2 == 0 { // Завершаем четные задачи
			_, err := app.TaskUseCase.ToggleTaskStatus(ctx, task.ID)
			testutils.AssertNoError(t, err, "Toggle task status should not return error")
		}
	}

	// Проверяем статистику (если доступна)
	if app.AnalyticsUseCase != nil {
		stats, err := app.AnalyticsUseCase.GetTasksStats(ctx)
		testutils.AssertNoError(t, err, "Get tasks stats should not return error")
		testutils.AssertNotNil(t, stats, "Stats should not be nil")
		testutils.AssertGreaterThan(t, stats.TotalTasks, 0, "Total tasks should be greater than 0")
	}
}

func TestTaskFlow_WithRealDatabase(t *testing.T) {
	// Настраиваем тестовый контейнер с реальной БД
	container := internal.SetupTestContainer(t)
	defer container.TeardownTestContainer(t)

	// Загружаем тестовые данные
	container.LoadTestData(t)

	// Получаем тестовое приложение
	app := container.GetTestApp()
	ctx := context.Background()

	// Создаем новую задачу поверх тестовых данных
	createReq := models.CreateTaskRequest{
		Title:       "Database Integration Test",
		Description: "Testing with real database",
		Priority:    models.PriorityMedium,
		DueDate:     nil,
	}

	createdTask, err := app.TaskUseCase.CreateTask(ctx, createReq)
	testutils.AssertNoError(t, err, "Create task with real DB should not return error")

	// Проверяем, что задача действительно сохранена в БД
	var taskCount int
	query := "SELECT COUNT(*) FROM tasks WHERE title = $1"
	err = container.TestDB.QueryRow(query, createReq.Title).Scan(&taskCount)
	testutils.AssertNoError(t, err, "Query task count should not return error")
	testutils.AssertEqual(t, 1, taskCount, "Task should be saved in database")

	// Обновляем задачу
	updateReq := models.UpdateTaskRequest{
		ID:          createdTask.ID,
		Title:       "Updated Database Integration Test",
		Description: "Updated with real database",
		Priority:    models.PriorityHigh,
	}

	updatedTask, err := app.TaskUseCase.UpdateTask(ctx, updateReq)
	testutils.AssertNoError(t, err, "Update task with real DB should not return error")

	// Проверяем в БД
	var dbTitle string
	query = "SELECT title FROM tasks WHERE id = $1"
	err = container.TestDB.QueryRow(query, updatedTask.ID).Scan(&dbTitle)
	testutils.AssertNoError(t, err, "Query updated task should not return error")
	testutils.AssertEqual(t, updateReq.Title, dbTitle, "Task should be updated in database")

	// Удаляем задачу
	err = app.TaskUseCase.DeleteTask(ctx, createdTask.ID)
	testutils.AssertNoError(t, err, "Delete task with real DB should not return error")

	// Проверяем, что задача удалена из БД
	err = container.TestDB.QueryRow(query, createdTask.ID).Scan(&dbTitle)
	testutils.AssertError(t, err, "Query deleted task should return error")
}

func TestTaskFlow_Concurrency(t *testing.T) {
	// Настраиваем тестовый контейнер
	container := internal.SetupTestContainer(t)
	defer container.TeardownTestContainer(t)

	// Очищаем данные
	container.ClearTestData(t)

	// Получаем тестовое приложение
	app := container.GetTestApp()
	ctx := context.Background()

	// Создаем задачу для конкурентных операций
	createReq := models.CreateTaskRequest{
		Title:       "Concurrency Test Task",
		Description: "Testing concurrent operations",
		Priority:    models.PriorityMedium,
	}

	task, err := app.TaskUseCase.CreateTask(ctx, createReq)
	testutils.AssertNoError(t, err, "Create task should not return error")

	// Запускаем конкурентные операции обновления
	done := make(chan bool, 2)

	// Горутина 1: обновление приоритета
	go func() {
		defer func() { done <- true }()
		updateReq := models.UpdateTaskRequest{
			ID:          task.ID,
			Title:       task.Title,
			Description: task.Description,
			Priority:    models.PriorityHigh,
		}
		_, err := app.TaskUseCase.UpdateTask(ctx, updateReq)
		if err != nil {
			t.Logf("Concurrent update 1 failed: %v", err)
		}
	}()

	// Горутина 2: переключение статуса
	go func() {
		defer func() { done <- true }()
		_, err := app.TaskUseCase.ToggleTaskStatus(ctx, task.ID)
		if err != nil {
			t.Logf("Concurrent toggle failed: %v", err)
		}
	}()

	// Ждем завершения обеих операций
	<-done
	<-done

	// Проверяем, что задача все еще существует и в валидном состоянии
	finalTask, err := app.TaskUseCase.GetTaskByID(ctx, task.ID)
	testutils.AssertNoError(t, err, "Get task after concurrent operations should not return error")
	testutils.AssertNotNil(t, finalTask, "Task should still exist after concurrent operations")
}

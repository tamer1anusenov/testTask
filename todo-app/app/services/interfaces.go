package services

import (
	"context"
	"todo-app/app/models"
)

// TaskService определяет интерфейс для сервиса управления задачами
type TaskService interface {
	// CreateTask создает новую задачу
	CreateTask(ctx context.Context, req models.CreateTaskRequest) (*models.Task, error)

	// GetAllTasks получает список всех задач с применением фильтров и сортировки
	GetAllTasks(ctx context.Context, filter models.TaskFilter, sort models.TaskSort) ([]*models.Task, error)

	// GetTaskByID получает задачу по ID
	GetTaskByID(ctx context.Context, id int) (*models.Task, error)

	// UpdateTask обновляет существующую задачу
	UpdateTask(ctx context.Context, req models.UpdateTaskRequest) (*models.Task, error)

	// DeleteTask удаляет задачу
	DeleteTask(ctx context.Context, id int) error

	// ToggleTaskStatus переключает статус задачи (active/completed)
	ToggleTaskStatus(ctx context.Context, id int) (*models.Task, error)

	// GetDashboardStats получает статистику для дашборда
	GetDashboardStats(ctx context.Context) (*models.DashboardStats, error)
}

// AppServices объединяет все сервисы приложения
type AppServices struct {
	TaskService TaskService
}

package usecases

import (
	"context"
	"todo-app/app/models"
)

// TaskUseCase определяет интерфейс для бизнес-логики управления задачами
type TaskUseCase interface {
	// CreateTask создает новую задачу с валидацией и бизнес-правилами
	CreateTask(ctx context.Context, req models.CreateTaskRequest) (*models.Task, error)

	// UpdateTask обновляет существующую задачу с проверкой существования
	UpdateTask(ctx context.Context, req models.UpdateTaskRequest) (*models.Task, error)

	// DeleteTask удаляет задачу с проверкой прав доступа
	DeleteTask(ctx context.Context, id int) error

	// ToggleTaskStatus переключает статус задачи (active/completed)
	ToggleTaskStatus(ctx context.Context, id int) (*models.Task, error)

	// GetTasks получает список задач с применением фильтров и сортировки
	GetTasks(ctx context.Context, filter models.TaskFilter, sort models.TaskSort) ([]*models.Task, error)

	// GetTaskByID получает задачу по ID с проверкой существования
	GetTaskByID(ctx context.Context, id int) (*models.Task, error)

	// GetTasksWithPagination получает задачи с пагинацией
	GetTasksWithPagination(ctx context.Context, filter models.TaskFilter, sort models.TaskSort, page, limit int) (*models.TaskListResponse, error)
}

// AnalyticsUseCase определяет интерфейс для аналитики и статистики задач
type AnalyticsUseCase interface {
	// GetTasksStats получает общую статистику по задачам
	GetTasksStats(ctx context.Context) (*models.TaskStats, error)

	// GetDashboardStats получает статистику для дашборда
	GetDashboardStats(ctx context.Context) (*models.DashboardStats, error)

	// GetOverdueTasks получает просроченные задачи
	GetOverdueTasks(ctx context.Context) ([]*models.Task, error)

	// GetHighPriorityTasks получает задачи с высоким приоритетом
	GetHighPriorityTasks(ctx context.Context) ([]*models.Task, error)

	// GetCompletionRates получает статистику выполнения задач
	GetCompletionRates(ctx context.Context, period string) (map[string]float64, error)

	// GetTasksByPriority группирует задачи по приоритетам
	GetTasksByPriority(ctx context.Context) (map[models.Priority][]*models.Task, error)
}

// ExportUseCase определяет интерфейс для экспорта данных
type ExportUseCase interface {
	// ExportTasksToCSV экспортирует задачи в формат CSV
	ExportTasksToCSV(ctx context.Context, filter models.TaskFilter) ([]byte, error)

	// ExportTasksToJSON экспортирует задачи в формат JSON
	ExportTasksToJSON(ctx context.Context, filter models.TaskFilter) ([]byte, error)

	// ExportTasksToPDF экспортирует задачи в формат PDF
	ExportTasksToPDF(ctx context.Context, filter models.TaskFilter) ([]byte, error)

	// GetExportableFields возвращает список полей доступных для экспорта
	GetExportableFields() []string
}

// UseCases объединяет все use case интерфейсы
type UseCases struct {
	Task      TaskUseCase
	Analytics AnalyticsUseCase
	Export    ExportUseCase
}

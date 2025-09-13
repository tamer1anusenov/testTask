package repository

import (
	"context"
	"todo-app/app/models"
)

// TaskRepository определяет интерфейс для работы с задачами в хранилище
type TaskRepository interface {
	// Create создает новую задачу и возвращает ее с заполненным ID
	Create(ctx context.Context, task *models.Task) (*models.Task, error)

	// GetAll получает список задач с учетом фильтров и сортировки
	GetAll(ctx context.Context, filter models.TaskFilter, sort models.TaskSort) ([]*models.Task, error)

	// GetByID получает задачу по ID
	GetByID(ctx context.Context, id int) (*models.Task, error)

	// Update обновляет существующую задачу
	Update(ctx context.Context, task *models.Task) (*models.Task, error)

	// Delete удаляет задачу по ID
	Delete(ctx context.Context, id int) error

	// MarkAsCompleted помечает задачу как выполненную
	MarkAsCompleted(ctx context.Context, id int) error

	// MarkAsActive помечает задачу как активную
	MarkAsActive(ctx context.Context, id int) error

	// GetTasksStats получает статистику по задачам
	GetTasksStats(ctx context.Context) (*models.TaskStats, error)

	// GetTasksCount получает количество задач с учетом фильтра
	GetTasksCount(ctx context.Context, filter models.TaskFilter) (int, error)

	// GetRecentTasks получает последние созданные задачи
	GetRecentTasks(ctx context.Context, limit int) ([]*models.Task, error)

	// GetUpcomingTasks получает задачи с ближайшими сроками
	GetUpcomingTasks(ctx context.Context, limit int) ([]*models.Task, error)
}

// SettingsRepository определяет интерфейс для работы с настройками приложения
type SettingsRepository interface {
	// GetSettings получает настройки приложения
	GetSettings(ctx context.Context) (*models.AppSettings, error)

	// UpdateSettings обновляет настройки приложения
	UpdateSettings(ctx context.Context, settings *models.AppSettings) error
}

// Repository объединяет все репозитории
type Repository struct {
	Task     TaskRepository
	Settings SettingsRepository
}

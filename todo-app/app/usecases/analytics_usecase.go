package usecases

import (
	"context"
	"fmt"
	"time"
	"todo-app/app/models"
	"todo-app/app/services"
)

// AnalyticsUseCaseImpl реализует интерфейс AnalyticsUseCase
type AnalyticsUseCaseImpl struct {
	taskService services.TaskService
}

// NewAnalyticsUseCase создает новый экземпляр AnalyticsUseCase
func NewAnalyticsUseCase(taskService services.TaskService) AnalyticsUseCase {
	return &AnalyticsUseCaseImpl{
		taskService: taskService,
	}
}

// GetTasksStats получает общую статистику по задачам
func (uc *AnalyticsUseCaseImpl) GetTasksStats(ctx context.Context) (*models.TaskStats, error) {
	// Получаем дашборд статистику и извлекаем из неё TaskStats
	dashboardStats, err := uc.taskService.GetDashboardStats(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get dashboard stats: %w", err)
	}

	return &dashboardStats.TaskStats, nil
}

// GetDashboardStats получает полную статистику для дашборда
func (uc *AnalyticsUseCaseImpl) GetDashboardStats(ctx context.Context) (*models.DashboardStats, error) {
	// Вызов сервисного слоя
	stats, err := uc.taskService.GetDashboardStats(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get dashboard stats: %w", err)
	}

	// Дополнительная бизнес-логика для аналитики
	// Например, можно добавить вычисление дополнительных метрик

	return stats, nil
}

// GetOverdueTasks получает просроченные задачи
func (uc *AnalyticsUseCaseImpl) GetOverdueTasks(ctx context.Context) ([]*models.Task, error) {
	// Создаем фильтр для получения активных задач
	filter := models.TaskFilter{
		Status:   models.TaskStatusActive,
		DateType: models.DateFilterOverdue,
	}

	// Сортировка по дате выполнения (просроченные раньше идут первыми)
	sort := models.TaskSort{
		Field: models.SortFieldDueDate,
		Order: models.SortOrderAsc,
	}

	// Получаем задачи через сервисный слой
	tasks, err := uc.taskService.GetAllTasks(ctx, filter, sort)
	if err != nil {
		return nil, fmt.Errorf("failed to get overdue tasks: %w", err)
	}

	// Дополнительная фильтрация на уровне бизнес-логики
	overdueTasks := make([]*models.Task, 0)
	now := time.Now()

	for _, task := range tasks {
		if task.DueDate != nil && task.DueDate.Before(now) && task.Status == models.TaskStatusActive {
			overdueTasks = append(overdueTasks, task)
		}
	}

	return overdueTasks, nil
}

// GetHighPriorityTasks получает задачи с высоким приоритетом
func (uc *AnalyticsUseCaseImpl) GetHighPriorityTasks(ctx context.Context) ([]*models.Task, error) {
	// Создаем фильтр для задач с высоким приоритетом
	filter := models.TaskFilter{
		Status:   models.TaskStatusActive, // Только активные задачи
		Priority: models.PriorityHigh,
	}

	// Сортировка по дате создания (новые задачи сначала)
	sort := models.TaskSort{
		Field: models.SortFieldCreatedAt,
		Order: models.SortOrderDesc,
	}

	// Получаем задачи через сервисный слой
	tasks, err := uc.taskService.GetAllTasks(ctx, filter, sort)
	if err != nil {
		return nil, fmt.Errorf("failed to get high priority tasks: %w", err)
	}

	return tasks, nil
}

// GetCompletionRates получает статистику выполнения задач за период
func (uc *AnalyticsUseCaseImpl) GetCompletionRates(ctx context.Context, period string) (map[string]float64, error) {
	rates := make(map[string]float64)

	// Определяем временной диапазон на основе периода
	var fromDate time.Time
	now := time.Now()

	switch period {
	case "week":
		fromDate = now.AddDate(0, 0, -7)
	case "month":
		fromDate = now.AddDate(0, -1, 0)
	case "year":
		fromDate = now.AddDate(-1, 0, 0)
	default:
		fromDate = now.AddDate(0, 0, -30) // По умолчанию за последние 30 дней
	}

	// Получаем все задачи за период
	allFilter := models.TaskFilter{
		DueFrom: &fromDate,
		DueTo:   &now,
	}

	allTasks, err := uc.taskService.GetAllTasks(ctx, allFilter, models.GetDefaultSort())
	if err != nil {
		return nil, fmt.Errorf("failed to get tasks for completion rate: %w", err)
	}

	// Получаем завершенные задачи за период
	completedFilter := models.TaskFilter{
		Status:  models.TaskStatusCompleted,
		DueFrom: &fromDate,
		DueTo:   &now,
	}

	completedTasks, err := uc.taskService.GetAllTasks(ctx, completedFilter, models.GetDefaultSort())
	if err != nil {
		return nil, fmt.Errorf("failed to get completed tasks for completion rate: %w", err)
	}

	// Вычисляем статистику
	totalTasks := len(allTasks)
	completedCount := len(completedTasks)

	if totalTasks > 0 {
		rates["completion_rate"] = float64(completedCount) / float64(totalTasks) * 100
	} else {
		rates["completion_rate"] = 0
	}

	rates["total_tasks"] = float64(totalTasks)
	rates["completed_tasks"] = float64(completedCount)
	rates["active_tasks"] = float64(totalTasks - completedCount)

	// Дополнительные метрики
	overdueTasks, err := uc.GetOverdueTasks(ctx)
	if err == nil {
		rates["overdue_tasks"] = float64(len(overdueTasks))
		if totalTasks > 0 {
			rates["overdue_rate"] = float64(len(overdueTasks)) / float64(totalTasks) * 100
		}
	}

	return rates, nil
}

// GetTasksByPriority группирует активные задачи по приоритетам
func (uc *AnalyticsUseCaseImpl) GetTasksByPriority(ctx context.Context) (map[models.Priority][]*models.Task, error) {
	// Получаем все активные задачи
	filter := models.TaskFilter{
		Status: models.TaskStatusActive,
	}

	sort := models.TaskSort{
		Field: models.SortFieldPriority,
		Order: models.SortOrderDesc, // Высокий приоритет сначала
	}

	tasks, err := uc.taskService.GetAllTasks(ctx, filter, sort)
	if err != nil {
		return nil, fmt.Errorf("failed to get tasks by priority: %w", err)
	}

	// Группируем задачи по приоритетам
	priorityGroups := make(map[models.Priority][]*models.Task)

	// Инициализируем группы для всех приоритетов
	priorityGroups[models.PriorityLow] = make([]*models.Task, 0)
	priorityGroups[models.PriorityMedium] = make([]*models.Task, 0)
	priorityGroups[models.PriorityHigh] = make([]*models.Task, 0)

	// Распределяем задачи по группам
	for _, task := range tasks {
		priorityGroups[task.Priority] = append(priorityGroups[task.Priority], task)
	}

	return priorityGroups, nil
}

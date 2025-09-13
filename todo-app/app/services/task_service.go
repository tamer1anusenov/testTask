package services

import (
	"context"
	"fmt"
	"time"
	"todo-app/app/models"
	"todo-app/app/repository"
	"todo-app/internal/validation"
)

// TaskServiceImpl реализует интерфейс TaskService
type TaskServiceImpl struct {
	repo      repository.TaskRepository
	validator *validation.TaskValidator
}

// NewTaskService создает новый экземпляр сервиса задач
func NewTaskService(repo repository.TaskRepository) TaskService {
	return &TaskServiceImpl{
		repo:      repo,
		validator: validation.NewTaskValidator(),
	}
}

// CreateTask создает новую задачу
func (s *TaskServiceImpl) CreateTask(ctx context.Context, req models.CreateTaskRequest) (*models.Task, error) {
	// Валидация запроса
	if err := s.validator.ValidateCreateTaskRequest(req); err != nil {
		return nil, fmt.Errorf("invalid create task request: %w", err)
	}

	// Установка значений по умолчанию и текущего времени
	now := time.Now()

	task := &models.Task{
		Title:       req.Title,
		Description: req.Description,
		Status:      models.TaskStatusActive, // Новая задача всегда активна
		Priority:    req.Priority,
		DueDate:     req.DueDate,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	// Если приоритет не указан, устанавливаем средний
	if task.Priority == "" {
		task.Priority = models.PriorityMedium
	}

	// Сохранение в репозитории
	createdTask, err := s.repo.Create(ctx, task)
	if err != nil {
		return nil, fmt.Errorf("failed to create task: %w", err)
	}

	return createdTask, nil
}

// GetAllTasks получает список всех задач с применением фильтров и сортировки
func (s *TaskServiceImpl) GetAllTasks(ctx context.Context, filter models.TaskFilter, sort models.TaskSort) ([]*models.Task, error) {
	// Валидация фильтра и сортировки
	if err := s.validator.ValidateTaskFilter(filter); err != nil {
		return nil, fmt.Errorf("invalid task filter: %w", err)
	}

	if err := s.validator.ValidateTaskSort(sort); err != nil {
		return nil, fmt.Errorf("invalid task sort: %w", err)
	}

	// Применение фильтров и сортировки через репозиторий
	tasks, err := s.repo.GetAll(ctx, filter, sort)
	if err != nil {
		return nil, fmt.Errorf("failed to get tasks: %w", err)
	}

	return tasks, nil
}

// GetTaskByID получает задачу по ID
func (s *TaskServiceImpl) GetTaskByID(ctx context.Context, id int) (*models.Task, error) {
	// Валидация ID
	if err := s.validator.ValidateID(id); err != nil {
		return nil, fmt.Errorf("invalid task ID: %w", err)
	}

	// Получение задачи из репозитория
	task, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get task by ID: %w", err)
	}

	return task, nil
}

// UpdateTask обновляет существующую задачу
func (s *TaskServiceImpl) UpdateTask(ctx context.Context, req models.UpdateTaskRequest) (*models.Task, error) {
	// Валидация запроса
	if err := s.validator.ValidateUpdateTaskRequest(req); err != nil {
		return nil, fmt.Errorf("invalid update task request: %w", err)
	}

	// Проверка существования задачи
	existingTask, err := s.repo.GetByID(ctx, req.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to find task for update: %w", err)
	}

	// Обновляем только измененные поля, сохраняя значения системных полей
	existingTask.Title = req.Title
	existingTask.Description = req.Description
	existingTask.Priority = req.Priority
	existingTask.DueDate = req.DueDate
	existingTask.UpdatedAt = time.Now()

	// Сохранение изменений в репозитории
	updatedTask, err := s.repo.Update(ctx, existingTask)
	if err != nil {
		return nil, fmt.Errorf("failed to update task: %w", err)
	}

	return updatedTask, nil
}

// DeleteTask удаляет задачу
func (s *TaskServiceImpl) DeleteTask(ctx context.Context, id int) error {
	// Валидация ID
	if err := s.validator.ValidateID(id); err != nil {
		return fmt.Errorf("invalid task ID: %w", err)
	}

	// Проверка существования задачи
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to find task for deletion: %w", err)
	}

	// Удаление задачи
	err = s.repo.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}

	return nil
}

// ToggleTaskStatus переключает статус задачи между активным и выполненным
func (s *TaskServiceImpl) ToggleTaskStatus(ctx context.Context, id int) (*models.Task, error) {
	// Валидация ID
	if err := s.validator.ValidateID(id); err != nil {
		return nil, fmt.Errorf("invalid task ID: %w", err)
	}

	// Получение текущей задачи
	task, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to find task for status toggle: %w", err)
	}

	// Переключение статуса
	now := time.Now()
	task.UpdatedAt = now

	if task.Status == models.TaskStatusActive {
		// Если задача была активной, делаем её выполненной
		task.Status = models.TaskStatusCompleted
		task.CompletedAt = &now

		// Обновление в репозитории
		err = s.repo.MarkAsCompleted(ctx, id)
	} else {
		// Если задача была выполненной, делаем её активной
		task.Status = models.TaskStatusActive
		task.CompletedAt = nil

		// Обновление в репозитории
		err = s.repo.MarkAsActive(ctx, id)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to toggle task status: %w", err)
	}

	// Получаем обновленную задачу из репозитория
	updatedTask, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get updated task: %w", err)
	}

	return updatedTask, nil
}

// GetDashboardStats получает статистику для дашборда
func (s *TaskServiceImpl) GetDashboardStats(ctx context.Context) (*models.DashboardStats, error) {
	// Получение статистики по задачам
	stats, err := s.repo.GetTasksStats(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get task stats: %w", err)
	}

	// Получение последних созданных задач
	recentTasks, err := s.repo.GetRecentTasks(ctx, 5) // Получаем 5 последних задач
	if err != nil {
		return nil, fmt.Errorf("failed to get recent tasks: %w", err)
	}

	// Получение предстоящих задач
	upcomingTasks, err := s.repo.GetUpcomingTasks(ctx, 5) // Получаем 5 ближайших задач
	if err != nil {
		return nil, fmt.Errorf("failed to get upcoming tasks: %w", err)
	}

	// Подготовка статистики по приоритетам и статусам
	priorityBreakdown := make(map[models.Priority]int)
	statusBreakdown := make(map[models.TaskStatus]int)

	// Базовые статусы всегда должны быть в результате
	statusBreakdown[models.TaskStatusActive] = 0
	statusBreakdown[models.TaskStatusCompleted] = 0

	// Базовые приоритеты всегда должны быть в результате
	priorityBreakdown[models.PriorityLow] = 0
	priorityBreakdown[models.PriorityMedium] = 0
	priorityBreakdown[models.PriorityHigh] = 0

	// Подготовка моделей ответа
	recentTasksResponse := make([]*models.TaskResponse, 0, len(recentTasks))
	upcomingTasksResponse := make([]*models.TaskResponse, 0, len(upcomingTasks))

	// Конвертация моделей Task в TaskResponse для недавних задач
	for _, task := range recentTasks {
		isOverdue := false
		if task.DueDate != nil && task.Status == models.TaskStatusActive && task.DueDate.Before(time.Now()) {
			isOverdue = true
		}

		recentTasksResponse = append(recentTasksResponse, &models.TaskResponse{
			ID:          task.ID,
			Title:       task.Title,
			Description: task.Description,
			Status:      task.Status,
			Priority:    task.Priority,
			DueDate:     task.DueDate,
			CreatedAt:   task.CreatedAt,
			UpdatedAt:   task.UpdatedAt,
			CompletedAt: task.CompletedAt,
			IsOverdue:   isOverdue,
		})
	}

	// Конвертация моделей Task в TaskResponse для предстоящих задач
	for _, task := range upcomingTasks {
		isOverdue := false
		if task.DueDate != nil && task.Status == models.TaskStatusActive && task.DueDate.Before(time.Now()) {
			isOverdue = true
		}

		upcomingTasksResponse = append(upcomingTasksResponse, &models.TaskResponse{
			ID:          task.ID,
			Title:       task.Title,
			Description: task.Description,
			Status:      task.Status,
			Priority:    task.Priority,
			DueDate:     task.DueDate,
			CreatedAt:   task.CreatedAt,
			UpdatedAt:   task.UpdatedAt,
			CompletedAt: task.CompletedAt,
			IsOverdue:   isOverdue,
		})
	}

	// Формирование итогового ответа
	dashboardStats := &models.DashboardStats{
		TaskStats:         *stats,
		PriorityBreakdown: priorityBreakdown,
		StatusBreakdown:   statusBreakdown,
		RecentTasks:       recentTasksResponse,
		UpcomingTasks:     upcomingTasksResponse,
	}

	return dashboardStats, nil
}

package usecases

import (
	"context"
	"fmt"
	"time"
	"todo-app/app/models"
	"todo-app/app/services"
	"todo-app/internal/validation"
)

// TaskUseCaseImpl реализует интерфейс TaskUseCase
type TaskUseCaseImpl struct {
	taskService services.TaskService
	validator   *validation.TaskValidator
}

// NewTaskUseCase создает новый экземпляр TaskUseCase
func NewTaskUseCase(taskService services.TaskService) TaskUseCase {
	return &TaskUseCaseImpl{
		taskService: taskService,
		validator:   validation.NewTaskValidator(),
	}
}

// CreateTask создает новую задачу с применением бизнес-правил
func (uc *TaskUseCaseImpl) CreateTask(ctx context.Context, req models.CreateTaskRequest) (*models.Task, error) {
	// Дополнительная валидация на уровне use case
	if err := uc.validator.ValidateCreateTaskRequest(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Бизнес-правила для создания задачи
	if req.DueDate != nil && req.DueDate.Before(time.Now().AddDate(0, 0, -1)) {
		return nil, fmt.Errorf("due date cannot be in the past")
	}

	// Если приоритет не указан, устанавливаем средний по умолчанию
	if req.Priority == "" {
		req.Priority = models.PriorityMedium
	}

	// Вызов сервисного слоя
	task, err := uc.taskService.CreateTask(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create task: %w", err)
	}

	return task, nil
}

// UpdateTask обновляет существующую задачу с проверками
func (uc *TaskUseCaseImpl) UpdateTask(ctx context.Context, req models.UpdateTaskRequest) (*models.Task, error) {
	// Валидация запроса
	if err := uc.validator.ValidateUpdateTaskRequest(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Проверяем, что задача существует
	existingTask, err := uc.taskService.GetTaskByID(ctx, req.ID)
	if err != nil {
		return nil, fmt.Errorf("task not found: %w", err)
	}

	// Бизнес-правила для обновления
	if req.DueDate != nil && req.DueDate.Before(time.Now().AddDate(0, 0, -1)) {
		return nil, fmt.Errorf("due date cannot be in the past")
	}

	// Если задача уже выполнена, не разрешаем изменять некоторые поля
	if existingTask.Status == models.TaskStatusCompleted {
		// Можно изменить только описание у выполненной задачи
		if req.Title != existingTask.Title || req.Priority != existingTask.Priority {
			return nil, fmt.Errorf("cannot modify title or priority of completed task")
		}
	}

	// Вызов сервисного слоя
	task, err := uc.taskService.UpdateTask(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to update task: %w", err)
	}

	return task, nil
}

// DeleteTask удаляет задачу с проверкой прав
func (uc *TaskUseCaseImpl) DeleteTask(ctx context.Context, id int) error {
	// Валидация ID
	if err := uc.validator.ValidateID(id); err != nil {
		return fmt.Errorf("invalid task ID: %w", err)
	}

	// Проверяем существование задачи
	task, err := uc.taskService.GetTaskByID(ctx, id)
	if err != nil {
		return fmt.Errorf("task not found: %w", err)
	}

	// Бизнес-правило: можно удалять только свои задачи
	// В будущем здесь можно добавить проверку прав доступа
	_ = task // Пока просто убираем предупреждение компилятора

	// Вызов сервисного слоя
	err = uc.taskService.DeleteTask(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}

	return nil
}

// ToggleTaskStatus переключает статус задачи с обновлением временных меток
func (uc *TaskUseCaseImpl) ToggleTaskStatus(ctx context.Context, id int) (*models.Task, error) {
	// Валидация ID
	if err := uc.validator.ValidateID(id); err != nil {
		return nil, fmt.Errorf("invalid task ID: %w", err)
	}

	// Получаем текущую задачу
	task, err := uc.taskService.GetTaskByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("task not found: %w", err)
	}

	// Бизнес-логика переключения статуса
	if task.Status == models.TaskStatusActive {
		// При завершении задачи проверяем, не просрочена ли она
		if task.DueDate != nil && task.DueDate.Before(time.Now()) {
			// Можно добавить специальную логику для просроченных задач
			// Например, отметить как "completed late"
		}
	} else if task.Status == models.TaskStatusCompleted {
		// При возврате в активное состояние очищаем completed_at
		// Это будет обработано в сервисном слое
	}

	// Вызов сервисного слоя
	updatedTask, err := uc.taskService.ToggleTaskStatus(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to toggle task status: %w", err)
	}

	return updatedTask, nil
}

// GetTasks получает список задач с применением бизнес-правил фильтрации
func (uc *TaskUseCaseImpl) GetTasks(ctx context.Context, filter models.TaskFilter, sort models.TaskSort) ([]*models.Task, error) {
	// Валидация фильтра и сортировки
	if err := uc.validator.ValidateTaskFilter(filter); err != nil {
		return nil, fmt.Errorf("invalid filter: %w", err)
	}

	if err := uc.validator.ValidateTaskSort(sort); err != nil {
		return nil, fmt.Errorf("invalid sort: %w", err)
	}

	// Применение бизнес-правил к фильтрам
	// Например, пользователь может видеть только свои задачи
	// В будущем здесь можно добавить фильтрацию по правам доступа

	// Вызов сервисного слоя
	tasks, err := uc.taskService.GetAllTasks(ctx, filter, sort)
	if err != nil {
		return nil, fmt.Errorf("failed to get tasks: %w", err)
	}

	return tasks, nil
}

// GetTaskByID получает задачу по ID с проверкой прав доступа
func (uc *TaskUseCaseImpl) GetTaskByID(ctx context.Context, id int) (*models.Task, error) {
	// Валидация ID
	if err := uc.validator.ValidateID(id); err != nil {
		return nil, fmt.Errorf("invalid task ID: %w", err)
	}

	// Вызов сервисного слоя
	task, err := uc.taskService.GetTaskByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("task not found: %w", err)
	}

	// Бизнес-правило: проверка прав доступа к задаче
	// В будущем здесь можно добавить проверку, может ли пользователь видеть эту задачу

	return task, nil
}

// GetTasksWithPagination получает задачи с пагинацией
func (uc *TaskUseCaseImpl) GetTasksWithPagination(ctx context.Context, filter models.TaskFilter, sort models.TaskSort, page, limit int) (*models.TaskListResponse, error) {
	// Валидация параметров пагинации
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20 // Значение по умолчанию
	}

	// Валидация фильтра и сортировки
	if err := uc.validator.ValidateTaskFilter(filter); err != nil {
		return nil, fmt.Errorf("invalid filter: %w", err)
	}

	if err := uc.validator.ValidateTaskSort(sort); err != nil {
		return nil, fmt.Errorf("invalid sort: %w", err)
	}

	// Получение всех задач (в будущем можно оптимизировать с пагинацией на уровне БД)
	allTasks, err := uc.taskService.GetAllTasks(ctx, filter, sort)
	if err != nil {
		return nil, fmt.Errorf("failed to get tasks: %w", err)
	}

	// Применение пагинации
	totalCount := len(allTasks)
	start := (page - 1) * limit
	end := start + limit

	if start >= totalCount {
		// Если страница выходит за пределы, возвращаем пустой список
		return &models.TaskListResponse{
			Tasks:      []*models.TaskResponse{},
			TotalCount: totalCount,
			Filter:     filter,
			Sort:       sort,
		}, nil
	}

	if end > totalCount {
		end = totalCount
	}

	pagedTasks := allTasks[start:end]

	// Конвертация в TaskResponse
	taskResponses := make([]*models.TaskResponse, 0, len(pagedTasks))
	for _, task := range pagedTasks {
		isOverdue := false
		if task.DueDate != nil && task.Status == models.TaskStatusActive && task.DueDate.Before(time.Now()) {
			isOverdue = true
		}

		taskResponses = append(taskResponses, &models.TaskResponse{
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

	return &models.TaskListResponse{
		Tasks:      taskResponses,
		TotalCount: totalCount,
		Filter:     filter,
		Sort:       sort,
	}, nil
}

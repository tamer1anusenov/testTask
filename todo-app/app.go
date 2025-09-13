package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"
	"todo-app/app/config"
	"todo-app/app/models"
	"todo-app/app/usecases"
	"todo-app/internal/utils"
)

// App struct
type App struct {
	ctx              context.Context
	db               *sql.DB
	config           *config.Config
	logger           *utils.Logger
	TaskUseCase      usecases.TaskUseCase
	AnalyticsUseCase usecases.AnalyticsUseCase
	ExportUseCase    usecases.ExportUseCase
}

// NewApp creates a new App application struct (for backward compatibility)
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods (for backward compatibility)
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	if a.logger != nil {
		a.logger.Info("Application started successfully")
	}
}

// Startup is the public method for Wails
func (a *App) Startup(ctx context.Context) {
	a.startup(ctx)
}

// Shutdown is called when the app is shutting down
func (a *App) Shutdown(ctx context.Context) {
	if a.logger != nil {
		a.logger.Info("Application shutting down")
	}
}

// Greet returns a greeting for the given name (example Wails method)
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

// GetAppInfo возвращает информацию о приложении
func (a *App) GetAppInfo() map[string]interface{} {
	info := map[string]interface{}{
		"name":    "Todo App",
		"version": "1.0.0",
	}

	if a.config != nil {
		info["name"] = a.config.App.Name
		info["version"] = a.config.App.Version
		info["environment"] = a.config.App.Environment
		info["debug"] = a.config.App.Debug
	}

	return info
}

// HealthCheck проверяет состояние приложения (Wails method)
func (a *App) HealthCheck() map[string]interface{} {
	result := map[string]interface{}{
		"status":    "ok",
		"timestamp": fmt.Sprintf("%d", time.Now().Unix()),
	}

	// Проверяем базу данных
	if a.db != nil {
		if err := a.db.Ping(); err != nil {
			result["database"] = "error"
			result["database_error"] = err.Error()
			result["status"] = "error"
		} else {
			result["database"] = "ok"
		}
	}

	return result
}

// === Task Management Methods (Wails bindings) ===

// CreateTask создает новую задачу
func (a *App) CreateTask(title, description string, priorityStr string, deadline string) (*models.Task, error) {
	if a.TaskUseCase == nil {
		return nil, fmt.Errorf("task use case not initialized")
	}

	// Конвертируем строку в Priority
	var priority models.Priority
	switch priorityStr {
	case "low":
		priority = models.PriorityLow
	case "medium":
		priority = models.PriorityMedium
	case "high":
		priority = models.PriorityHigh
	default:
		priority = models.PriorityMedium
	}

	// Парсим дедлайн если он предоставлен
	var dueDate *time.Time
	if deadline != "" {
		parsedDate, err := time.Parse("2006-01-02", deadline)
		if err != nil {
			return nil, fmt.Errorf("invalid deadline format, expected YYYY-MM-DD: %w", err)
		}
		dueDate = &parsedDate
	}

	req := models.CreateTaskRequest{
		Title:       title,
		Description: description,
		Priority:    priority,
		DueDate:     dueDate,
	}

	// Log task creation in the terminal
	fmt.Printf("Task created via frontend: Title='%s', Description='%s', Priority='%s', Deadline='%s'\n", title, description, priorityStr, deadline)

	return a.TaskUseCase.CreateTask(a.ctx, req)
}

// GetAllTasks возвращает все задачи
func (a *App) GetAllTasks() ([]*models.Task, error) {
	if a.TaskUseCase == nil {
		return nil, fmt.Errorf("task use case not initialized")
	}

	// Исключаем архивные задачи по умолчанию
	filter := models.TaskFilter{
		Archived: false,
	}
	sort := models.TaskSort{
		Field: models.SortFieldCreatedAt,
		Order: models.SortOrderDesc,
	}

	return a.TaskUseCase.GetTasks(a.ctx, filter, sort)
}

// GetTasksByStatus возвращает задачи по статусу
func (a *App) GetTasksByStatus(status string) ([]*models.Task, error) {
	if a.TaskUseCase == nil {
		return nil, fmt.Errorf("task use case not initialized")
	}

	var taskStatus models.TaskStatus
	switch status {
	case "active":
		taskStatus = models.TaskStatusActive
	case "completed":
		taskStatus = models.TaskStatusCompleted
	default:
		taskStatus = models.TaskStatusActive
	}

	filter := models.TaskFilter{
		Status:   taskStatus,
		Archived: false, // Исключаем архивные задачи
	}
	sort := models.TaskSort{
		Field: models.SortFieldCreatedAt,
		Order: models.SortOrderDesc,
	}

	return a.TaskUseCase.GetTasks(a.ctx, filter, sort)
}

// UpdateTask обновляет задачу
func (a *App) UpdateTask(id int, title, description, priorityStr string, deadline string) (*models.Task, error) {
	if a.TaskUseCase == nil {
		return nil, fmt.Errorf("task use case not initialized")
	}

	var priority models.Priority
	switch priorityStr {
	case "low":
		priority = models.PriorityLow
	case "medium":
		priority = models.PriorityMedium
	case "high":
		priority = models.PriorityHigh
	default:
		priority = models.PriorityMedium
	}

	// Парсим дедлайн если он предоставлен
	var dueDate *time.Time
	if deadline != "" {
		parsedDate, err := time.Parse("2006-01-02", deadline)
		if err != nil {
			return nil, fmt.Errorf("invalid deadline format, expected YYYY-MM-DD: %w", err)
		}
		dueDate = &parsedDate
	}

	req := models.UpdateTaskRequest{
		ID:          id,
		Title:       title,
		Description: description,
		Priority:    priority,
		DueDate:     dueDate,
	}

	return a.TaskUseCase.UpdateTask(a.ctx, req)
}

// DeleteTask удаляет задачу
func (a *App) DeleteTask(id int) error {
	if a.TaskUseCase == nil {
		return fmt.Errorf("task use case not initialized")
	}

	return a.TaskUseCase.DeleteTask(a.ctx, id)
}

// ToggleTaskStatus переключает статус задачи
func (a *App) ToggleTaskStatus(id int) (*models.Task, error) {
	if a.TaskUseCase == nil {
		return nil, fmt.Errorf("task use case not initialized")
	}

	return a.TaskUseCase.ToggleTaskStatus(a.ctx, id)
}

// GetTaskByID получает задачу по ID
func (a *App) GetTaskByID(id int) (*models.Task, error) {
	if a.TaskUseCase == nil {
		return nil, fmt.Errorf("task use case not initialized")
	}

	return a.TaskUseCase.GetTaskByID(a.ctx, id)
}

// === Analytics Methods ===

// GetTasksStats возвращает статистику по задачам
func (a *App) GetTasksStats() interface{} {
	if a.AnalyticsUseCase == nil {
		return utils.ErrorResponse(fmt.Errorf("analytics use case not initialized"))
	}

	stats, err := a.AnalyticsUseCase.GetTasksStats(a.ctx)
	return utils.WailsResponse(stats, err)
}

// GetDashboardStats возвращает статистику для дашборда
func (a *App) GetDashboardStats() interface{} {
	if a.AnalyticsUseCase == nil {
		return utils.ErrorResponse(fmt.Errorf("analytics use case not initialized"))
	}

	stats, err := a.AnalyticsUseCase.GetDashboardStats(a.ctx)
	return utils.WailsResponse(stats, err)
}

// === Priority-based Methods ===

// GetTasksByPriority возвращает задачи определенного приоритета
func (a *App) GetTasksByPriority(priorityStr string) ([]*models.Task, error) {
	if a.TaskUseCase == nil {
		return nil, fmt.Errorf("task use case not initialized")
	}

	var priority models.Priority
	switch priorityStr {
	case "low":
		priority = models.PriorityLow
	case "medium":
		priority = models.PriorityMedium
	case "high":
		priority = models.PriorityHigh
	default:
		return nil, fmt.Errorf("invalid priority: %s", priorityStr)
	}

	filter := models.TaskFilter{
		Priority: priority,
		Archived: false, // Исключаем архивные задачи
	}
	sort := models.TaskSort{
		Field: models.SortFieldCreatedAt,
		Order: models.SortOrderDesc,
	}

	return a.TaskUseCase.GetTasks(a.ctx, filter, sort)
}

// === Archive Methods ===

// ArchiveTask отправляет задачу в архив
func (a *App) ArchiveTask(id int) (*models.Task, error) {
	if a.TaskUseCase == nil {
		return nil, fmt.Errorf("task use case not initialized")
	}

	// Получаем текущую задачу
	task, err := a.TaskUseCase.GetTaskByID(a.ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get task: %w", err)
	}

	// Проверяем, что задача выполнена
	if task.Status != models.TaskStatusCompleted {
		return nil, fmt.Errorf("task must be completed before archiving")
	}

	// Обновляем задачу с archived = true
	req := models.UpdateTaskRequest{
		ID:          id,
		Title:       task.Title,
		Description: task.Description,
		Priority:    task.Priority,
		DueDate:     task.DueDate,
		Archived:    true,
	}

	return a.TaskUseCase.UpdateTask(a.ctx, req)
}

// GetArchivedTasks возвращает все архивные задачи
func (a *App) GetArchivedTasks() ([]*models.Task, error) {
	if a.TaskUseCase == nil {
		return nil, fmt.Errorf("task use case not initialized")
	}

	filter := models.TaskFilter{
		Archived: true,
	}
	sort := models.TaskSort{
		Field: models.SortFieldUpdatedAt,
		Order: models.SortOrderDesc,
	}

	return a.TaskUseCase.GetTasks(a.ctx, filter, sort)
}

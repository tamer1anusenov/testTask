package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"todo-app/app/models"

	_ "github.com/lib/pq"
)

// postgresTaskRepository реализует TaskRepository для PostgreSQL
type postgresTaskRepository struct {
	db *sql.DB
}

// NewPostgresTaskRepository создает новый PostgreSQL репозиторий для задач
func NewPostgresTaskRepository(db *sql.DB) TaskRepository {
	return &postgresTaskRepository{db: db}
}

// Create создает новую задачу
func (r *postgresTaskRepository) Create(ctx context.Context, task *models.Task) (*models.Task, error) {
	query := `
        INSERT INTO tasks (title, description, status, priority, due_date, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
        RETURNING id, created_at, updated_at`

	now := time.Now()
	task.CreatedAt = now
	task.UpdatedAt = now

	err := r.db.QueryRowContext(ctx, query,
		task.Title,
		task.Description,
		task.Status,
		task.Priority,
		task.DueDate,
		task.CreatedAt,
		task.UpdatedAt,
	).Scan(&task.ID, &task.CreatedAt, &task.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to create task: %w", err)
	}

	return task, nil
}

// GetAll получает список задач с фильтрацией и сортировкой
func (r *postgresTaskRepository) GetAll(ctx context.Context, filter models.TaskFilter, sort models.TaskSort) ([]*models.Task, error) {
	whereClause, args := r.buildWhereClause(filter)
	orderClause := r.buildOrderClause(sort)

	query := fmt.Sprintf(`
        SELECT id, title, description, status, priority, due_date, created_at, updated_at, completed_at
        FROM tasks
        %s
        %s`, whereClause, orderClause)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get tasks: %w", err)
	}
	defer rows.Close()

	var tasks []*models.Task
	for rows.Next() {
		task := &models.Task{}
		err := rows.Scan(
			&task.ID,
			&task.Title,
			&task.Description,
			&task.Status,
			&task.Priority,
			&task.DueDate,
			&task.CreatedAt,
			&task.UpdatedAt,
			&task.CompletedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan task: %w", err)
		}
		tasks = append(tasks, task)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return tasks, nil
}

// GetByID получает задачу по ID
func (r *postgresTaskRepository) GetByID(ctx context.Context, id int) (*models.Task, error) {
	query := `
        SELECT id, title, description, status, priority, due_date, created_at, updated_at, completed_at
        FROM tasks 
        WHERE id = $1`

	task := &models.Task{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&task.Status,
		&task.Priority,
		&task.DueDate,
		&task.CreatedAt,
		&task.UpdatedAt,
		&task.CompletedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("task with id %d not found", id)
		}
		return nil, fmt.Errorf("failed to get task: %w", err)
	}

	return task, nil
}

// Update обновляет задачу
func (r *postgresTaskRepository) Update(ctx context.Context, task *models.Task) (*models.Task, error) {
	query := `
        UPDATE tasks 
        SET title = $2, description = $3, priority = $4, due_date = $5, updated_at = $6
        WHERE id = $1
        RETURNING updated_at`

	task.UpdatedAt = time.Now()

	err := r.db.QueryRowContext(ctx, query,
		task.ID,
		task.Title,
		task.Description,
		task.Priority,
		task.DueDate,
		task.UpdatedAt,
	).Scan(&task.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("task with id %d not found", task.ID)
		}
		return nil, fmt.Errorf("failed to update task: %w", err)
	}

	return task, nil
}

// Delete удаляет задачу
func (r *postgresTaskRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM tasks WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("task with id %d not found", id)
	}

	return nil
}

// MarkAsCompleted помечает задачу как выполненную
func (r *postgresTaskRepository) MarkAsCompleted(ctx context.Context, id int) error {
	query := `
        UPDATE tasks 
        SET status = $2, completed_at = $3, updated_at = $4
        WHERE id = $1`

	now := time.Now()

	result, err := r.db.ExecContext(ctx, query, id, models.TaskStatusCompleted, &now, now)
	if err != nil {
		return fmt.Errorf("failed to mark task as completed: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("task with id %d not found", id)
	}

	return nil
}

// MarkAsActive помечает задачу как активную
func (r *postgresTaskRepository) MarkAsActive(ctx context.Context, id int) error {
	query := `
        UPDATE tasks 
        SET status = $2, completed_at = NULL, updated_at = $3
        WHERE id = $1`

	now := time.Now()

	result, err := r.db.ExecContext(ctx, query, id, models.TaskStatusActive, now)
	if err != nil {
		return fmt.Errorf("failed to mark task as active: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("task with id %d not found", id)
	}

	return nil
}

// GetTasksStats получает статистику по задачам
func (r *postgresTaskRepository) GetTasksStats(ctx context.Context) (*models.TaskStats, error) {
	query := `
        SELECT 
            COUNT(*) as total_tasks,
            COUNT(CASE WHEN status = 'active' THEN 1 END) as active_tasks,
            COUNT(CASE WHEN status = 'completed' THEN 1 END) as completed_tasks,
            COUNT(CASE WHEN status = 'active' AND due_date IS NOT NULL AND due_date < NOW() THEN 1 END) as overdue_tasks,
            COUNT(CASE WHEN status = 'active' AND DATE(due_date) = CURRENT_DATE THEN 1 END) as today_tasks,
            COUNT(CASE WHEN status = 'active' AND due_date BETWEEN NOW() AND NOW() + INTERVAL '7 days' THEN 1 END) as week_tasks
        FROM tasks`

	stats := &models.TaskStats{}
	err := r.db.QueryRowContext(ctx, query).Scan(
		&stats.TotalTasks,
		&stats.ActiveTasks,
		&stats.CompletedTasks,
		&stats.OverdueTasks,
		&stats.TodayTasks,
		&stats.WeekTasks,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get tasks stats: %w", err)
	}

	return stats, nil
}

// GetTasksCount получает количество задач с учетом фильтра
func (r *postgresTaskRepository) GetTasksCount(ctx context.Context, filter models.TaskFilter) (int, error) {
	whereClause, args := r.buildWhereClause(filter)

	query := fmt.Sprintf(`SELECT COUNT(*) FROM tasks %s`, whereClause)

	var count int
	err := r.db.QueryRowContext(ctx, query, args...).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get tasks count: %w", err)
	}

	return count, nil
}

// GetRecentTasks получает последние созданные задачи
func (r *postgresTaskRepository) GetRecentTasks(ctx context.Context, limit int) ([]*models.Task, error) {
	query := `
        SELECT id, title, description, status, priority, due_date, created_at, updated_at, completed_at
        FROM tasks
        ORDER BY created_at DESC
        LIMIT $1`

	rows, err := r.db.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get recent tasks: %w", err)
	}
	defer rows.Close()

	var tasks []*models.Task
	for rows.Next() {
		task := &models.Task{}
		err := rows.Scan(
			&task.ID,
			&task.Title,
			&task.Description,
			&task.Status,
			&task.Priority,
			&task.DueDate,
			&task.CreatedAt,
			&task.UpdatedAt,
			&task.CompletedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan task: %w", err)
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

// GetUpcomingTasks получает задачи с ближайшими сроками
func (r *postgresTaskRepository) GetUpcomingTasks(ctx context.Context, limit int) ([]*models.Task, error) {
	query := `
        SELECT id, title, description, status, priority, due_date, created_at, updated_at, completed_at
        FROM tasks
        WHERE status = 'active' AND due_date IS NOT NULL AND due_date >= NOW()
        ORDER BY due_date ASC
        LIMIT $1`

	rows, err := r.db.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get upcoming tasks: %w", err)
	}
	defer rows.Close()

	var tasks []*models.Task
	for rows.Next() {
		task := &models.Task{}
		err := rows.Scan(
			&task.ID,
			&task.Title,
			&task.Description,
			&task.Status,
			&task.Priority,
			&task.DueDate,
			&task.CreatedAt,
			&task.UpdatedAt,
			&task.CompletedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan task: %w", err)
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

// buildWhereClause строит WHERE условие и возвращает аргументы
func (r *postgresTaskRepository) buildWhereClause(filter models.TaskFilter) (string, []interface{}) {
	var conditions []string
	var args []interface{}
	argIndex := 1

	// Фильтр по статусу
	if filter.Status != "" {
		conditions = append(conditions, fmt.Sprintf("status = $%d", argIndex))
		args = append(args, filter.Status)
		argIndex++
	}

	// Фильтр по приоритету
	if filter.Priority != "" {
		conditions = append(conditions, fmt.Sprintf("priority = $%d", argIndex))
		args = append(args, filter.Priority)
		argIndex++
	}

	// Поиск по тексту
	if filter.Search != "" {
		conditions = append(conditions, fmt.Sprintf("(title ILIKE $%d OR description ILIKE $%d)", argIndex, argIndex))
		searchPattern := "%" + filter.Search + "%"
		args = append(args, searchPattern)
		argIndex++
	}

	// Фильтр по дате
	switch filter.DateType {
	case models.DateFilterToday:
		conditions = append(conditions, "DATE(due_date) = CURRENT_DATE")
	case models.DateFilterWeek:
		conditions = append(conditions, "due_date BETWEEN NOW() AND NOW() + INTERVAL '7 days'")
	case models.DateFilterOverdue:
		conditions = append(conditions, "status = 'active' AND due_date IS NOT NULL AND due_date < NOW()")
	}

	// Диапазон дат
	if filter.DueFrom != nil {
		conditions = append(conditions, fmt.Sprintf("due_date >= $%d", argIndex))
		args = append(args, filter.DueFrom)
		argIndex++
	}

	if filter.DueTo != nil {
		conditions = append(conditions, fmt.Sprintf("due_date <= $%d", argIndex))
		args = append(args, filter.DueTo)
		argIndex++
	}

	if len(conditions) == 0 {
		return "", args
	}

	return "WHERE " + strings.Join(conditions, " AND "), args
}

// buildOrderClause строит ORDER BY условие
func (r *postgresTaskRepository) buildOrderClause(sort models.TaskSort) string {
	var orderField string

	switch sort.Field {
	case models.SortFieldTitle:
		orderField = "title"
	case models.SortFieldPriority:
		orderField = "CASE priority WHEN 'high' THEN 1 WHEN 'medium' THEN 2 WHEN 'low' THEN 3 END"
	case models.SortFieldDueDate:
		orderField = "due_date"
	case models.SortFieldStatus:
		orderField = "status"
	default:
		orderField = "created_at"
	}

	orderDirection := "DESC"
	if sort.Order == models.SortOrderAsc {
		orderDirection = "ASC"
	}

	return fmt.Sprintf("ORDER BY %s %s", orderField, orderDirection)
}

// postgresSettingsRepository реализует SettingsRepository для PostgreSQL
type postgresSettingsRepository struct {
	db *sql.DB
}

// NewPostgresSettingsRepository создает новый PostgreSQL репозиторий для настроек
func NewPostgresSettingsRepository(db *sql.DB) SettingsRepository {
	return &postgresSettingsRepository{db: db}
}

// GetSettings получает настройки приложения
func (r *postgresSettingsRepository) GetSettings(ctx context.Context) (*models.AppSettings, error) {
	query := `
        SELECT theme, language, notifications_on, auto_save 
        FROM app_settings 
        WHERE id = 1`

	settings := &models.AppSettings{}
	err := r.db.QueryRowContext(ctx, query).Scan(
		&settings.Theme,
		&settings.Language,
		&settings.NotificationsOn,
		&settings.AutoSave,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			// Возвращаем настройки по умолчанию
			return &models.AppSettings{
				Theme:           "light",
				Language:        "en",
				NotificationsOn: true,
				AutoSave:        true,
			}, nil
		}
		return nil, fmt.Errorf("failed to get settings: %w", err)
	}

	return settings, nil
}

// UpdateSettings обновляет настройки приложения
func (r *postgresSettingsRepository) UpdateSettings(ctx context.Context, settings *models.AppSettings) error {
	query := `
        INSERT INTO app_settings (id, theme, language, notifications_on, auto_save, updated_at)
        VALUES (1, $1, $2, $3, $4, $5)
        ON CONFLICT (id) 
        DO UPDATE SET 
            theme = EXCLUDED.theme,
            language = EXCLUDED.language,
            notifications_on = EXCLUDED.notifications_on,
            auto_save = EXCLUDED.auto_save,
            updated_at = EXCLUDED.updated_at`

	_, err := r.db.ExecContext(ctx, query,
		settings.Theme,
		settings.Language,
		settings.NotificationsOn,
		settings.AutoSave,
		time.Now(),
	)

	if err != nil {
		return fmt.Errorf("failed to update settings: %w", err)
	}

	return nil
}

// NewRepository создает новый репозиторий со всеми зависимостями
func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		Task:     NewPostgresTaskRepository(db),
		Settings: NewPostgresSettingsRepository(db),
	}
}

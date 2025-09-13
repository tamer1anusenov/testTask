package models

import (
	"database/sql/driver"
	"fmt"
	"time"
)

// Task представляет основную модель задачи
type Task struct {
	ID          int        `json:"id" db:"id"`
	Title       string     `json:"title" db:"title"`
	Description string     `json:"description" db:"description"`
	Status      TaskStatus `json:"status" db:"status"`
	Priority    Priority   `json:"priority" db:"priority"`
	DueDate     *time.Time `json:"due_date" db:"due_date"`
	Archived    bool       `json:"archived" db:"archived"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
	CompletedAt *time.Time `json:"completed_at" db:"completed_at"`
}

// TaskStatus представляет статус задачи
type TaskStatus string

const (
	TaskStatusActive    TaskStatus = "active"
	TaskStatusCompleted TaskStatus = "completed"
)

// Value реализует driver.Valuer для TaskStatus
func (ts TaskStatus) Value() (driver.Value, error) {
	return string(ts), nil
}

// Scan реализует sql.Scanner для TaskStatus
func (ts *TaskStatus) Scan(value interface{}) error {
	if value == nil {
		*ts = TaskStatusActive
		return nil
	}
	switch s := value.(type) {
	case string:
		*ts = TaskStatus(s)
	case []byte:
		*ts = TaskStatus(s)
	default:
		return fmt.Errorf("cannot scan %T into TaskStatus", value)
	}
	return nil
}

// Priority представляет приоритет задачи
type Priority string

const (
	PriorityLow    Priority = "low"
	PriorityMedium Priority = "medium"
	PriorityHigh   Priority = "high"
)

// Value реализует driver.Valuer для Priority
func (p Priority) Value() (driver.Value, error) {
	return string(p), nil
}

// Scan реализует sql.Scanner для Priority
func (p *Priority) Scan(value interface{}) error {
	if value == nil {
		*p = PriorityMedium
		return nil
	}
	switch s := value.(type) {
	case string:
		*p = Priority(s)
	case []byte:
		*p = Priority(s)
	default:
		return fmt.Errorf("cannot scan %T into Priority", value)
	}
	return nil
}

// TaskStats представляет статистику задач
type TaskStats struct {
	TotalTasks     int `json:"total_tasks"`
	ActiveTasks    int `json:"active_tasks"`
	CompletedTasks int `json:"completed_tasks"`
	OverdueTasks   int `json:"overdue_tasks"`
	TodayTasks     int `json:"today_tasks"`
	WeekTasks      int `json:"week_tasks"`
}

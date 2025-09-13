package models

import "time"

// TaskResponse представляет ответ с информацией о задаче
type TaskResponse struct {
	ID          int        `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      TaskStatus `json:"status"`
	Priority    Priority   `json:"priority"`
	DueDate     *time.Time `json:"due_date"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	CompletedAt *time.Time `json:"completed_at"`
	IsOverdue   bool       `json:"is_overdue"`
}

// TaskListResponse представляет ответ со списком задач
type TaskListResponse struct {
	Tasks      []*TaskResponse `json:"tasks"`
	TotalCount int             `json:"total_count"`
	Filter     TaskFilter      `json:"filter"`
	Sort       TaskSort        `json:"sort"`
}

// DashboardStats представляет статистику для дашборда
type DashboardStats struct {
	TaskStats         TaskStats          `json:"task_stats"`
	PriorityBreakdown map[Priority]int   `json:"priority_breakdown"`
	StatusBreakdown   map[TaskStatus]int `json:"status_breakdown"`
	RecentTasks       []*TaskResponse    `json:"recent_tasks"`
	UpcomingTasks     []*TaskResponse    `json:"upcoming_tasks"`
}

// AppSettings представляет настройки приложения
type AppSettings struct {
	Theme           string `json:"theme" validate:"oneof=light dark"`
	Language        string `json:"language" validate:"oneof=en ru"`
	NotificationsOn bool   `json:"notifications_on"`
	AutoSave        bool   `json:"auto_save"`
}

// ErrorResponse представляет ответ с ошибкой
type ErrorResponse struct {
	Error   string `json:"error"`
	Code    string `json:"code"`
	Details string `json:"details,omitempty"`
}

// SuccessResponse представляет успешный ответ
type SuccessResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

package models

import "time"

// CreateTaskRequest представляет запрос на создание задачи
type CreateTaskRequest struct {
	Title       string     `json:"title" validate:"required,min=1,max=255"`
	Description string     `json:"description" validate:"max=1000"`
	Priority    Priority   `json:"priority" validate:"oneof=low medium high"`
	DueDate     *time.Time `json:"due_date"`
}

// UpdateTaskRequest представляет запрос на обновление задачи
type UpdateTaskRequest struct {
	ID          int        `json:"id" validate:"required,gt=0"`
	Title       string     `json:"title" validate:"required,min=1,max=255"`
	Description string     `json:"description" validate:"max=1000"`
	Priority    Priority   `json:"priority" validate:"oneof=low medium high"`
	DueDate     *time.Time `json:"due_date"`
	Archived    bool       `json:"archived"`
}

// ToggleTaskStatusRequest представляет запрос на изменение статуса задачи
type ToggleTaskStatusRequest struct {
	ID int `json:"id" validate:"required,gt=0"`
}

// DeleteTaskRequest представляет запрос на удаление задачи
type DeleteTaskRequest struct {
	ID int `json:"id" validate:"required,gt=0"`
}

// GetTaskRequest представляет запрос на получение задачи по ID
type GetTaskRequest struct {
	ID int `json:"id" validate:"required,gt=0"`
}

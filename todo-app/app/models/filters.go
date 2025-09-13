package models

import "time"

// TaskFilter представляет фильтры для поиска задач
type TaskFilter struct {
	Status   TaskStatus `json:"status"`    // all, active, completed
	Priority Priority   `json:"priority"`  // all, low, medium, high
	DateType DateFilter `json:"date_type"` // all, today, week, overdue
	Search   string     `json:"search"`    // поиск по заголовку и описанию
	DueFrom  *time.Time `json:"due_from"`  // задачи с даты
	DueTo    *time.Time `json:"due_to"`    // задачи до даты
	Archived bool       `json:"archived"`  // показывать архивные задачи
}

// DateFilter представляет типы фильтрации по дате
type DateFilter string

const (
	DateFilterAll     DateFilter = "all"
	DateFilterToday   DateFilter = "today"
	DateFilterWeek    DateFilter = "week"
	DateFilterOverdue DateFilter = "overdue"
)

// TaskSort представляет параметры сортировки задач
type TaskSort struct {
	Field SortField `json:"field"` // created_at, priority, due_date, title
	Order SortOrder `json:"order"` // asc, desc
}

// SortField представляет поле для сортировки
type SortField string

const (
	SortFieldCreatedAt SortField = "created_at"
	SortFieldUpdatedAt SortField = "updated_at"
	SortFieldPriority  SortField = "priority"
	SortFieldDueDate   SortField = "due_date"
	SortFieldTitle     SortField = "title"
	SortFieldStatus    SortField = "status"
)

// SortOrder представляет направление сортировки
type SortOrder string

const (
	SortOrderAsc  SortOrder = "asc"
	SortOrderDesc SortOrder = "desc"
)

// GetDefaultFilter возвращает фильтр по умолчанию
func GetDefaultFilter() TaskFilter {
	return TaskFilter{
		Status:   TaskStatusActive,
		Priority: "", // все приоритеты
		DateType: DateFilterAll,
		Search:   "",
	}
}

// GetDefaultSort возвращает сортировку по умолчанию
func GetDefaultSort() TaskSort {
	return TaskSort{
		Field: SortFieldCreatedAt,
		Order: SortOrderDesc,
	}
}

// IsValidStatus проверяет валидность статуса
func IsValidStatus(status string) bool {
	return status == string(TaskStatusActive) ||
		status == string(TaskStatusCompleted) ||
		status == ""
}

// IsValidPriority проверяет валидность приоритета
func IsValidPriority(priority string) bool {
	return priority == string(PriorityLow) ||
		priority == string(PriorityMedium) ||
		priority == string(PriorityHigh) ||
		priority == ""
}

// IsValidDateFilter проверяет валидность фильтра по дате
func IsValidDateFilter(dateType string) bool {
	return dateType == string(DateFilterAll) ||
		dateType == string(DateFilterToday) ||
		dateType == string(DateFilterWeek) ||
		dateType == string(DateFilterOverdue)
}

// IsValidSortField проверяет валидность поля сортировки
func IsValidSortField(field string) bool {
	return field == string(SortFieldCreatedAt) ||
		field == string(SortFieldUpdatedAt) ||
		field == string(SortFieldPriority) ||
		field == string(SortFieldDueDate) ||
		field == string(SortFieldTitle) ||
		field == string(SortFieldStatus)
}

// IsValidSortOrder проверяет валидность направления сортировки
func IsValidSortOrder(order string) bool {
	return order == string(SortOrderAsc) ||
		order == string(SortOrderDesc)
}

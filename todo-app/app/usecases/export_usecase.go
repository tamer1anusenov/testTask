package usecases

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"strings"
	"time"
	"todo-app/app/models"
	"todo-app/app/services"
)

// ExportUseCaseImpl реализует интерфейс ExportUseCase
type ExportUseCaseImpl struct {
	taskService services.TaskService
}

// NewExportUseCase создает новый экземпляр ExportUseCase
func NewExportUseCase(taskService services.TaskService) ExportUseCase {
	return &ExportUseCaseImpl{
		taskService: taskService,
	}
}

// ExportTasksToCSV экспортирует задачи в формат CSV
func (uc *ExportUseCaseImpl) ExportTasksToCSV(ctx context.Context, filter models.TaskFilter) ([]byte, error) {
	// Получаем задачи с учетом фильтра
	sort := models.TaskSort{
		Field: models.SortFieldCreatedAt,
		Order: models.SortOrderDesc,
	}

	tasks, err := uc.taskService.GetAllTasks(ctx, filter, sort)
	if err != nil {
		return nil, fmt.Errorf("failed to get tasks for CSV export: %w", err)
	}

	// Создаем CSV буфер
	var csvContent strings.Builder
	writer := csv.NewWriter(&csvContent)

	// Записываем заголовки
	headers := []string{
		"ID",
		"Title",
		"Description",
		"Status",
		"Priority",
		"Due Date",
		"Created At",
		"Updated At",
		"Completed At",
		"Is Overdue",
	}

	if err := writer.Write(headers); err != nil {
		return nil, fmt.Errorf("failed to write CSV headers: %w", err)
	}

	// Записываем данные задач
	for _, task := range tasks {
		record := uc.taskToCSVRecord(task)
		if err := writer.Write(record); err != nil {
			return nil, fmt.Errorf("failed to write CSV record: %w", err)
		}
	}

	// Финализируем запись
	writer.Flush()
	if err := writer.Error(); err != nil {
		return nil, fmt.Errorf("CSV writer error: %w", err)
	}

	return []byte(csvContent.String()), nil
}

// ExportTasksToJSON экспортирует задачи в формат JSON
func (uc *ExportUseCaseImpl) ExportTasksToJSON(ctx context.Context, filter models.TaskFilter) ([]byte, error) {
	// Получаем задачи с учетом фильтра
	sort := models.TaskSort{
		Field: models.SortFieldCreatedAt,
		Order: models.SortOrderDesc,
	}

	tasks, err := uc.taskService.GetAllTasks(ctx, filter, sort)
	if err != nil {
		return nil, fmt.Errorf("failed to get tasks for JSON export: %w", err)
	}

	// Конвертируем в структуру для экспорта
	exportData := struct {
		ExportedAt time.Time         `json:"exported_at"`
		Filter     models.TaskFilter `json:"filter"`
		Count      int               `json:"count"`
		Tasks      []*models.Task    `json:"tasks"`
	}{
		ExportedAt: time.Now(),
		Filter:     filter,
		Count:      len(tasks),
		Tasks:      tasks,
	}

	// Сериализуем в JSON
	jsonData, err := json.MarshalIndent(exportData, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal tasks to JSON: %w", err)
	}

	return jsonData, nil
}

// ExportTasksToPDF экспортирует задачи в формат PDF
func (uc *ExportUseCaseImpl) ExportTasksToPDF(ctx context.Context, filter models.TaskFilter) ([]byte, error) {
	// Получаем задачи с учетом фильтра
	sort := models.TaskSort{
		Field: models.SortFieldCreatedAt,
		Order: models.SortOrderDesc,
	}

	tasks, err := uc.taskService.GetAllTasks(ctx, filter, sort)
	if err != nil {
		return nil, fmt.Errorf("failed to get tasks for PDF export: %w", err)
	}

	// Простая HTML-генерация для PDF (в реальном проекте лучше использовать специальную библиотеку)
	htmlContent := uc.generateHTMLReport(tasks, filter)

	// В реальном проекте здесь должна быть конвертация HTML в PDF
	// Например, с использованием библиотеки wkhtmltopdf или chromedp
	// Пока возвращаем HTML как байты
	return []byte(htmlContent), fmt.Errorf("PDF generation not implemented yet - returning HTML content")
}

// GetExportableFields возвращает список полей доступных для экспорта
func (uc *ExportUseCaseImpl) GetExportableFields() []string {
	return []string{
		"id",
		"title",
		"description",
		"status",
		"priority",
		"due_date",
		"created_at",
		"updated_at",
		"completed_at",
		"is_overdue",
	}
}

// taskToCSVRecord конвертирует задачу в запись CSV
func (uc *ExportUseCaseImpl) taskToCSVRecord(task *models.Task) []string {
	// Определяем, просрочена ли задача
	isOverdue := "false"
	if task.DueDate != nil && task.Status == models.TaskStatusActive && task.DueDate.Before(time.Now()) {
		isOverdue = "true"
	}

	// Форматируем даты
	dueDateStr := ""
	if task.DueDate != nil {
		dueDateStr = task.DueDate.Format("2006-01-02 15:04:05")
	}

	completedAtStr := ""
	if task.CompletedAt != nil {
		completedAtStr = task.CompletedAt.Format("2006-01-02 15:04:05")
	}

	return []string{
		fmt.Sprintf("%d", task.ID),
		task.Title,
		task.Description,
		string(task.Status),
		string(task.Priority),
		dueDateStr,
		task.CreatedAt.Format("2006-01-02 15:04:05"),
		task.UpdatedAt.Format("2006-01-02 15:04:05"),
		completedAtStr,
		isOverdue,
	}
}

// generateHTMLReport генерирует HTML отчет для PDF
func (uc *ExportUseCaseImpl) generateHTMLReport(tasks []*models.Task, filter models.TaskFilter) string {
	var html strings.Builder

	html.WriteString(`<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Tasks Report</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 20px; }
        h1 { color: #333; }
        .header { margin-bottom: 20px; }
        .filter-info { background: #f5f5f5; padding: 10px; margin-bottom: 20px; }
        table { width: 100%; border-collapse: collapse; }
        th, td { border: 1px solid #ddd; padding: 8px; text-align: left; }
        th { background-color: #f2f2f2; }
        .status-active { color: #007bff; }
        .status-completed { color: #28a745; }
        .priority-high { color: #dc3545; font-weight: bold; }
        .priority-medium { color: #ffc107; }
        .priority-low { color: #6c757d; }
        .overdue { background-color: #ffebee; }
    </style>
</head>
<body>`)

	// Заголовок отчета
	html.WriteString(fmt.Sprintf(`
    <div class="header">
        <h1>Tasks Report</h1>
        <p>Generated on: %s</p>
        <p>Total tasks: %d</p>
    </div>`, time.Now().Format("2006-01-02 15:04:05"), len(tasks)))

	// Информация о фильтре
	html.WriteString(`<div class="filter-info">`)
	html.WriteString(`<h3>Applied Filters:</h3>`)
	if filter.Status != "" {
		html.WriteString(fmt.Sprintf(`<p>Status: %s</p>`, filter.Status))
	}
	if filter.Priority != "" {
		html.WriteString(fmt.Sprintf(`<p>Priority: %s</p>`, filter.Priority))
	}
	if filter.Search != "" {
		html.WriteString(fmt.Sprintf(`<p>Search: %s</p>`, filter.Search))
	}
	html.WriteString(`</div>`)

	// Таблица задач
	html.WriteString(`
    <table>
        <thead>
            <tr>
                <th>ID</th>
                <th>Title</th>
                <th>Status</th>
                <th>Priority</th>
                <th>Due Date</th>
                <th>Created</th>
            </tr>
        </thead>
        <tbody>`)

	for _, task := range tasks {
		statusClass := "status-" + strings.ToLower(string(task.Status))
		priorityClass := "priority-" + strings.ToLower(string(task.Priority))

		rowClass := ""
		if task.DueDate != nil && task.Status == models.TaskStatusActive && task.DueDate.Before(time.Now()) {
			rowClass = "overdue"
		}

		dueDateStr := ""
		if task.DueDate != nil {
			dueDateStr = task.DueDate.Format("2006-01-02")
		}

		html.WriteString(fmt.Sprintf(`
            <tr class="%s">
                <td>%d</td>
                <td>%s</td>
                <td class="%s">%s</td>
                <td class="%s">%s</td>
                <td>%s</td>
                <td>%s</td>
            </tr>`,
			rowClass,
			task.ID,
			task.Title,
			statusClass, task.Status,
			priorityClass, task.Priority,
			dueDateStr,
			task.CreatedAt.Format("2006-01-02")))
	}

	html.WriteString(`
        </tbody>
    </table>
</body>
</html>`)

	return html.String()
}

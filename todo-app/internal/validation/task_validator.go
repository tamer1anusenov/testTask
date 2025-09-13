package validation

import (
	"fmt"
	"time"

	"todo-app/app/models"

	"github.com/go-playground/validator/v10"
)

// TaskValidator представляет валидатор для задач
type TaskValidator struct {
	validator *validator.Validate
}

// NewTaskValidator создает новый валидатор задач
func NewTaskValidator() *TaskValidator {
	v := validator.New()

	// Регистрация кастомных валидаторов
	v.RegisterValidation("priority", validatePriority)
	v.RegisterValidation("status", validateStatus)
	v.RegisterValidation("future_date", validateFutureDate)

	return &TaskValidator{
		validator: v,
	}
}

// ValidateCreateTaskRequest валидирует запрос создания задачи
func (tv *TaskValidator) ValidateCreateTaskRequest(req models.CreateTaskRequest) error {
	if err := tv.validator.Struct(req); err != nil {
		return formatValidationError(err)
	}

	// Дополнительные проверки
	if req.DueDate != nil && req.DueDate.Before(time.Now()) {
		return fmt.Errorf("дата выполнения не может быть в прошлом")
	}

	return nil
}

// ValidateUpdateTaskRequest валидирует запрос обновления задачи
func (tv *TaskValidator) ValidateUpdateTaskRequest(req models.UpdateTaskRequest) error {
	if err := tv.validator.Struct(req); err != nil {
		return formatValidationError(err)
	}

	// Дополнительные проверки
	if req.DueDate != nil && req.DueDate.Before(time.Now()) {
		return fmt.Errorf("дата выполнения не может быть в прошлом")
	}

	return nil
}

// ValidateTaskFilter валидирует фильтр задач
func (tv *TaskValidator) ValidateTaskFilter(filter models.TaskFilter) error {
	if filter.Status != "" && !models.IsValidStatus(string(filter.Status)) {
		return fmt.Errorf("некорректный статус: %s", filter.Status)
	}

	if filter.Priority != "" && !models.IsValidPriority(string(filter.Priority)) {
		return fmt.Errorf("некорректный приоритет: %s", filter.Priority)
	}

	if !models.IsValidDateFilter(string(filter.DateType)) {
		return fmt.Errorf("некорректный фильтр по дате: %s", filter.DateType)
	}

	// Проверка диапазона дат
	if filter.DueFrom != nil && filter.DueTo != nil {
		if filter.DueFrom.After(*filter.DueTo) {
			return fmt.Errorf("дата 'с' не может быть позже даты 'по'")
		}
	}

	return nil
}

// ValidateTaskSort валидирует параметры сортировки
func (tv *TaskValidator) ValidateTaskSort(sort models.TaskSort) error {
	if !models.IsValidSortField(string(sort.Field)) {
		return fmt.Errorf("некорректное поле сортировки: %s", sort.Field)
	}

	if !models.IsValidSortOrder(string(sort.Order)) {
		return fmt.Errorf("некорректный порядок сортировки: %s", sort.Order)
	}

	return nil
}

// ValidateID валидирует ID задачи
func (tv *TaskValidator) ValidateID(id int) error {
	if id <= 0 {
		return fmt.Errorf("ID должен быть положительным числом")
	}
	return nil
}

// Кастомные валидаторы
func validatePriority(fl validator.FieldLevel) bool {
	priority := fl.Field().String()
	return models.IsValidPriority(priority)
}

func validateStatus(fl validator.FieldLevel) bool {
	status := fl.Field().String()
	return models.IsValidStatus(status)
}

func validateFutureDate(fl validator.FieldLevel) bool {
	if fl.Field().IsNil() {
		return true // nil допустим
	}

	date := fl.Field().Interface().(*time.Time)
	if date == nil {
		return true
	}

	return date.After(time.Now())
}

// formatValidationError форматирует ошибки валидации
func formatValidationError(err error) error {
	validationErrors := err.(validator.ValidationErrors)

	for _, fieldError := range validationErrors {
		switch fieldError.Tag() {
		case "required":
			return fmt.Errorf("поле '%s' обязательно для заполнения", fieldError.Field())
		case "min":
			return fmt.Errorf("поле '%s' должно содержать минимум %s символов", fieldError.Field(), fieldError.Param())
		case "max":
			return fmt.Errorf("поле '%s' должно содержать максимум %s символов", fieldError.Field(), fieldError.Param())
		case "gt":
			return fmt.Errorf("поле '%s' должно быть больше %s", fieldError.Field(), fieldError.Param())
		case "oneof":
			return fmt.Errorf("поле '%s' должно быть одним из: %s", fieldError.Field(), fieldError.Param())
		case "priority":
			return fmt.Errorf("некорректный приоритет в поле '%s'", fieldError.Field())
		case "status":
			return fmt.Errorf("некорректный статус в поле '%s'", fieldError.Field())
		case "future_date":
			return fmt.Errorf("дата в поле '%s' должна быть в будущем", fieldError.Field())
		default:
			return fmt.Errorf("ошибка валидации поля '%s': %s", fieldError.Field(), fieldError.Tag())
		}
	}

	return fmt.Errorf("ошибка валидации")
}

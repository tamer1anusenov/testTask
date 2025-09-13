package utils

import (
	"fmt"
	"runtime"
)

// Определение кастомных типов ошибок
type ErrorType string

const (
	ErrorTypeValidation   ErrorType = "VALIDATION_ERROR"
	ErrorTypeNotFound     ErrorType = "NOT_FOUND"
	ErrorTypeUnauthorized ErrorType = "UNAUTHORIZED"
	ErrorTypeForbidden    ErrorType = "FORBIDDEN"
	ErrorTypeConflict     ErrorType = "CONFLICT"
	ErrorTypeInternal     ErrorType = "INTERNAL_ERROR"
	ErrorTypeDatabase     ErrorType = "DATABASE_ERROR"
	ErrorTypeExternal     ErrorType = "EXTERNAL_SERVICE_ERROR"
	ErrorTypeBadRequest   ErrorType = "BAD_REQUEST"
	ErrorTypeTimeout      ErrorType = "TIMEOUT"
)

// AppError представляет кастомную ошибку приложения
type AppError struct {
	Type        ErrorType `json:"type"`
	Message     string    `json:"message"`
	Code        int       `json:"code,omitempty"`
	Details     string    `json:"details,omitempty"`
	StackTrace  string    `json:"stack_trace,omitempty"`
	OriginalErr error     `json:"-"`
}

// Error реализует интерфейс error
func (e *AppError) Error() string {
	if e.Details != "" {
		return fmt.Sprintf("%s: %s (%s)", e.Type, e.Message, e.Details)
	}
	return fmt.Sprintf("%s: %s", e.Type, e.Message)
}

// Unwrap возвращает оригинальную ошибку для поддержки errors.Is и errors.As
func (e *AppError) Unwrap() error {
	return e.OriginalErr
}

// WithStackTrace добавляет stack trace к ошибке
func (e *AppError) WithStackTrace() *AppError {
	e.StackTrace = getStackTrace()
	return e
}

// WithDetails добавляет дополнительные детали к ошибке
func (e *AppError) WithDetails(details string) *AppError {
	e.Details = details
	return e
}

// WithCode добавляет код ошибки
func (e *AppError) WithCode(code int) *AppError {
	e.Code = code
	return e
}

// NewError создает новую AppError
func NewError(errorType ErrorType, message string) *AppError {
	return &AppError{
		Type:    errorType,
		Message: message,
	}
}

// NewErrorWithCause создает новую AppError с оригинальной ошибкой
func NewErrorWithCause(errorType ErrorType, message string, originalErr error) *AppError {
	return &AppError{
		Type:        errorType,
		Message:     message,
		OriginalErr: originalErr,
	}
}

// Предопределенные конструкторы ошибок

// NewValidationError создает ошибку валидации
func NewValidationError(message string) *AppError {
	return NewError(ErrorTypeValidation, message).WithCode(400)
}

// NewNotFoundError создает ошибку "не найдено"
func NewNotFoundError(resource string) *AppError {
	return NewError(ErrorTypeNotFound, fmt.Sprintf("%s not found", resource)).WithCode(404)
}

// NewUnauthorizedError создает ошибку неавторизованного доступа
func NewUnauthorizedError(message string) *AppError {
	if message == "" {
		message = "Unauthorized access"
	}
	return NewError(ErrorTypeUnauthorized, message).WithCode(401)
}

// NewForbiddenError создает ошибку запрещенного доступа
func NewForbiddenError(message string) *AppError {
	if message == "" {
		message = "Access forbidden"
	}
	return NewError(ErrorTypeForbidden, message).WithCode(403)
}

// NewConflictError создает ошибку конфликта
func NewConflictError(message string) *AppError {
	return NewError(ErrorTypeConflict, message).WithCode(409)
}

// NewInternalError создает внутреннюю ошибку сервера
func NewInternalError(message string) *AppError {
	if message == "" {
		message = "Internal server error"
	}
	return NewError(ErrorTypeInternal, message).WithCode(500).WithStackTrace()
}

// NewDatabaseError создает ошибку базы данных
func NewDatabaseError(message string, originalErr error) *AppError {
	return NewErrorWithCause(ErrorTypeDatabase, message, originalErr).WithCode(500).WithStackTrace()
}

// NewExternalServiceError создает ошибку внешнего сервиса
func NewExternalServiceError(service, message string, originalErr error) *AppError {
	fullMessage := fmt.Sprintf("External service '%s' error: %s", service, message)
	return NewErrorWithCause(ErrorTypeExternal, fullMessage, originalErr).WithCode(502)
}

// NewBadRequestError создает ошибку неверного запроса
func NewBadRequestError(message string) *AppError {
	return NewError(ErrorTypeBadRequest, message).WithCode(400)
}

// NewTimeoutError создает ошибку таймаута
func NewTimeoutError(operation string) *AppError {
	message := fmt.Sprintf("Operation '%s' timed out", operation)
	return NewError(ErrorTypeTimeout, message).WithCode(408)
}

// WrapError оборачивает стандартную ошибку в AppError
func WrapError(err error, errorType ErrorType, message string) *AppError {
	if err == nil {
		return nil
	}

	// Если ошибка уже является AppError, возвращаем её
	if appErr, ok := err.(*AppError); ok {
		return appErr
	}

	return NewErrorWithCause(errorType, message, err).WithStackTrace()
}

// WrapValidationError оборачивает ошибку как ошибку валидации
func WrapValidationError(err error, message string) *AppError {
	return WrapError(err, ErrorTypeValidation, message).WithCode(400)
}

// WrapDatabaseError оборачивает ошибку как ошибку базы данных
func WrapDatabaseError(err error, operation string) *AppError {
	message := fmt.Sprintf("Database operation failed: %s", operation)
	return WrapError(err, ErrorTypeDatabase, message).WithCode(500)
}

// IsErrorType проверяет, является ли ошибка определенного типа
func IsErrorType(err error, errorType ErrorType) bool {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Type == errorType
	}
	return false
}

// GetErrorCode возвращает код ошибки
func GetErrorCode(err error) int {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Code
	}
	return 500 // Неизвестная ошибка по умолчанию
}

// GetUserFriendlyMessage возвращает пользовательское сообщение об ошибке
func GetUserFriendlyMessage(err error) string {
	if appErr, ok := err.(*AppError); ok {
		switch appErr.Type {
		case ErrorTypeValidation:
			return "Please check your input and try again"
		case ErrorTypeNotFound:
			return "The requested resource was not found"
		case ErrorTypeUnauthorized:
			return "Please log in to continue"
		case ErrorTypeForbidden:
			return "You don't have permission to perform this action"
		case ErrorTypeConflict:
			return "This action conflicts with existing data"
		case ErrorTypeTimeout:
			return "The operation took too long to complete. Please try again"
		case ErrorTypeExternal:
			return "A third-party service is currently unavailable"
		default:
			return "Something went wrong. Please try again later"
		}
	}
	return "An unexpected error occurred"
}

// getStackTrace получает stack trace
func getStackTrace() string {
	const depth = 32
	var pcs [depth]uintptr
	n := runtime.Callers(3, pcs[:])
	frames := runtime.CallersFrames(pcs[:n])

	var trace string
	for {
		frame, more := frames.Next()
		trace += fmt.Sprintf("%s:%d %s\n", frame.File, frame.Line, frame.Function)
		if !more {
			break
		}
	}
	return trace
}

// LogError логирует ошибку (простая реализация, может быть расширена)
func LogError(err error, context string) {
	if appErr, ok := err.(*AppError); ok {
		fmt.Printf("[ERROR] %s: %s\n", context, appErr.Error())
		if appErr.StackTrace != "" {
			fmt.Printf("Stack trace:\n%s\n", appErr.StackTrace)
		}
	} else {
		fmt.Printf("[ERROR] %s: %s\n", context, err.Error())
	}
}

// LogErrorWithDetails логирует ошибку с дополнительными деталями
func LogErrorWithDetails(err error, context string, details map[string]interface{}) {
	fmt.Printf("[ERROR] %s: %s\n", context, err.Error())
	for key, value := range details {
		fmt.Printf("  %s: %v\n", key, value)
	}

	if appErr, ok := err.(*AppError); ok && appErr.StackTrace != "" {
		fmt.Printf("Stack trace:\n%s\n", appErr.StackTrace)
	}
}

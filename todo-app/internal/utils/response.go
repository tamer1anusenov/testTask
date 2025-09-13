package utils

import (
	"encoding/json"
	"fmt"
)

// StandardResponse представляет стандартную структуру ответа
type StandardResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Message string      `json:"message,omitempty"`
	Code    int         `json:"code,omitempty"`
}

// PaginationMeta содержит метаданные для пагинации
type PaginationMeta struct {
	CurrentPage int  `json:"current_page"`
	PageSize    int  `json:"page_size"`
	TotalItems  int  `json:"total_items"`
	TotalPages  int  `json:"total_pages"`
	HasNext     bool `json:"has_next"`
	HasPrev     bool `json:"has_previous"`
}

// PaginatedResponse представляет ответ с пагинацией
type PaginatedResponse struct {
	Success bool           `json:"success"`
	Data    interface{}    `json:"data"`
	Meta    PaginationMeta `json:"meta"`
	Error   string         `json:"error,omitempty"`
}

// SuccessResponse создает успешный ответ с данными
func SuccessResponse(data interface{}) StandardResponse {
	return StandardResponse{
		Success: true,
		Data:    data,
	}
}

// SuccessResponseWithMessage создает успешный ответ с данными и сообщением
func SuccessResponseWithMessage(data interface{}, message string) StandardResponse {
	return StandardResponse{
		Success: true,
		Data:    data,
		Message: message,
	}
}

// ErrorResponse создает ответ с ошибкой
func ErrorResponse(err error) StandardResponse {
	return StandardResponse{
		Success: false,
		Error:   err.Error(),
	}
}

// ErrorResponseWithCode создает ответ с ошибкой и кодом
func ErrorResponseWithCode(err error, code int) StandardResponse {
	return StandardResponse{
		Success: false,
		Error:   err.Error(),
		Code:    code,
	}
}

// ErrorResponseWithMessage создает ответ с кастомным сообщением об ошибке
func ErrorResponseWithMessage(message string) StandardResponse {
	return StandardResponse{
		Success: false,
		Error:   message,
	}
}

// ValidationErrorResponse создает ответ для ошибок валидации
func ValidationErrorResponse(validationErrors map[string]string) StandardResponse {
	return StandardResponse{
		Success: false,
		Error:   "Validation failed",
		Data:    validationErrors,
		Code:    400,
	}
}

// NotFoundResponse создает ответ для случаев "не найдено"
func NotFoundResponse(resource string) StandardResponse {
	return StandardResponse{
		Success: false,
		Error:   fmt.Sprintf("%s not found", resource),
		Code:    404,
	}
}

// UnauthorizedResponse создает ответ для неавторизованных запросов
func UnauthorizedResponse() StandardResponse {
	return StandardResponse{
		Success: false,
		Error:   "Unauthorized access",
		Code:    401,
	}
}

// ForbiddenResponse создает ответ для запрещенных операций
func ForbiddenResponse() StandardResponse {
	return StandardResponse{
		Success: false,
		Error:   "Access forbidden",
		Code:    403,
	}
}

// InternalErrorResponse создает ответ для внутренних ошибок сервера
func InternalErrorResponse() StandardResponse {
	return StandardResponse{
		Success: false,
		Error:   "Internal server error",
		Code:    500,
	}
}

// PaginatedSuccessResponse создает успешный ответ с пагинацией
func PaginatedSuccessResponse(data interface{}, meta PaginationMeta) PaginatedResponse {
	return PaginatedResponse{
		Success: true,
		Data:    data,
		Meta:    meta,
	}
}

// PaginatedErrorResponse создает ответ с ошибкой для пагинированных запросов
func PaginatedErrorResponse(err error) PaginatedResponse {
	return PaginatedResponse{
		Success: false,
		Error:   err.Error(),
		Meta:    PaginationMeta{},
	}
}

// CalculatePaginationMeta вычисляет метаданные пагинации
func CalculatePaginationMeta(currentPage, pageSize, totalItems int) PaginationMeta {
	totalPages := (totalItems + pageSize - 1) / pageSize
	if totalPages == 0 {
		totalPages = 1
	}

	hasNext := currentPage < totalPages
	hasPrev := currentPage > 1

	return PaginationMeta{
		CurrentPage: currentPage,
		PageSize:    pageSize,
		TotalItems:  totalItems,
		TotalPages:  totalPages,
		HasNext:     hasNext,
		HasPrev:     hasPrev,
	}
}

// ToJSON конвертирует ответ в JSON строку
func (r StandardResponse) ToJSON() (string, error) {
	jsonBytes, err := json.Marshal(r)
	if err != nil {
		return "", fmt.Errorf("failed to marshal response to JSON: %w", err)
	}
	return string(jsonBytes), nil
}

// ToJSON конвертирует пагинированный ответ в JSON строку
func (r PaginatedResponse) ToJSON() (string, error) {
	jsonBytes, err := json.Marshal(r)
	if err != nil {
		return "", fmt.Errorf("failed to marshal paginated response to JSON: %w", err)
	}
	return string(jsonBytes), nil
}

// WailsResponse адаптирует StandardResponse для Wails (без error поля в случае успеха)
func WailsResponse(data interface{}, err error) interface{} {
	if err != nil {
		return ErrorResponse(err)
	}
	return SuccessResponse(data)
}

// WailsResponseWithMessage адаптирует ответ с сообщением для Wails
func WailsResponseWithMessage(data interface{}, message string, err error) interface{} {
	if err != nil {
		return ErrorResponse(err)
	}
	return SuccessResponseWithMessage(data, message)
}

// ResponseBuilder предоставляет fluent API для создания ответов
type ResponseBuilder struct {
	response StandardResponse
}

// NewResponseBuilder создает новый ResponseBuilder
func NewResponseBuilder() *ResponseBuilder {
	return &ResponseBuilder{
		response: StandardResponse{},
	}
}

// Success устанавливает успешный статус
func (rb *ResponseBuilder) Success() *ResponseBuilder {
	rb.response.Success = true
	return rb
}

// Error устанавливает статус ошибки
func (rb *ResponseBuilder) Error(err error) *ResponseBuilder {
	rb.response.Success = false
	rb.response.Error = err.Error()
	return rb
}

// Data устанавливает данные ответа
func (rb *ResponseBuilder) Data(data interface{}) *ResponseBuilder {
	rb.response.Data = data
	return rb
}

// Message устанавливает сообщение
func (rb *ResponseBuilder) Message(message string) *ResponseBuilder {
	rb.response.Message = message
	return rb
}

// Code устанавливает код ответа
func (rb *ResponseBuilder) Code(code int) *ResponseBuilder {
	rb.response.Code = code
	return rb
}

// Build возвращает построенный ответ
func (rb *ResponseBuilder) Build() StandardResponse {
	return rb.response
}

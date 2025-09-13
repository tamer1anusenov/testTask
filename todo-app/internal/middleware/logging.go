package middleware

import (
	"fmt"
	"net/http"
	"time"

	"todo-app/internal/utils"
)

// LoggingConfig содержит настройки middleware для логирования
type LoggingConfig struct {
	Logger           *utils.Logger
	LogRequestBody   bool
	LogResponseBody  bool
	SkipPaths        []string
	LogHeaders       []string
	SensitiveHeaders []string // Заголовки, которые нужно скрывать в логах
}

// RequestInfo содержит информацию о запросе для логирования
type RequestInfo struct {
	Method        string
	Path          string
	RemoteAddr    string
	UserAgent     string
	ContentLength int64
	Headers       map[string]string
	Body          string
	StartTime     time.Time
}

// ResponseInfo содержит информацию об ответе для логирования
type ResponseInfo struct {
	StatusCode    int
	ContentLength int
	Headers       map[string]string
	Body          string
	Duration      time.Duration
}

// responseWriter обертка для захвата информации об ответе
type responseWriter struct {
	http.ResponseWriter
	statusCode    int
	contentLength int
	body          []byte
}

// WriteHeader захватывает код статуса
func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// Write захватывает тело ответа
func (rw *responseWriter) Write(b []byte) (int, error) {
	if rw.statusCode == 0 {
		rw.statusCode = http.StatusOK
	}

	rw.contentLength += len(b)
	if cap(rw.body) < len(rw.body)+len(b) {
		// Ограничиваем размер логируемого тела ответа
		maxBodySize := 1024 // 1KB
		if len(rw.body) < maxBodySize {
			remaining := maxBodySize - len(rw.body)
			if len(b) > remaining {
				rw.body = append(rw.body, b[:remaining]...)
			} else {
				rw.body = append(rw.body, b...)
			}
		}
	} else {
		rw.body = append(rw.body, b...)
	}

	return rw.ResponseWriter.Write(b)
}

// DefaultLoggingConfig возвращает настройки логирования по умолчанию
func DefaultLoggingConfig() LoggingConfig {
	return LoggingConfig{
		Logger:          utils.DefaultLogger(),
		LogRequestBody:  false,
		LogResponseBody: false,
		SkipPaths: []string{
			"/health",
			"/metrics",
			"/favicon.ico",
		},
		LogHeaders: []string{
			"Content-Type",
			"Accept",
			"User-Agent",
		},
		SensitiveHeaders: []string{
			"Authorization",
			"Cookie",
			"Set-Cookie",
			"X-API-Key",
		},
	}
}

// DevelopmentLoggingConfig возвращает подробные настройки для разработки
func DevelopmentLoggingConfig() LoggingConfig {
	return LoggingConfig{
		Logger:          utils.DefaultLogger(),
		LogRequestBody:  true,
		LogResponseBody: true,
		SkipPaths: []string{
			"/health",
		},
		LogHeaders: []string{
			"Content-Type",
			"Accept",
			"User-Agent",
			"Origin",
			"Referer",
		},
		SensitiveHeaders: []string{
			"Authorization",
			"Cookie",
			"Set-Cookie",
			"X-API-Key",
		},
	}
}

// Logging создает middleware для логирования HTTP запросов
func Logging(config LoggingConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Проверяем, нужно ли пропустить логирование для этого пути
			if shouldSkipPath(r.URL.Path, config.SkipPaths) {
				next.ServeHTTP(w, r)
				return
			}

			startTime := time.Now()

			// Собираем информацию о запросе
			requestInfo := collectRequestInfo(r, config, startTime)

			// Логируем входящий запрос
			logIncomingRequest(config.Logger, requestInfo)

			// Создаем обертку для ResponseWriter
			rw := &responseWriter{
				ResponseWriter: w,
				statusCode:     0,
				contentLength:  0,
				body:           make([]byte, 0),
			}

			// Выполняем запрос
			next.ServeHTTP(rw, r)

			// Собираем информацию об ответе
			responseInfo := collectResponseInfo(rw, config, time.Since(startTime))

			// Логируем завершенный запрос
			logCompletedRequest(config.Logger, requestInfo, responseInfo)
		})
	}
}

// collectRequestInfo собирает информацию о запросе
func collectRequestInfo(r *http.Request, config LoggingConfig, startTime time.Time) RequestInfo {
	info := RequestInfo{
		Method:        r.Method,
		Path:          r.URL.Path,
		RemoteAddr:    getClientIP(r),
		UserAgent:     r.UserAgent(),
		ContentLength: r.ContentLength,
		Headers:       make(map[string]string),
		StartTime:     startTime,
	}

	// Собираем заголовки
	for _, headerName := range config.LogHeaders {
		if value := r.Header.Get(headerName); value != "" {
			if isSensitiveHeader(headerName, config.SensitiveHeaders) {
				info.Headers[headerName] = "[REDACTED]"
			} else {
				info.Headers[headerName] = value
			}
		}
	}

	// Читаем тело запроса если нужно
	if config.LogRequestBody && r.ContentLength > 0 && r.ContentLength < 1024 {
		// TODO: Реализовать чтение тела запроса
		// Нужно быть осторожным, чтобы не "израсходовать" тело запроса
		info.Body = "[BODY_READING_NOT_IMPLEMENTED]"
	}

	return info
}

// collectResponseInfo собирает информацию об ответе
func collectResponseInfo(rw *responseWriter, config LoggingConfig, duration time.Duration) ResponseInfo {
	info := ResponseInfo{
		StatusCode:    rw.statusCode,
		ContentLength: rw.contentLength,
		Headers:       make(map[string]string),
		Duration:      duration,
	}

	// Собираем заголовки ответа
	for _, headerName := range config.LogHeaders {
		if value := rw.Header().Get(headerName); value != "" {
			info.Headers[headerName] = value
		}
	}

	// Добавляем тело ответа если нужно
	if config.LogResponseBody && len(rw.body) > 0 {
		info.Body = string(rw.body)
	}

	return info
}

// logIncomingRequest логирует входящий запрос
func logIncomingRequest(logger *utils.Logger, info RequestInfo) {
	fields := map[string]interface{}{
		"method":         info.Method,
		"path":           info.Path,
		"remote_addr":    info.RemoteAddr,
		"user_agent":     info.UserAgent,
		"content_length": info.ContentLength,
	}

	// Добавляем заголовки
	for key, value := range info.Headers {
		fields["header_"+key] = value
	}

	// Добавляем тело если есть
	if info.Body != "" {
		fields["body"] = info.Body
	}

	logger.Info("Incoming request", fields)
}

// logCompletedRequest логирует завершенный запрос
func logCompletedRequest(logger *utils.Logger, requestInfo RequestInfo, responseInfo ResponseInfo) {
	fields := map[string]interface{}{
		"method":        requestInfo.Method,
		"path":          requestInfo.Path,
		"remote_addr":   requestInfo.RemoteAddr,
		"status_code":   responseInfo.StatusCode,
		"duration_ms":   responseInfo.Duration.Milliseconds(),
		"response_size": responseInfo.ContentLength,
	}

	// Добавляем заголовки ответа
	for key, value := range responseInfo.Headers {
		fields["response_header_"+key] = value
	}

	// Добавляем тело ответа если есть
	if responseInfo.Body != "" {
		fields["response_body"] = responseInfo.Body
	}

	// Определяем уровень логирования на основе статус кода
	message := fmt.Sprintf("%s %s - %d (%dms)",
		requestInfo.Method, requestInfo.Path, responseInfo.StatusCode, responseInfo.Duration.Milliseconds())

	if responseInfo.StatusCode >= 500 {
		logger.Error(message, fields)
	} else if responseInfo.StatusCode >= 400 {
		logger.Warn(message, fields)
	} else {
		logger.Info(message, fields)
	}
}

// shouldSkipPath проверяет, нужно ли пропустить логирование для пути
func shouldSkipPath(path string, skipPaths []string) bool {
	for _, skipPath := range skipPaths {
		if path == skipPath {
			return true
		}
	}
	return false
}

// isSensitiveHeader проверяет, является ли заголовок чувствительным
func isSensitiveHeader(headerName string, sensitiveHeaders []string) bool {
	for _, sensitive := range sensitiveHeaders {
		if headerName == sensitive {
			return true
		}
	}
	return false
}

// getClientIP пытается получить реальный IP клиента
func getClientIP(r *http.Request) string {
	// Проверяем заголовки в порядке приоритета

	// X-Forwarded-For (может содержать список IP)
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		// Берем первый IP из списка
		for i := 0; i < len(xff); i++ {
			if xff[i] == ',' {
				return xff[:i]
			}
		}
		return xff
	}

	// X-Real-IP
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}

	// X-Forwarded
	if xf := r.Header.Get("X-Forwarded"); xf != "" {
		return xf
	}

	// Fallback к RemoteAddr
	return r.RemoteAddr
}

// RequestIDMiddleware добавляет уникальный ID к каждому запросу
func RequestIDMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := generateRequestID()

			// Добавляем ID в контекст запроса
			// В Go 1.7+ можно использовать context.WithValue

			// Добавляем ID в заголовок ответа
			w.Header().Set("X-Request-ID", requestID)

			next.ServeHTTP(w, r)
		})
	}
}

// generateRequestID генерирует уникальный ID запроса
func generateRequestID() string {
	// Простая реализация на основе времени
	// В продакшене лучше использовать UUID
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

// StructuredLogging создает middleware для структурированного логирования
func StructuredLogging(logger *utils.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			startTime := time.Now()

			// Создаем контекстный логгер для этого запроса
			contextLogger := logger.WithFields(map[string]interface{}{
				"request_id": generateRequestID(),
				"method":     r.Method,
				"path":       r.URL.Path,
				"remote_ip":  getClientIP(r),
			})

			contextLogger.Info("Request started")

			// Создаем обертку для ResponseWriter
			rw := &responseWriter{
				ResponseWriter: w,
				statusCode:     0,
				contentLength:  0,
				body:           make([]byte, 0),
			}

			// Выполняем запрос
			next.ServeHTTP(rw, r)

			// Логируем завершение
			duration := time.Since(startTime)
			contextLogger = logger.WithFields(map[string]interface{}{
				"status_code": rw.statusCode,
				"duration_ms": duration.Milliseconds(),
				"size_bytes":  rw.contentLength,
			})

			if rw.statusCode >= 500 {
				contextLogger.Error("Request completed with server error")
			} else if rw.statusCode >= 400 {
				contextLogger.Warn("Request completed with client error")
			} else {
				contextLogger.Info("Request completed successfully")
			}
		})
	}
}

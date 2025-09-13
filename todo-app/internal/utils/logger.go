package utils

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// LogLevel представляет уровень логирования
type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
	FATAL
)

// String возвращает строковое представление уровня лога
func (l LogLevel) String() string {
	switch l {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARN:
		return "WARN"
	case ERROR:
		return "ERROR"
	case FATAL:
		return "FATAL"
	default:
		return "UNKNOWN"
	}
}

// Logger представляет кастомный логгер
type Logger struct {
	level      LogLevel
	output     io.Writer
	debugLog   *log.Logger
	infoLog    *log.Logger
	warnLog    *log.Logger
	errorLog   *log.Logger
	fatalLog   *log.Logger
	jsonFormat bool
}

// LoggerConfig содержит конфигурацию логгера
type LoggerConfig struct {
	Level      LogLevel
	Output     io.Writer
	JSONFormat bool
	LogFile    string
}

// NewLogger создает новый логгер
func NewLogger(config LoggerConfig) (*Logger, error) {
	output := config.Output
	if output == nil {
		output = os.Stdout
	}

	// Если указан файл логов, создаем MultiWriter
	if config.LogFile != "" {
		// Создаем директорию для логов если она не существует
		logDir := filepath.Dir(config.LogFile)
		if err := os.MkdirAll(logDir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create log directory: %w", err)
		}

		file, err := os.OpenFile(config.LogFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return nil, fmt.Errorf("failed to open log file: %w", err)
		}
		output = io.MultiWriter(config.Output, file)
	}

	logger := &Logger{
		level:      config.Level,
		output:     output,
		jsonFormat: config.JSONFormat,
	}

	// Создаем логгеры для каждого уровня
	flags := log.LstdFlags
	if !config.JSONFormat {
		flags = log.LstdFlags | log.Lshortfile
	}

	logger.debugLog = log.New(output, "", flags)
	logger.infoLog = log.New(output, "", flags)
	logger.warnLog = log.New(output, "", flags)
	logger.errorLog = log.New(output, "", flags)
	logger.fatalLog = log.New(output, "", flags)

	return logger, nil
}

// DefaultLogger создает логгер с настройками по умолчанию
func DefaultLogger() *Logger {
	logger, _ := NewLogger(LoggerConfig{
		Level:      INFO,
		Output:     os.Stdout,
		JSONFormat: false,
	})
	return logger
}

// formatMessage форматирует сообщение в зависимости от формата
func (l *Logger) formatMessage(level LogLevel, message string, fields map[string]interface{}) string {
	if l.jsonFormat {
		return l.formatJSON(level, message, fields)
	}
	return l.formatText(level, message, fields)
}

// formatJSON форматирует сообщение в JSON формате
func (l *Logger) formatJSON(level LogLevel, message string, fields map[string]interface{}) string {
	// Простая JSON сериализация (можно заменить на encoding/json)
	json := fmt.Sprintf(`{"timestamp":"%s","level":"%s","message":"%s"`,
		time.Now().UTC().Format(time.RFC3339), level.String(), message)

	for key, value := range fields {
		json += fmt.Sprintf(`,"%s":"%v"`, key, value)
	}

	// Добавляем caller info
	_, file, line, ok := runtime.Caller(3)
	if ok {
		json += fmt.Sprintf(`,"caller":"%s:%d"`, filepath.Base(file), line)
	}

	json += "}"
	return json
}

// formatText форматирует сообщение в текстовом формате
func (l *Logger) formatText(level LogLevel, message string, fields map[string]interface{}) string {
	text := fmt.Sprintf("[%s] %s", level.String(), message)

	if len(fields) > 0 {
		var fieldStrings []string
		for key, value := range fields {
			fieldStrings = append(fieldStrings, fmt.Sprintf("%s=%v", key, value))
		}
		text += " " + strings.Join(fieldStrings, " ")
	}

	return text
}

// shouldLog проверяет, нужно ли логировать сообщение данного уровня
func (l *Logger) shouldLog(level LogLevel) bool {
	return level >= l.level
}

// Debug логирует сообщение уровня DEBUG
func (l *Logger) Debug(message string, fields ...map[string]interface{}) {
	if !l.shouldLog(DEBUG) {
		return
	}

	var fieldMap map[string]interface{}
	if len(fields) > 0 {
		fieldMap = fields[0]
	}

	l.debugLog.Println(l.formatMessage(DEBUG, message, fieldMap))
}

// Info логирует сообщение уровня INFO
func (l *Logger) Info(message string, fields ...map[string]interface{}) {
	if !l.shouldLog(INFO) {
		return
	}

	var fieldMap map[string]interface{}
	if len(fields) > 0 {
		fieldMap = fields[0]
	}

	l.infoLog.Println(l.formatMessage(INFO, message, fieldMap))
}

// Warn логирует сообщение уровня WARN
func (l *Logger) Warn(message string, fields ...map[string]interface{}) {
	if !l.shouldLog(WARN) {
		return
	}

	var fieldMap map[string]interface{}
	if len(fields) > 0 {
		fieldMap = fields[0]
	}

	l.warnLog.Println(l.formatMessage(WARN, message, fieldMap))
}

// Error логирует сообщение уровня ERROR
func (l *Logger) Error(message string, fields ...map[string]interface{}) {
	if !l.shouldLog(ERROR) {
		return
	}

	var fieldMap map[string]interface{}
	if len(fields) > 0 {
		fieldMap = fields[0]
	}

	l.errorLog.Println(l.formatMessage(ERROR, message, fieldMap))
}

// Fatal логирует сообщение уровня FATAL и завершает программу
func (l *Logger) Fatal(message string, fields ...map[string]interface{}) {
	var fieldMap map[string]interface{}
	if len(fields) > 0 {
		fieldMap = fields[0]
	}

	l.fatalLog.Println(l.formatMessage(FATAL, message, fieldMap))
	os.Exit(1)
}

// LogError логирует ошибку с дополнительным контекстом
func (l *Logger) LogError(err error, context string, fields ...map[string]interface{}) {
	var fieldMap map[string]interface{}
	if len(fields) > 0 {
		fieldMap = fields[0]
	} else {
		fieldMap = make(map[string]interface{})
	}

	fieldMap["error"] = err.Error()
	fieldMap["context"] = context

	l.Error("Error occurred", fieldMap)
}

// LogDBQuery логирует запрос к базе данных
func (l *Logger) LogDBQuery(query string, duration time.Duration, err error) {
	fields := map[string]interface{}{
		"query":    query,
		"duration": duration.String(),
	}

	if err != nil {
		fields["error"] = err.Error()
		l.Error("Database query failed", fields)
	} else {
		l.Debug("Database query executed", fields)
	}
}

// LogRequest логирует HTTP запрос (для будущего использования)
func (l *Logger) LogRequest(method, path string, statusCode int, duration time.Duration, userID string) {
	fields := map[string]interface{}{
		"method":      method,
		"path":        path,
		"status_code": statusCode,
		"duration":    duration.String(),
	}

	if userID != "" {
		fields["user_id"] = userID
	}

	level := INFO
	if statusCode >= 400 {
		level = WARN
	}
	if statusCode >= 500 {
		level = ERROR
	}

	message := fmt.Sprintf("%s %s - %d", method, path, statusCode)

	switch level {
	case INFO:
		l.Info(message, fields)
	case WARN:
		l.Warn(message, fields)
	case ERROR:
		l.Error(message, fields)
	}
}

// LogAPICall логирует вызов внешнего API
func (l *Logger) LogAPICall(service, endpoint string, method string, statusCode int, duration time.Duration, err error) {
	fields := map[string]interface{}{
		"service":     service,
		"endpoint":    endpoint,
		"method":      method,
		"status_code": statusCode,
		"duration":    duration.String(),
	}

	if err != nil {
		fields["error"] = err.Error()
		l.Error("External API call failed", fields)
	} else {
		l.Info("External API call completed", fields)
	}
}

// LogUserAction логирует действие пользователя
func (l *Logger) LogUserAction(userID, action string, resourceType, resourceID string) {
	fields := map[string]interface{}{
		"user_id":       userID,
		"action":        action,
		"resource_type": resourceType,
		"resource_id":   resourceID,
	}

	l.Info("User action performed", fields)
}

// WithFields создает временный логгер с предустановленными полями
func (l *Logger) WithFields(fields map[string]interface{}) *ContextLogger {
	return &ContextLogger{
		logger: l,
		fields: fields,
	}
}

// ContextLogger предоставляет логгер с предустановленным контекстом
type ContextLogger struct {
	logger *Logger
	fields map[string]interface{}
}

// Debug логирует DEBUG сообщение с контекстом
func (cl *ContextLogger) Debug(message string) {
	cl.logger.Debug(message, cl.fields)
}

// Info логирует INFO сообщение с контекстом
func (cl *ContextLogger) Info(message string) {
	cl.logger.Info(message, cl.fields)
}

// Warn логирует WARN сообщение с контекстом
func (cl *ContextLogger) Warn(message string) {
	cl.logger.Warn(message, cl.fields)
}

// Error логирует ERROR сообщение с контекстом
func (cl *ContextLogger) Error(message string) {
	cl.logger.Error(message, cl.fields)
}

// Глобальный логгер
var defaultLogger = DefaultLogger()

// SetDefaultLogger устанавливает глобальный логгер
func SetDefaultLogger(logger *Logger) {
	defaultLogger = logger
}

// Глобальные функции для удобства использования
func Debug(message string, fields ...map[string]interface{}) {
	defaultLogger.Debug(message, fields...)
}

func Info(message string, fields ...map[string]interface{}) {
	defaultLogger.Info(message, fields...)
}

func Warn(message string, fields ...map[string]interface{}) {
	defaultLogger.Warn(message, fields...)
}

func Error(message string, fields ...map[string]interface{}) {
	defaultLogger.Error(message, fields...)
}

func Fatal(message string, fields ...map[string]interface{}) {
	defaultLogger.Fatal(message, fields...)
}

func GlobalLogError(err error, context string, fields ...map[string]interface{}) {
	defaultLogger.LogError(err, context, fields...)
}

func LogDBQuery(query string, duration time.Duration, err error) {
	defaultLogger.LogDBQuery(query, duration, err)
}

func LogRequest(method, path string, statusCode int, duration time.Duration, userID string) {
	defaultLogger.LogRequest(method, path, statusCode, duration, userID)
}

package utils

import (
	"fmt"
	"time"
)

// TimeFormats содержит различные форматы времени
var TimeFormats = struct {
	ISO8601   string
	RFC3339   string
	DateOnly  string
	TimeOnly  string
	DateTime  string
	Timestamp string
}{
	ISO8601:   "2006-01-02T15:04:05Z07:00",
	RFC3339:   time.RFC3339,
	DateOnly:  "2006-01-02",
	TimeOnly:  "15:04:05",
	DateTime:  "2006-01-02 15:04:05",
	Timestamp: "2006-01-02T15:04:05.000Z",
}

// GetCurrentTime возвращает текущее время в UTC
func GetCurrentTime() time.Time {
	return time.Now().UTC()
}

// FormatTimeISO8601 форматирует время в ISO8601
func FormatTimeISO8601(t time.Time) string {
	return t.Format(TimeFormats.ISO8601)
}

// FormatTimeCustom форматирует время в указанном формате
func FormatTimeCustom(t time.Time, format string) string {
	return t.Format(format)
}

// ParseTimeISO8601 парсит время из строки в формате ISO8601
func ParseTimeISO8601(timeStr string) (time.Time, error) {
	t, err := time.Parse(TimeFormats.ISO8601, timeStr)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to parse time '%s': %w", timeStr, err)
	}
	return t, nil
}

// ParseTimeCustom парсит время из строки в указанном формате
func ParseTimeCustom(timeStr, format string) (time.Time, error) {
	t, err := time.Parse(format, timeStr)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to parse time '%s' with format '%s': %w", timeStr, format, err)
	}
	return t, nil
}

// TimestampToString конвертирует Unix timestamp в строку
func TimestampToString(timestamp int64) string {
	return time.Unix(timestamp, 0).UTC().Format(TimeFormats.ISO8601)
}

// StringToTimestamp конвертирует строку в Unix timestamp
func StringToTimestamp(timeStr string) (int64, error) {
	t, err := ParseTimeISO8601(timeStr)
	if err != nil {
		return 0, err
	}
	return t.Unix(), nil
}

// DaysDifference вычисляет разницу в днях между двумя датами
func DaysDifference(start, end time.Time) int {
	diff := end.Sub(start)
	return int(diff.Hours() / 24)
}

// HoursDifference вычисляет разницу в часах между двумя датами
func HoursDifference(start, end time.Time) int {
	diff := end.Sub(start)
	return int(diff.Hours())
}

// MinutesDifference вычисляет разницу в минутах между двумя датами
func MinutesDifference(start, end time.Time) int {
	diff := end.Sub(start)
	return int(diff.Minutes())
}

// IsToday проверяет, является ли указанная дата сегодняшней
func IsToday(t time.Time) bool {
	now := GetCurrentTime()
	return t.Year() == now.Year() && t.Month() == now.Month() && t.Day() == now.Day()
}

// IsThisWeek проверяет, находится ли дата в текущей неделе
func IsThisWeek(t time.Time) bool {
	now := GetCurrentTime()

	// Получаем начало недели (понедельник)
	weekday := int(now.Weekday())
	if weekday == 0 { // Воскресенье
		weekday = 7
	}
	startOfWeek := now.AddDate(0, 0, -(weekday - 1)).Truncate(24 * time.Hour)
	endOfWeek := startOfWeek.AddDate(0, 0, 7)

	return t.After(startOfWeek) && t.Before(endOfWeek)
}

// GetStartOfDay возвращает начало дня для указанной даты
func GetStartOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

// GetEndOfDay возвращает конец дня для указанной даты
func GetEndOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 999999999, t.Location())
}

// GetStartOfWeek возвращает начало недели (понедельник) для указанной даты
func GetStartOfWeek(t time.Time) time.Time {
	weekday := int(t.Weekday())
	if weekday == 0 { // Воскресенье
		weekday = 7
	}
	return GetStartOfDay(t.AddDate(0, 0, -(weekday - 1)))
}

// GetEndOfWeek возвращает конец недели (воскресенье) для указанной даты
func GetEndOfWeek(t time.Time) time.Time {
	startOfWeek := GetStartOfWeek(t)
	return GetEndOfDay(startOfWeek.AddDate(0, 0, 6))
}

// TimeZoneOffset возвращает смещение часового пояса в часах
func TimeZoneOffset(t time.Time) int {
	_, offset := t.Zone()
	return offset / 3600 // конвертируем секунды в часы
}

// ConvertToTimezone конвертирует время в указанный часовой пояс
func ConvertToTimezone(t time.Time, timezone string) (time.Time, error) {
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to load timezone '%s': %w", timezone, err)
	}
	return t.In(loc), nil
}

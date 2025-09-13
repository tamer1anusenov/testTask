package integration

import (
	"fmt"
	"os"
	"testing"
)

// TestMain служит точкой входа для всех интеграционных тестов
func TestMain(m *testing.M) {
	fmt.Println("=== Starting Todo App Integration Tests ===")

	// Проверяем переменные окружения для тестирования
	if os.Getenv("TEST_DATABASE_URL") == "" {
		fmt.Println("Warning: TEST_DATABASE_URL not set, using default test database")
	}

	// Запускаем тесты
	code := m.Run()

	fmt.Println("=== Integration Tests Completed ===")

	// Завершаем с соответствующим кодом выхода
	os.Exit(code)
}

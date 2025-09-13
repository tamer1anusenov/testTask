package internal

import (
	"database/sql"
	"testing"
	"todo-app/app/config"
	"todo-app/app/usecases"
	"todo-app/internal/testutils"

	_ "github.com/lib/pq"
)

// TestApp содержит все компоненты приложения для тестирования
type TestApp struct {
	TaskUseCase      usecases.TaskUseCase
	AnalyticsUseCase usecases.AnalyticsUseCase
	ExportUseCase    usecases.ExportUseCase
}

// TestContainer управляет тестовой средой
type TestContainer struct {
	TestDB *sql.DB
	cfg    *config.Config
	app    *TestApp
}

// SetupTestContainer создает тестовый контейнер с тестовой БД
func SetupTestContainer(t *testing.T) *TestContainer {
	// Загружаем тестовую конфигурацию
	cfg := &config.Config{
		Database: config.DatabaseConfig{
			Host:     "localhost",
			Port:     5432,
			User:     "postgres",
			Password: "postgres",
			DBName:   "todo_db_test",
			SSLMode:  "disable",
		},
		App: config.AppConfig{
			Debug: true,
		},
	}

	// Подключаемся к тестовой БД
	dsn := cfg.GetDatabaseDSN()
	db, err := sql.Open("postgres", dsn)
	testutils.AssertNoError(t, err, "Failed to connect to test database")

	// Проверяем соединение
	err = db.Ping()
	testutils.AssertNoError(t, err, "Failed to ping test database")

	// Создаем компоненты вручную для тестирования
	// TODO: Здесь нужно будет создать реальные usecase после их реализации
	testApp := &TestApp{
		// TaskUseCase:      nil, // Заполним после создания реальных usecase
		// AnalyticsUseCase: nil,
		// ExportUseCase:    nil,
	}

	return &TestContainer{
		TestDB: db,
		cfg:    cfg,
		app:    testApp,
	}
}

// TeardownTestContainer очищает тестовую среду
func (tc *TestContainer) TeardownTestContainer(t *testing.T) {
	if tc.TestDB != nil {
		err := tc.TestDB.Close()
		testutils.AssertNoError(t, err, "Failed to close test database")
	}
}

// GetTestApp возвращает тестовое приложение
func (tc *TestContainer) GetTestApp() *TestApp {
	return tc.app
}

// ClearTestData очищает все тестовые данные
func (tc *TestContainer) ClearTestData(t *testing.T) {
	queries := []string{
		"TRUNCATE TABLE tasks RESTART IDENTITY CASCADE",
	}

	for _, query := range queries {
		_, err := tc.TestDB.Exec(query)
		testutils.AssertNoError(t, err, "Failed to clear test data: "+query)
	}
}

// LoadTestData загружает тестовые данные из файла
func (tc *TestContainer) LoadTestData(t *testing.T) {
	testDataSQL := `
		INSERT INTO tasks (title, description, status, priority, created_at, updated_at, due_date) VALUES
		('Test Task 1', 'First test task', 'active', 'high', NOW(), NOW(), NULL),
		('Test Task 2', 'Second test task', 'completed', 'medium', NOW(), NOW(), '2024-12-25 10:00:00'),
		('Test Task 3', 'Third test task', 'active', 'low', NOW(), NOW(), '2024-12-30 15:00:00'),
		('Test Task 4', 'Fourth test task', 'active', 'high', NOW(), NOW(), NULL),
		('Test Task 5', 'Fifth test task', 'completed', 'medium', NOW(), NOW(), '2024-12-20 09:00:00');
	`

	_, err := tc.TestDB.Exec(testDataSQL)
	testutils.AssertNoError(t, err, "Failed to load test data")
}

// ExecuteSQL выполняет произвольный SQL запрос (для сложных тестов)
func (tc *TestContainer) ExecuteSQL(t *testing.T, query string, args ...interface{}) {
	_, err := tc.TestDB.Exec(query, args...)
	testutils.AssertNoError(t, err, "Failed to execute SQL: "+query)
}

// QueryRow выполняет запрос и возвращает одну строку
func (tc *TestContainer) QueryRow(query string, args ...interface{}) *sql.Row {
	return tc.TestDB.QueryRow(query, args...)
}

// Query выполняет запрос и возвращает несколько строк
func (tc *TestContainer) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return tc.TestDB.Query(query, args...)
}

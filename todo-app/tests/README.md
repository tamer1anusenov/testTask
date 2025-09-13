# README: Система тестирования Todo App

## 📁 Структура тестов

```
tests/
├── integration/           # Интеграционные тесты
│   ├── main_test.go      # Главный файл тестов
│   └── task_flow_test.go # Тесты бизнес-процессов
├── internal/             # Внутренние утилиты тестирования
│   └── test_container.go # Управление тестовой средой
└── run_tests.sh         # Скрипт запуска всех тестов
```

## 🧪 Типы тестов

### 1. Unit-тесты (Модульные)
- **Расположение**: `app/repository/*_test.go`, `app/services/*_test.go`
- **Назначение**: Тестирование отдельных компонентов
- **Технологии**: go-sqlmock для моков БД
- **Запуск**: `go test ./app/repository/... ./app/services/...`

### 2. Integration-тесты (Интеграционные)
- **Расположение**: `tests/integration/`
- **Назначение**: Тестирование взаимодействия компонентов
- **Технологии**: Реальная тестовая БД PostgreSQL
- **Запуск**: `go test ./tests/integration/...`

## 🚀 Запуск тестов

### Быстрый запуск всех тестов
```bash
cd todo-app
./tests/run_tests.sh
```

### Запуск отдельных типов тестов
```bash
# Unit-тесты только
go test ./app/repository/... ./app/services/...

# Интеграционные тесты только
go test ./tests/integration/...

# С покрытием кода
go test -cover ./...

# С детальным выводом
go test -v ./...
```

## 🗄️ Тестовая база данных

### Настройка
Перед запуском интеграционных тестов убедитесь, что:

1. **PostgreSQL установлен и запущен**
2. **Создана тестовая база данных**:
```sql
CREATE DATABASE todo_db_test;
CREATE USER test_user WITH PASSWORD 'test_password';
GRANT ALL PRIVILEGES ON DATABASE todo_db_test TO test_user;
```

3. **Применены миграции**:
```bash
# Из корня проекта
migrate -path database/migrations -database "postgres://test_user:test_password@localhost/todo_db_test?sslmode=disable" up
```

### Переменные окружения
```bash
export TEST_DB_HOST=localhost
export TEST_DB_PORT=5432
export TEST_DB_USER=test_user
export TEST_DB_PASSWORD=test_password
export TEST_DB_NAME=todo_db_test
export TEST_DB_SSLMODE=disable
```

## 📊 Покрытие кода

### Генерация отчета
```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

### Просмотр покрытия
```bash
# В терминале
go tool cover -func=coverage.out

# В браузере
open coverage.html
```

## 🏗️ Архитектура тестирования

### TestContainer
Основной компонент для управления тестовой средой:
- Создание изолированной тестовой БД
- Загрузка тестовых данных
- Очистка данных между тестами
- Управление жизненным циклом тестов

### Тестовые данные
- **Статичные данные**: `database/seeds/test_data.sql`
- **Динамические данные**: Создаются в процессе тестирования
- **Очистка**: Автоматическая между тестами

### Утилиты тестирования
- **Assertions**: `internal/testutils/assertions.go`
- **Моки**: Использование sqlmock для БД
- **Хелперы**: Вспомогательные функции для тестов

## ✅ Типы тестовых сценариев

### Модульные тесты
1. **Repository Layer**:
   - CRUD операции с моком БД
   - Валидация SQL запросов
   - Обработка ошибок

2. **Service Layer**:
   - Бизнес-логика с мок-репозиториями
   - Валидация входных данных
   - Обработка исключительных ситуаций

### Интеграционные тесты
1. **Task Flow Tests**:
   - Полный жизненный цикл задачи
   - Создание → Обновление → Удаление
   - Переключение статусов

2. **Database Integration**:
   - Работа с реальной БД
   - Транзакции и откаты
   - Проверка целостности данных

3. **Concurrency Tests**:
   - Параллельные операции
   - Race conditions
   - Блокировки БД

## 🐛 Отладка тестов

### Логирование
```bash
# Включить детальные логи
export LOG_LEVEL=debug

# Запуск с подробным выводом
go test -v -run TestSpecificTest ./tests/integration/
```

### Изоляция тестов
```bash
# Запуск конкретного теста
go test -run TestTaskFlow_CreateUpdateDeleteFlow ./tests/integration/

# Запуск тестов с таймаутом
go test -timeout 30s ./tests/integration/
```

## 📝 Написание новых тестов

### Unit-тест для repository
```go
func TestTaskRepository_Create(t *testing.T) {
    db, mock := testutils.SetupMockDB(t)
    defer db.Close()
    
    repo := NewTaskRepository(db)
    
    // Настройка мока
    mock.ExpectExec("INSERT INTO tasks").
        WithArgs("Test", "Description", "active", "high").
        WillReturnResult(sqlmock.NewResult(1, 1))
    
    // Выполнение теста
    task, err := repo.Create(context.Background(), createReq)
    
    // Проверки
    testutils.AssertNoError(t, err)
    testutils.AssertNotNil(t, task)
}
```

### Integration-тест
```go
func TestTaskFlow_NewFeature(t *testing.T) {
    container := internal.SetupTestContainer(t)
    defer container.TeardownTestContainer(t)
    
    // Очистка данных
    container.ClearTestData(t)
    
    // Тестирование нового функционала
    app := container.GetTestApp()
    // ... test logic
}
```

## 🔧 Конфигурация CI/CD

### GitHub Actions пример
```yaml
name: Tests
on: [push, pull_request]
jobs:
  test:
    runs-on: ubuntu-latest
    services:
      postgres:
        image: postgres:13
        env:
          POSTGRES_PASSWORD: postgres
          POSTGRES_DB: todo_db_test
        options: >-
          --health-cmd pg_isready
          --health-interval 10s
          --health-timeout 5s
          --health-retries 5
    steps:
      - uses: actions/checkout@v2
      - uses: actions/setup-go@v2
        with:
          go-version: 1.21
      - run: ./tests/run_tests.sh
```

## 📈 Метрики качества

### Требования к покрытию
- **Unit-тесты**: Минимум 80% покрытия
- **Integration-тесты**: Покрытие основных user stories
- **Критический код**: 100% покрытия (auth, payments, etc.)

### Performance тесты
```bash
# Бенчмарки
go test -bench=. ./app/...

# Профилирование
go test -cpuprofile=cpu.prof -memprofile=mem.prof ./app/...
```

## 🚨 Troubleshooting

### Частые проблемы

1. **"Database connection failed"**
   - Проверьте, что PostgreSQL запущен
   - Убедитесь в правильности credentials
   - Проверьте сетевое подключение

2. **"Import cycle detected"**
   - Используйте internal/testutils для общих утилит
   - Избегайте циклических зависимостей между пакетами

3. **"Test timeout"**
   - Увеличьте таймаут: `go test -timeout 60s`
   - Оптимизируйте медленные тесты

4. **"Resource leak"**
   - Всегда закрывайте соединения с БД
   - Используйте defer для cleanup

### Получение помощи
- Изучите логи тестов: `go test -v`
- Проверьте состояние БД: подключитесь через psql
- Запустите тесты по отдельности для изоляции проблем

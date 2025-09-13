# Todo App - Фаза 9: Сборка проекта

## Архитектура

Проект организован по принципу Clean Architecture с Dependency Injection:

```
main.go                 # Точка входа, настройка Wails
app.go                  # Wails bindings и API методы
app/
├── container.go        # DI контейнер
├── app.go             # Основная структура приложения
├── config/            # Конфигурация
├── models/            # Модели данных
├── repository/        # Слой данных
├── services/          # Бизнес-сервисы
└── usecases/          # Use cases (бизнес-логика)
internal/
├── utils/             # Утилиты (DB, Logger, Response, Errors)
└── middleware/        # HTTP middleware
```

## Запуск приложения

1. **Настройка конфигурации:**
```bash
cp .env.example .env
# Отредактируйте .env файл под ваши настройки
```

2. **Настройка базы данных:**
```bash
# Создайте базу данных PostgreSQL
createdb todo_db

# Миграции выполнятся автоматически при запуске
```

3. **Запуск в режиме разработки:**
```bash
wails dev
```

4. **Сборка для продакшена:**
```bash
wails build
```

## Основные компоненты

### main.go
- Загрузка конфигурации из переменных окружения
- Инициализация DI контейнера  
- Выполнение миграций БД
- Настройка и запуск Wails приложения

### app/container.go (DI Container)
- Инициализация всех зависимостей в правильном порядке
- Настройка подключения к БД
- Создание Repository → Services → UseCases цепочки
- Управление жизненным циклом ресурсов

### app.go (Wails Bindings)
Методы для фронтенда:
- `CreateTask(title, description, priority)` - создание задачи
- `GetAllTasks()` - получение всех задач
- `GetTasksByStatus(status)` - фильтрация по статусу
- `UpdateTask(id, ...)` - обновление задачи  
- `DeleteTask(id)` - удаление задачи
- `ToggleTaskStatus(id)` - переключение статуса
- `GetTasksStats()` - статистика
- `GetDashboardStats()` - данные для дашборда

## Конфигурация

Приложение настраивается через переменные окружения:

### База данных
- `DB_HOST` - хост БД (localhost)
- `DB_PORT` - порт БД (5432)  
- `DB_USER` - пользователь БД
- `DB_PASSWORD` - пароль БД
- `DB_NAME` - имя БД
- `DB_SSLMODE` - SSL режим (disable/require)

### Логирование  
- `LOG_LEVEL` - уровень логов (debug/info/warn/error)
- `LOG_JSON_FORMAT` - JSON формат (true/false)
- `LOG_FILE` - файл для логов

### Wails окно
- `WAILS_TITLE` - заголовок окна
- `WAILS_WIDTH` - ширина окна
- `WAILS_HEIGHT` - высота окна

## Dependency Injection Flow

```
main.go → Container.NewContainer()
    ↓
1. Logger initialization
    ↓  
2. Database connection + migrations
    ↓
3. Repositories (TaskRepository)
    ↓
4. Services (TaskService) 
    ↓
5. UseCases (TaskUseCase, AnalyticsUseCase, ExportUseCase)
    ↓
6. App creation with all dependencies
    ↓
7. Wails binding and startup
```

## Health Check

Приложение включает health check методы:
- `Container.HealthCheck()` - проверка всех компонентов
- `App.HealthCheck()` - проверка через Wails API

## Graceful Shutdown

При завершении работы:
- Закрывается подключение к БД
- Освобождаются все ресурсы  
- Логируется процесс завершения

## Следующие шаги

- Интеграция с фронтендом через Wails bindings
- Добавление middleware для HTTP API (если понадобится)
- Расширение аналитики и экспорта
- Добавление тестов

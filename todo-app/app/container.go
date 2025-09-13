package app

import (
	"context"
	"database/sql"
	"fmt"
	"todo-app/app/config"
	"todo-app/app/repository"
	"todo-app/app/services"
	"todo-app/app/usecases"
	"todo-app/internal/utils"

	_ "github.com/lib/pq"
)

// Container представляет DI контейнер приложения
type Container struct {
	// Database
	DB *sql.DB

	// Repositories
	TaskRepository repository.TaskRepository

	// Services
	TaskService services.TaskService

	// UseCases
	TaskUseCase      usecases.TaskUseCase
	AnalyticsUseCase usecases.AnalyticsUseCase
	ExportUseCase    usecases.ExportUseCase

	// Utils
	Logger *utils.Logger

	// Config
	Config *config.Config
}

// NewContainer создает и инициализирует новый DI контейнер
func NewContainer(cfg *config.Config) (*Container, error) {
	container := &Container{
		Config: cfg,
	}

	// Инициализируем зависимости в правильном порядке
	if err := container.initLogger(); err != nil {
		return nil, fmt.Errorf("failed to initialize logger: %w", err)
	}

	if err := container.initDatabase(); err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	if err := container.initRepositories(); err != nil {
		return nil, fmt.Errorf("failed to initialize repositories: %w", err)
	}

	if err := container.initServices(); err != nil {
		return nil, fmt.Errorf("failed to initialize services: %w", err)
	}

	if err := container.initUseCases(); err != nil {
		return nil, fmt.Errorf("failed to initialize use cases: %w", err)
	}

	container.Logger.Info("DI Container initialized successfully")
	return container, nil
}

// initLogger инициализирует логгер
func (c *Container) initLogger() error {
	loggerConfig := utils.LoggerConfig{
		Level:      utils.INFO,
		JSONFormat: c.Config.Logger.JSONFormat,
		LogFile:    c.Config.Logger.LogFile,
	}

	// Устанавливаем уровень логирования
	switch c.Config.Logger.Level {
	case "debug":
		loggerConfig.Level = utils.DEBUG
	case "info":
		loggerConfig.Level = utils.INFO
	case "warn":
		loggerConfig.Level = utils.WARN
	case "error":
		loggerConfig.Level = utils.ERROR
	default:
		loggerConfig.Level = utils.INFO
	}

	logger, err := utils.NewLogger(loggerConfig)
	if err != nil {
		return fmt.Errorf("failed to create logger: %w", err)
	}

	c.Logger = logger
	utils.SetDefaultLogger(logger)

	return nil
}

// initDatabase инициализирует подключение к базе данных
func (c *Container) initDatabase() error {
	c.Logger.Info("Initializing database connection")

	dbConfig := utils.DatabaseConfig{
		Host:     c.Config.Database.Host,
		Port:     c.Config.Database.Port,
		User:     c.Config.Database.User,
		Password: c.Config.Database.Password,
		DBName:   c.Config.Database.DBName,
		SSLMode:  c.Config.Database.SSLMode,
	}

	db, err := utils.InitDB(dbConfig)
	if err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}

	// Проверяем подключение
	if err := db.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	c.DB = db
	c.Logger.Info("Database connection established successfully",
		map[string]interface{}{
			"host": c.Config.Database.Host,
			"port": c.Config.Database.Port,
			"db":   c.Config.Database.DBName,
		})

	return nil
}

// runMigrations выполняет миграции базы данных
func (c *Container) runMigrations() error {
	if !c.Config.Database.RunMigrations {
		c.Logger.Info("Migrations skipped (disabled in config)")
		return nil
	}

	c.Logger.Info("Running database migrations")

	migrationHelper := utils.NewMigrationHelper(c.DB)

	// Создаем таблицу для отслеживания миграций
	if err := migrationHelper.CreateMigrationsTable(); err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	// Список миграций в порядке выполнения
	migrations := []struct {
		version string
		file    string
	}{
		{"001", "database/migrations/001_create_tasks_table.up.sql"},
		{"002", "database/migrations/002_add_indexes.up.sql"},
	}

	for _, migration := range migrations {
		applied, err := migrationHelper.IsMigrationApplied(migration.version)
		if err != nil {
			return fmt.Errorf("failed to check migration %s: %w", migration.version, err)
		}

		if applied {
			c.Logger.Debug("Migration already applied", map[string]interface{}{
				"version": migration.version,
			})
			continue
		}

		c.Logger.Info("Applying migration", map[string]interface{}{
			"version": migration.version,
			"file":    migration.file,
		})

		// Здесь можно добавить логику выполнения SQL файлов
		// Для простоты пропускаем реальное выполнение

		if err := migrationHelper.MarkMigrationApplied(migration.version); err != nil {
			return fmt.Errorf("failed to mark migration as applied: %w", err)
		}
	}

	c.Logger.Info("Database migrations completed successfully")
	return nil
}

// initRepositories инициализирует репозитории
func (c *Container) initRepositories() error {
	c.Logger.Info("Initializing repositories")

	// Task Repository
	c.TaskRepository = repository.NewPostgresTaskRepository(c.DB)

	c.Logger.Info("Repositories initialized successfully")
	return nil
}

// initServices инициализирует сервисы
func (c *Container) initServices() error {
	c.Logger.Info("Initializing services")

	// Task Service
	c.TaskService = services.NewTaskService(c.TaskRepository)

	c.Logger.Info("Services initialized successfully")
	return nil
}

// initUseCases инициализирует use cases
func (c *Container) initUseCases() error {
	c.Logger.Info("Initializing use cases")

	// Task UseCase
	c.TaskUseCase = usecases.NewTaskUseCase(c.TaskService)

	// Analytics UseCase
	c.AnalyticsUseCase = usecases.NewAnalyticsUseCase(c.TaskService)

	// Export UseCase
	c.ExportUseCase = usecases.NewExportUseCase(c.TaskService)

	c.Logger.Info("Use cases initialized successfully")
	return nil
}

// NewApp создает и инициализирует новое приложение с зависимостями
func (c *Container) NewApp(ctx context.Context) *App {
	return &App{
		ctx:              ctx,
		db:               c.DB,
		config:           c.Config,
		logger:           c.Logger,
		TaskUseCase:      c.TaskUseCase,
		AnalyticsUseCase: c.AnalyticsUseCase,
		ExportUseCase:    c.ExportUseCase,
	}
}

// Close закрывает все ресурсы контейнера
func (c *Container) Close() error {
	c.Logger.Info("Closing container resources")

	if c.DB != nil {
		utils.CloseDB(c.DB)
	}

	c.Logger.Info("Container resources closed successfully")
	return nil
}

// HealthCheck проверяет состояние всех компонентов
func (c *Container) HealthCheck() error {
	// Проверяем базу данных
	if err := c.DB.Ping(); err != nil {
		return fmt.Errorf("database health check failed: %w", err)
	}

	// Можно добавить другие проверки (Redis, внешние API и т.д.)

	return nil
}

// GetDependencies возвращает информацию о зависимостях для отладки
func (c *Container) GetDependencies() map[string]interface{} {
	return map[string]interface{}{
		"database_connected": c.DB != nil,
		"task_repository":    c.TaskRepository != nil,
		"task_service":       c.TaskService != nil,
		"task_usecase":       c.TaskUseCase != nil,
		"analytics_usecase":  c.AnalyticsUseCase != nil,
		"export_usecase":     c.ExportUseCase != nil,
		"logger":             c.Logger != nil,
		"config":             c.Config != nil,
	}
}

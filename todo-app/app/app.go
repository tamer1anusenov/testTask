package app

import (
	"context"
	"database/sql"
	"todo-app/app/config"
	"todo-app/app/usecases"
	"todo-app/internal/utils"
)

// App представляет основную структуру приложения
type App struct {
	ctx              context.Context
	db               *sql.DB
	config           *config.Config
	logger           *utils.Logger
	TaskUseCase      usecases.TaskUseCase
	AnalyticsUseCase usecases.AnalyticsUseCase
	ExportUseCase    usecases.ExportUseCase
}

// GetContext возвращает контекст приложения
func (a *App) GetContext() context.Context {
	return a.ctx
}

// GetDB возвращает подключение к базе данных
func (a *App) GetDB() *sql.DB {
	return a.db
}

// GetConfig возвращает конфигурацию приложения
func (a *App) GetConfig() *config.Config {
	return a.config
}

// GetLogger возвращает логгер приложения
func (a *App) GetLogger() *utils.Logger {
	return a.logger
}

// Startup вызывается при запуске приложения
func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
	if a.logger != nil {
		a.logger.Info("Application started successfully")
	}
}

// Shutdown вызывается при завершении приложения
func (a *App) Shutdown(ctx context.Context) {
	if a.logger != nil {
		a.logger.Info("Application shutting down")
	}
}

// HealthCheck проверяет состояние приложения
func (a *App) HealthCheck() error {
	if a.db != nil {
		return a.db.Ping()
	}
	return nil
}

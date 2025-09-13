package main

import (
	"context"
	"embed"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"todo-app/app"
	"todo-app/app/config"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// Загружаем конфигурацию
	cfg, err := loadConfiguration()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Выводим конфигурацию (без чувствительных данных)
	if cfg.IsDevelopment() {
		cfg.Print()
	}

	// Создаем DI контейнер
	container, err := app.NewContainer(cfg)
	if err != nil {
		log.Fatalf("Failed to create container: %v", err)
	}

	// Настраиваем graceful shutdown
	defer func() {
		if err := container.Close(); err != nil {
			log.Printf("Error closing container: %v", err)
		}
	}()

	// Выполняем миграции если необходимо
	if err := runMigrations(container); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Создаем экземпляр приложения для Wails
	ctx := context.Background()
	wailsApp := NewApp()

	// Инициализируем Wails App с backend dependencies
	wailsApp.ctx = ctx
	wailsApp.db = container.DB
	wailsApp.config = cfg
	wailsApp.logger = container.Logger
	wailsApp.TaskUseCase = container.TaskUseCase
	wailsApp.AnalyticsUseCase = container.AnalyticsUseCase
	wailsApp.ExportUseCase = container.ExportUseCase

	// Настраиваем Wails опции
	wailsOptions := buildWailsOptions(cfg, wailsApp)

	// Запускаем приложение
	if err := wails.Run(wailsOptions); err != nil {
		log.Fatalf("Error running application: %v", err)
	}
}

// loadConfiguration загружает конфигурацию приложения
func loadConfiguration() (*config.Config, error) {
	// Пробуем загрузить из переменных окружения
	cfg, err := config.LoadFromEnv()
	if err != nil {
		return nil, fmt.Errorf("failed to load config from environment: %w", err)
	}

	// Валидируем конфигурацию
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return cfg, nil
}

// runMigrations выполняет миграции базы данных
func runMigrations(container *app.Container) error {
	// Проверяем health состояние перед миграциями
	if err := container.HealthCheck(); err != nil {
		return fmt.Errorf("health check failed before migrations: %w", err)
	}

	// В container.go уже есть логика миграций, но здесь можно добавить
	// дополнительную обработку если необходимо

	return nil
}

// buildWailsOptions создает опции для Wails приложения
func buildWailsOptions(cfg *config.Config, appInstance *App) *options.App {
	return &options.App{
		Title:  cfg.Wails.Title,
		Width:  cfg.Wails.Width,
		Height: cfg.Wails.Height,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        appInstance.Startup,
		OnShutdown:       appInstance.Shutdown,
		Bind: []interface{}{
			appInstance,
		},
		// Дополнительные настройки окна
		WindowStartState:  getWindowStartState(cfg),
		DisableResize:     !cfg.Wails.Window.Resizable,
		Fullscreen:        cfg.Wails.Window.Fullscreen,
		AlwaysOnTop:       cfg.Wails.Window.AlwaysOnTop,
		HideWindowOnClose: cfg.Wails.Window.HideWindowOnClose,
	}
}

// getWindowStartState определяет начальное состояние окна
func getWindowStartState(cfg *config.Config) options.WindowStartState {
	if cfg.Wails.Window.Maximized {
		return options.Maximised
	}
	if cfg.Wails.Window.Fullscreen {
		return options.Fullscreen
	}
	return options.Normal
}

// setupGracefulShutdown настраивает graceful shutdown
func setupGracefulShutdown(container *app.Container) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		log.Println("Received shutdown signal")

		if err := container.Close(); err != nil {
			log.Printf("Error during shutdown: %v", err)
		}

		os.Exit(0)
	}()
}

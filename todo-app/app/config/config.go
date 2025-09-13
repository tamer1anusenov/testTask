package config

import (
	"fmt"
	"os"
	"strconv"
)

// Config представляет конфигурацию приложения
type Config struct {
	App      AppConfig      `yaml:"app"`
	Database DatabaseConfig `yaml:"database"`
	Logger   LoggerConfig   `yaml:"logger"`
	Wails    WailsConfig    `yaml:"wails"`
}

// AppConfig содержит настройки приложения
type AppConfig struct {
	Name        string `yaml:"name"`
	Version     string `yaml:"version"`
	Environment string `yaml:"environment"`
	Debug       bool   `yaml:"debug"`
}

// DatabaseConfig содержит настройки базы данных
type DatabaseConfig struct {
	Host          string `yaml:"host"`
	Port          int    `yaml:"port"`
	User          string `yaml:"user"`
	Password      string `yaml:"password"`
	DBName        string `yaml:"dbname"`
	SSLMode       string `yaml:"sslmode"`
	RunMigrations bool   `yaml:"run_migrations"`
}

// LoggerConfig содержит настройки логирования
type LoggerConfig struct {
	Level      string `yaml:"level"`
	JSONFormat bool   `yaml:"json_format"`
	LogFile    string `yaml:"log_file"`
}

// WailsConfig содержит настройки Wails приложения
type WailsConfig struct {
	Title  string       `yaml:"title"`
	Width  int          `yaml:"width"`
	Height int          `yaml:"height"`
	Window WindowConfig `yaml:"window"`
}

// WindowConfig содержит настройки окна приложения
type WindowConfig struct {
	Resizable         bool `yaml:"resizable"`
	Maximized         bool `yaml:"maximized"`
	Fullscreen        bool `yaml:"fullscreen"`
	AlwaysOnTop       bool `yaml:"always_on_top"`
	HideWindowOnClose bool `yaml:"hide_window_on_close"`
}

// NewDefaultConfig создает конфигурацию по умолчанию
func NewDefaultConfig() *Config {
	return &Config{
		App: AppConfig{
			Name:        "Todo App",
			Version:     "1.0.0",
			Environment: "development",
			Debug:       true,
		},
		Database: DatabaseConfig{
			Host:          "localhost",
			Port:          5432,
			User:          "todo_user",
			Password:      "1234",
			DBName:        "todo_db",
			SSLMode:       "disable",
			RunMigrations: true,
		},
		Logger: LoggerConfig{
			Level:      "info",
			JSONFormat: false,
			LogFile:    "",
		},
		Wails: WailsConfig{
			Title:  "Todo App",
			Width:  1024,
			Height: 768,
			Window: WindowConfig{
				Resizable:         true,
				Maximized:         false,
				Fullscreen:        false,
				AlwaysOnTop:       false,
				HideWindowOnClose: false,
			},
		},
	}
}

// LoadFromEnv загружает конфигурацию из переменных окружения
func LoadFromEnv() (*Config, error) {
	config := NewDefaultConfig()

	// App settings
	if env := os.Getenv("APP_NAME"); env != "" {
		config.App.Name = env
	}
	if env := os.Getenv("APP_VERSION"); env != "" {
		config.App.Version = env
	}
	if env := os.Getenv("APP_ENVIRONMENT"); env != "" {
		config.App.Environment = env
	}
	if env := os.Getenv("APP_DEBUG"); env != "" {
		if debug, err := strconv.ParseBool(env); err == nil {
			config.App.Debug = debug
		}
	}

	// Database settings
	if env := os.Getenv("DB_HOST"); env != "" {
		config.Database.Host = env
	}
	if env := os.Getenv("DB_PORT"); env != "" {
		if port, err := strconv.Atoi(env); err == nil {
			config.Database.Port = port
		}
	}
	if env := os.Getenv("DB_USER"); env != "" {
		config.Database.User = env
	}
	if env := os.Getenv("DB_PASSWORD"); env != "" {
		config.Database.Password = env
	}
	if env := os.Getenv("DB_NAME"); env != "" {
		config.Database.DBName = env
	}
	if env := os.Getenv("DB_SSLMODE"); env != "" {
		config.Database.SSLMode = env
	}
	if env := os.Getenv("DB_RUN_MIGRATIONS"); env != "" {
		if runMigrations, err := strconv.ParseBool(env); err == nil {
			config.Database.RunMigrations = runMigrations
		}
	}

	// Logger settings
	if env := os.Getenv("LOG_LEVEL"); env != "" {
		config.Logger.Level = env
	}
	if env := os.Getenv("LOG_JSON_FORMAT"); env != "" {
		if jsonFormat, err := strconv.ParseBool(env); err == nil {
			config.Logger.JSONFormat = jsonFormat
		}
	}
	if env := os.Getenv("LOG_FILE"); env != "" {
		config.Logger.LogFile = env
	}

	// Wails settings
	if env := os.Getenv("WAILS_TITLE"); env != "" {
		config.Wails.Title = env
	}
	if env := os.Getenv("WAILS_WIDTH"); env != "" {
		if width, err := strconv.Atoi(env); err == nil {
			config.Wails.Width = width
		}
	}
	if env := os.Getenv("WAILS_HEIGHT"); env != "" {
		if height, err := strconv.Atoi(env); err == nil {
			config.Wails.Height = height
		}
	}

	return config, nil
}

// Validate проверяет корректность конфигурации
func (c *Config) Validate() error {
	if c.App.Name == "" {
		return fmt.Errorf("app name cannot be empty")
	}

	if c.Database.Host == "" {
		return fmt.Errorf("database host cannot be empty")
	}

	if c.Database.Port <= 0 || c.Database.Port > 65535 {
		return fmt.Errorf("database port must be between 1 and 65535")
	}

	if c.Database.User == "" {
		return fmt.Errorf("database user cannot be empty")
	}

	if c.Database.DBName == "" {
		return fmt.Errorf("database name cannot be empty")
	}

	validLogLevels := map[string]bool{
		"debug": true,
		"info":  true,
		"warn":  true,
		"error": true,
		"fatal": true,
	}

	if !validLogLevels[c.Logger.Level] {
		return fmt.Errorf("invalid log level: %s", c.Logger.Level)
	}

	if c.Wails.Width <= 0 {
		return fmt.Errorf("window width must be positive")
	}

	if c.Wails.Height <= 0 {
		return fmt.Errorf("window height must be positive")
	}

	return nil
}

// IsDevelopment проверяет, находится ли приложение в режиме разработки
func (c *Config) IsDevelopment() bool {
	return c.App.Environment == "development"
}

// IsProduction проверяет, находится ли приложение в продакшене
func (c *Config) IsProduction() bool {
	return c.App.Environment == "production"
}

// GetDatabaseDSN возвращает строку подключения к базе данных
func (c *Config) GetDatabaseDSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Database.Host,
		c.Database.Port,
		c.Database.User,
		c.Database.Password,
		c.Database.DBName,
		c.Database.SSLMode,
	)
}

// Print выводит конфигурацию в консоль (без чувствительных данных)
func (c *Config) Print() {
	fmt.Printf("App Configuration:\n")
	fmt.Printf("  Name: %s\n", c.App.Name)
	fmt.Printf("  Version: %s\n", c.App.Version)
	fmt.Printf("  Environment: %s\n", c.App.Environment)
	fmt.Printf("  Debug: %t\n", c.App.Debug)
	fmt.Printf("Database Configuration:\n")
	fmt.Printf("  Host: %s\n", c.Database.Host)
	fmt.Printf("  Port: %d\n", c.Database.Port)
	fmt.Printf("  User: %s\n", c.Database.User)
	fmt.Printf("  Database: %s\n", c.Database.DBName)
	fmt.Printf("  SSL Mode: %s\n", c.Database.SSLMode)
	fmt.Printf("  Run Migrations: %t\n", c.Database.RunMigrations)
	fmt.Printf("Logger Configuration:\n")
	fmt.Printf("  Level: %s\n", c.Logger.Level)
	fmt.Printf("  JSON Format: %t\n", c.Logger.JSONFormat)
	fmt.Printf("  Log File: %s\n", c.Logger.LogFile)
	fmt.Printf("Wails Configuration:\n")
	fmt.Printf("  Title: %s\n", c.Wails.Title)
	fmt.Printf("  Size: %dx%d\n", c.Wails.Width, c.Wails.Height)
}

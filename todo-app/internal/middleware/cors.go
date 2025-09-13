package middleware

import (
	"fmt"
	"net/http"
	"strings"
)

// CORSConfig содержит настройки CORS
type CORSConfig struct {
	AllowedOrigins     []string
	AllowedMethods     []string
	AllowedHeaders     []string
	ExposedHeaders     []string
	AllowCredentials   bool
	MaxAge             int
	OptionsPassthrough bool
}

// DefaultCORSConfig возвращает настройки CORS по умолчанию
func DefaultCORSConfig() CORSConfig {
	return CORSConfig{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
			http.MethodHead,
			http.MethodOptions,
		},
		AllowedHeaders: []string{
			"Accept",
			"Authorization",
			"Content-Type",
			"X-CSRF-Token",
			"X-Requested-With",
		},
		ExposedHeaders:     []string{},
		AllowCredentials:   false,
		MaxAge:             86400, // 24 часа
		OptionsPassthrough: false,
	}
}

// DevelopmentCORSConfig возвращает настройки CORS для разработки
func DevelopmentCORSConfig() CORSConfig {
	return CORSConfig{
		AllowedOrigins: []string{
			"http://localhost:3000",
			"http://localhost:5173",
			"http://127.0.0.1:3000",
			"http://127.0.0.1:5173",
			"wails://wails",
		},
		AllowedMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
			http.MethodHead,
			http.MethodOptions,
		},
		AllowedHeaders: []string{
			"Accept",
			"Authorization",
			"Content-Type",
			"X-CSRF-Token",
			"X-Requested-With",
			"Origin",
			"Cache-Control",
			"X-File-Name",
		},
		ExposedHeaders: []string{
			"Content-Length",
			"Content-Type",
			"Date",
		},
		AllowCredentials:   true,
		MaxAge:             300, // 5 минут для разработки
		OptionsPassthrough: false,
	}
}

// ProductionCORSConfig возвращает безопасные настройки CORS для продакшена
func ProductionCORSConfig(allowedOrigins []string) CORSConfig {
	return CORSConfig{
		AllowedOrigins: allowedOrigins,
		AllowedMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
			http.MethodOptions,
		},
		AllowedHeaders: []string{
			"Accept",
			"Authorization",
			"Content-Type",
			"X-Requested-With",
		},
		ExposedHeaders:     []string{},
		AllowCredentials:   true,
		MaxAge:             86400, // 24 часа
		OptionsPassthrough: false,
	}
}

// CORS создает middleware для обработки CORS
func CORS(config CORSConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")

			// Проверяем, разрешен ли origin
			if isOriginAllowed(origin, config.AllowedOrigins) {
				w.Header().Set("Access-Control-Allow-Origin", origin)
			} else if len(config.AllowedOrigins) == 1 && config.AllowedOrigins[0] == "*" {
				w.Header().Set("Access-Control-Allow-Origin", "*")
			}

			// Устанавливаем разрешенные методы
			if len(config.AllowedMethods) > 0 {
				w.Header().Set("Access-Control-Allow-Methods", strings.Join(config.AllowedMethods, ", "))
			}

			// Устанавливаем разрешенные заголовки
			if len(config.AllowedHeaders) > 0 {
				w.Header().Set("Access-Control-Allow-Headers", strings.Join(config.AllowedHeaders, ", "))
			}

			// Устанавливаем exposed заголовки
			if len(config.ExposedHeaders) > 0 {
				w.Header().Set("Access-Control-Expose-Headers", strings.Join(config.ExposedHeaders, ", "))
			}

			// Устанавливаем credentials
			if config.AllowCredentials {
				w.Header().Set("Access-Control-Allow-Credentials", "true")
			}

			// Устанавливаем max age
			if config.MaxAge > 0 {
				w.Header().Set("Access-Control-Max-Age", fmt.Sprintf("%d", config.MaxAge))
			}

			// Обрабатываем preflight запросы
			if r.Method == http.MethodOptions {
				if config.OptionsPassthrough {
					next.ServeHTTP(w, r)
				} else {
					w.WriteHeader(http.StatusNoContent)
				}
				return
			}

			// Продолжаем выполнение цепочки middleware
			next.ServeHTTP(w, r)
		})
	}
}

// isOriginAllowed проверяет, разрешен ли origin
func isOriginAllowed(origin string, allowedOrigins []string) bool {
	for _, allowedOrigin := range allowedOrigins {
		if allowedOrigin == "*" {
			return true
		}
		if allowedOrigin == origin {
			return true
		}
		// Поддержка wildcard поддоменов (например, *.example.com)
		if strings.HasPrefix(allowedOrigin, "*.") {
			domain := allowedOrigin[2:]
			if strings.HasSuffix(origin, "."+domain) || origin == domain {
				return true
			}
		}
	}
	return false
}

// WebSocketCORS создает специализированный CORS middleware для WebSocket соединений
func WebSocketCORS(allowedOrigins []string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")

			if !isOriginAllowed(origin, allowedOrigins) {
				http.Error(w, "Origin not allowed", http.StatusForbidden)
				return
			}

			// Для WebSocket соединений не нужно устанавливать CORS заголовки
			// так как браузер их не проверяет для WebSocket
			next.ServeHTTP(w, r)
		})
	}
}

// SecurityHeaders добавляет дополнительные заголовки безопасности
func SecurityHeaders() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Защита от XSS
			w.Header().Set("X-Content-Type-Options", "nosniff")
			w.Header().Set("X-Frame-Options", "DENY")
			w.Header().Set("X-XSS-Protection", "1; mode=block")

			// Контроль referrer
			w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")

			// Content Security Policy (базовый)
			w.Header().Set("Content-Security-Policy", "default-src 'self'")

			next.ServeHTTP(w, r)
		})
	}
}

// APIKeyValidator создает middleware для проверки API ключей
func APIKeyValidator(validKeys []string, headerName string) func(http.Handler) http.Handler {
	if headerName == "" {
		headerName = "X-API-Key"
	}

	keyMap := make(map[string]bool)
	for _, key := range validKeys {
		keyMap[key] = true
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			apiKey := r.Header.Get(headerName)

			if apiKey == "" {
				http.Error(w, "API key required", http.StatusUnauthorized)
				return
			}

			if !keyMap[apiKey] {
				http.Error(w, "Invalid API key", http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// RateLimiter создает простой rate limiter middleware
type RateLimiterConfig struct {
	RequestsPerMinute int
	BurstSize         int
}

// Примечание: Это упрощенная реализация rate limiter'а
// В продакшене рекомендуется использовать более сложные решения
// с поддержкой Redis или других внешних хранилищ
func RateLimiter(config RateLimiterConfig) func(http.Handler) http.Handler {
	// Здесь должна быть реализация rate limiter'а
	// Для простоты оставляем заглушку
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// TODO: Реализовать rate limiting
			// clientIP := r.RemoteAddr
			// if isRateLimited(clientIP, config) {
			//     http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			//     return
			// }

			next.ServeHTTP(w, r)
		})
	}
}

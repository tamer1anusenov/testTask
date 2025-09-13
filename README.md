# testTask

# Todo App - Desktop приложение для управления задачами

## Описание
Кросс-платформенное desktop приложение на Wails для управления списком задач и мини дэшбордом.

## Технологии
- Backend: Go + Wails
- Frontend: HTML/CSS/JavaScript
- База данных: PostgreSQL
- Архитектура: Repository → Service → UseCase

## Установка и запуск

### Требования
- Go 1.19+
- Node.js 16+
- Wails CLI
- PostgreSQL 12+

### Инструкции по запуску

#### Быстрый старт с Makefile (рекомендуется)
```bash
# Клонировать репозиторий
git clone <https://github.com/tamer1anusenov/testTask.git>
cd testTask/todo-app

# Полная настройка одной командой (PostgreSQL + БД + зависимости)
make full-setup

# Создать .env файл (см. секцию "Конфигурация базы данных")
# Затем запустить приложение
make dev
```

#### Ручная настройка

#### 1. Подготовка PostgreSQL
```bash
# Автоматически с Makefile
make setup-db

# Или вручную:
sudo systemctl start postgresql
sudo systemctl enable postgresql

# Создайте базу данных и пользователя (используйте те же данные, что в .env файле)
sudo -u postgres psql
CREATE DATABASE todoapp;
CREATE USER todouser WITH ENCRYPTED PASSWORD '1234';
GRANT ALL PRIVILEGES ON DATABASE todoapp TO todouser;
\q
```

#### 2. Настройка проекта
```bash
# Клонировать репозиторий
git clone <repository-url>
cd testTask/todo-app

# Установить зависимости Go
go mod download

# Установить зависимости frontend
cd frontend
npm install
cd ..
```

#### 3. Конфигурация базы данных
Создайте файл `.env` в корне проекта `todo-app/` и укажите данные для подключения к PostgreSQL:
```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=todouser
DB_PASSWORD=1234
DB_NAME=todoapp
DB_SSLMODE=disable
```

#### 4. Установка Wails CLI (если не установлен)
```bash
go install github.com/wailsapp/wails/v2/cmd/wails@latest
```

#### 5. Запуск приложения

**С использованием Makefile:**
```bash
# Для разработки
make dev

# Для сборки production версии
make build

# Проверить статус PostgreSQL
make check-postgres

# Перезапустить PostgreSQL
make restart-postgres
```

**Или напрямую с Wails:**
```bash
# Для разработки
wails dev

# Для сборки production версии
wails build
```

### Доступные команды Makefile
```bash
make help           # Показать все доступные команды
make setup-db       # Настроить PostgreSQL и создать БД
make start-postgres # Запустить PostgreSQL
make stop-postgres  # Остановить PostgreSQL
make dev           # Запустить в режиме разработки
make build         # Собрать production версию
make clean         # Очистить build артефакты
make test-db       # Проверить подключение к БД
make full-setup    # Полная настройка для новых пользователей
```

#### Проверка подключения к базе данных
Приложение автоматически применит миграции при первом запуске. Если возникают проблемы с подключением, проверьте:
- PostgreSQL запущен: `sudo systemctl status postgresql`
- Данные в `.env` файле корректны
- Пользователь имеет права доступа к базе данных


## Video
[Watch Demo](https://github.com/user-attachments/assets/789a6f4c-1951-4291-91b6-def33db13033)

## Checklist выполненных функций

### Интерфейс пользователя 
- Текстовое поле для ввода новой задачи
- Кнопка добавления задачи
- Отображение списка задач
- CSS стилизация
- Значки и цвета для статусов
- Адаптивная верстка
- Переключение светлой/темной темы

### ✅ Добавление задач 
- Добавление новых задач
- Валидация ввода
- Добавление даты и времени выполнения
- Установка приоритета задачи

### ✅ Удаление задач 
- Удаление задач из списка
- Подтверждение удаления (модальное окно)

### ✅ Управление выполнением задач 
- Отметка задачи как выполненной
- Зачеркивание текста выполненных задач
- Перемещение в раздел "Выполненные задачи"
- Возврат в "Активные задачи"

### ✅ Сохранение состояния 
- Сохранение состояния при закрытии
- Загрузка состояния при запуске
- Использование PostgreSQL
- Архитектура repo → service → usecase

### ✅ Фильтрация и сортировка 
- Фильтрация по статусу
- Сортировка по дате добавления
- Сортировка по приоритету
- Фильтрация по дате

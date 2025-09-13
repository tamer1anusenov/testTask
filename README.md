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

#### 1. Подготовка PostgreSQL
```bash
# Убедитесь, что PostgreSQL запущен и активен
sudo systemctl start postgresql
sudo systemctl enable postgresql

# Создайте базу данных и пользователя
sudo -u postgres psql
CREATE DATABASE todoapp;
CREATE USER todouser WITH ENCRYPTED PASSWORD 'your_password';
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
DB_PASSWORD=your_password
DB_NAME=todoapp
DB_SSLMODE=disable
```

#### 4. Установка Wails CLI (если не установлен)
```bash
go install github.com/wailsapp/wails/v2/cmd/wails@latest
```

#### 5. Запуск приложения

**Для разработки:**
```bash
wails dev
```

**Для сборки production версии:**
```bash
wails build
```

#### Проверка подключения к базе данных
Приложение автоматически применит миграции при первом запуске. Если возникают проблемы с подключением, проверьте:
- PostgreSQL запущен: `sudo systemctl status postgresql`
- Данные в `.env` файле корректны
- Пользователь имеет права доступа к базе данных

## Checklist выполненных функций

### ✅ Интерфейс пользователя 
- ✅ Текстовое поле для ввода новой задачи
- ✅ Кнопка добавления задачи
- ✅ Отображение списка задач
- ✅ CSS стилизация
- ✅ Значки и цвета для статусов
- ✅ Адаптивная верстка
- ✅ Переключение светлой/темной темы

### ✅ Добавление задач 
- ✅ Добавление новых задач
- ✅ Валидация ввода
- ✅ Добавление даты и времени выполнения
- ✅ Установка приоритета задачи

### ✅ Удаление задач 
- ✅ Удаление задач из списка
- ✅ Подтверждение удаления (модальное окно)

### ✅ Управление выполнением задач 
- ✅ Отметка задачи как выполненной
- ✅ Зачеркивание текста выполненных задач
- ✅ Перемещение в раздел "Выполненные задачи"
- ✅ Возврат в "Активные задачи"

### ✅ Сохранение состояния 
- ✅ Сохранение состояния при закрытии
- ✅ Загрузка состояния при запуске
- ✅ Использование PostgreSQL
- ✅ Архитектура repo → service → usecase

### ✅ Фильтрация и сортировка 
- ✅ Фильтрация по статусу
- ✅ Сортировка по дате добавления
- ✅ Сортировка по приоритету
- ✅ Фильтрация по дате

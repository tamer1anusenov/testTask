-- Добавляем поле archived в таблицу tasks
ALTER TABLE tasks ADD COLUMN archived BOOLEAN NOT NULL DEFAULT FALSE;

-- Создаем индекс для поля archived
CREATE INDEX idx_tasks_archived ON tasks(archived);

-- Создаем составной индекс для статуса и архива
CREATE INDEX idx_tasks_status_archived ON tasks(status, archived);

-- Создаем составной индекс для приоритета и архива
CREATE INDEX idx_tasks_priority_archived ON tasks(priority, archived);

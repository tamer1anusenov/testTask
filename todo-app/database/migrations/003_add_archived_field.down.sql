-- Удаляем индексы
DROP INDEX IF EXISTS idx_tasks_priority_archived;
DROP INDEX IF EXISTS idx_tasks_status_archived;
DROP INDEX IF EXISTS idx_tasks_archived;

-- Удаляем поле archived
ALTER TABLE tasks DROP COLUMN archived;

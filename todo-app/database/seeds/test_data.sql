-- Test data for todo app
-- This file contains sample data for testing purposes

-- Clear existing data
TRUNCATE tasks RESTART IDENTITY CASCADE;

-- Insert test tasks
INSERT INTO tasks (title, description, status, priority, created_at, updated_at) VALUES
('Test Task 1', 'This is a test task for unit testing', 'active', 'medium', NOW(), NOW()),
('Test Task 2', 'Another test task with high priority', 'active', 'high', NOW(), NOW()),
('Completed Test Task', 'This task is already completed', 'completed', 'low', NOW() - INTERVAL '1 day', NOW()),
('Task with Due Date', 'This task has a due date', 'active', 'high', NOW(), NOW()),
('Overdue Task', 'This task is overdue', 'active', 'medium', NOW() - INTERVAL '3 days', NOW() - INTERVAL '3 days'),
('Long Description Task', 'This task has a very long description to test how the application handles longer text content. It should be displayed properly in the UI and stored correctly in the database without any issues.', 'active', 'low', NOW(), NOW()),
('Task for Update Test', 'This task will be used to test update operations', 'active', 'medium', NOW(), NOW()),
('Task for Delete Test', 'This task will be used to test delete operations', 'active', 'low', NOW(), NOW()),
('Task for Toggle Test', 'This task will be used to test status toggle operations', 'active', 'high', NOW(), NOW()),
('Analytics Test Task 1', 'Task for analytics testing - completed yesterday', 'completed', 'medium', NOW() - INTERVAL '1 day', NOW() - INTERVAL '1 day'),
('Analytics Test Task 2', 'Task for analytics testing - completed today', 'completed', 'high', NOW(), NOW()),
('Weekly Analytics Task', 'Task for weekly analytics testing', 'completed', 'low', NOW() - INTERVAL '5 days', NOW() - INTERVAL '5 days');

-- Update some tasks with due dates
UPDATE tasks SET due_date = NOW() + INTERVAL '3 days' WHERE title = 'Task with Due Date';
UPDATE tasks SET due_date = NOW() - INTERVAL '1 day' WHERE title = 'Overdue Task';
UPDATE tasks SET completed_at = NOW() WHERE status = 'completed';

-- Add some tasks with various creation dates for analytics
INSERT INTO tasks (title, description, status, priority, created_at, updated_at, completed_at) VALUES
('Week Ago Task', 'Task created a week ago', 'completed', 'medium', NOW() - INTERVAL '7 days', NOW() - INTERVAL '6 days', NOW() - INTERVAL '6 days'),
('Month Ago Task', 'Task created a month ago', 'completed', 'low', NOW() - INTERVAL '30 days', NOW() - INTERVAL '29 days', NOW() - INTERVAL '29 days'),
('Recent Active Task', 'Recently created active task', 'active', 'high', NOW() - INTERVAL '2 hours', NOW() - INTERVAL '2 hours', NULL);

-- Commit the changes
COMMIT;

UPDATE tasks SET description = '' WHERE description IS NULL;

ALTER TABLE tasks MODIFY description TEXT NOT NULL;

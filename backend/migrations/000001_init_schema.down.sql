DROP INDEX IF EXISTS idx_wishes_status;
DROP INDEX IF EXISTS idx_wishes_child_id;
DROP INDEX IF EXISTS idx_tasks_status;
DROP INDEX IF EXISTS idx_tasks_parent_id;
DROP INDEX IF EXISTS idx_tasks_child_id;
DROP INDEX IF EXISTS idx_children_parent_id;

DROP TABLE IF EXISTS wishes;
DROP TABLE IF EXISTS tasks;
DROP TABLE IF EXISTS children;
DROP TABLE IF EXISTS parents;
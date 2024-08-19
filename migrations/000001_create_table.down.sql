-- Drop the tables in reverse order of dependencies
DROP TABLE IF EXISTS audit_logs;
DROP TABLE IF EXISTS black_list;
DROP TABLE IF EXISTS employees;
DROP TABLE IF EXISTS hr;

-- Drop the ENUM type
DROP TYPE IF EXISTS action;

-- Database initialization script
-- This script creates the database and user for the BezBase application

-- Create database
SELECT 'CREATE DATABASE bezbase'
WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'bezbase')\gexec

-- Create user (optional, you can use existing user)
DO $$
BEGIN
    IF NOT EXISTS (SELECT FROM pg_catalog.pg_roles WHERE rolname = 'bezbase_user') THEN
        CREATE USER bezbase_user WITH PASSWORD 'bezbase_password';
    END IF;
END
$$;

-- Grant privileges
GRANT ALL PRIVILEGES ON DATABASE bezbase TO bezbase_user;

-- Connect to the database
\c bezbase;

-- Grant schema privileges
GRANT ALL ON SCHEMA public TO bezbase_user;
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO bezbase_user;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO bezbase_user;


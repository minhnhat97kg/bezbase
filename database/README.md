# Database Setup Instructions

## Prerequisites
- PostgreSQL installed and running
- psql command-line tool available

## Setup Steps

### 1. Create Database and User
Run the initialization script to create the database and user:

```bash
psql -U postgres -f database/init.sql
```

### 2. Run Migrations
Apply the database migrations:

```bash
psql -U bezbase_user -d bezbase -f database/migrations/001_create_users_table.sql
```

### 3. Environment Configuration
Copy the environment file and update the database connection string:

```bash
cp backend/.env.example backend/.env
```

Update the `DATABASE_URL` in the `.env` file:
```
DATABASE_URL=postgres://bezbase_user:bezbase_password@localhost/bezbase?sslmode=disable
```

## Database Schema

### Users Table
- `id`: Primary key (SERIAL)
- `email`: Unique email address (VARCHAR(255))
- `password`: Hashed password (VARCHAR(255))
- `first_name`: User's first name (VARCHAR(100))
- `last_name`: User's last name (VARCHAR(100))
- `created_at`: Creation timestamp (TIMESTAMP)
- `updated_at`: Last update timestamp (TIMESTAMP)

## Connection Details
- **Database**: bezbase
- **User**: bezbase_user
- **Password**: bezbase_password
- **Host**: localhost
- **Port**: 5432

## Useful Commands

### Connect to database
```bash
psql -U bezbase_user -d bezbase -h localhost
```

### Check tables
```sql
\dt
```

### Check user table structure
```sql
\d users
```

### View all users
```sql
SELECT id, email, first_name, last_name, created_at FROM users;
```

## Security Notes
- Change the default password in production
- Use environment variables for sensitive data
- Consider using connection pooling for production
- Enable SSL in production environments


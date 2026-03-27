# 🛒 UniBazar — University Marketplace




---

## Installation

### 1. Clone the repository

download

### 2. Install Go dependencies

```bash
go mod tidy
```

### 3. Install Node dependencies (for CSS build)

```bash
npm install
```

This installs `tailwindcss` and `@tailwindcss/cli` as defined in `package.json`.

### 4. Configure the database connection

Edit `internal/db/db.go` and update the connection string: and make sure that you have provided the required port and sometimes there is possiblity of blocked port or dynamic port and you can check them on the computer management 

```go
// SQL username + password
connStr := "server=localhost;user id=sa;password=YourPassword123;database=unibazar;encrypt=disable"

// OR — Windows Authentication
connStr := "server=localhost;database=unibazar;trusted_connection=yes;encrypt=disable"
```

---

## Database Setup

Run the migration files in order from `000001` to `000023`.

### Option A — SSMS (Windows)

1. Open SSMS → connect to your SQL Server instance
2. Right-click **Databases → New Database** → name it `unibazar`
3. copy the content from db\schema.sql to create schema
4. then copy step by step as single query from db\trigger.sql
5. run the application
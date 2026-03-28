# Database Setup Guide

## Quick Start

### 1. Start PostgreSQL with Docker

```bash
docker-compose up -d
```

This will start a PostgreSQL container with the littleLabor database. The migrations are automatically applied on startup.

### 2. Verify the Connection

```bash
# Check container is running
docker-compose ps

# Check logs
docker-compose logs postgres
```

### 3. Build and Run the Application

```bash
go build -o bin/littleLabor ./cmd/server
./bin/littleLabor
```

The server will connect to the database and start listening on `http://localhost:8080`

### 4. Test the API

```bash
# Health check
curl http://localhost:8080/health

# Get all users
curl http://localhost:8080/api/users

# Get all chores
curl http://localhost:8080/api/chores

# Get all rewards
curl http://localhost:8080/api/rewards
```

## Database Schema

The migrations create the following tables:

- **users**: User accounts with points system
- **chores**: Task definitions with point values
- **user_chores**: Many-to-many relationship between users and chores
- **completions**: Logs of completed chores with points earned
- **rewards**: Available rewards in the system
- **user_rewards**: Earned/redeemed rewards per user

## Stopping the Database

```bash
docker-compose down
```

To also remove the database volume:

```bash
docker-compose down -v
```

## Environment Variables

Copy `.env.example` to `.env` and modify as needed. The application reads from environment variables with these defaults:

- `DB_HOST` - localhost
- `DB_PORT` - 5432
- `DB_USER` - littlelabor
- `DB_PASSWORD` - littlelabor_dev_password
- `DB_NAME` - littlelabor
- `DB_SSLMODE` - disable

## Directory Structure

```
migrations/
├── 001_initial_schema.sql    # Initial schema with all tables
internal/db/
├── connection.go             # Database connection setup
├── models.go                 # Data models (User, Chore, etc.)
├── user_repository.go        # User CRUD operations
├── chore_repository.go       # Chore and completion operations
└── reward_repository.go      # Reward operations
```

## Next Steps

- [ ] Implement authentication/authorization
- [ ] Create POST endpoints for creating users, chores, and rewards
- [ ] Add middleware for request validation
- [ ] Implement gamification logic (streaks, achievements)
- [ ] Create front-end interface

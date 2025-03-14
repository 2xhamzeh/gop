# Running the API

## Prerequisites

- Docker

Note: Docker files are using go version 1.24.0, if your local go version is newer, you might run into issues, to fix them change the version used by the docker files or the version specified in your go.mod to 1.24.0

## Setup and Running

1. Set up environment variables:

```bash
cp .env.example .env
```

2. Start services:

```bash
docker compose up -d
```

This command will:

- start a postgreSQL database container.
- run a migration script in a separate container against the database.
- run the API container and connect to the database.

Your API will be available at `http://localhost:<SERVER_PORT>`, 8080 is default in .env.example

## Environment Variables

| Variable     | Purpose                      |
| ------------ | ---------------------------- |
| PGHOST       | PostgreSQL Host              |
| PGUSER       | PostgreSQL user              |
| PGPASSWORD   | PostgreSQL password          |
| PGDATABASE   | PostgreSQL name              |
| PGSSLMODE    | PostgreSQL SSL mode          |
| JWT_SECRET   | JWT signing secret           |
| JWT_DURATION | JWT token duration           |
| SERVER_HOST  | The host name of your server |
| SERVER_PORT  | API server port              |

## Maintenance Commands

### Logs

```bash
docker compose logs -f # -f will show logs in real time
```

### Service Management

```bash
# Stop services
docker compose down

# Rebuild
docker compose build
```

### Database Operations

```bash
# Check database connectivity
docker compose exec db pg_isready

# Reset database
docker compose down -v
docker compose up --build -d
```

## Troubleshooting

If you see warnings about orphaned containers:

```bash
docker compose up --remove-orphans
```

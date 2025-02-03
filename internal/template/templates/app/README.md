# Getting Started

## Prerequisites

- Docker and Docker Compose

## Setup and Running

1. Set up environment variables:

```bash
cp .env.example .env
```

2. Start services:

```bash
docker compose up -d
```

3. Run database migrations:

```bash
docker compose --profile tools run --rm migrate up
```

Your API will be available at `http://localhost:<SERVER_PORT>`

## Environment Variables

| Variable     | Purpose                 | Default     |
| ------------ | ----------------------- | ----------- |
| APP_ENV      | Application environment | development |
| DB_USER      | Database username       | -           |
| DB_PASSWORD  | Database password       | -           |
| DB_NAME      | Database name           | -           |
| DB_SSLMODE   | PostgreSQL SSL mode     | disable     |
| JWT_SECRET   | JWT signing secret      | -           |
| JWT_DURATION | JWT token duration      | 1h          |
| SERVER_PORT  | API server port         | 8080        |

## Maintenance Commands

### Logs

```bash
# All services
docker compose logs -f

# Single service
docker compose logs -f [app|db]
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
docker compose exec db pg_isready -U $DB_USER -d $DB_NAME

# Reset database
docker compose down
docker volume rm ${PWD##*/}_postgres_data
docker compose up -d
docker compose --profile tools run --rm migrate up
```

## Troubleshooting

If you see warnings about orphaned containers:

```bash
docker compose up --remove-orphans
```

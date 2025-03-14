# Database Migrations

This directory contains the migration tool for managing database schema changes.

## Prerequisites

- PostgreSQL
- `golang-migrate` CLI tool (for creating new migrations)

## Creating New Migrations

To create a new migration:

```bash
migrate create -ext sql -dir migrations -seq name_of_your_migration
```

This will create two files:

- `migrations/XXXXXX_name_of_your_migration.up.sql`
- `migrations/XXXXXX_name_of_your_migration.down.sql`

The `up.sql` file should contain the changes you want to make, and the `down.sql` file should contain the commands to reverse these changes.

## Running Migrations

The `cmd/migrate/main.go` file is meant for running the database migrations when deploying the application (can be dockerized).
It can be used in development. It can also be extended to support steps if needed. However using the migrate CLI during development might be better depending on your needs.

All commands must be run from the project root directory.

### Using migrate CLI

```bash
migrate -source file://migrations -database postgresql://user:password@host/db?sslmode=disable [command]
```

### Using the go file

#### Configuration

The migration tool reads the configuration from the environmental variables.
By default it reads the .env file if you have one, otherwise make sure to export the variables needed.

Required environment variables:

```
PGHOST
PGUSER
PGPASSWORD
PGDATABASE
PGSSLMODE
```

#### Apply Migrations

```bash
# Using configuration from env variables
go run cmd/migrate/main.go up
```

#### Rollback Migrations

```bash
go run cmd/migrate/main.go down
```

#### Check Migration Status

```bash
go run cmd/migrate/main.go version
```

## Migration Files

- All migration files are stored in the `migrations/` directory at the project root
- Files are executed in order based on their version number (prefix)
- Each migration should be atomic and self-contained
- Always test migrations, especially `down` migrations, before applying to production

## Best Practices

1. **Keep Migrations Small**

   - Each migration should make a small, focused change
   - This makes it easier to test and rollback if needed

2. **Test Both Directions**

   - Always test both `up` and `down` migrations
   - Ensure `down` migrations completely reverse the `up` migrations

3. **Never Edit Existing Migrations**

   - Once a migration has been applied to any environment, treat it as immutable
   - If you need to make changes, create a new migration

4. **Use Transactions**
   - Wrap complex migrations in transactions
   - This ensures database consistency

## Example Migration

```sql
-- 000001_create_users.up.sql
BEGIN;
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) NOT NULL
);
COMMIT;

-- 000001_create_users.down.sql
BEGIN;
DROP TABLE users;
COMMIT;
```

## Troubleshooting

If you encounter "dirty" database state:

1. Check the current version: `go run cmd/migrate/main.go version`
2. Fix any issues in your migration files
3. You may need to manually fix the database state

For more information about the migration library, refer to the [golang-migrate documentation](https://github.com/golang-migrate/migrate).

# gop

A simple Go project generator that creates projects from templates.

## Installation

```bash
go install github.com/2xhamzeh/gop@latest
```

## Usage

1. Create and enter your project directory:

   ```bash
   mkdir my-project && cd my-project
   ```

2. Generate project from template:

   ```bash
   gop [template-name] [module-name]
   ```

## Available Templates

- `empty` - Empty project
- `rest` - REST API with authentication, postgreSQL, Docker files and more

### REST

The `rest` template is currently the largest template provided by this tool.
It has:

- docker images and docker compose setup
- PostgreSQL
- database migrations
- JWT token authentication
- user setup
- simple validator
- configuration setup using environmental variables
- chi router
- middlewares
- central logging for requests and errors with request ID
- domain errors setup and conversion of domain errors to http errors on response
- standard response format and json helpers
- context flow
- dependency injection
- Handler > Service > Repository
- and more...

Unlike other tools, this is not meant to be a framework. This is just a project template you can use, modify or reference while developing a REST API. (still being developed)

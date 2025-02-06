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

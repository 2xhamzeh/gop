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

2. Initialize a Go module:

   ```bash
   go mod init <module-name>
   ```

3. Generate project from template:

   ```bash
   gop <template-name>
   ```

4. Install dependencies:
   ```bash
   go mod tidy
   ```

## Available Templates

- `empty` - Empty project
- `app` - REST API with authentication, postgreSQL, Docker files and more

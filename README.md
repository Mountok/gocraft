# GoCraft

GoCraft is a CLI generator for production-ready Go project structure.

It creates a clean, ordinary Go service foundation without forcing a framework or runtime dependency on the generated project.

## Install

System install on macOS/Linux:

```sh
curl -fsSL https://raw.githubusercontent.com/Mountok/gocraft/main/install.sh | sh
```

This installs `gocraft` to `/usr/local/bin`, so it works from any directory.

If `/usr/local/bin` requires admin rights, the installer will use `sudo`.

Uninstall:

```sh
curl -fsSL https://raw.githubusercontent.com/Mountok/gocraft/main/uninstall.sh | sh
```

Manual Go install is also supported, but it requires Go's bin directory to be in `PATH`:

```sh
go install github.com/Mountok/gocraft/cmd/gocraft@latest
```

## Usage

Create a new service:

```sh
gocraft new user-service
```

Create a Gin service:

```sh
gocraft new user-service gin
```

Create a chi service:

```sh
gocraft new user-service chi
```

Create a Fiber service:

```sh
gocraft new user-service fiber
```

Run interactive project creation:

```sh
gocraft new
```

In a terminal, this opens a colored wizard with arrow-key selection for framework and architecture.

Check installed version:

```sh
gocraft version
gocraft check-update
```

Generated structure:

```text
user-service/
|-- cmd/user-service/
|-- internal/
|   |-- config/
|   |-- handler/
|   |-- models/
|   |-- repository/
|   `-- service/
|-- migrations/
|-- configs/
|-- pkg/
|-- .env
|-- Makefile
|-- go.mod
`-- README.md
```

The generated service uses `net/http` by default, structured logging from `log/slog`, environment-based config, and a `GET /health` endpoint.

Add a layered resource inside a generated project:

```sh
gocraft make resource user
```

This creates model, repository, service, and handler skeleton files.

## Supported Today

- Current version: `v1.4.0`
- `gocraft new <name>`
- `gocraft new <name> gin`
- `gocraft new <name> chi`
- `gocraft new <name> fiber`
- interactive `gocraft new` with arrow-key TUI
- `--router nethttp`
- `--router gin`
- `--router chi`
- `--router fiber`
- `--arch layered`
- `gocraft make resource <name>`
- `gocraft version`
- automatic update warning when a newer GitHub tag is available

Flags for Echo, Clean Architecture, database support, and ORM support are intentionally rejected until those releases are implemented.

## Roadmap

Version 1.0:

- `net/http`
- Layered Architecture
- project structure generation

Version 2.0:

- Fiber
- Echo
- Clean Architecture

Version 3.0:

- PostgreSQL
- MySQL
- GORM
- SQLX
- migrations

Version 4.0:

- CRUD generation
- OpenAPI
- Swagger
- JWT
- middleware

Version 5.0:

- custom templates
- plugins
- template marketplace
- microservice generation

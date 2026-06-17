# Changelog

## v2.1.0

### Added

- Architecture-specific generated README content.
- Clean Architecture resource generation for `gocraft make resource <name>`.
- Generation checks for all router and architecture combinations.

### Changed

- Improved architecture descriptions in the interactive TUI.

## v2.0.0

### Added

- `--arch clean` for Clean Architecture project generation.
- Interactive architecture selection now includes Clean Architecture.

## v1.5.0

### Added

- `gocraft new <name> echo` for Echo project generation.
- Interactive framework selection now includes Echo.

## v1.4.0

### Added

- `gocraft new <name> fiber` for Fiber project generation.
- Interactive framework selection now includes Fiber.

## v1.3.0

### Added

- Colored interactive TUI for `gocraft new`.
- Arrow-key selection for framework and architecture.

### Changed

- Non-terminal interactive input remains available as a fallback for pipes and tests.

## v1.2.0

### Added

- `gocraft new <name> chi` for chi project generation.
- Interactive framework selection now includes chi.

## v1.1.0

### Added

- `gocraft new <name> gin` for Gin project generation.
- `gocraft new <name> default` as an explicit default `net/http` shortcut.
- Interactive `gocraft new` prompt for project name and framework selection.
- `gocraft version` and `gocraft check-update`.
- Automatic update warning when a newer GitHub tag is available.
- Uninstall script via `uninstall.sh`.

## v1.0.0

First GoCraft version.

### Added

- `gocraft new <name>` for generating a layered Go project.
- `net/http` server scaffold with `GET /health`.
- Environment-based config.
- `log/slog` structured logging.
- `Makefile`, `.env`, `go.mod`, `README.md`, `configs`, `migrations`, and `pkg` generation.
- `gocraft make resource <name>` for model, repository, service, and handler skeletons.
- System installer via `install.sh`.

### Next

- Gin router support.

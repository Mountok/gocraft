# Roadmap

This file tracks what is already done and what should be built next.

## Released

### v1.0.0

- [x] `net/http` project generation
- [x] Layered Architecture
- [x] base project structure
- [x] health endpoint
- [x] config, logging, Makefile, README, `.env`

### v1.1.0

- [x] Gin project generation
- [x] interactive `gocraft new`
- [x] `gocraft version`
- [x] `gocraft check-update`
- [x] auto update warning
- [x] install/uninstall scripts

### v1.2.0

- [x] chi project generation
- [x] chi in interactive framework selection

### v1.3.0

- [x] colored TUI wizard
- [x] arrow-key framework selection
- [x] arrow-key architecture selection
- [x] non-terminal fallback input

### v1.4.0

- [x] Fiber project generation
- [x] Fiber in interactive framework selection

### v1.5.0

- [x] Echo project generation
- [x] Echo in interactive framework selection

### v2.0.0

- [x] Clean Architecture generation
- [x] clean architecture in interactive selection
- [x] clean structure with domain, usecase, interface, infrastructure

### v2.1.0 - Architecture Polish

- [x] improve generated README per architecture
- [x] add architecture-specific resource generation
- [x] add compile checks for all router and architecture combinations
- [x] improve TUI descriptions for architecture selection

### v3.0.0 - Database Support

- [x] PostgreSQL option
- [x] database config in `.env`
- [x] database connection package
- [x] Docker Compose generation
- [x] initial migrations directory and examples

## Next Versions

### v3.1.0 - Database Polish

- [ ] wire DB connection into generated main behind an optional health check
- [ ] improve migration command documentation
- [ ] add database-specific generated README sections per architecture

### v3.2.0 - MySQL Support

- [ ] MySQL option
- [ ] MySQL Docker Compose generation
- [ ] MySQL database connection package

### v3.3.0 - ORM Support

- [ ] GORM option
- [ ] SQLX option
- [ ] repository templates for selected DB/ORM
- [ ] migration commands in Makefile

### v4.0.0 - API Tooling

- [ ] CRUD generation improvements
- [ ] OpenAPI generation
- [ ] Swagger UI setup
- [ ] JWT auth skeleton
- [ ] middleware templates

### v5.0.0 - Templates and Plugins

- [ ] custom template directory
- [ ] project template config file
- [ ] plugin system
- [ ] template marketplace concept
- [ ] microservice generation mode

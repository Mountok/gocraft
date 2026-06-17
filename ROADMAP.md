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

## Next Versions

### v3.0.0 - Database Support

- [ ] PostgreSQL option
- [ ] MySQL option
- [ ] database config in `.env`
- [ ] database connection package
- [ ] Docker Compose generation
- [ ] initial migrations directory and examples

### v3.1.0 - ORM Support

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

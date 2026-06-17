# GoCraft

![GitHub stars](https://img.shields.io/github/stars/Mountok/gocraft?style=for-the-badge)
![Latest release](https://img.shields.io/github/v/tag/Mountok/gocraft?label=release&style=for-the-badge)
![License](https://img.shields.io/github/license/Mountok/gocraft?style=for-the-badge)
![Go version](https://img.shields.io/github/go-mod/go-version/Mountok/gocraft?style=for-the-badge)

GoCraft is a CLI generator for production-ready Go project structure.

It creates clean, ordinary Go services without forcing a framework runtime or locking the project to the generator.

## Status

| Item | Value |
| --- | --- |
| Current version | `v2.0.0` |
| License | MIT |
| Default server | `net/http` |
| Architectures | Layered, Clean Architecture |
| Frameworks | `net/http`, Gin, chi, Fiber, Echo |
| Installer | macOS/Linux shell installer |
| Update checks | Built in |

## Install

System install on macOS/Linux:

```sh
curl -fsSL https://raw.githubusercontent.com/Mountok/gocraft/main/install.sh | sh
```

This installs `gocraft` to `/usr/local/bin`, so it works from any directory. If `/usr/local/bin` requires admin rights, the installer uses `sudo`.

Uninstall:

```sh
curl -fsSL https://raw.githubusercontent.com/Mountok/gocraft/main/uninstall.sh | sh
```

Manual Go install is also supported, but it requires Go's bin directory to be in `PATH`:

```sh
go install github.com/Mountok/gocraft/cmd/gocraft@latest
```

## Quick Start

Create a default `net/http` service:

```sh
gocraft new user-service
```

Create a service with a framework:

```sh
gocraft new user-service gin
gocraft new user-service chi
gocraft new user-service fiber
gocraft new user-service echo
```

Create a Clean Architecture service:

```sh
gocraft new --arch clean user-service
```

Open the interactive wizard:

```sh
gocraft new
```

In a terminal, the wizard uses colors and arrow-key selection for framework and architecture.

## Generated Project

Layered architecture output:

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

Clean Architecture output:

```text
user-service/
|-- cmd/user-service/
|-- internal/
|   |-- config/
|   |-- domain/
|   |-- usecase/
|   |-- interface/http/
|   `-- infrastructure/repository/
|-- migrations/
|-- configs/
|-- pkg/
|-- .env
|-- Makefile
|-- go.mod
`-- README.md
```

Generated services include environment config, `log/slog` structured logging, `Makefile`, README, and `GET /health`.

## Commands

| Command | Description |
| --- | --- |
| `gocraft new <name>` | Create a default `net/http` layered project |
| `gocraft new <name> gin` | Create a Gin project |
| `gocraft new <name> chi` | Create a chi project |
| `gocraft new <name> fiber` | Create a Fiber project |
| `gocraft new <name> echo` | Create an Echo project |
| `gocraft new --arch clean <name>` | Create a Clean Architecture project |
| `gocraft make resource <name>` | Generate model, repository, service, and handler skeletons |
| `gocraft version` | Show installed version |
| `gocraft check-update` | Check the latest GitHub version |

## Feature Matrix

| Feature | Status |
| --- | --- |
| Project generator | Done |
| Layered Architecture | Done |
| Clean Architecture | Done |
| `net/http` | Done |
| Gin | Done |
| chi | Done |
| Fiber | Done |
| Echo | Done |
| Colored TUI wizard | Done |
| Resource generation | Basic |
| PostgreSQL/MySQL | Planned |
| GORM/SQLX | Planned |
| Migrations | Planned |
| OpenAPI/Swagger | Planned |
| JWT/middleware | Planned |
| Custom templates/plugins | Planned |

## Updates

GoCraft checks for newer GitHub tags before generation commands and prints an update command when a newer release exists.

Manual check:

```sh
gocraft check-update
```

Disable automatic checks for one command:

```sh
GOCRAFT_SKIP_UPDATE_CHECK=1 gocraft new user-service
```

## Project Links

| Document | Description |
| --- | --- |
| [ROADMAP.md](ROADMAP.md) | Version checklist and planned work |
| [CHANGELOG.md](CHANGELOG.md) | Release history |
| [LICENSE](LICENSE) | MIT license |

## License

GoCraft is released under the MIT License.

package generator

import (
	"go/parser"
	"go/token"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestNewProjectCreatesLayeredHTTPProject(t *testing.T) {
	tmp := t.TempDir()
	oldWD, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(oldWD); err != nil {
			t.Fatal(err)
		}
	})
	if err := os.Chdir(tmp); err != nil {
		t.Fatal(err)
	}

	if err := NewProject(ProjectOptions{Name: "user-service", Router: "nethttp", Arch: "layered"}); err != nil {
		t.Fatalf("NewProject() error = %v", err)
	}

	for _, path := range []string{
		"user-service/cmd/user-service/main.go",
		"user-service/internal/config/config.go",
		"user-service/internal/handler/health_handler.go",
		"user-service/internal/service/health_service.go",
		"user-service/internal/repository/health_repository.go",
		"user-service/internal/models/health.go",
		"user-service/migrations/.gitkeep",
		"user-service/configs/config.example.yaml",
		"user-service/.env",
		"user-service/Makefile",
		"user-service/go.mod",
		"user-service/README.md",
	} {
		if _, err := os.Stat(filepath.Join(tmp, path)); err != nil {
			t.Fatalf("expected generated file %s: %v", path, err)
		}
	}
}

func TestNewProjectCreatesGinProject(t *testing.T) {
	tmp := t.TempDir()
	oldWD, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(oldWD); err != nil {
			t.Fatal(err)
		}
	})
	if err := os.Chdir(tmp); err != nil {
		t.Fatal(err)
	}

	if err := NewProject(ProjectOptions{Name: "api", Router: "gin", Arch: "layered"}); err != nil {
		t.Fatalf("NewProject() error = %v", err)
	}

	goMod, err := os.ReadFile(filepath.Join(tmp, "api", "go.mod"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(goMod), "github.com/gin-gonic/gin") {
		t.Fatal("expected Gin dependency in go.mod")
	}
}

func TestNewProjectCreatesChiProject(t *testing.T) {
	tmp := t.TempDir()
	oldWD, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(oldWD); err != nil {
			t.Fatal(err)
		}
	})
	if err := os.Chdir(tmp); err != nil {
		t.Fatal(err)
	}

	if err := NewProject(ProjectOptions{Name: "api", Router: "chi", Arch: "layered"}); err != nil {
		t.Fatalf("NewProject() error = %v", err)
	}

	goMod, err := os.ReadFile(filepath.Join(tmp, "api", "go.mod"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(goMod), "github.com/go-chi/chi/v5") {
		t.Fatal("expected chi dependency in go.mod")
	}
}

func TestNewProjectCreatesFiberProject(t *testing.T) {
	tmp := t.TempDir()
	oldWD, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(oldWD); err != nil {
			t.Fatal(err)
		}
	})
	if err := os.Chdir(tmp); err != nil {
		t.Fatal(err)
	}

	if err := NewProject(ProjectOptions{Name: "api", Router: "fiber", Arch: "layered"}); err != nil {
		t.Fatalf("NewProject() error = %v", err)
	}

	goMod, err := os.ReadFile(filepath.Join(tmp, "api", "go.mod"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(goMod), "github.com/gofiber/fiber/v2") {
		t.Fatal("expected Fiber dependency in go.mod")
	}
}

func TestNewProjectCreatesEchoProject(t *testing.T) {
	tmp := t.TempDir()
	oldWD, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(oldWD); err != nil {
			t.Fatal(err)
		}
	})
	if err := os.Chdir(tmp); err != nil {
		t.Fatal(err)
	}

	if err := NewProject(ProjectOptions{Name: "api", Router: "echo", Arch: "layered"}); err != nil {
		t.Fatalf("NewProject() error = %v", err)
	}

	goMod, err := os.ReadFile(filepath.Join(tmp, "api", "go.mod"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(goMod), "github.com/labstack/echo/v4") {
		t.Fatal("expected Echo dependency in go.mod")
	}
}

func TestNewProjectCreatesCleanProject(t *testing.T) {
	tmp := t.TempDir()
	oldWD, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(oldWD); err != nil {
			t.Fatal(err)
		}
	})
	if err := os.Chdir(tmp); err != nil {
		t.Fatal(err)
	}

	if err := NewProject(ProjectOptions{Name: "api", Router: "nethttp", Arch: "clean"}); err != nil {
		t.Fatalf("NewProject() error = %v", err)
	}

	for _, path := range []string{
		"api/internal/domain/health.go",
		"api/internal/usecase/health_usecase.go",
		"api/internal/interface/http/health_handler.go",
		"api/internal/infrastructure/repository/health_repository.go",
	} {
		if _, err := os.Stat(filepath.Join(tmp, path)); err != nil {
			t.Fatalf("expected generated file %s: %v", path, err)
		}
	}
}

func TestNewProjectCreatesPostgresFiles(t *testing.T) {
	tmp := t.TempDir()
	oldWD, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(oldWD); err != nil {
			t.Fatal(err)
		}
	})
	if err := os.Chdir(tmp); err != nil {
		t.Fatal(err)
	}

	if err := NewProject(ProjectOptions{Name: "api", Router: "nethttp", Arch: "layered", DB: "postgres"}); err != nil {
		t.Fatalf("NewProject() error = %v", err)
	}

	for _, path := range []string{
		"api/internal/database/postgres.go",
		"api/docker-compose.yml",
		"api/migrations/000001_init.up.sql",
		"api/migrations/000001_init.down.sql",
	} {
		if _, err := os.Stat(filepath.Join(tmp, path)); err != nil {
			t.Fatalf("expected generated file %s: %v", path, err)
		}
	}

	goMod, err := os.ReadFile(filepath.Join(tmp, "api", "go.mod"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(goMod), "github.com/jackc/pgx/v5") {
		t.Fatal("expected pgx dependency in go.mod")
	}
}

func TestNewProjectCreatesCleanPostgresFiles(t *testing.T) {
	tmp := t.TempDir()
	oldWD, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(oldWD); err != nil {
			t.Fatal(err)
		}
	})
	if err := os.Chdir(tmp); err != nil {
		t.Fatal(err)
	}

	if err := NewProject(ProjectOptions{Name: "api", Router: "nethttp", Arch: "clean", DB: "pg"}); err != nil {
		t.Fatalf("NewProject() error = %v", err)
	}
	if _, err := os.Stat(filepath.Join(tmp, "api/internal/infrastructure/database/postgres.go")); err != nil {
		t.Fatalf("expected clean postgres database file: %v", err)
	}
}

func TestNewProjectRejectsUnsupportedRouter(t *testing.T) {
	err := NewProject(ProjectOptions{Name: "api", Router: "fasthttp", Arch: "layered"})
	if err == nil {
		t.Fatal("expected unsupported router error")
	}
}

func TestNewProjectRejectsUnsupportedArchitecture(t *testing.T) {
	err := NewProject(ProjectOptions{Name: "api", Router: "nethttp", Arch: "hexagonal"})
	if err == nil {
		t.Fatal("expected unsupported architecture error")
	}
}

func TestNewProjectRejectsUnsupportedDatabase(t *testing.T) {
	err := NewProject(ProjectOptions{Name: "api", Router: "nethttp", Arch: "layered", DB: "mysql"})
	if err == nil {
		t.Fatal("expected unsupported database error")
	}
}

func TestNewResourceDoesNotOverwriteExistingFiles(t *testing.T) {
	tmp := t.TempDir()
	if err := os.WriteFile(filepath.Join(tmp, "go.mod"), []byte("module example.com/app\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(tmp, "internal/models"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(tmp, "internal/models/user.go"), []byte("package models\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	err := NewResource(tmp, "user")
	if err == nil {
		t.Fatal("expected overwrite error")
	}
}

func TestNewResourceCreatesCleanArchitectureFiles(t *testing.T) {
	tmp := t.TempDir()
	oldWD, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(oldWD); err != nil {
			t.Fatal(err)
		}
	})
	if err := os.Chdir(tmp); err != nil {
		t.Fatal(err)
	}
	if err := NewProject(ProjectOptions{Name: "api", Router: "nethttp", Arch: "clean"}); err != nil {
		t.Fatalf("NewProject() error = %v", err)
	}
	if err := NewResource(filepath.Join(tmp, "api"), "user"); err != nil {
		t.Fatalf("NewResource() error = %v", err)
	}

	for _, path := range []string{
		"api/internal/domain/user.go",
		"api/internal/usecase/user_usecase.go",
		"api/internal/interface/http/user_handler.go",
		"api/internal/infrastructure/repository/user_repository.go",
	} {
		if _, err := os.Stat(filepath.Join(tmp, path)); err != nil {
			t.Fatalf("expected generated file %s: %v", path, err)
		}
	}
}

func TestAllRouterArchitectureCombinationsGenerateParseableGo(t *testing.T) {
	routers := []string{"nethttp", "gin", "chi", "fiber", "echo"}
	architectures := []string{"layered", "clean"}

	for _, router := range routers {
		for _, arch := range architectures {
			name := router + "-" + arch
			t.Run(name, func(t *testing.T) {
				tmp := t.TempDir()
				oldWD, err := os.Getwd()
				if err != nil {
					t.Fatal(err)
				}
				t.Cleanup(func() {
					if err := os.Chdir(oldWD); err != nil {
						t.Fatal(err)
					}
				})
				if err := os.Chdir(tmp); err != nil {
					t.Fatal(err)
				}
				if err := NewProject(ProjectOptions{Name: name, Router: router, Arch: arch}); err != nil {
					t.Fatalf("NewProject() error = %v", err)
				}
				assertParseableGo(t, filepath.Join(tmp, name))
				if router == "nethttp" {
					cmd := exec.Command("go", "test", "./...")
					cmd.Dir = filepath.Join(tmp, name)
					if output, err := cmd.CombinedOutput(); err != nil {
						t.Fatalf("generated net/http project does not compile: %v\n%s", err, output)
					}
				}
			})
		}
	}
}

func TestAllPostgresRouterArchitectureCombinationsGenerateParseableGo(t *testing.T) {
	routers := []string{"nethttp", "gin", "chi", "fiber", "echo"}
	architectures := []string{"layered", "clean"}

	for _, router := range routers {
		for _, arch := range architectures {
			name := router + "-" + arch + "-postgres"
			t.Run(name, func(t *testing.T) {
				tmp := t.TempDir()
				oldWD, err := os.Getwd()
				if err != nil {
					t.Fatal(err)
				}
				t.Cleanup(func() {
					if err := os.Chdir(oldWD); err != nil {
						t.Fatal(err)
					}
				})
				if err := os.Chdir(tmp); err != nil {
					t.Fatal(err)
				}
				if err := NewProject(ProjectOptions{Name: name, Router: router, Arch: arch, DB: "postgres"}); err != nil {
					t.Fatalf("NewProject() error = %v", err)
				}
				assertParseableGo(t, filepath.Join(tmp, name))
			})
		}
	}
}

func assertParseableGo(t *testing.T, root string) {
	t.Helper()
	fset := token.NewFileSet()
	if err := filepath.WalkDir(root, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() || filepath.Ext(path) != ".go" {
			return nil
		}
		_, err = parser.ParseFile(fset, path, nil, parser.AllErrors)
		return err
	}); err != nil {
		t.Fatalf("generated Go files are not parseable: %v", err)
	}
}

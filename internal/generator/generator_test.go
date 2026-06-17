package generator

import (
	"os"
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

func TestNewProjectRejectsUnsupportedRouter(t *testing.T) {
	err := NewProject(ProjectOptions{Name: "api", Router: "fasthttp", Arch: "layered"})
	if err == nil {
		t.Fatal("expected unsupported router error")
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

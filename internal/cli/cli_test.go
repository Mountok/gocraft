package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestRunHelp(t *testing.T) {
	t.Setenv("GOCRAFT_SKIP_UPDATE_CHECK", "1")
	var stdout, stderr bytes.Buffer
	if err := Run([]string{"help"}, nil, &stdout, &stderr); err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if stdout.Len() == 0 {
		t.Fatal("expected help output")
	}
}

func TestRunNew(t *testing.T) {
	t.Setenv("GOCRAFT_SKIP_UPDATE_CHECK", "1")
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

	var stdout, stderr bytes.Buffer
	if err := Run([]string{"new", "api"}, nil, &stdout, &stderr); err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if _, err := os.Stat(filepath.Join(tmp, "api", "cmd", "api", "main.go")); err != nil {
		t.Fatalf("expected generated main.go: %v", err)
	}
}

func TestRunNewWithPositionalFramework(t *testing.T) {
	t.Setenv("GOCRAFT_SKIP_UPDATE_CHECK", "1")
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

	var stdout, stderr bytes.Buffer
	if err := Run([]string{"new", "api", "chi"}, nil, &stdout, &stderr); err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if _, err := os.Stat(filepath.Join(tmp, "api", "cmd", "api", "main.go")); err != nil {
		t.Fatalf("expected generated main.go: %v", err)
	}
}

func TestRunNewInteractive(t *testing.T) {
	t.Setenv("GOCRAFT_SKIP_UPDATE_CHECK", "1")
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

	var stdout, stderr bytes.Buffer
	stdin := bytes.NewBufferString("api\n3\n")
	if err := Run([]string{"new"}, stdin, &stdout, &stderr); err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if _, err := os.Stat(filepath.Join(tmp, "api", "cmd", "api", "main.go")); err != nil {
		t.Fatalf("expected generated main.go: %v", err)
	}
}

func TestRunVersion(t *testing.T) {
	var stdout, stderr bytes.Buffer
	if err := Run([]string{"version"}, nil, &stdout, &stderr); err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if stdout.String() == "" {
		t.Fatal("expected version output")
	}
}

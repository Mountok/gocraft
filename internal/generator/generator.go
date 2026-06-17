package generator

import (
	"bytes"
	"errors"
	"fmt"
	"go/format"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"
)

const defaultModulePrefix = "github.com/example"

var validName = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_-]*$`)

type ProjectOptions struct {
	Name   string
	Router string
	Arch   string
	DB     string
	ORM    string
}

type fileSpec struct {
	Path     string
	Template string
	Go       bool
}

type templateData struct {
	Name        string
	PackageName string
	ModulePath  string
	Resource    string
	TypeName    string
}

func NewProject(options ProjectOptions) error {
	if err := validateProjectOptions(options); err != nil {
		return err
	}
	if _, err := os.Stat(options.Name); err == nil {
		return fmt.Errorf("target directory %q already exists", options.Name)
	} else if !errors.Is(err, os.ErrNotExist) {
		return err
	}

	data := templateData{
		Name:        options.Name,
		PackageName: packageName(options.Name),
		ModulePath:  defaultModulePrefix + "/" + options.Name,
	}

	for _, dir := range []string{
		"cmd/" + options.Name,
		"internal/config",
		"internal/handler",
		"internal/models",
		"internal/repository",
		"internal/service",
		"configs",
		"migrations",
		"pkg",
	} {
		if err := os.MkdirAll(filepath.Join(options.Name, dir), 0o755); err != nil {
			return err
		}
	}

	files := []fileSpec{
		{Path: "go.mod", Template: projectGoModTemplate},
		{Path: ".env", Template: envTemplate},
		{Path: "Makefile", Template: makefileTemplate},
		{Path: "README.md", Template: projectReadmeTemplate},
		{Path: "cmd/{{.Name}}/main.go", Template: mainTemplate, Go: true},
		{Path: "internal/config/config.go", Template: configTemplate, Go: true},
		{Path: "internal/handler/health_handler.go", Template: healthHandlerTemplate, Go: true},
		{Path: "internal/service/health_service.go", Template: healthServiceTemplate, Go: true},
		{Path: "internal/repository/health_repository.go", Template: healthRepositoryTemplate, Go: true},
		{Path: "internal/models/health.go", Template: healthModelTemplate, Go: true},
		{Path: "configs/config.example.yaml", Template: configExampleTemplate},
		{Path: "migrations/.gitkeep", Template: ""},
		{Path: "pkg/.gitkeep", Template: ""},
	}

	for _, file := range files {
		if err := writeTemplate(options.Name, file, data); err != nil {
			return err
		}
	}

	return nil
}

func NewResource(root, name string) error {
	if !validName.MatchString(name) {
		return fmt.Errorf("invalid resource name %q", name)
	}
	modulePath, err := readModulePath(filepath.Join(root, "go.mod"))
	if err != nil {
		return err
	}

	data := templateData{
		ModulePath: modulePath,
		Resource:   packageName(name),
		TypeName:   typeName(name),
	}

	files := []fileSpec{
		{Path: "internal/models/{{.Resource}}.go", Template: resourceModelTemplate, Go: true},
		{Path: "internal/repository/{{.Resource}}_repository.go", Template: resourceRepositoryTemplate, Go: true},
		{Path: "internal/service/{{.Resource}}_service.go", Template: resourceServiceTemplate, Go: true},
		{Path: "internal/handler/{{.Resource}}_handler.go", Template: resourceHandlerTemplate, Go: true},
	}

	for _, file := range files {
		if err := writeTemplate(root, file, data); err != nil {
			return err
		}
	}

	return nil
}

func validateProjectOptions(options ProjectOptions) error {
	if !validName.MatchString(options.Name) {
		return fmt.Errorf("invalid project name %q", options.Name)
	}
	if options.Router == "" {
		options.Router = "nethttp"
	}
	if options.Router != "nethttp" {
		return fmt.Errorf("router %q is not supported yet; use nethttp", options.Router)
	}
	if options.Arch == "" {
		options.Arch = "layered"
	}
	if options.Arch != "layered" {
		return fmt.Errorf("architecture %q is not supported yet; use layered", options.Arch)
	}
	if options.DB != "" {
		return fmt.Errorf("database support %q is planned for a future release", options.DB)
	}
	if options.ORM != "" {
		return fmt.Errorf("ORM support %q is planned for a future release", options.ORM)
	}
	return nil
}

func writeTemplate(root string, file fileSpec, data templateData) error {
	relPath, err := render(file.Path, data)
	if err != nil {
		return err
	}
	target := filepath.Join(root, relPath)
	if _, err := os.Stat(target); err == nil {
		return fmt.Errorf("file %q already exists", target)
	} else if !errors.Is(err, os.ErrNotExist) {
		return err
	}

	content, err := render(file.Template, data)
	if err != nil {
		return err
	}
	output := []byte(content)
	if file.Go {
		formatted, err := format.Source(output)
		if err != nil {
			return fmt.Errorf("format %s: %w", relPath, err)
		}
		output = formatted
	}

	if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
		return err
	}
	return os.WriteFile(target, output, 0o644)
}

func render(source string, data templateData) (string, error) {
	tpl, err := template.New("file").Parse(source)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if err := tpl.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func packageName(name string) string {
	name = strings.ReplaceAll(name, "-", "_")
	return strings.ToLower(name)
}

func typeName(name string) string {
	parts := regexp.MustCompile(`[-_]+`).Split(name, -1)
	for i, part := range parts {
		if part == "" {
			continue
		}
		parts[i] = strings.ToUpper(part[:1]) + strings.ToLower(part[1:])
	}
	return strings.Join(parts, "")
}

func readModulePath(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("read go.mod: %w", err)
	}
	for _, line := range strings.Split(string(content), "\n") {
		fields := strings.Fields(line)
		if len(fields) == 2 && fields[0] == "module" {
			return fields[1], nil
		}
	}
	return "", fmt.Errorf("module path not found in %s", path)
}

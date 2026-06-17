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
	Name         string
	PackageName  string
	ModulePath   string
	Router       string
	Arch         string
	IsClean      bool
	HasDeps      bool
	UsesGin      bool
	UsesChi      bool
	UsesFiber    bool
	UsesEcho     bool
	UsesPostgres bool
	Resource     string
	TypeName     string
}

func NewProject(options ProjectOptions) error {
	router, err := normalizeRouter(options.Router)
	if err != nil {
		return err
	}
	options.Router = router
	if options.Arch == "" {
		options.Arch = "layered"
	}
	options.DB = normalizeDB(options.DB)

	if err := validateProjectOptions(options); err != nil {
		return err
	}
	if _, err := os.Stat(options.Name); err == nil {
		return fmt.Errorf("target directory %q already exists", options.Name)
	} else if !errors.Is(err, os.ErrNotExist) {
		return err
	}

	data := templateData{
		Name:         options.Name,
		PackageName:  packageName(options.Name),
		ModulePath:   defaultModulePrefix + "/" + options.Name,
		Router:       options.Router,
		Arch:         options.Arch,
		IsClean:      options.Arch == "clean",
		HasDeps:      options.Router != "nethttp" || options.DB == "postgres",
		UsesGin:      options.Router == "gin",
		UsesChi:      options.Router == "chi",
		UsesFiber:    options.Router == "fiber",
		UsesEcho:     options.Router == "echo",
		UsesPostgres: options.DB == "postgres",
	}

	dirs := []string{
		"cmd/" + options.Name,
		"internal/config",
		"internal/handler",
		"internal/models",
		"internal/repository",
		"internal/service",
		"configs",
		"migrations",
		"pkg",
	}
	files := layeredFiles(options.Router)
	if options.Arch == "clean" {
		dirs = []string{
			"cmd/" + options.Name,
			"internal/config",
			"internal/domain",
			"internal/usecase",
			"internal/interface/http",
			"internal/infrastructure/repository",
			"configs",
			"migrations",
			"pkg",
		}
		files = cleanFiles(options.Router)
	}
	if options.DB == "postgres" {
		dirs = append(dirs, postgresDirs(options.Arch)...)
		files = append(files, postgresFiles(options.Arch)...)
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(filepath.Join(options.Name, dir), 0o755); err != nil {
			return err
		}
	}

	for _, file := range files {
		if err := writeTemplate(options.Name, file, data); err != nil {
			return err
		}
	}

	return nil
}

func postgresDirs(arch string) []string {
	if arch == "clean" {
		return []string{"internal/infrastructure/database"}
	}
	return []string{"internal/database"}
}

func postgresFiles(arch string) []fileSpec {
	dbPath := "internal/database/postgres.go"
	if arch == "clean" {
		dbPath = "internal/infrastructure/database/postgres.go"
	}
	return []fileSpec{
		{Path: dbPath, Template: postgresTemplate, Go: true},
		{Path: "docker-compose.yml", Template: dockerComposeTemplate},
		{Path: "migrations/000001_init.up.sql", Template: migrationUpTemplate},
		{Path: "migrations/000001_init.down.sql", Template: migrationDownTemplate},
	}
}

func layeredFiles(router string) []fileSpec {
	mainTpl := netHTTPMainTemplate
	healthHandlerTpl := netHTTPHealthHandlerTemplate
	if router == "gin" {
		mainTpl = ginMainTemplate
		healthHandlerTpl = ginHealthHandlerTemplate
	} else if router == "chi" {
		mainTpl = chiMainTemplate
	} else if router == "fiber" {
		mainTpl = fiberMainTemplate
		healthHandlerTpl = fiberHealthHandlerTemplate
	} else if router == "echo" {
		mainTpl = echoMainTemplate
		healthHandlerTpl = echoHealthHandlerTemplate
	}

	return []fileSpec{
		{Path: "go.mod", Template: projectGoModTemplate},
		{Path: ".env", Template: envTemplate},
		{Path: "Makefile", Template: makefileTemplate},
		{Path: "README.md", Template: projectReadmeTemplate},
		{Path: "cmd/{{.Name}}/main.go", Template: mainTpl, Go: true},
		{Path: "internal/config/config.go", Template: configTemplate, Go: true},
		{Path: "internal/handler/health_handler.go", Template: healthHandlerTpl, Go: true},
		{Path: "internal/service/health_service.go", Template: healthServiceTemplate, Go: true},
		{Path: "internal/repository/health_repository.go", Template: healthRepositoryTemplate, Go: true},
		{Path: "internal/models/health.go", Template: healthModelTemplate, Go: true},
		{Path: "configs/config.example.yaml", Template: configExampleTemplate},
		{Path: "migrations/.gitkeep", Template: ""},
		{Path: "pkg/.gitkeep", Template: ""},
	}
}

func cleanFiles(router string) []fileSpec {
	mainTpl := cleanNetHTTPMainTemplate
	healthHandlerTpl := cleanNetHTTPHealthHandlerTemplate
	if router == "gin" {
		mainTpl = cleanGinMainTemplate
		healthHandlerTpl = cleanGinHealthHandlerTemplate
	} else if router == "chi" {
		mainTpl = cleanChiMainTemplate
	} else if router == "fiber" {
		mainTpl = cleanFiberMainTemplate
		healthHandlerTpl = cleanFiberHealthHandlerTemplate
	} else if router == "echo" {
		mainTpl = cleanEchoMainTemplate
		healthHandlerTpl = cleanEchoHealthHandlerTemplate
	}

	return []fileSpec{
		{Path: "go.mod", Template: projectGoModTemplate},
		{Path: ".env", Template: envTemplate},
		{Path: "Makefile", Template: makefileTemplate},
		{Path: "README.md", Template: projectReadmeTemplate},
		{Path: "cmd/{{.Name}}/main.go", Template: mainTpl, Go: true},
		{Path: "internal/config/config.go", Template: configTemplate, Go: true},
		{Path: "internal/domain/health.go", Template: cleanHealthDomainTemplate, Go: true},
		{Path: "internal/usecase/health_usecase.go", Template: cleanHealthUsecaseTemplate, Go: true},
		{Path: "internal/interface/http/health_handler.go", Template: healthHandlerTpl, Go: true},
		{Path: "internal/infrastructure/repository/health_repository.go", Template: cleanHealthRepositoryTemplate, Go: true},
		{Path: "configs/config.example.yaml", Template: configExampleTemplate},
		{Path: "migrations/.gitkeep", Template: ""},
		{Path: "pkg/.gitkeep", Template: ""},
	}
}

func NewResource(root, name string) error {
	if !validName.MatchString(name) {
		return fmt.Errorf("invalid resource name %q", name)
	}
	modulePath, err := readModulePath(filepath.Join(root, "go.mod"))
	if err != nil {
		return err
	}
	arch := "layered"
	if _, err := os.Stat(filepath.Join(root, "internal", "domain")); err == nil {
		arch = "clean"
	} else if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}

	data := templateData{
		ModulePath: modulePath,
		Arch:       arch,
		IsClean:    arch == "clean",
		Resource:   packageName(name),
		TypeName:   typeName(name),
	}

	files := []fileSpec{
		{Path: "internal/models/{{.Resource}}.go", Template: resourceModelTemplate, Go: true},
		{Path: "internal/repository/{{.Resource}}_repository.go", Template: resourceRepositoryTemplate, Go: true},
		{Path: "internal/service/{{.Resource}}_service.go", Template: resourceServiceTemplate, Go: true},
		{Path: "internal/handler/{{.Resource}}_handler.go", Template: resourceHandlerTemplate, Go: true},
	}
	if arch == "clean" {
		files = []fileSpec{
			{Path: "internal/domain/{{.Resource}}.go", Template: cleanResourceDomainTemplate, Go: true},
			{Path: "internal/usecase/{{.Resource}}_usecase.go", Template: cleanResourceUsecaseTemplate, Go: true},
			{Path: "internal/interface/http/{{.Resource}}_handler.go", Template: cleanResourceHandlerTemplate, Go: true},
			{Path: "internal/infrastructure/repository/{{.Resource}}_repository.go", Template: cleanResourceRepositoryTemplate, Go: true},
		}
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
	if options.Arch == "" {
		options.Arch = "layered"
	}
	if options.Arch != "layered" && options.Arch != "clean" {
		return fmt.Errorf("architecture %q is not supported yet; use layered or clean", options.Arch)
	}
	if options.DB != "" && options.DB != "postgres" {
		return fmt.Errorf("database support %q is not supported yet; use postgres", options.DB)
	}
	if options.ORM != "" {
		return fmt.Errorf("ORM support %q is planned for a future release", options.ORM)
	}
	return nil
}

func normalizeDB(db string) string {
	switch strings.ToLower(strings.TrimSpace(db)) {
	case "", "none", "no":
		return ""
	case "postgres", "postgresql", "pg":
		return "postgres"
	default:
		return strings.ToLower(strings.TrimSpace(db))
	}
}

func normalizeRouter(router string) (string, error) {
	switch strings.ToLower(strings.TrimSpace(router)) {
	case "", "default", "nethttp", "net/http", "http":
		return "nethttp", nil
	case "gin":
		return "gin", nil
	case "chi":
		return "chi", nil
	case "fiber":
		return "fiber", nil
	case "echo":
		return "echo", nil
	default:
		return "", fmt.Errorf("framework %q is not supported yet; use default, gin, chi, fiber, or echo", router)
	}
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

package cli

import (
	"bufio"
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/manifoldco/promptui"

	"github.com/Mountok/gocraft/internal/generator"
	"github.com/Mountok/gocraft/internal/update"
	"github.com/Mountok/gocraft/internal/version"
)

const usage = `GoCraft - Go project architecture generator

Usage:
  gocraft new
  gocraft new <name> [default|gin|chi|fiber|echo]
  gocraft new [--router default|nethttp|gin|chi|fiber|echo] [--arch layered|clean] <name>
  gocraft make resource <name>
  gocraft version
  gocraft check-update

Examples:
  gocraft new user-service
  gocraft new user-service gin
  gocraft new user-service chi
  gocraft new user-service fiber
  gocraft new user-service echo
  gocraft make resource user
`

func Run(args []string, stdin io.Reader, stdout, stderr io.Writer) error {
	if len(args) == 0 {
		_, _ = io.WriteString(stdout, usage)
		return nil
	}

	switch args[0] {
	case "help", "--help", "-h":
		_, _ = io.WriteString(stdout, usage)
		return nil
	case "new":
		update.WarnIfOutdated(stderr)
		return runNew(args[1:], stdin, stdout, stderr)
	case "make":
		update.WarnIfOutdated(stderr)
		return runMake(args[1:], stdout, stderr)
	case "version", "--version", "-v":
		fmt.Fprintf(stdout, "gocraft %s\n", version.Current())
		return nil
	case "check-update":
		return runCheckUpdate(stdout)
	default:
		return fmt.Errorf("unknown command %q\n\n%s", args[0], usage)
	}
}

func runCheckUpdate(stdout io.Writer) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := update.Check(ctx)
	if err != nil {
		return err
	}
	if result.Outdated {
		fmt.Fprintf(stdout, "update available: %s -> %s\n", result.Current, result.Latest)
		fmt.Fprintln(stdout, "Update: curl -fsSL https://raw.githubusercontent.com/Mountok/gocraft/main/install.sh | sh")
		return nil
	}
	fmt.Fprintf(stdout, "gocraft is up to date: %s\n", result.Current)
	return nil
}

func runNew(args []string, stdin io.Reader, stdout, stderr io.Writer) error {
	flags := flag.NewFlagSet("new", flag.ContinueOnError)
	flags.SetOutput(stderr)
	router := flags.String("router", "", "HTTP framework: default, nethttp, gin, chi, fiber, echo")
	arch := flags.String("arch", "layered", "architecture: layered, clean")
	db := flags.String("db", "", "database: postgres")
	orm := flags.String("orm", "", "ORM support: planned")

	if err := flags.Parse(args); err != nil {
		return err
	}
	archOverride := ""
	flags.Visit(func(f *flag.Flag) {
		if f.Name == "arch" {
			archOverride = *arch
		}
	})
	if flags.NArg() == 0 {
		return runInteractiveNew(stdin, stdout, archOverride, *db, *orm)
	}
	if flags.NArg() > 2 {
		return errors.New("usage: gocraft new <name> [default|gin|chi|fiber|echo]")
	}
	if flags.NArg() == 2 && *router != "" {
		return errors.New("use either positional framework or --router, not both")
	}
	if flags.NArg() == 2 {
		*router = flags.Arg(1)
	}

	options := generator.ProjectOptions{
		Name:   flags.Arg(0),
		Router: *router,
		Arch:   *arch,
		DB:     *db,
		ORM:    *orm,
	}
	if err := generator.NewProject(options); err != nil {
		return err
	}

	fmt.Fprintf(stdout, "created GoCraft project %s\n", options.Name)
	return nil
}

func runInteractiveNew(stdin io.Reader, stdout io.Writer, arch, db, orm string) error {
	if isTerminal(stdin) {
		return runTUIInteractiveNew(stdout, arch, db, orm)
	}
	return runFallbackInteractiveNew(stdin, stdout, arch, db, orm)
}

func runTUIInteractiveNew(stdout io.Writer, arch, db, orm string) error {
	fmt.Fprintln(stdout, "\033[36mGoCraft project wizard\033[0m")
	fmt.Fprintln(stdout, "\033[90mUse arrow keys to select options. Press Enter to confirm.\033[0m")

	prompt := promptui.Prompt{
		Label: "\033[35mProject name\033[0m",
		Validate: func(input string) error {
			if strings.TrimSpace(input) == "" {
				return errors.New("project name is required")
			}
			return nil
		},
	}
	name, err := prompt.Run()
	if err != nil {
		return err
	}

	frameworks := []choice{
		{Label: "default", Description: "net/http standard library", Value: "default"},
		{Label: "gin", Description: "Gin web framework", Value: "gin"},
		{Label: "chi", Description: "chi lightweight router", Value: "chi"},
		{Label: "fiber", Description: "Fiber web framework", Value: "fiber"},
		{Label: "echo", Description: "Echo web framework", Value: "echo"},
	}
	_, framework, err := selectChoice("Framework", frameworks)
	if err != nil {
		return err
	}

	architectures := []choice{
		{Label: "layered", Description: "simple service layout: handler -> service -> repository", Value: "layered"},
		{Label: "clean", Description: "strict boundaries: domain -> usecase -> interface -> infrastructure", Value: "clean"},
	}
	_, selectedArch, err := selectChoice("Architecture", architectures)
	if err != nil {
		return err
	}
	if arch != "" {
		selectedArch.Value = arch
	}

	options := generator.ProjectOptions{
		Name:   strings.TrimSpace(name),
		Router: framework.Value,
		Arch:   selectedArch.Value,
		DB:     db,
		ORM:    orm,
	}
	if err := generator.NewProject(options); err != nil {
		return err
	}

	fmt.Fprintf(stdout, "\033[32mcreated GoCraft project %s\033[0m\n", options.Name)
	return nil
}

func runFallbackInteractiveNew(stdin io.Reader, stdout io.Writer, arch, db, orm string) error {
	reader := bufio.NewReader(stdin)

	fmt.Fprint(stdout, "Project name: ")
	name, err := reader.ReadString('\n')
	if err != nil && !errors.Is(err, io.EOF) {
		return err
	}
	name = strings.TrimSpace(name)
	if name == "" {
		return errors.New("project name is required")
	}

	fmt.Fprintln(stdout, "Framework:")
	fmt.Fprintln(stdout, "  1) default (net/http)")
	fmt.Fprintln(stdout, "  2) gin")
	fmt.Fprintln(stdout, "  3) chi")
	fmt.Fprintln(stdout, "  4) fiber")
	fmt.Fprintln(stdout, "  5) echo")
	fmt.Fprint(stdout, "Select framework [1]: ")
	framework, err := reader.ReadString('\n')
	if err != nil && !errors.Is(err, io.EOF) {
		return err
	}
	framework = strings.TrimSpace(framework)
	if framework == "" || framework == "1" {
		framework = "default"
	} else if framework == "2" {
		framework = "gin"
	} else if framework == "3" {
		framework = "chi"
	} else if framework == "4" {
		framework = "fiber"
	} else if framework == "5" {
		framework = "echo"
	}

	fmt.Fprintln(stdout, "Architecture:")
	fmt.Fprintln(stdout, "  1) layered")
	fmt.Fprintln(stdout, "  2) clean")
	fmt.Fprint(stdout, "Select architecture [1]: ")
	selectedArch, err := reader.ReadString('\n')
	if err != nil && !errors.Is(err, io.EOF) {
		return err
	}
	selectedArch = strings.TrimSpace(selectedArch)
	if selectedArch == "" || selectedArch == "1" {
		selectedArch = "layered"
	} else if selectedArch == "2" {
		selectedArch = "clean"
	}
	if arch != "" && arch != "layered" {
		selectedArch = arch
	}

	options := generator.ProjectOptions{
		Name:   name,
		Router: framework,
		Arch:   selectedArch,
		DB:     db,
		ORM:    orm,
	}
	if err := generator.NewProject(options); err != nil {
		return err
	}

	fmt.Fprintf(stdout, "created GoCraft project %s\n", options.Name)
	return nil
}

type choice struct {
	Label       string
	Description string
	Value       string
}

func selectChoice(label string, items []choice) (int, choice, error) {
	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}:",
		Active:   "{{ \">\" | cyan }} {{ .Label | cyan }} {{ .Description | faint }}",
		Inactive: "  {{ .Label }} {{ .Description | faint }}",
		Selected: "{{ \"OK\" | green }} {{ .Label | green }}",
	}
	selectPrompt := promptui.Select{
		Label:     "\033[35m" + label + "\033[0m",
		Items:     items,
		Templates: templates,
		Size:      len(items),
	}
	index, _, err := selectPrompt.Run()
	if err != nil {
		return 0, choice{}, err
	}
	return index, items[index], nil
}

func isTerminal(input io.Reader) bool {
	file, ok := input.(*os.File)
	if !ok {
		return false
	}
	info, err := file.Stat()
	if err != nil {
		return false
	}
	return info.Mode()&os.ModeCharDevice != 0
}

func runMake(args []string, stdout, stderr io.Writer) error {
	if len(args) == 0 {
		return errors.New("usage: gocraft make resource <name>")
	}
	if args[0] != "resource" {
		return fmt.Errorf("unknown make target %q", args[0])
	}

	flags := flag.NewFlagSet("make resource", flag.ContinueOnError)
	flags.SetOutput(stderr)
	if err := flags.Parse(args[1:]); err != nil {
		return err
	}
	if flags.NArg() != 1 {
		return errors.New("usage: gocraft make resource <name>")
	}

	if err := generator.NewResource(".", flags.Arg(0)); err != nil {
		return err
	}

	fmt.Fprintf(stdout, "created resource %s\n", flags.Arg(0))
	return nil
}

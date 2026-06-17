package cli

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"strings"

	"github.com/Mountok/gocraft/internal/generator"
)

const usage = `GoCraft - Go project architecture generator

Usage:
  gocraft new
  gocraft new <name> [default|gin]
  gocraft new <name> [--router default|nethttp|gin] [--arch layered]
  gocraft make resource <name>

Examples:
  gocraft new user-service
  gocraft new user-service gin
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
		return runNew(args[1:], stdin, stdout, stderr)
	case "make":
		return runMake(args[1:], stdout, stderr)
	default:
		return fmt.Errorf("unknown command %q\n\n%s", args[0], usage)
	}
}

func runNew(args []string, stdin io.Reader, stdout, stderr io.Writer) error {
	flags := flag.NewFlagSet("new", flag.ContinueOnError)
	flags.SetOutput(stderr)
	router := flags.String("router", "", "HTTP framework: default, nethttp, gin")
	arch := flags.String("arch", "layered", "architecture: layered")
	db := flags.String("db", "", "database support: planned")
	orm := flags.String("orm", "", "ORM support: planned")

	if err := flags.Parse(args); err != nil {
		return err
	}
	if flags.NArg() == 0 {
		return runInteractiveNew(stdin, stdout, *arch, *db, *orm)
	}
	if flags.NArg() > 2 {
		return errors.New("usage: gocraft new <name> [default|gin]")
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
	}

	options := generator.ProjectOptions{
		Name:   name,
		Router: framework,
		Arch:   arch,
		DB:     db,
		ORM:    orm,
	}
	if err := generator.NewProject(options); err != nil {
		return err
	}

	fmt.Fprintf(stdout, "created GoCraft project %s\n", options.Name)
	return nil
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

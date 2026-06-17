package cli

import (
	"errors"
	"flag"
	"fmt"
	"io"

	"github.com/Mountok/gocraft/internal/generator"
)

const usage = `GoCraft - Go project architecture generator

Usage:
  gocraft new <name> [--router nethttp] [--arch layered]
  gocraft make resource <name>

Examples:
  gocraft new user-service
  gocraft make resource user
`

func Run(args []string, stdout, stderr io.Writer) error {
	if len(args) == 0 {
		_, _ = io.WriteString(stdout, usage)
		return nil
	}

	switch args[0] {
	case "help", "--help", "-h":
		_, _ = io.WriteString(stdout, usage)
		return nil
	case "new":
		return runNew(args[1:], stdout, stderr)
	case "make":
		return runMake(args[1:], stdout, stderr)
	default:
		return fmt.Errorf("unknown command %q\n\n%s", args[0], usage)
	}
}

func runNew(args []string, stdout, stderr io.Writer) error {
	flags := flag.NewFlagSet("new", flag.ContinueOnError)
	flags.SetOutput(stderr)
	router := flags.String("router", "nethttp", "HTTP router: nethttp")
	arch := flags.String("arch", "layered", "architecture: layered")
	db := flags.String("db", "", "database support: planned")
	orm := flags.String("orm", "", "ORM support: planned")

	if err := flags.Parse(args); err != nil {
		return err
	}
	if flags.NArg() != 1 {
		return errors.New("usage: gocraft new <name> [--router nethttp] [--arch layered]")
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

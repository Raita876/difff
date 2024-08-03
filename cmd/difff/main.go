package main

import (
	"log"
	"os"

	"github.com/Raita876/difff/internal/difff"
	"github.com/urfave/cli/v2"
)

var (
	version string
	name    string
)

type Options struct{}

type Option interface {
	apply(*Options)
}

func (options *Options) Set(opts ...Option) {
	for _, o := range opts {
		o.apply(options)
	}
}

func run(source, target string, o *Options) error {
	return difff.Run(source, target)
}

func Run(c *cli.Context) error {
	source := c.Args().Get(0)
	target := c.Args().Get(1)

	o := &Options{}
	o.Set()

	return run(source, target, o)
}

func main() {
	cli.VersionFlag = &cli.BoolFlag{
		Name:    "version",
		Aliases: []string{"V"},
		Usage:   "print only the version",
	}

	app := &cli.App{
		Version:   version,
		Name:      name,
		Usage:     "This CLI compares files located in two directories and outputs the differences.",
		UsageText: "difff <source_path> <target_path>",
		Action:    Run,
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

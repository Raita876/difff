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

type Options struct {
	Format difff.FormatType
}

type Option interface {
	apply(*Options)
}

type formatOptions difff.FormatType

func (fo formatOptions) apply(o *Options) {
	o.Format = difff.FormatType(fo)
}

func (options *Options) Set(opts ...Option) {
	for _, o := range opts {
		o.apply(options)
	}
}

func run(source, target string, o *Options) error {
	return difff.Run(source, target, o.Format)
}

func Run(c *cli.Context) error {
	source := c.Args().Get(0)
	target := c.Args().Get(1)

	o := &Options{}
	o.Set(
		formatOptions(c.String("format")),
	)

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
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "format",
				Aliases: []string{"f"},
				Value:   "JSON",
				Usage:   "specify the output format. support: JSON, YAML, XML",
			},
		},
		Action: Run,
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

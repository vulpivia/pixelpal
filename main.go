package main

import (
	"os"

	"github.com/Acid147/pixelpal/commands"
	"github.com/Acid147/pixelpal/terminal"
	"github.com/mkideal/cli"
)

func main() {
	root := &cli.Command{
		Desc: "PixelPal is a command line utility for validating and converting the palette of image (especially pixel art) files.",
		Fn:   commands.Root,
		Argv: func() interface{} { return new(terminal.RootOptions) },
	}

	convert := &cli.Command{
		Name: "convert",
		Desc: "Convert an image (or multiple) to a target palette.",
		Fn:   commands.Convert,
		Argv: func() interface{} { return new(terminal.ConvertOptions) },
	}

	validate := &cli.Command{
		Name: "validate",
		Desc: "Check if an image (or multiple) uses a palette.",
		Fn:   commands.Validate,
		Argv: func() interface{} { return new(terminal.ValidateOptions) },
	}

	find := &cli.Command{
		CanSubRoute: true,
		Name:        "find",
		Desc:        "Find the palette used by an image (or multiple).",
		Fn:          commands.Find,
		Argv:        func() interface{} { return new(terminal.FindOptions) },
	}

	if err := cli.Root(
		root,
		cli.Tree(convert),
		cli.Tree(validate),
		cli.Tree(find),
	).Run(os.Args[1:]); err != nil {
		terminal.PrintError(err)
		os.Exit(1)
	}
}

package commands

import (
	"github.com/Acid147/pixelpal/terminal"
	"github.com/logrusorgru/aurora/v3"
	"github.com/mkideal/cli"
)

func Root(context *cli.Context) error {
	options := context.Argv().(*terminal.RootOptions)

	if options.Version {
		context.String("0.2.0")
		return nil
	}

	context.String(aurora.Bold("Usage").String() + ":\n")
	context.String("\n")
	context.String("  pixelpal <command>\n")
	context.String("  pixelpal -h | --help\n")
	context.String("  pixelpal -v | --version\n")
	context.String("\n")
	context.WriteUsage()
	return nil
}

package terminal

import (
	"log"
	"strings"

	"github.com/logrusorgru/aurora/v3"
	"github.com/mattn/go-colorable"
)

type HelpOption struct {
	Help bool `cli:"h,help" usage:"Show help"`
}

type PaletteOption struct {
	Palette string `cli:"p,palette" usage:"Name or path (.png) of the palette to use"`
}

type RootOptions struct {
	HelpOption
	Version bool `cli:"v,version" usage:"Show version information"`
}

type ConvertOptions struct {
	HelpOption
	PaletteOption
	Output   string `cli:"o,output" usage:"Output file(s)"`
	Contrast int    `cli:"c,contrast" usage:"Tolerance of the contrast in percent" dft:"0"`
	L        int    `cli:"l" usage:"Tolerance of the L* component in the CIELAB color space" dft:"0"`
	A        int    `cli:"a" usage:"Tolerance of the a* component in the CIELAB color space" dft:"0"`
	B        int    `cli:"b" usage:"Tolerance of the b* component in the CIELAB color space" dft:"0"`
	Step     int    `cli:"s,stepsize" usage:"Step size of each tolerance" dft:"1"`
	X        int    `cli:"x" usage:"Tile size in x direction" dft:"0"`
	Y        int    `cli:"y" usage:"Tile size in y direction" dft:"0"`
}

type ValidateOptions struct {
	HelpOption
	PaletteOption
}

type FindOptions struct {
	HelpOption
}

func PrintError(err error) {
	str := err.Error()
	str = strings.ToUpper(str[:1]) + str[1:]

	log.SetOutput(colorable.NewColorableStdout())
	log.SetFlags(0)

	log.Printf(aurora.Red("Error: %s.").String(), str)
	log.Printf("Try 'pixelpal --help' for more information.")
}

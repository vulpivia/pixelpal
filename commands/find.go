package commands

import (
	"image"
	"time"

	"github.com/Acid147/pixelpal/data"
	"github.com/Acid147/pixelpal/io"
	"github.com/Acid147/pixelpal/terminal"
	"github.com/logrusorgru/aurora/v3"
	"github.com/mkideal/cli"
	"github.com/pegasus-toolset/color"
)

func Find(context *cli.Context) error {
	options := context.Argv().(*terminal.FindOptions)

	if options.Help {
		context.String(aurora.Bold("Usage").String() + ":\n")
		context.String("\n")
		context.String("  pixelpal find <file>...\n")
		context.String("  pixelpal find -h | --help\n")
		context.String("\n")
		context.WriteUsage()
		return nil
	}

	start := time.Now()

	context.String("Loading input file(s)...\n")
	images, err := io.LoadInputFiles(context.Args())
	if err != nil {
		return err
	}
	context.String(aurora.Green("Done!").String() + "\n")

	// Find palette for a single image
	if len(images) == 1 {
		palette := findPalette(images[0])
		if len(palette) != 0 {
			context.String(aurora.Blue("Image uses palette '"+palette[0]+"'.").String() + "\n")
		} else {
			context.String(aurora.BrightBlack("Image uses an unknown palette.").String() + "\n")
		}

		context.String("Time taken: " + time.Since(start).String())
		return nil
	}

	paths := context.Args()
	for i, img := range images {
		palette := findPalette(img)
		if len(palette) != 0 {
			context.String(aurora.Blue("Image '"+paths[i]+"' uses palette '"+palette[0]+"'.").String() + "\n")
		} else {
			context.String(aurora.BrightBlack("Image '"+paths[i]+"' uses an unknown palette.").String() + "\n")
		}
	}

	context.String("Time taken: " + time.Since(start).String())
	return nil
}

func findPalette(img image.Image) []string {
	// Get image dimensions
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	var matches []string
	for name, palette := range data.Palettes {
		match := true
	out:
		for x := 0; x < width; x++ {
			for y := 0; y < height; y++ {
				c := color.RGBModel.Convert(img.At(x, y)).(color.Color)
				if !palette.Contains(c) {
					match = false
					break out
				}
			}
		}

		if match {
			matches = append(matches, name)
		}
	}

	return matches
}

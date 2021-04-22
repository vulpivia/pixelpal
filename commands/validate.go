package commands

import (
	"errors"
	"image"
	"time"

	"github.com/Acid147/pixelpal/io"
	"github.com/Acid147/pixelpal/terminal"
	"github.com/logrusorgru/aurora/v3"
	"github.com/mkideal/cli"
	"github.com/pegasus-toolset/color"
)

func Validate(context *cli.Context) error {
	options := context.Argv().(*terminal.ValidateOptions)

	if options.Help {
		context.String(aurora.Bold("Usage").String() + ":\n")
		context.String("\n")
		context.String("  pixelpal validate [options] <file>...\n")
		context.String("  pixelpal validate -h | --help\n")
		context.String("\n")
		context.WriteUsage()
		return nil
	}

	if options.Palette == "" {
		return errors.New("the option '--palette' is required")
	}

	start := time.Now()

	context.String("Loading input file(s)...\n")
	images, err := io.LoadInputFiles(context.Args())
	if err != nil {
		return err
	}
	context.String(aurora.Green("Done!").String() + "\n")

	context.String("Loading palette...\n")
	palettePath := options.Palette
	palette, err := io.LoadPalette(palettePath)
	if err != nil {
		return err
	}
	context.String(aurora.Green("Done!").String() + "\n")

	// Validate single image
	if len(images) == 1 {
		if validateImage(images[0], palette) {
			context.String(aurora.Blue("Image uses palette '"+palettePath+"'.").String() + "\n")
		} else {
			context.String(aurora.BrightBlack("Image doesn't use palette '"+palettePath+"'.").String() + "\n")
		}

		context.String("Time taken: " + time.Since(start).String())
		return nil
	}

	paths := context.Args()
	for i, img := range images {
		if validateImage(img, palette) {
			context.String(aurora.Blue("Image '"+paths[i]+"' uses palette '"+palettePath+"'.").String() + "\n")
		} else {
			context.String(aurora.BrightBlack("Image '"+paths[i]+"' doesn't use palette '"+palettePath+"'.").String() + "\n")
		}
	}

	context.String("Time taken: " + time.Since(start).String())
	return nil
}

func validateImage(img image.Image, palette color.Palette) bool {
	// Get image dimensions
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			c := color.RGBModel.Convert(img.At(x, y)).(color.Color)
			if !palette.Contains(c) {
				return false
			}
		}
	}

	return true
}

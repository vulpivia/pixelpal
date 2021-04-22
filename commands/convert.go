package commands

import (
	"errors"
	"image"
	"image/color"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/Acid147/pixelpal/io"
	"github.com/Acid147/pixelpal/terminal"
	"github.com/logrusorgru/aurora/v3"
	"github.com/mkideal/cli"

	col "github.com/pegasus-toolset/color"
	pal "github.com/pegasus-toolset/color/palette"
)

type step uint8

const (
	Contrast step = iota
	L
	A
	B
)

// Convert an image to a palette.
func Convert(context *cli.Context) error {
	options := context.Argv().(*terminal.ConvertOptions)

	if options.Help {
		context.String(aurora.Bold("Usage").String() + ":\n")
		context.String("\n")
		context.String("  pixelpal convert [options] <file>...\n")
		context.String("  pixelpal convert -h | --help\n")
		context.String("\n")
		context.WriteUsage()
		context.String("\n")
		context.String("If the options " + aurora.Bold("-x").String() + " and/or " + aurora.Bold("-y").String() + " are set and not 0, the input image is split into multiple tiles. Each tile is treated as a separate image, meaning that each tile conversion may use different calculated offsets for contrast, L*, a*, and b*.")
		return nil
	}

	if options.Palette == "" {
		return errors.New("the option '--palette' is required")
	}

	if options.Output == "" {
		return errors.New("the option '--output' is required")
	}

	// Measure time taken
	start := time.Now()

	// Parse option flags and use them as parameter
	contrastTolerance, tolerances := getTolerances(options)
	stepSize := float64(options.Step)

	context.String("Loading input file(s)...\n")
	inputImages, err := io.LoadInputFiles(context.Args())
	if err != nil {
		return err
	}
	context.String(aurora.Green("Done!").String() + "\n")

	context.String("Loading palette...\n")
	palette, err := io.LoadPalette(options.Palette)
	if err != nil {
		return err
	}
	context.String(aurora.Green("Done!").String() + "\n")

	// Get output file name and remove file extension
	outputPath := options.Output
	// Remove ".png" suffix
	if strings.HasSuffix(outputPath, ".png") {
		outputPath = outputPath[:len(outputPath)-4]
	}

	// Save single image
	if len(inputImages) == 1 {
		context.String("Converting image...\n")

		err := convertAndSave(inputImages[0], palette, contrastTolerance, tolerances, stepSize, options.X, options.Y, outputPath)
		if err != nil {
			return err
		}

		context.String(aurora.Green("Done!").String() + "\n")
		context.String("Time taken: " + time.Since(start).String())
		return nil
	}

	// Save multiple images
	for i, inputImage := range inputImages {
		context.String("Converting image #" + strconv.Itoa(i+1) + "...\n")

		path := outputPath + strconv.Itoa(i+1) + ".png"
		err := convertAndSave(inputImage, palette, contrastTolerance, tolerances, stepSize, options.X, options.Y, path)
		if err != nil {
			return err
		}

		context.String(aurora.Green("Done!").String() + "\n")
	}

	context.String("Time taken: " + time.Since(start).String())
	return nil
}

// Convert an image to a target palette, adjusting the colors of the input image
// by modifiers for contrast, L*, a*, and b* values before converting.
func convertImage(img image.Image, palette col.Palette, contrast float64, modifier col.Lab) *image.Paletted {
	// Get image dimensions
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	// Create output image
	result := image.NewPaletted(bounds, palette.ToStdPalette())

	// Fill output image
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			c := col.LabModel.Convert(img.At(x, y)).(col.Lab)

			c.L = (c.L-50)*(contrast+1) + 50
			c.A *= contrast + 1
			c.B *= contrast + 1

			c.L += modifier.L
			c.A += modifier.A
			c.B += modifier.B

			result.Set(x, y, convertColor(c, palette))
		}
	}

	return result
}

// Convert a single color to its nearest color in a palette.
func convertColor(c col.Lab, palette col.Palette) color.Color {
	minDistance := math.MaxFloat64

	var result color.Color
	for _, paletteColor := range palette {
		distance := c.DistanceTo(paletteColor)
		if distance < minDistance {
			minDistance = distance
			result = paletteColor
		}
	}

	return result
}

// calculateModifiers calculates the values by which contrast, L*, a*, and b*
// are shifted before converting an image.
func calculateModifiers(img image.Image, palette col.Palette, contrastTolerance float64, tolerances col.Lab, stepSize float64) (float64, col.Lab) {
	sourcePalette := pal.FromImage(img)

	var currentContrast float64
	current := col.Lab{}
	var bestContrast float64
	best := col.Lab{L: math.MaxFloat64, A: math.MaxFloat64, B: math.MaxFloat64}

	// Go through every possible contrast modifier value to find the best
	// contrast.
	colorCount := 0
	for currentContrast = -contrastTolerance; currentContrast <= contrastTolerance; currentContrast += stepSize / 100 {
		colorCount, bestContrast = calculateModifier(currentContrast, current, bestContrast, best, Contrast, colorCount, sourcePalette, palette)
	}

	// Go through every possible L* modifier value to find the best L*.
	colorCount = 0
	for current.L = -tolerances.L; current.L <= tolerances.L; current.L += stepSize {
		colorCount, best.L = calculateModifier(currentContrast, current, bestContrast, best, L, colorCount, sourcePalette, palette)
	}

	// Go through every possible a* modifier value to find the best a*.
	colorCount = 0
	for current.A = -tolerances.A; current.A <= tolerances.A; current.A += stepSize {
		colorCount, best.A = calculateModifier(currentContrast, current, bestContrast, best, A, colorCount, sourcePalette, palette)
	}

	// Go through every possible b* modifier value to find the best b*.
	colorCount = 0
	for current.B = -tolerances.B; current.B <= tolerances.B; current.B += stepSize {
		colorCount, best.B = calculateModifier(currentContrast, current, bestContrast, best, B, colorCount, sourcePalette, palette)
	}

	return bestContrast, best
}

// Calculate the modifier of the value for contrast, L*, a*, or b*.
//
// If the modifier in the parameter "current" is better than the current best,
// it returns the new best value for contrast, L*, a*, or b* and the count of
// resulting colors using the target palette. If not, it returns the old best
// value (e.g. best.contrast if the contrast value modifier is calculated).
//
// "step" decides which modifier value is calculated. For calculating L*, the
// best contrast modifier must already be set. For calculating a*, the best
// contrast and L* modifiers must already be set. For calculating b*, the best
// contrast, L*, and a* modifiers must already be set.
func calculateModifier(currentContrast float64, current col.Lab, bestContrast float64, best col.Lab, step step, colorCount int, sourcePalette col.Palette, targetPalette col.Palette) (int, float64) {
	var lab col.Lab
	var usedColors col.Palette

	// For each color in the source image, find the optimal target color by
	// looping through all colors in the target palette.
	for _, sourceColor := range sourcePalette {
		bestDistance := math.MaxFloat64
		bestIndex := 0
		for i, targetColor := range targetPalette {
			lab = col.LabModel.Convert(sourceColor).(col.Lab)

			// Apply contrast modifier
			if step == Contrast {
				lab.L = (lab.L-50)*(currentContrast+1) + 50
				lab.A *= currentContrast + 1
				lab.B *= currentContrast + 1
			} else {
				lab.L = (lab.L-50)*(bestContrast+1) + 50
				lab.A *= bestContrast + 1
				lab.B *= bestContrast + 1
			}

			// Apply L*, a*, and b* modifiers
			switch step {
			case L:
				lab.L += current.L
			case A:
				lab.L += best.L
				lab.A += current.A
			case B:
				lab.L += best.L
				lab.A += best.A
				lab.B += current.B
			}

			// If a closer target color is found, use that one as the target
			// color.
			distanceToTargetColor := lab.DistanceTo(targetColor)
			if bestDistance > distanceToTargetColor {
				bestDistance = distanceToTargetColor
				bestIndex = i
			}
		}

		if !usedColors.Contains(targetPalette[bestIndex]) {
			usedColors = append(usedColors, targetPalette[bestIndex])
		}
	}

	// Check if an improvement has been made (more colors of the target palette
	// have been used, meaning there is less loss of color information when
	// converting the image)
	//
	// If not, return the same modifier value that existed before.
	moreColors := colorCount < len(usedColors)
	switch step {
	case Contrast:
		if moreColors && math.Abs(currentContrast) < math.Abs(bestContrast) {
			return len(usedColors), currentContrast
		}
		return colorCount, bestContrast
	case L:
		if moreColors && math.Abs(current.L) < math.Abs(best.L) {
			return len(usedColors), current.L
		}
		return colorCount, best.L
	case A:
		if moreColors && math.Abs(current.A) < math.Abs(best.A) {
			return len(usedColors), current.A
		}
		return colorCount, best.A
	case B:
		if moreColors && math.Abs(current.B) < math.Abs(best.B) {
			return len(usedColors), current.B
		}
		return colorCount, best.B
	}

	// Step must be Contrast, L, A, or B. An invalid step results in a panic.
	panic("Step outside of valid range!")
}

// Split an image into tiles of width x and height y.
func split(img image.Image, x, y int) []image.Image {
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()
	if x == 0 {
		x = width
	}
	if y == 0 {
		y = height
	}

	var tiles []image.Image

	for currentX := 0; currentX < width; currentX += x {
		// Smaller image if image width wasn't a multiple of tile width
		actualX := x
		if width-currentX < x {
			actualX = width - currentX
		}
		for currentY := 0; currentY < height; currentY += y {
			// Smaller image if image height wasn't a multiple of tile height
			actualY := y
			if height-currentY < y {
				actualY = height - currentY
			}

			tile := image.NewRGBA(image.Rect(0, 0, x, y))
			for pixelX := 0; pixelX < actualX; pixelX++ {
				for pixelY := 0; pixelY < actualY; pixelY++ {
					c := img.At(currentX+pixelX, currentY+pixelY)
					tile.Set(pixelX, pixelY, c)
				}
			}
			tiles = append(tiles, tile)
		}
	}

	return tiles
}

// Join tiles into a full image of size bounds.
func join(tiles []image.Image, bounds image.Rectangle) image.Image {
	tileBounds := tiles[0].Bounds()
	tileWidth := tileBounds.Dx()
	tileHeight := tileBounds.Dy()

	img := image.NewRGBA(bounds)
	imgWidth := bounds.Dx()
	imgHeight := bounds.Dy()
	for x := 0; x < imgWidth; x += tileWidth {
		for y := 0; y < imgHeight; y += tileHeight {
			tile := tiles[y/tileHeight+(x/tileWidth)*(int(math.Ceil(float64(imgHeight)/float64(tileHeight))))]
			currentTileBounds := tile.Bounds()
			currentTileWidth := currentTileBounds.Dx()
			currentTileHeight := currentTileBounds.Dy()
			for pixelX := 0; pixelX < currentTileWidth; pixelX++ {
				for pixelY := 0; pixelY < currentTileHeight; pixelY++ {
					c := tile.At(pixelX, pixelY)
					img.Set(x+pixelX, y+pixelY, c)
				}
			}
		}
	}
	return img
}

// Convert the command line flags into usable values.
func getTolerances(options *terminal.ConvertOptions) (float64, col.Lab) {
	result := col.Lab{
		L: float64(options.L),
		A: float64(options.A),
		B: float64(options.B),
	}
	return float64(options.Contrast) / 100, result
}

// Convert one image and save it to an output file.
func convertAndSave(image image.Image, palette col.Palette, contrastTolerance float64, tolerances col.Lab, stepSize float64, tileWidth, tileHeight int, path string) error {
	tiles := split(image, tileWidth, tileHeight)

	for i, tile := range tiles {
		contrast, modifiers := calculateModifiers(tile, palette, contrastTolerance, tolerances, stepSize)
		tiles[i] = convertImage(tile, palette, contrast, modifiers)
	}

	image = join(tiles, image.Bounds())

	return io.SaveImage(path+".png", image)
}

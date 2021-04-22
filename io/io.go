package io

import (
	"errors"
	"image"
	"image/png"
	"os"
	"strings"

	"github.com/Acid147/pixelpal/data"
	"github.com/pegasus-toolset/color"
)

func LoadInputFiles(paths []string) ([]image.Image, error) {
	if len(paths) == 0 {
		return nil, errors.New("no input file specified")
	}

	var images []image.Image
	for _, path := range paths {
		img, err := loadImage(path)
		if err != nil {
			return nil, err
		}

		images = append(images, img)
	}

	return images, nil
}

func LoadPalette(name string) (color.Palette, error) {

	p := data.Palettes[strings.ToLower(name)]
	if p != nil {
		return p, nil
	}

	img, err := loadImage(name)
	if err != nil {
		return nil, err
	}

	var palette color.Palette
	bounds := img.Bounds()
	for x := 0; x < bounds.Dx(); x++ {
		for y := 0; y < bounds.Dy(); y++ {
			c := color.RGBModel.Convert(img.At(x, y)).(color.RGB)

			if palette.Contains(c) {
				continue
			}

			palette = append(palette, c)
		}
	}
	return palette, nil
}

func SaveImage(path string, img image.Image) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	err = png.Encode(file, img)
	if err != nil {
		return err
	}

	return nil
}

func loadImage(path string) (image.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	result, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	return result, nil
}

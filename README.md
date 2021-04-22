# PixelPal

**PixelPal** is a command line utility for validating and converting the palette
of image (especially pixel art) files.

When using third party pixel art assets, the palettes of different assets may
not match. To unify the look of these assets, you might want to convert their
color palettes to match.

## Usage

PixelPal can be run using a command-line interface:

```sh
pixelpal convert [options] <file>...
pixelpal validate [options] <file>...
pixelpal find <file>...
```

### Convert

The `convert` command converts an image (or multiple images) to a target
palette.

Valid options are:

- `-h`, `--help`: Show help
- `-p`, `--palette`: Name or path (.png) of the palette to use (required)
- `-o`, `--output`: Output file(s) (required)
- `-c`, `--contrast`: Tolerance of the contrast in percent (default value is
  `0`)
- `-l`: Tolerance of the L* component in the CIELAB color space (default value
  is `0`)
- `-a`: Tolerance of the a* component in the CIELAB color space (default value
  is `0`)
- `-b`: Tolerance of the b* component in the CIELAB color space (default value
  is `0`)
- `-s`, `--stepsize`: Step size of each tolerance (default value is `1`)

Tolerances represent a spectrum in which the source image can be moved to get a
result that better utilizes the target palette, although the higher a tolerance,
the more the resulting image may differ from the source image. If for example a
contrast tolerance of 10% was chosen, the contrast of the input image may be
increased or decreased by up to 10% to make better use of all the colors in the
target palette.

The L\*a\*b\* component tolerances control the following:

- **L\***: A lower L\* value darkens the image, a higher value lightens it.
- **a\***: A lower a\* value shifts all colors towards green, a higher value
  shifts them towards red.
- **b\***: A lower b\* value shifts all colors towards blue, a higher value
  shifts them towards yellow.

The step size controls the intervals in which the input image is manipulated
according to the tolerances. A step size of `5` and a contrast tolerance of
10% means that the following contrasts are used:

- 90%
- 95%
- 100%
- 105%
- 110%

### Validate

The `validate` command checks if an image (or multiple images) uses a palette.

Valid options are:

- `-h`, `--help`: Show help
- `-p`, `--palette`: Name or path (.png) of the palette to use (required)

### Find

The `find` command finds the palette used by an image (or multiple images).

Valid options are:

- `-h`, `--help`: Show help

## Contributing

Pull requests are welcome.

### Project Structure

The commands of PixelPal are defined in [`main.go`](main.go). The
options for each command are defined in
[`terminal/terminal.go`](terminal/terminal.go). The program then gives
control to the chosen command in `commands/` (e.g. `commands/convert.go`) to the
function with the same name as the command.

Each command function is responsible for its help text, validating required
options and of course executing the command itself.

#### Convert

The `convert` command executes the following steps:

- It first loads the input files (input images and palette)
- Then it starts converting the input image (or loops through each input image,
  if multiple are given)
  - The input image is split into tiles (see function `split`)
  - For each tile, modifiers are calculated (see function `calculateModifiers`)
  - Each tile gets converted according to its modifier values (see function
    `convertImage`)
  - The tiles are joined together again (see function `join`)
  - The output image gets saved

#### Validate

The `validate` command executes the following steps:

- It first loads the input files (input images and palette)
- Then it starts validating the input image (or loops through each input image,
  if multiple are given)
- For each image, it goes through each pixel and checks if the given palette
  contains the color of the pixel

#### Find

The `find` command executes the following steps:

- It first loads the input images
- Then it goes through each input image and tries to find a matching palette by
  going through each known palette and checking each pixel of the image against
  that palette
- If at least one matching palette has been found, the name of the first palette
  that matches the colors of the input image is printed to the terminal.

## Credits

PixelPal is written by [Norman Rauschen](https://github.com/Acid147).

The project uses the following open source packages:

- [cli](https://github.com/mkideal/cli)
- [color](https://github.com/pegasus-toolset/color)

## License

This project is released under the [Unlicense](LICENSE.md).

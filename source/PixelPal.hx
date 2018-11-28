import sys.io.File;
import tink.cli.Rest;

@alias(false)
class PixelPal
{
    public var help:Bool;
    public var version:Bool;
    public var palette:String;
    public var output:String;

    public function new() {}

    @:defaultCommand
    public function run(rest:Rest<String>)
    {
        if (help)
        {
            runHelp();
            return;
        }

        if (version)
        {
            runVersion();
            return;
        }

        if (palette == null)
        {
            Sys.println("The palette option is required.");
            return;
        }

        if (output != null)
        {
            runConvert(rest);
            return;
        }

        runValidate(rest);
    }

    function runHelp()
    {
        Sys.println("Usage:");
        Sys.println("  pixelpal [OPTIONS] [ARGS]");
        Sys.println("");
        Sys.println("Options:");
        Sys.println("  -h, --help: Help");
        Sys.println("  -v, --version: Version");
        Sys.println("  -p PALETTE, --palette=PALETTE: Name of the palette to use");
        Sys.println("  -o OUTPUT, --output=OUTPUT: Convert to palette and save to output file");
        Sys.println("");
        Sys.println("Arguments:");
        Sys.println("  INPUT: Input file");
    }

    function runVersion()
    {
        Sys.println("pixelpal 0.1.0");
    }

    function runConvert(rest:Rest<String>)
    {
        var images = loadInputFiles(rest);
        if (images == null)
        {
            return;
        }

        var palette = Palette.fromPNG(palette);
        if (palette == null)
        {
            Sys.println("Unable to read palette file '" + palette + "'");
            return;
        }

        for (i in 0...images.length)
        {
            images[i].convert(palette);

            if (images.length > 1)
            {
                if (!images[i].save(output + i))
                {
                    Sys.println("Unable to write file '" + output + i + "'");
                }
            }
            else
            {
                if (!images[i].save(output))
                {
                    Sys.println("Unable to write file '" + output + "'");
                }
            }
        }

        Sys.println("Conversion finished.");
    }

    function runValidate(rest:Rest<String>)
    {
        var images = loadInputFiles(rest);
        if (images == null)
        {
            return;
        }

        var palette = Palette.fromPNG(palette);
        if (palette == null)
        {
            Sys.println("Unable to read palette file '" + palette + "'");
            return;
        }

        var errorCount = 0;

        for (i in 0...images.length)
        {
            if (!images[i].validate(palette))
            {
                Sys.println("File '" + rest[i] + "' contains colors outside of the specified palette.");
                errorCount++;
            }
        }

        Sys.println("Validation finished. " + errorCount + " of " + images.length + " files contain colors outside of the specified palette.");
    }

    function loadInputFiles(rest:Rest<String>):Array<Image>
    {
        if (rest.length == 0)
        {
            Sys.println("No input file specified");
            Sys.println("Try 'pixelpal --help' for more information.");
            return null;
        }

        var images = [];

        for (path in rest)
        {
            var image = new Image(path);
            if (image.empty) {
                Sys.println("Unable to read file '" + path + "'");
                return null;
            }
            images.push(image);
        }

        return images;
    }
}
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
    }

    function runValidate(rest:Rest<String>)
    {
        var images = loadInputFiles(rest);
        if (images == null)
        {
            return;
        }
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
import format.png.Reader;
import format.png.Tools;
import haxe.io.BytesData;
import sys.io.File;

class Palette
{
    var colors:Array<Color>;

    /**
        Get palette information from a PNG file.

        @param path path of the PNG file
        @return the palette
    **/
    public static function fromPNG(path:String):Palette
    {
        var data;
        var width;
        var height;

        try
        {
            var handle = File.read(path, true);
            var d = new Reader(handle).read();
            var hdr = Tools.getHeader(d);

            data = Tools.extract32(d).getData();
            width = hdr.width;
            height = hdr.height;

            handle.close();
        }
        catch (error:Dynamic)
        {
            return null;
        }

        var colors:Array<Color> = [];

        for (x in 0...width)
        {
            for (y in 0...height)
            {
                var b = data[x * 4 + y * width * 4];
                var g = data[x * 4 + y * width * 4 + 1];
                var r = data[x * 4 + y * width * 4 + 2];
                var a = data[x * 4 + y * width * 4 + 3];

                var color = new Color(r, g, b);
                if (colors.filter(function(c) return c.r == color.r && c.g == color.g && c.b == color.b).length == 0)
                {
                    colors.push(color);
                }
            }
        }

        return new Palette(colors);
    }

    public function new(colors:Array<Color>)
    {
        this.colors = colors;
    }

    public function contains(color:Color):Bool
    {
        return colors.filter(function(c) return c.r == color.r && c.g == color.g && c.b == color.b).length > 0;
    }
}
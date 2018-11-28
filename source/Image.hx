import format.png.Reader;
import format.png.Tools;
import format.png.Writer;
import haxe.io.Bytes;
import haxe.io.BytesData;
import sys.io.File;

class Image
{
    /**
        True if the image has been loaded successfully.
    **/
    public var empty:Bool;

    var bytes:Bytes;
    var data:BytesData;
    var width:Int;
    var height:Int;

    public function new(path:String)
    {
        try
        {
            var handle = File.read(path, true);
            var d = new Reader(handle).read();
            var hdr = Tools.getHeader(d);

            bytes = Tools.extract32(d);
            data = bytes.getData();
            width = hdr.width;
            height = hdr.height;

            handle.close();
        }
        catch (error:Dynamic)
        {
            empty = true;
            return;
        }

        empty = false;
    }

    /**
        Check if the image contains colors outside of a palette.

        @param palette the palette that contains all valid colors
        @return true if all colors of the image are contained in the palette
    **/
    public function validate(palette:Palette):Bool
    {
        for (x in 0...width)
        {
            for (y in 0...height)
            {
                var b = data[x * 4 + y * width * 4];
                var g = data[x * 4 + y * width * 4 + 1];
                var r = data[x * 4 + y * width * 4 + 2];
                var a = data[x * 4 + y * width * 4 + 3];

                if (a == 0)
                {
                    continue;
                }

                var color = new Color(r, g, b);
                if (!palette.contains(color))
                {
                    return false;
                }
            }
        }

        return true;
    }

    /**
        Convert all pixels to colors contained in the palette.

        @param palette the palette that contains the valid colors
    **/
    public function convert(palette:Palette)
    {
        for (x in 0...width)
        {
            for (y in 0...height)
            {
                var b = data[x * 4 + y * width * 4];
                var g = data[x * 4 + y * width * 4 + 1];
                var r = data[x * 4 + y * width * 4 + 2];
                var a = data[x * 4 + y * width * 4 + 3];

                if (a == 0)
                {
                    continue;
                }

                var color = new Color(r, g, b).convert(palette);

                data[x * 4 + y * width * 4] = Std.int(color.b);
                data[x * 4 + y * width * 4 + 1] = Std.int(color.g);
                data[x * 4 + y * width * 4 + 2] = Std.int(color.r);
                data[x * 4 + y * width * 4 + 3] = 255;
            }
        }
    }

    /**
        Save the image.

        @param path the path of the output file
        @return true if saving was successful
    **/
    public function save(path:String):Bool
    {
        for (i in 0...data.length)
        {
            bytes.set(i, data[i]);
        }

        try
        {
            var out = File.write(path, true);
            new Writer(out).write(Tools.build32BGRA(width, height, bytes));
        }
        catch (error:Dynamic)
        {
            return false;
        }

        return true;
    }
}
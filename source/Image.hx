import format.png.Reader;
import format.png.Tools;
import haxe.io.BytesData;
import sys.io.File;

class Image
{
    /**
        True if the image has been loaded successfully.
    **/
    public var empty:Bool;

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

            data = Tools.extract32(d).getData();
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

    public function validate(palette:Palette)
    {
        for (x in 0...width)
        {
            for (y in 0...height)
            {
                var b = data[x * 4 + y * width * 4];
                var g = data[x * 4 + y * width * 4 + 1];
                var r = data[x * 4 + y * width * 4 + 2];
                var a = data[x * 4 + y * width * 4 + 3];

                var color = new Color(r, g, b);
                if (!palette.contains(color))
                {
                    return false;
                }
            }
        }

        return true;
    }
}
import format.png.Reader;
import format.png.Tools;
import haxe.io.BytesData;
import sys.io.File;

class Image
{
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
        catch(error:Dynamic)
        {
            empty = true;
            return;
        }

        empty = false;
    }
}
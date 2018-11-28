class Color
{
    public var r:Float;
    public var g:Float;
    public var b:Float;

    /**
        Convert a color from RGB to LAB.

        @param r red, from 0-1
        @param g green, from 0-1
        @param b blue, from 0-1
        @return the LAB color
    **/
    public static function rgb2lab(r:Float, g:Float, b:Float):Color
    {
        r = (r > 0.04045) ? Math.pow((r + 0.055) / 1.055, 2.4) : r / 12.92;
        g = (g > 0.04045) ? Math.pow((g + 0.055) / 1.055, 2.4) : g / 12.92;
        b = (b > 0.04045) ? Math.pow((b + 0.055) / 1.055, 2.4) : b / 12.92;

        var x = (r * 0.4124 + g * 0.3576 + b * 0.1805) / 0.95047;
        var y = (r * 0.2126 + g * 0.7152 + b * 0.0722) / 1.00000;
        var z = (r * 0.0193 + g * 0.1192 + b * 0.9505) / 1.08883;

        x = (x > 0.008856) ? Math.pow(x, 1/3) : (7.787 * x) + 16/116;
        y = (y > 0.008856) ? Math.pow(y, 1/3) : (7.787 * y) + 16/116;
        z = (z > 0.008856) ? Math.pow(z, 1/3) : (7.787 * z) + 16/116;

        return new Color((116 * y) - 16, 500 * (x - y), 200 * (y - z));
    }

    public function new(r:Float, g:Float, b:Float)
    {
        this.r = r;
        this.g = g;
        this.b = b;
    }

    public function convert(palette:Palette):Color
    {
        var diff:Float = 1000;
        var outputColor = new Color(0, 0, 0);

        for (paletteColor in palette.colors)
        {
            var colorR = r / 255;
            var colorG = g / 255;
            var colorB = b / 255;
            var paletteR = paletteColor.r / 255;
            var paletteG = paletteColor.g / 255;
            var paletteB = paletteColor.b / 255;

            var colorLAB = Color.rgb2lab(colorR, colorG, colorB);
            var paletteLAB = Color.rgb2lab(paletteR, paletteG, paletteB);

            var currentDiff = Math.sqrt(Math.pow(colorLAB.r - paletteLAB.r, 2) + Math.pow(colorLAB.g - paletteLAB.g, 2) + Math.pow(colorLAB.b - paletteLAB.b, 2));
            if (currentDiff < diff)
            {
                diff = currentDiff;
                outputColor = paletteColor;
            }
        }

        return outputColor;
    }
}
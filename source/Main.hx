import tink.Cli;

class Main
{
    static function main()
    {
        Cli.process(Sys.args(), new PixelPal()).handle(Cli.exit);
    }
}
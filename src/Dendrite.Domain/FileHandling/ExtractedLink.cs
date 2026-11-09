namespace Dendrite.Domain.FileHandling;

public sealed record ExtractedLink
{
    public string Raw { get; }
    public string Display { get; }
    public string TargetSlug { get; }
    public int Line { get; }
    public int Col { get; }

    public ExtractedLink(string raw, string display, string targetSlug, int line, int col)
    {
        Raw = raw;
        Display = display;
        TargetSlug = targetSlug;
        Line = line;
        Col = col;
    }
}

namespace Dendrite.Infrastructure.Persistence.Rows;

internal sealed record LinkRow
{
    public string Id { get; }
    public string TargetSlug { get; }
    public string Raw { get; }
    public string? Display { get; }
    public int Line { get; }
    public int Col { get; }

    public LinkRow(string id, string targetSlug, string raw, string? display, int line, int col)
    {
        Id = id;
        TargetSlug = targetSlug;
        Raw = raw;
        Display = display;
        Line = line;
        Col = col;
    }
}

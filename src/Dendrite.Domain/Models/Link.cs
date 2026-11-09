namespace Dendrite.Domain.Models;

/// <summary>A wiki link found in a note. Owned by the <see cref="Note"/> aggregate.</summary>
public sealed class Link
{
    public string Id { get; }
    public string TargetSlug { get; }
    public string Raw { get; }
    public string Display { get; }
    public int Line { get; }
    public int Col { get; }

    public Link(string id, string targetSlug, string raw, string display, int line, int col)
    {
        Id = id;
        TargetSlug = targetSlug;
        Raw = raw;
        Display = display;
        Line = line;
        Col = col;
    }

    public static Link Create(string targetSlug, string raw, string display, int line, int col) =>
        new(Guid.CreateVersion7().ToString(), targetSlug, raw, display, line, col);
}

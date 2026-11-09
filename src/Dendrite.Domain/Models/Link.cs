namespace Dendrite.Domain.Models;

#pragma warning disable IDE0290 // Use primary constructor
public sealed class Link
{
    public string Id { get; }
    public string FromNoteId { get; }
    public string TargetSlug { get; }
    public string Raw { get; }
    public string Display { get; }
    public int Line { get; }
    public int Col { get; }

    public Link(string id, string fromNoteId, string targetSlug, string raw, string display, int line, int col)
    {
        Id = id;
        FromNoteId = fromNoteId;
        TargetSlug = targetSlug;
        Raw = raw;
        Display = display;
        Line = line;
        Col = col;
    }

    public static Link Create(string fromNoteId, string targetSlug, string raw, string display, int line, int col) =>
        new(Guid.CreateVersion7().ToString(), fromNoteId, targetSlug, raw, display, line, col);
}

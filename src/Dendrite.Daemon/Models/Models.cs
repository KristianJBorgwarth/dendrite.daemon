namespace Dendrite.Daemon.Models;

public sealed record Note(string Id, string Path, string? Title, string Slug, string CreatedAt, string UpdatedAt)
{
    public static Note Create(string path, string title, string slug)
    {
        var now = DateTime.Now.ToString("yyyy-MM-dd HH:mm:ss");
        return new Note(Guid.CreateVersion7().ToString(), path, title, slug, now, now);
    }
}

public sealed record Link(string Id, string FromNoteId, string TargetSlug, string Raw, string Display, int Line, int Col)
{
    public static Link Create(string fromNoteId, string targetSlug, string raw, string display, int line, int col) =>
        new(Guid.CreateVersion7().ToString(), fromNoteId, targetSlug, raw, display, line, col);
}

public sealed record Tag(string Name);

public sealed record NoteTag(string NoteId, string TagId);

public sealed record CustomFrontMatter(string NoteId, string Key, string Value);

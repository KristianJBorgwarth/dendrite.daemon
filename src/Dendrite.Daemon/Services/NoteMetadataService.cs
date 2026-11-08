using Dendrite.Daemon.FileHandling;
using Dendrite.Daemon.Models;
using Dendrite.Daemon.Persistence.Repositories;

namespace Dendrite.Daemon.Services;

/// <summary>Writes and removes the index rows (tags, links, custom frontmatter) that hang off a note.</summary>
public sealed class NoteMetadataService(
    NoteRepository notes,
    TagRepository tags,
    LinkRepository links,
    CfeRepository cfe)
{
    public void WriteTags(string noteId, IReadOnlyList<string> tagNames)
    {
        var names = tagNames.Distinct().ToList();
        var existing = tags.GetByNames(names);
        var created = names
            .Where(n => existing.All(t => t.Name != n))
            .Select(n => new Tag(n))
            .ToList();

        tags.Insert(created);
        tags.InsertNoteTags(names.Select(n => new NoteTag(noteId, n)));
    }

    public void WriteLinks(string noteId, IReadOnlyList<ExtractedLink> extracted) =>
        links.Insert(extracted.Select(l => Link.Create(noteId, l.TargetSlug, l.Raw, l.Display, l.Line, l.Col)));

    public void WriteCustomFrontMatter(string noteId, IReadOnlyList<KeyValuePair<string, string>> entries) =>
        cfe.InsertRange(entries.Select(e => new CustomFrontMatter(noteId, e.Key, e.Value)));

    public void DeleteMetadata(string noteId)
    {
        links.Delete(noteId);
        tags.DeleteNoteTags(noteId);
        cfe.Delete(noteId);
    }

    public void DeleteNote(string noteId)
    {
        DeleteMetadata(noteId);
        notes.Delete(noteId);
    }
}

using Dendrite.Application.Abstractions;
using Dendrite.Domain.FileHandling;
using Dendrite.Domain.Models;

namespace Dendrite.Application.Services;

public sealed class IndexRebuilder(
    IUnitOfWork uow,
    INoteRepository notes,
    ILinkRepository links,
    ITagRepository tags,
    IIndexRepository index,
    ICfeRepository cfe,
    IVaultFiles files,
    VaultContext vault)
{
    public void Rebuild()
    {
        var root = vault.Config.VaultPath;
        var parsed = files.ListNotePaths()
            .Select(path => NoteFileParser.Parse(root, path, files.ReadText(path)))
            .ToList();

        using var tx = uow.Begin();

        index.Wipe();

        foreach (var file in parsed)
        {
            var note = Note.Create(file.Path, file.Title, file.Slug);
            notes.Insert(note);
            tags.Insert(file.FrontMatter.Tags.Select(t => new Tag(t)));
            tags.InsertNoteTags(file.FrontMatter.Tags.Select(t => new NoteTag(note.Id, t)));
            links.Insert(file.Links.Select(l => Link.Create(note.Id, l.TargetSlug, l.Raw, l.Display, l.Line, l.Col)));
            cfe.InsertRange(file.FrontMatter.Custom.Select(e => new CustomFrontMatter(note.Id, e.Key, e.Value)));
        }

        tx.Commit();
    }
}

using Dendrite.Daemon.FileHandling;
using Dendrite.Daemon.Models;
using Dendrite.Daemon.Persistence;
using Dendrite.Daemon.Persistence.Repositories;

namespace Dendrite.Daemon.Services;

public sealed class IndexRebuilder(
    UnitOfWork uow,
    NoteRepository notes,
    LinkRepository links,
    TagRepository tags,
    IndexRepository index,
    CfeRepository cfe,
    VaultStore vaultStore)
{
    public void Rebuild(string vaultRoot)
    {
        var files = ReadFiles(vaultRoot);

        using var tx = uow.Begin();

        index.Wipe();

        foreach (var file in files)
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

    private List<NoteFile> ReadFiles(string vaultRoot)
    {
        var ignored = vaultStore.Config.EffectiveExcludes;
        var files = new List<NoteFile>();
        Walk(vaultRoot);
        return files;

        void Walk(string dir)
        {
            foreach (var entry in new DirectoryInfo(dir).EnumerateFileSystemInfos().OrderBy(e => e.Name, StringComparer.Ordinal))
            {
                if (entry is DirectoryInfo)
                {
                    if (entry.LinkTarget is not null)
                        continue;

                    if (ignored.Contains(entry.Name))
                    {
                        Log.Debug($"Skipping directory during index rebuild: {entry.FullName}");
                        continue;
                    }
                    Walk(entry.FullName);
                }
                else if (entry.Extension == ".md")
                {
                    files.Add(NoteFileReader.Read(vaultRoot, entry.FullName));
                }
            }
        }
    }
}

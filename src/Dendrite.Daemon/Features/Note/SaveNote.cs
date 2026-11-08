using Dendrite.Daemon.FileHandling;
using Dendrite.Daemon.Persistence;
using Dendrite.Daemon.Persistence.Repositories;
using Dendrite.Daemon.Services;
using Mediator;

namespace Dendrite.Daemon.Features.Note;

public sealed record SaveNote(string Path) : ICommand<Unit>;

public sealed class SaveNoteHandler(
    UnitOfWork uow,
    NoteRepository notes,
    NoteMetadataService metadata,
    VaultStore vaultStore) : ICommandHandler<SaveNote, Unit>
{
    public ValueTask<Unit> Handle(SaveNote command, CancellationToken cancellationToken)
    {
        var file = NoteFileReader.Read(vaultStore.Config.VaultPath, command.Path);

        using var tx = uow.Begin();

        var note = notes.GetBySlug(file.Slug);
        if (note is null)
        {
            note = Models.Note.Create(file.Path, file.Title, file.Slug);
            notes.Insert(note);
        }
        else
        {
            metadata.DeleteMetadata(note.Id);
            notes.Update(note.Id, file.Path, file.Title, file.Slug);
        }

        metadata.WriteTags(note.Id, file.FrontMatter.Tags);
        metadata.WriteLinks(note.Id, file.Links);
        metadata.WriteCustomFrontMatter(note.Id, file.FrontMatter.Custom);

        tx.Commit();
        return ValueTask.FromResult(Unit.Value);
    }
}

using Dendrite.Application.Abstractions;
using Dendrite.Application.Services;
using Dendrite.Domain.FileHandling;
using Dendrite.Domain.Models;
using Mediator;

namespace Dendrite.Application.Features.Notes;

public sealed record SaveNote(string Path) : ICommand<Unit>;

public sealed class SaveNoteHandler(
    IUnitOfWork uow,
    INoteRepository notes,
    NoteMetadataService metadata,
    IVaultFiles files,
    VaultContext vault) : ICommandHandler<SaveNote, Unit>
{
    public ValueTask<Unit> Handle(SaveNote command, CancellationToken cancellationToken)
    {
        var file = NoteFileParser.Parse(vault.Config.VaultPath, command.Path, files.ReadText(command.Path));

        using var tx = uow.Begin();

        var note = notes.GetBySlug(file.Slug);
        if (note is null)
        {
            note = Note.Create(file.Path, file.Title, file.Slug);
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

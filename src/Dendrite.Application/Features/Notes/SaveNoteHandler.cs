using Dendrite.Application.Abstractions;
using Dendrite.Domain.FileHandling;
using Dendrite.Domain.Models;
using Mediator;

namespace Dendrite.Application.Features.Notes;

public sealed class SaveNoteHandler : ICommandHandler<SaveNote, Unit>
{
    private readonly IUnitOfWork _uow;
    private readonly INoteRepository _notes;
    private readonly IVaultFiles _files;
    private readonly VaultContext _vault;

    public SaveNoteHandler(IUnitOfWork uow, INoteRepository notes, IVaultFiles files, VaultContext vault)
    {
        _uow = uow;
        _notes = notes;
        _files = files;
        _vault = vault;
    }

    public ValueTask<Unit> Handle(SaveNote command, CancellationToken cancellationToken)
    {
        var file = NoteFileParser.Parse(_vault.Config.VaultPath, command.Path, _files.ReadText(command.Path));

        using var tx = _uow.Begin();

        var note = _notes.GetBySlug(file.Slug);
        if (note is null)
        {
            _notes.Add(Note.Create(file));
        }
        else
        {
            note.Reindex(file);
            _notes.Save(note);
        }

        tx.Commit();
        return ValueTask.FromResult(Unit.Value);
    }
}

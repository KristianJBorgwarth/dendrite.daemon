using Dendrite.Daemon.Persistence;
using Dendrite.Daemon.Persistence.Repositories;
using Dendrite.Daemon.Services;
using Mediator;

namespace Dendrite.Daemon.Features.Note;

public sealed record DeleteNote(string Path) : ICommand<Unit>;

public sealed class DeleteNoteHandler(UnitOfWork uow, NoteRepository notes, NoteMetadataService metadata)
    : ICommandHandler<DeleteNote, Unit>
{
    public ValueTask<Unit> Handle(DeleteNote command, CancellationToken cancellationToken)
    {
        using var tx = uow.Begin();

        var note = notes.GetByPath(command.Path)
                   ?? throw new InvalidOperationException($"note not found: {command.Path}");

        metadata.DeleteNote(note.Id);

        // A failed commit throws, so the error reaches the client instead of being swallowed.
        tx.Commit();
        return ValueTask.FromResult(Unit.Value);
    }
}

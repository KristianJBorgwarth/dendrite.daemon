using Dendrite.Application.Abstractions;
using Dendrite.Domain.Models;
using Mediator;

namespace Dendrite.Application.Features.Notes;

public sealed class DeleteNoteHandler : ICommandHandler<DeleteNote, Unit>
{
    private readonly IUnitOfWork _uow;
    private readonly INoteRepository _notes;

    public DeleteNoteHandler(IUnitOfWork uow, INoteRepository notes)
    {
        _uow = uow;
        _notes = notes;
    }

    public ValueTask<Unit> Handle(DeleteNote command, CancellationToken cancellationToken)
    {
        using var tx = _uow.Begin();

        var note = _notes.GetByPath(command.Path)
                   ?? throw new InvalidOperationException($"note not found: {command.Path}");

        _notes.Remove(note);

        tx.Commit();
        return ValueTask.FromResult(Unit.Value);
    }
}

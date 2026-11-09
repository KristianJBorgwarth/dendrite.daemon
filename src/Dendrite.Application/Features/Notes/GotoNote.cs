using Dendrite.Application.Abstractions;
using Dendrite.Domain.FileHandling;
using Mediator;

namespace Dendrite.Application.Features.Notes;

public sealed record GotoNote(string Link) : IQuery<GotoNoteResult?>;

public sealed class GotoNoteHandler(INoteRepository notes) : IQueryHandler<GotoNote, GotoNoteResult?>
{
    public ValueTask<GotoNoteResult?> Handle(GotoNote query, CancellationToken cancellationToken)
    {
        var parsed = LinkParser.Parse(query.Link);

        GotoNoteResult? result = parsed.Kind switch
        {
            LinkKind.Note => notes.GetBySlug(parsed.Target) is { } note ? new GotoNoteResult(note.Path, "note") : null,
            LinkKind.Url => new GotoNoteResult(parsed.Target, "url"),
            _ => throw new InvalidOperationException("unsupported link type"),
        };

        return ValueTask.FromResult(result);
    }
}

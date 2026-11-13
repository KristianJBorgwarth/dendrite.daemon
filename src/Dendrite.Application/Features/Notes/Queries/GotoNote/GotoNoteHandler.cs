using Dendrite.Application.Abstractions;
using Dendrite.Domain.FileHandling;
using Mediator;

namespace Dendrite.Application.Features.Notes;

public sealed class GotoNoteHandler : IQueryHandler<GotoNote, GotoNoteResult?>
{
    private readonly INoteQueries _queries;

    public GotoNoteHandler(INoteQueries queries)
    {
        _queries = queries;
    }

    public ValueTask<GotoNoteResult?> Handle(GotoNote query, CancellationToken cancellationToken)
    {
        var parsed = LinkParser.Parse(query.Link);

        GotoNoteResult? result = parsed.Kind switch
        {
            LinkKind.Note => _queries.FindPathBySlug(parsed.Target) is { } path ? new GotoNoteResult(path, "note") : null,
            LinkKind.Url => new GotoNoteResult(parsed.Target, "url"),
            _ => throw new InvalidOperationException("unsupported link type"),
        };

        return ValueTask.FromResult(result);
    }
}

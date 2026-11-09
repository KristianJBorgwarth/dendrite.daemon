using Dendrite.Application.Abstractions;
using Mediator;

namespace Dendrite.Application.Features.Notes;

public sealed class GetNotesByTagHandler : IQueryHandler<GetNotesByTag, List<NoteDto>>
{
    private readonly INoteQueries _queries;

    public GetNotesByTagHandler(INoteQueries queries)
    {
        _queries = queries;
    }

    public ValueTask<List<NoteDto>> Handle(GetNotesByTag query, CancellationToken cancellationToken) =>
        ValueTask.FromResult(_queries.GetNotesByTag(query.Tag));
}

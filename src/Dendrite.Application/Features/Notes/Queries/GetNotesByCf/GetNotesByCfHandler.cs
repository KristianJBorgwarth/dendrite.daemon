using Dendrite.Application.Abstractions;
using Mediator;

namespace Dendrite.Application.Features.Notes;

public sealed class GetNotesByCfHandler : IQueryHandler<GetNotesByCf, List<NoteDto>>
{
    private readonly INoteQueries _queries;

    public GetNotesByCfHandler(INoteQueries queries)
    {
        _queries = queries;
    }

    public ValueTask<List<NoteDto>> Handle(GetNotesByCf query, CancellationToken cancellationToken) =>
        ValueTask.FromResult(_queries.GetNotesByCustomFrontMatter(query.Key, query.Value));
}

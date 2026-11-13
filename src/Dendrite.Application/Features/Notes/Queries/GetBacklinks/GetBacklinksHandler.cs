using Dendrite.Application.Abstractions;
using Mediator;

namespace Dendrite.Application.Features.Notes;

public sealed class GetBacklinksHandler : IQueryHandler<GetBacklinks, List<BacklinkDto>>
{
    private readonly INoteQueries _queries;

    public GetBacklinksHandler(INoteQueries queries)
    {
        _queries = queries;
    }

    public ValueTask<List<BacklinkDto>> Handle(GetBacklinks query, CancellationToken cancellationToken) =>
        ValueTask.FromResult(_queries.GetBacklinks(query.Path));
}

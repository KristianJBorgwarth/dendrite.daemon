using Dendrite.Application.Abstractions;
using Mediator;

namespace Dendrite.Application.Features.Completion;

public sealed class CompleteTagHandler : IQueryHandler<CompleteTag, List<string>>
{
    private readonly INoteQueries _queries;

    public CompleteTagHandler(INoteQueries queries)
    {
        _queries = queries;
    }

    public ValueTask<List<string>> Handle(CompleteTag query, CancellationToken cancellationToken)
    {
        var names = _queries.SearchTags(query.Query);
        Log.Debug($"Got tags for '{query.Query}': {names.Count}");
        return ValueTask.FromResult(names);
    }
}

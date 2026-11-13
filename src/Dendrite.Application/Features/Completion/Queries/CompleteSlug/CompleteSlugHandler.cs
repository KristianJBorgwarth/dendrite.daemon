using Dendrite.Application.Abstractions;
using Mediator;

namespace Dendrite.Application.Features.Completion;

public sealed class CompleteSlugHandler(INoteQueries queries) : IQueryHandler<CompleteSlug, List<string>>
{
    private readonly INoteQueries _queries = queries;

    public ValueTask<List<string>> Handle(CompleteSlug query, CancellationToken cancellationToken) =>
        ValueTask.FromResult(_queries.SearchNotes(query.Query).Select(n => n.Slug).ToList());
}

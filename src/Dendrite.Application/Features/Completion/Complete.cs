using Dendrite.Application.Abstractions;
using Mediator;

namespace Dendrite.Application.Features.Completion;

public sealed record CompleteTag(string Query) : IQuery<List<string>>;

public sealed class CompleteTagHandler(ITagRepository tags) : IQueryHandler<CompleteTag, List<string>>
{
    public ValueTask<List<string>> Handle(CompleteTag query, CancellationToken cancellationToken)
    {
        var names = tags.Search(query.Query).Select(t => t.Name).ToList();
        Log.Debug($"Got tags for '{query.Query}': {names.Count}");
        return ValueTask.FromResult(names);
    }
}

public sealed record CompleteSlug(string Query) : IQuery<List<string>>;

public sealed class CompleteSlugHandler(INoteRepository notes) : IQueryHandler<CompleteSlug, List<string>>
{
    public ValueTask<List<string>> Handle(CompleteSlug query, CancellationToken cancellationToken) =>
        ValueTask.FromResult(notes.Search(query.Query).Select(n => n.Slug).ToList());
}

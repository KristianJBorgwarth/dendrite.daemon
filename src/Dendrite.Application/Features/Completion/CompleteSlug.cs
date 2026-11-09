using Mediator;

namespace Dendrite.Application.Features.Completion;

public sealed record CompleteSlug : IQuery<List<string>>
{
    public string Query { get; }

    public CompleteSlug(string query)
    {
        Query = query;
    }
}

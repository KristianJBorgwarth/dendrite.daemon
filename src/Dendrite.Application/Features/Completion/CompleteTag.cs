using Mediator;

namespace Dendrite.Application.Features.Completion;

public sealed record CompleteTag : IQuery<List<string>>
{
    public string Query { get; }

    public CompleteTag(string query)
    {
        Query = query;
    }
}

using Dendrite.Application.Vaults;
using Mediator;

namespace Dendrite.Application.Features.Completion;

public sealed record CompleteTag : IQuery<List<string>>, IVaultRequest
{
    public string Query { get; }
    public string? Vault { get; }

    public CompleteTag(string query, string? vault = null)
    {
        Query = query;
        Vault = vault;
    }
}

using Dendrite.Application.Vaults;
using Mediator;

namespace Dendrite.Application.Features.Vault;

public sealed record RebuildIndex : ICommand<Unit>, IVaultRequest
{
    public string? Vault { get; }

    public RebuildIndex(string? vault = null)
    {
        Vault = vault;
    }
}

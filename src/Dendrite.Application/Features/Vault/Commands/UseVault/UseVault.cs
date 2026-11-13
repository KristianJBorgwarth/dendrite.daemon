using Mediator;

namespace Dendrite.Application.Features.Vault;

public sealed record UseVault : ICommand<Unit>
{
    public string Vault { get; }

    public UseVault(string vault)
    {
        Vault = vault;
    }
}

using Mediator;

namespace Dendrite.Application.Features.Vault;

public sealed record RemoveVault : ICommand<Unit>
{
    public string Vault { get; }

    public RemoveVault(string vault)
    {
        Vault = vault;
    }
}

using Mediator;

namespace Dendrite.Application.Features.Vault;

/// <summary>Makes a vault the one requests fall back to when they name none.</summary>
public sealed record UseVault : ICommand<Unit>
{
    public string Vault { get; }

    public UseVault(string vault)
    {
        Vault = vault;
    }
}

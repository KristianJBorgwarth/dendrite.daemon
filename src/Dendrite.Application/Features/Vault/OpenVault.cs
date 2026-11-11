using Mediator;

namespace Dendrite.Application.Features.Vault;

/// <summary>Registers an existing folder as a vault from its <c>.dendrite/vault.json</c>, or re-points the vault of that name to it.</summary>
public sealed record OpenVault : ICommand<Unit>
{
    public string Path { get; }

    public OpenVault(string path)
    {
        Path = path;
    }
}

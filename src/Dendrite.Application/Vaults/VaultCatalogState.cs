using Dendrite.Domain.Vault;

namespace Dendrite.Application.Vaults;

/// <summary>What the daemon remembers about vaults between runs.</summary>
public sealed class VaultCatalogState
{
    public string? Active { get; init; }

    public List<VaultConfig> Vaults { get; init; } = [];
}

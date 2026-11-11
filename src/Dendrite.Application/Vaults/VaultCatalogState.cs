using Dendrite.Domain.Vault;

namespace Dendrite.Application.Vaults;

public sealed class VaultCatalogState
{
    public string? Active { get; init; }

    public List<VaultConfig> Vaults { get; init; } = [];
}

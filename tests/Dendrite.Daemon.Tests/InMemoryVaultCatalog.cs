using Dendrite.Application.Abstractions;
using Dendrite.Application.Vaults;

namespace Dendrite.Daemon.Tests;

public sealed class InMemoryVaultCatalog : IVaultCatalog
{
    private VaultCatalogState _state = new();

    public VaultCatalogState Load() => _state;

    public void Save(VaultCatalogState state) => _state = state;
}

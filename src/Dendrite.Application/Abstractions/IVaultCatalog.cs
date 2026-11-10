using Dendrite.Application.Vaults;

namespace Dendrite.Application.Abstractions;

/// <summary>Remembers registered vaults and the active one across daemon restarts.</summary>
public interface IVaultCatalog
{
    VaultCatalogState Load();

    void Save(VaultCatalogState state);
}

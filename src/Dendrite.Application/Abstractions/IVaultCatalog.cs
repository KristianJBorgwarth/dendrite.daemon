using Dendrite.Application.Vaults;

namespace Dendrite.Application.Abstractions;

public interface IVaultCatalog
{
    VaultCatalogState Load();

    void Save(VaultCatalogState state);
}

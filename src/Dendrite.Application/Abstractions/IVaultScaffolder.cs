using Dendrite.Domain.Vault;

namespace Dendrite.Application.Abstractions;

public interface IVaultScaffolder
{
    bool IsEmptyOrMissing(string vaultPath);

    bool LooksLikeVault(VaultConfig config);

    void Scaffold(VaultConfig config);
}

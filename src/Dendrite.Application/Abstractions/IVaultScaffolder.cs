using Dendrite.Domain.Vault;

namespace Dendrite.Application.Abstractions;

/// <summary>Lays out a new vault on disk.</summary>
public interface IVaultScaffolder
{
    /// <summary>True when the directory is missing or has no entries.</summary>
    bool IsEmptyOrMissing(string vaultPath);

    /// <summary>True when the directory already contains the vault's templates directory.</summary>
    bool LooksLikeVault(VaultConfig config);

    /// <summary>Creates the vault and templates directories and any missing default templates. Never overwrites files.</summary>
    void Scaffold(VaultConfig config);
}

using Dendrite.Domain.Vault;

namespace Dendrite.Application.Abstractions;

/// <summary>The vault's own settings, kept inside the vault folder so they travel with it.</summary>
public interface IVaultManifest
{
    /// <summary>Reads the settings stored in the folder, or null when it has none. The path is the folder itself.</summary>
    VaultConfig? Read(string vaultPath);

    void Write(VaultConfig config);
}

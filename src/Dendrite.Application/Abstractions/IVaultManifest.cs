using Dendrite.Domain.Vault;

namespace Dendrite.Application.Abstractions;

public interface IVaultManifest
{
    VaultConfig? Read(string vaultPath);

    void Write(VaultConfig config);
}

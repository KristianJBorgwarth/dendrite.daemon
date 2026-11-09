using Dendrite.Domain.Vault;

namespace Dendrite.Application;

public sealed class VaultContext
{
    private VaultConfig? _config;

    public VaultConfig Config =>
        _config ?? throw new InvalidOperationException("vault is not initialized — call vault/init first");

    public void Initialize(VaultConfig config) => _config = config;
}

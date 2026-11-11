using Dendrite.Domain.Vault;

namespace Dendrite.Application.Vaults;

public sealed class VaultContext
{
    private VaultConfig? _config;

    public VaultConfig Config =>
        _config ?? throw new InvalidOperationException("no vault is initialized — call vault/init first");

    public void Select(VaultConfig config) => _config = config;
}

using Mediator;

namespace Dendrite.Application.Features.Vault;

public sealed record InitializeVault : ICommand<Unit>
{
    public VaultConfigDto Config { get; }

    public InitializeVault(VaultConfigDto config)
    {
        Config = config;
    }
}

using Mediator;

namespace Dendrite.Application.Features.Vault;

public sealed record CreateVault : ICommand<Unit>
{
    public VaultConfigDto Config { get; }

    public CreateVault(VaultConfigDto config)
    {
        Config = config;
    }
}

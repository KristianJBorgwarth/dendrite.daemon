using Dendrite.Application.Abstractions;
using Dendrite.Application.Vaults;
using Mediator;

namespace Dendrite.Application.Features.Vault;

public sealed class UseVaultHandler : ICommandHandler<UseVault, Unit>
{
    private readonly VaultRegistry _registry;
    private readonly IVaultCatalog _catalog;

    public UseVaultHandler(VaultRegistry registry, IVaultCatalog catalog)
    {
        _registry = registry;
        _catalog = catalog;
    }

    public ValueTask<Unit> Handle(UseVault command, CancellationToken cancellationToken)
    {
        _registry.Use(command.Vault);
        _catalog.Save(_registry.ToState());
        return ValueTask.FromResult(Unit.Value);
    }
}

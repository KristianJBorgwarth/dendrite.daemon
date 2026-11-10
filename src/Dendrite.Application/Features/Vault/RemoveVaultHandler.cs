using Dendrite.Application.Abstractions;
using Dendrite.Application.Vaults;
using Mediator;

namespace Dendrite.Application.Features.Vault;

/// <summary>Unregisters a vault and closes its index. The index file and the notes stay on disk.</summary>
public sealed class RemoveVaultHandler : ICommandHandler<RemoveVault, Unit>
{
    private readonly VaultRegistry _registry;
    private readonly IIndexDatabase _indexes;

    public RemoveVaultHandler(VaultRegistry registry, IIndexDatabase indexes)
    {
        _registry = registry;
        _indexes = indexes;
    }

    public ValueTask<Unit> Handle(RemoveVault command, CancellationToken cancellationToken)
    {
        if (!_registry.Remove(command.Vault))
            throw new InvalidOperationException($"unknown vault '{command.Vault}'");

        _indexes.Close(command.Vault);
        return ValueTask.FromResult(Unit.Value);
    }
}

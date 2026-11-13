using Dendrite.Application.Abstractions;
using Dendrite.Application.Services;
using Dendrite.Application.Vaults;
using Mediator;

namespace Dendrite.Application.Features.Vault;

public sealed class InitializeVaultHandler : ICommandHandler<InitializeVault, Unit>
{
    private readonly VaultInitializer _initializer;

    public InitializeVaultHandler(VaultInitializer initializer)
    {
        _initializer = initializer;
    }

    public ValueTask<Unit> Handle(InitializeVault command, CancellationToken cancellationToken)
    {
        _initializer.Initialize(VaultInitializer.ToConfig(command.Config));
        return ValueTask.FromResult(Unit.Value);
    }
}

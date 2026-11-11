using Dendrite.Application.Abstractions;
using Dendrite.Application.Services;
using Dendrite.Application.Vaults;
using Mediator;

namespace Dendrite.Application.Features.Vault;

/// <summary>Scaffolds a new vault, registers it and builds its first index.</summary>
public sealed class CreateVaultHandler : ICommandHandler<CreateVault, Unit>
{
    private readonly VaultInitializer _initializer;
    private readonly IVaultScaffolder _scaffolder;
    private readonly VaultRegistry _registry;

    public CreateVaultHandler(VaultInitializer initializer, IVaultScaffolder scaffolder, VaultRegistry registry)
    {
        _initializer = initializer;
        _scaffolder = scaffolder;
        _registry = registry;
    }

    public ValueTask<Unit> Handle(CreateVault command, CancellationToken cancellationToken)
    {
        var config = VaultInitializer.ToConfig(command.Config.WithDefaults());
        VaultRegistry.ValidateName(config.VaultName);

        var alreadyRegistered = _registry.CheckAvailability(config);
        if (!alreadyRegistered && !_scaffolder.IsEmptyOrMissing(config.VaultPath) && !_scaffolder.LooksLikeVault(config))
            throw new InvalidOperationException(
                $"'{config.VaultPath}' is not empty and is not a vault — refusing to create one there");

        _scaffolder.Scaffold(config);
        _initializer.Initialize(config);
        return ValueTask.FromResult(Unit.Value);
    }
}

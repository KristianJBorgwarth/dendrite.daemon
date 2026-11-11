using Dendrite.Application.Abstractions;
using Dendrite.Application.Services;
using Dendrite.Application.Vaults;
using Mediator;

namespace Dendrite.Application.Features.Vault;

public sealed class CreateVaultHandler : ICommandHandler<CreateVault, Unit>
{
    private readonly VaultInitializer _initializer;
    private readonly IVaultScaffolder _scaffolder;
    private readonly VaultRegistry _registry;
    private readonly IVaultManifest _manifest;

    public CreateVaultHandler(
        VaultInitializer initializer,
        IVaultScaffolder scaffolder,
        VaultRegistry registry,
        IVaultManifest manifest)
    {
        _initializer = initializer;
        _scaffolder = scaffolder;
        _registry = registry;
        _manifest = manifest;
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
        _manifest.Write(config);
        _initializer.Initialize(config, activate: true);
        return ValueTask.FromResult(Unit.Value);
    }
}

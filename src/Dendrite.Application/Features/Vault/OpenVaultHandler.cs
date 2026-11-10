using Dendrite.Application.Abstractions;
using Dendrite.Application.Services;
using Dendrite.Application.Vaults;
using Mediator;

namespace Dendrite.Application.Features.Vault;

public sealed class OpenVaultHandler : ICommandHandler<OpenVault, Unit>
{
    private readonly IVaultManifest _manifest;
    private readonly VaultRegistry _registry;
    private readonly VaultInitializer _initializer;
    private readonly VaultRelocator _relocator;
    private readonly IVaultCatalog _catalog;

    public OpenVaultHandler(
        IVaultManifest manifest,
        VaultRegistry registry,
        VaultInitializer initializer,
        VaultRelocator relocator,
        IVaultCatalog catalog)
    {
        _manifest = manifest;
        _registry = registry;
        _initializer = initializer;
        _relocator = relocator;
        _catalog = catalog;
    }

    public ValueTask<Unit> Handle(OpenVault command, CancellationToken cancellationToken)
    {
        var path = VaultRelocator.Normalize(command.Path);
        if (!Directory.Exists(path))
            throw new InvalidOperationException($"'{path}' is not a directory");

        var config = _manifest.Read(path)
            ?? throw new InvalidOperationException(
                $"'{path}' has no .dendrite/vault.json — create it with vault/create or register it with vault/init");

        if (_registry.Contains(config.VaultName))
        {
            _relocator.Relocate(config.VaultName, path);
            _registry.Use(config.VaultName);
            _catalog.Save(_registry.ToState());
        }
        else
        {
            _initializer.Initialize(config, activate: true);
        }

        return ValueTask.FromResult(Unit.Value);
    }
}

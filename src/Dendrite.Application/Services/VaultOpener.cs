using Dendrite.Application.Abstractions;
using Dendrite.Application.Vaults;

namespace Dendrite.Application.Services;

/// <summary>Registers a folder as a vault from its own manifest, or re-points the vault of that name to it.</summary>
public sealed class VaultOpener
{
    private readonly IVaultManifest _manifest;
    private readonly VaultRegistry _registry;
    private readonly VaultInitializer _initializer;
    private readonly IIndexDatabase _indexes;
    private readonly IVaultCatalog _catalog;

    public VaultOpener(
        IVaultManifest manifest,
        VaultRegistry registry,
        VaultInitializer initializer,
        IIndexDatabase indexes,
        IVaultCatalog catalog)
    {
        _manifest = manifest;
        _registry = registry;
        _initializer = initializer;
        _indexes = indexes;
        _catalog = catalog;
    }

    public void Open(string folder, bool activate)
    {
        var path = Path.GetFullPath(folder).TrimEnd('/');
        var config = _manifest.Read(path)
            ?? throw new InvalidOperationException($"'{path}' has no .dendrite/vault.json");

        if (!_registry.Contains(config.VaultName))
        {
            _initializer.Initialize(config, activate);
            return;
        }

        var old = _registry.Resolve(config.VaultName, null);
        if (old.VaultPath != path)
        {
            _registry.Replace(config);
            _indexes.Close(config.VaultName);
            _indexes.Open(config.VaultName, path, old.VaultPath);
        }

        if (activate)
            _registry.Use(config.VaultName);

        _catalog.Save(_registry.ToState());
    }
}

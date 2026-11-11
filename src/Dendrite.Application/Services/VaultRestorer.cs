using Dendrite.Application.Abstractions;
using Dendrite.Application.Vaults;

namespace Dendrite.Application.Services;

public sealed class VaultRestorer
{
    private readonly IVaultCatalog _catalog;
    private readonly VaultRegistry _registry;
    private readonly IIndexDatabase _indexes;
    private readonly IVaultManifest _manifest;

    public VaultRestorer(IVaultCatalog catalog, VaultRegistry registry, IIndexDatabase indexes, IVaultManifest manifest)
    {
        _catalog = catalog;
        _registry = registry;
        _indexes = indexes;
        _manifest = manifest;
    }

    public void Restore()
    {
        var state = _catalog.Load();

        foreach (var remembered in state.Vaults)
        {
            var config = remembered;
            try
            {
                if (!Directory.Exists(config.VaultPath))
                {
                    Log.Warn($"vault '{config.VaultName}' points at {config.VaultPath}, which does not exist — open its new folder to repair it");
                    _registry.Register(config);
                    continue;
                }

                config = _manifest.Read(config.VaultPath) ?? config;
                _registry.Register(config);
                _indexes.Open(config.VaultName, config.VaultPath);
            }
            catch (Exception ex)
            {
                _registry.Remove(config.VaultName);
                Log.Warn($"could not restore vault '{config.VaultName}': {ex.Message}");
            }
        }

        if (state.Active is not null && _registry.Contains(state.Active))
            _registry.Use(state.Active);
    }
}

using Dendrite.Application.Abstractions;
using Dendrite.Application.Vaults;

namespace Dendrite.Application.Services;

/// <summary>Re-registers the remembered vaults and the active one when the daemon starts.</summary>
public sealed class VaultRestorer
{
    private readonly IVaultCatalog _catalog;
    private readonly VaultRegistry _registry;
    private readonly IIndexDatabase _indexes;

    public VaultRestorer(IVaultCatalog catalog, VaultRegistry registry, IIndexDatabase indexes)
    {
        _catalog = catalog;
        _registry = registry;
        _indexes = indexes;
    }

    public void Restore()
    {
        var state = _catalog.Load();

        foreach (var config in state.Vaults)
        {
            try
            {
                if (!Directory.Exists(config.VaultPath))
                {
                    Log.Warn($"skipping vault '{config.VaultName}': {config.VaultPath} does not exist");
                    continue;
                }

                _registry.Register(config);
                _indexes.Open(config.VaultName);
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

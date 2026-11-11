using Dendrite.Application.Abstractions;
using Dendrite.Application.Features.Vault;
using Dendrite.Application.Vaults;
using Dendrite.Domain.Vault;

namespace Dendrite.Application.Services;

/// <summary>Registers a vault, opens its index and builds the index when it is empty.</summary>
public sealed class VaultInitializer
{
    private readonly VaultRegistry _registry;
    private readonly VaultContext _vault;
    private readonly IIndexDatabase _indexDatabase;
    private readonly INoteQueries _queries;
    private readonly IndexRebuilder _rebuilder;
    private readonly IVaultCatalog _catalog;
    private readonly IVaultManifest _manifest;

    public VaultInitializer(
        VaultRegistry registry,
        VaultContext vault,
        IIndexDatabase indexDatabase,
        INoteQueries queries,
        IndexRebuilder rebuilder,
        IVaultCatalog catalog,
        IVaultManifest manifest)
    {
        _registry = registry;
        _vault = vault;
        _indexDatabase = indexDatabase;
        _queries = queries;
        _rebuilder = rebuilder;
        _catalog = catalog;
        _manifest = manifest;
    }

    public static VaultConfig ToConfig(VaultConfigDto c)
    {
        if (string.IsNullOrEmpty(c.VaultPath))
            throw new InvalidOperationException("vault_path is empty — check your dendrite.nvim config");

        var daily = c.DailyNotes ?? new DailyNotesConfig();
        var vaultPath = c.VaultPath.Length > 1 ? c.VaultPath.TrimEnd('/', '\\') : c.VaultPath;

        return new VaultConfig(
            c.VaultName,
            vaultPath,
            c.TemplatesDir,
            c.ExcludeIndexFiles ?? [],
            c.OverrideDefaultIgnores,
            new DailyNotes(daily.Dir, daily.FilenameFormat, daily.TemplateName));
    }

    public void Initialize(VaultConfig config, bool activate = false)
    {
        var added = _registry.Register(config);
        try
        {
            _vault.Select(config);
            _indexDatabase.Open(config.VaultName, config.VaultPath);
            if (_manifest.Read(config.VaultPath) is null)
                _manifest.Write(config);
        }
        catch
        {
            if (added)
                _registry.Remove(config.VaultName);
            throw;
        }

        if (activate)
            _registry.Use(config.VaultName);

        _catalog.Save(_registry.ToState());

        if (_queries.CountNotes() > 0)
        {
            Log.Debug("Vault already initialized, skipping index rebuild");
            return;
        }

        _rebuilder.Rebuild();
    }
}

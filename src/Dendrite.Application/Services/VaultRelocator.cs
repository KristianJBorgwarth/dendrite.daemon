using Dendrite.Application.Abstractions;
using Dendrite.Application.Vaults;

namespace Dendrite.Application.Services;

/// <summary>Moves a registered vault to a new folder without losing its index.</summary>
public sealed class VaultRelocator
{
    private readonly VaultRegistry _registry;
    private readonly IIndexDatabase _indexes;
    private readonly IVaultCatalog _catalog;

    public VaultRelocator(VaultRegistry registry, IIndexDatabase indexes, IVaultCatalog catalog)
    {
        _registry = registry;
        _indexes = indexes;
        _catalog = catalog;
    }

    public static string Normalize(string path)
    {
        var full = Path.GetFullPath(path);
        return full.Length > 1 ? full.TrimEnd('/', '\\') : full;
    }

    public void Relocate(string name, string newPath)
    {
        var old = _registry.Resolve(name, null);
        var target = Normalize(newPath);

        if (old.VaultPath == target)
            return;

        if (!Directory.Exists(target))
            throw new InvalidOperationException($"'{target}' is not a directory");

        var owner = _registry.FindByPath(target);
        if (owner is not null)
            throw new InvalidOperationException($"'{target}' already belongs to vault '{owner.VaultName}'");

        _registry.Replace(old.WithPath(target));
        _indexes.RelocateNotes(name, old.VaultPath, target);
        _catalog.Save(_registry.ToState());
    }
}

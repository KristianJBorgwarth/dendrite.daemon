using Dendrite.Domain.Vault;

namespace Dendrite.Application.Vaults;

/// <summary>The vaults known to the daemon, keyed by vault name (the vault's identity).</summary>
public sealed class VaultRegistry
{
    private readonly Dictionary<string, VaultConfig> _vaults = new(StringComparer.Ordinal);

    public IReadOnlyCollection<VaultConfig> All => _vaults.Values;

    /// <summary>The vault requests fall back to when they name none and no path places them in one.</summary>
    public string? Active { get; private set; }

    public bool Contains(string name) => _vaults.ContainsKey(name);

    public void Use(string name)
    {
        if (!_vaults.ContainsKey(name))
            throw new InvalidOperationException($"unknown vault '{name}'");

        Active = name;
    }

    public VaultCatalogState ToState() => new() { Active = Active, Vaults = [.. _vaults.Values] };

    /// <summary>Adds or replaces a vault. Returns true when the vault was not registered before.</summary>
    public bool Register(VaultConfig config)
    {
        ValidateName(config.VaultName);

        if (_vaults.TryGetValue(config.VaultName, out var existing))
        {
            if (existing.VaultPath != config.VaultPath)
                throw new InvalidOperationException(
                    $"vault '{config.VaultName}' is already registered at '{existing.VaultPath}'");

            _vaults[config.VaultName] = config;
            return false;
        }

        _vaults[config.VaultName] = config;
        return true;
    }

    public bool Remove(string name)
    {
        if (Active == name)
            Active = null;

        return _vaults.Remove(name);
    }

    /// <summary>Points an already registered vault at a new folder.</summary>
    public void Replace(VaultConfig config)
    {
        if (!_vaults.ContainsKey(config.VaultName))
            throw new InvalidOperationException($"unknown vault '{config.VaultName}'");

        _vaults[config.VaultName] = config;
    }

    /// <summary>The vault registered at exactly this folder, if any.</summary>
    public VaultConfig? FindByPath(string vaultPath) =>
        _vaults.Values.FirstOrDefault(v => v.VaultPath == vaultPath);

    /// <summary>True when the vault is already registered at the same path; throws when the name belongs to another path.</summary>
    public bool CheckAvailability(VaultConfig config)
    {
        if (!_vaults.TryGetValue(config.VaultName, out var existing))
            return false;

        if (existing.VaultPath != config.VaultPath)
            throw new InvalidOperationException(
                $"vault '{config.VaultName}' is already registered at '{existing.VaultPath}'");

        return true;
    }

    /// <summary>True when the path lies inside a registered vault.</summary>
    public bool Owns(string path) => Owner(System.IO.Path.GetFullPath(path)) is not null;

    private VaultConfig? Owner(string fullPath) => _vaults.Values
        .Where(v => IsInside(v.VaultPath, fullPath))
        .OrderByDescending(v => v.VaultPath.Length)
        .FirstOrDefault();

    /// <summary>Picks the vault by explicit name, then by the vault root containing the path, then the active vault, then the only vault.</summary>
    public VaultConfig Resolve(string? vault, string? path)
    {
        if (!string.IsNullOrEmpty(vault))
        {
            return _vaults.TryGetValue(vault, out var named)
                ? named
                : throw new InvalidOperationException($"unknown vault '{vault}'");
        }

        if (!string.IsNullOrEmpty(path))
        {
            var owner = Owner(System.IO.Path.GetFullPath(path));
            if (owner is not null)
                return owner;
        }

        if (Active is not null && _vaults.TryGetValue(Active, out var active))
            return active;

        return _vaults.Count switch
        {
            0 => throw new InvalidOperationException("no vault is initialized — call vault/init first"),
            1 => _vaults.Values.First(),
            _ => throw new InvalidOperationException("several vaults are initialized — pass 'vault' to pick one"),
        };
    }

    private static bool IsInside(string root, string fullPath)
    {
        var normalized = System.IO.Path.GetFullPath(root).TrimEnd('/', '\\');
        return fullPath.StartsWith(normalized + System.IO.Path.DirectorySeparatorChar, StringComparison.Ordinal);
    }

    public static void ValidateName(string name)
    {
        if (string.IsNullOrWhiteSpace(name)
            || name.Contains('/')
            || name.Contains('\\')
            || name.Contains("..", StringComparison.Ordinal))
        {
            throw new InvalidOperationException($"invalid vault name '{name}'");
        }
    }
}

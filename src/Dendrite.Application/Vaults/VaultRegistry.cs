using Dendrite.Domain.Vault;

namespace Dendrite.Application.Vaults;

/// <summary>The vaults known to the daemon, keyed by vault name (the vault's identity).</summary>
public sealed class VaultRegistry
{
    private readonly Dictionary<string, VaultConfig> _vaults = new(StringComparer.Ordinal);

    public IReadOnlyCollection<VaultConfig> All => _vaults.Values;

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

    public void Remove(string name) => _vaults.Remove(name);

    /// <summary>Picks the vault by explicit name, then by the vault root containing the path, then the only vault.</summary>
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
            var full = System.IO.Path.GetFullPath(path);
            var owner = _vaults.Values
                .Where(v => IsInside(v.VaultPath, full))
                .OrderByDescending(v => v.VaultPath.Length)
                .FirstOrDefault();
            if (owner is not null)
                return owner;
        }

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

    private static void ValidateName(string name)
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

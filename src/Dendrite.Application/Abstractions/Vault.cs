using Dendrite.Domain.Vault;

namespace Dendrite.Application.Abstractions;

/// <summary>Reads the vault's files on disk. Resolution is relative to the initialized vault.</summary>
public interface IVaultFiles
{
    string ReadText(string path);

    string ReadTemplate(string templateName);

    /// <summary>Full paths of every indexable markdown file, honouring the vault's excludes.</summary>
    IReadOnlyList<string> ListNotePaths();
}

/// <summary>Opens (and migrates) the index database that belongs to a vault.</summary>
public interface IIndexDatabase
{
    void Open(string vaultName);
}

public sealed class VaultContext
{
    private VaultConfig? _config;

    public VaultConfig Config =>
        _config ?? throw new InvalidOperationException("vault is not initialized — call vault/init first");

    public void Initialize(VaultConfig config) => _config = config;
}

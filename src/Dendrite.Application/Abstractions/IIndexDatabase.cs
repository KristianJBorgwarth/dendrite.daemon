namespace Dendrite.Application.Abstractions;

/// <summary>Opens (and migrates) the index database that belongs to a vault.</summary>
public interface IIndexDatabase
{
    void Open(string vaultName);

    /// <summary>Closes the vault's index. The file stays on disk.</summary>
    void Close(string vaultName);

    long CountNotes(string vaultName);
}

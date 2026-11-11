namespace Dendrite.Application.Abstractions;

/// <summary>Opens (and migrates) the index database that belongs to a vault.</summary>
public interface IIndexDatabase
{
    /// <summary>
    /// Opens the index stored inside the vault folder. If the folder moved since the index was written, the stored note
    /// paths are rewritten to the current location. <paramref name="previousRoot"/> is the location to assume for an
    /// index that has not recorded one yet.
    /// </summary>
    void Open(string vaultName, string vaultPath, string? previousRoot = null);

    /// <summary>Closes the vault's index. The file stays on disk.</summary>
    void Close(string vaultName);

    /// <summary>The number of indexed notes, or 0 when the vault's index is not open.</summary>
    long CountNotes(string vaultName);
}

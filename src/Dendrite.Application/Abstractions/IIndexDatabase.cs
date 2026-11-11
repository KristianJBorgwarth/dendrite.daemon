namespace Dendrite.Application.Abstractions;

public interface IIndexDatabase
{
    void Open(string vaultName, string vaultPath, string? previousRoot = null);

    void Close(string vaultName);

    long CountNotes(string vaultName);
}

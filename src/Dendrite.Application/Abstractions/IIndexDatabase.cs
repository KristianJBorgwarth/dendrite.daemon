namespace Dendrite.Application.Abstractions;

/// <summary>Opens (and migrates) the index database that belongs to a vault.</summary>
public interface IIndexDatabase
{
    void Open(string vaultName);
}

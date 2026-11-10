using Dendrite.Application.Abstractions;
using Dendrite.Application.Vaults;
using Microsoft.Data.Sqlite;

namespace Dendrite.Infrastructure.Persistence;

/// <summary>Holds one index per vault and exposes the one belonging to the vault selected for the current request.</summary>
public sealed class Database : IIndexDatabase, IDisposable
{
    private readonly VaultContext _context;
    private readonly Dictionary<string, VaultIndex> _indexes = new(StringComparer.Ordinal);

    public Database(VaultContext context)
    {
        _context = context;
    }

    public SqliteTransaction? Transaction => Current.Transaction;

    public SqliteConnection Connection => Current.Connection;

    private VaultIndex Current => For(_context.Config.VaultName);

    public static string DataRoot()
    {
        var dataHome = Environment.GetEnvironmentVariable("XDG_DATA_HOME");
        if (string.IsNullOrEmpty(dataHome))
            dataHome = Path.Combine(Environment.GetFolderPath(Environment.SpecialFolder.UserProfile), ".local", "share");

        return Path.Combine(dataHome, "dendrite");
    }

    public static string IndexPathForVault(string vaultName) => Path.Combine(DataRoot(), vaultName, "index.db");

    public void Open(string vaultName) => Initialize(vaultName, IndexPathForVault(vaultName));

    public void Initialize(string vaultName, string dbPath)
    {
        var index = new VaultIndex(dbPath);

        if (_indexes.Remove(vaultName, out var previous))
            previous.Dispose();

        _indexes[vaultName] = index;
    }

    public void Close(string vaultName)
    {
        if (_indexes.Remove(vaultName, out var index))
            index.Dispose();
    }

    public long CountNotes(string vaultName)
    {
        using var cmd = For(vaultName).Connection.CreateCommand();
        cmd.CommandText = "SELECT COUNT(*) FROM note";
        return (long)cmd.ExecuteScalar()!;
    }

    public void RelocateNotes(string vaultName, string oldRoot, string newRoot)
    {
        using var cmd = For(vaultName).Connection.CreateCommand();
        cmd.CommandText = "UPDATE note SET path = @newRoot || substr(path, length(@oldRoot) + 1) WHERE substr(path, 1, length(@oldRoot) + 1) = @oldPrefix";
        cmd.Parameters.AddWithValue("@newRoot", newRoot);
        cmd.Parameters.AddWithValue("@oldRoot", oldRoot);
        cmd.Parameters.AddWithValue("@oldPrefix", oldRoot + "/");
        cmd.ExecuteNonQuery();
    }

    /// <summary>The connection of a specific vault, regardless of the current selection.</summary>
    public SqliteConnection ConnectionFor(string vaultName) => For(vaultName).Connection;

    public SqliteTransaction BeginTransaction() => Current.BeginTransaction();

    public void EndTransaction() => Current.EndTransaction();

    public void Dispose()
    {
        foreach (var index in _indexes.Values)
            index.Dispose();
    }

    private VaultIndex For(string vaultName) =>
        _indexes.TryGetValue(vaultName, out var index)
            ? index
            : throw new InvalidOperationException($"index for vault '{vaultName}' is not open");
}

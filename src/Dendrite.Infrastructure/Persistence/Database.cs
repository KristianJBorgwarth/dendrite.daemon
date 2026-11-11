using Dendrite.Application.Abstractions;
using Dendrite.Application.Vaults;
using Microsoft.Data.Sqlite;

namespace Dendrite.Infrastructure.Persistence;

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

    public static string IndexPath(string vaultPath) => Path.Combine(vaultPath, ".dendrite", "index.db");

    public static string LegacyIndexPath(string vaultName) => Path.Combine(DataRoot(), vaultName, "index.db");

    public void Open(string vaultName, string vaultPath, string? previousRoot = null)
    {
        var dbPath = IndexPath(vaultPath);
        AdoptLegacyIndex(vaultName, dbPath);
        Initialize(vaultName, dbPath);

        var stored = ReadRoot(vaultName);
        var from = stored ?? previousRoot ?? vaultPath;
        if (from != vaultPath)
            RelocateNotes(vaultName, from, vaultPath);

        WriteRoot(vaultName, vaultPath);
    }

    private static void AdoptLegacyIndex(string vaultName, string dbPath)
    {
        var legacy = LegacyIndexPath(vaultName);
        if (File.Exists(dbPath) || !File.Exists(legacy))
            return;

        Directory.CreateDirectory(Path.GetDirectoryName(dbPath)!);
        foreach (var suffix in new[] { "", "-wal", "-shm" })
        {
            if (File.Exists(legacy + suffix))
                File.Move(legacy + suffix, dbPath + suffix);
        }
    }

    private string? ReadRoot(string vaultName)
    {
        using var cmd = For(vaultName).Connection.CreateCommand();
        cmd.CommandText = "SELECT value FROM vault_meta WHERE key = 'root'";
        return cmd.ExecuteScalar() as string;
    }

    private void WriteRoot(string vaultName, string root)
    {
        using var cmd = For(vaultName).Connection.CreateCommand();
        cmd.CommandText = "INSERT INTO vault_meta (key, value) VALUES ('root', @root) ON CONFLICT(key) DO UPDATE SET value = @root";
        cmd.Parameters.AddWithValue("@root", root);
        cmd.ExecuteNonQuery();
    }

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
        if (!_indexes.ContainsKey(vaultName))
            return 0;

        using var cmd = For(vaultName).Connection.CreateCommand();
        cmd.CommandText = "SELECT COUNT(*) FROM note";
        return (long)cmd.ExecuteScalar()!;
    }

    private void RelocateNotes(string vaultName, string oldRoot, string newRoot)
    {
        using var cmd = For(vaultName).Connection.CreateCommand();
        cmd.CommandText = "UPDATE note SET path = @newRoot || substr(path, length(@oldRoot) + 1) WHERE substr(path, 1, length(@oldRoot) + 1) = @oldPrefix";
        cmd.Parameters.AddWithValue("@newRoot", newRoot);
        cmd.Parameters.AddWithValue("@oldRoot", oldRoot);
        cmd.Parameters.AddWithValue("@oldPrefix", oldRoot + "/");
        cmd.ExecuteNonQuery();
    }

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
            : throw new InvalidOperationException($"index for vault '{vaultName}' is not open — is its folder missing? Open the folder as a vault again");
}

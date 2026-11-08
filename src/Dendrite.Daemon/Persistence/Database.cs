using System.Reflection;
using Microsoft.Data.Sqlite;

namespace Dendrite.Daemon.Persistence;

public sealed class Database : IDisposable
{
    private SqliteConnection? _connection;

    public SqliteTransaction? Transaction { get; private set; }

    public SqliteConnection Connection =>
        _connection ?? throw new InvalidOperationException("database is not initialized — call vault/init first");

    public static string IndexPathForVault(string vaultName)
    {
        var dataHome = Environment.GetEnvironmentVariable("XDG_DATA_HOME");
        if (string.IsNullOrEmpty(dataHome))
            dataHome = Path.Combine(Environment.GetFolderPath(Environment.SpecialFolder.UserProfile), ".local", "share");

        return Path.Combine(dataHome, "dendrite", vaultName, "index.db");
    }

    public void Initialize(string dbPath)
    {
        _connection?.Dispose();
        Directory.CreateDirectory(Path.GetDirectoryName(dbPath)!);

        var connection = new SqliteConnection($"Data Source={dbPath}");
        connection.Open();
        Migrator.Apply(connection);
        _connection = connection;
    }

    public SqliteTransaction BeginTransaction()
    {
        if (Transaction is not null)
            throw new InvalidOperationException("a transaction is already active");

        return Transaction = Connection.BeginTransaction();
    }

    public void EndTransaction()
    {
        Transaction?.Dispose();
        Transaction = null;
    }

    public void Dispose()
    {
        EndTransaction();
        _connection?.Dispose();
    }
}

internal static class Migrator
{
    public static void Apply(SqliteConnection connection)
    {
        Execute(connection, "CREATE TABLE IF NOT EXISTS schema_migrations (version TEXT PRIMARY KEY);");

        var assembly = Assembly.GetExecutingAssembly();
        var migrations = assembly.GetManifestResourceNames()
            .Where(n => n.EndsWith(".sql", StringComparison.Ordinal))
            .Select(n => (Resource: n, Version: VersionOf(n)))
            .OrderBy(m => m.Version, StringComparer.Ordinal);

        foreach (var (resource, version) in migrations)
        {
            if (IsApplied(connection, version))
                continue;

            using var reader = new StreamReader(assembly.GetManifestResourceStream(resource)!);
            using var tx = connection.BeginTransaction();
            Execute(connection, reader.ReadToEnd(), tx);
            Execute(connection, $"INSERT INTO schema_migrations (version) VALUES ('{version}');", tx);
            tx.Commit();

            Log.Info($"applied migration {version}");
        }
    }

    // "Dendrite.Daemon.Persistence.Migrations.001_note.sql" -> "001_note.sql"
    private static string VersionOf(string resourceName)
    {
        var withoutExt = resourceName[..^".sql".Length];
        return withoutExt[(withoutExt.LastIndexOf('.') + 1)..] + ".sql";
    }

    private static bool IsApplied(SqliteConnection connection, string version)
    {
        using var cmd = connection.CreateCommand();
        cmd.CommandText = "SELECT COUNT(*) FROM schema_migrations WHERE version = $v";
        cmd.Parameters.AddWithValue("$v", version);
        return (long)cmd.ExecuteScalar()! > 0;
    }

    private static void Execute(SqliteConnection connection, string sql, SqliteTransaction? tx = null)
    {
        using var cmd = connection.CreateCommand();
        cmd.Transaction = tx;
        cmd.CommandText = sql;
        cmd.ExecuteNonQuery();
    }
}

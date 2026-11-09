using Dendrite.Application.Abstractions;
using Microsoft.Data.Sqlite;

namespace Dendrite.Infrastructure.Persistence;

public sealed class Database : IIndexDatabase, IDisposable
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

    public void Open(string vaultName) => Initialize(IndexPathForVault(vaultName));

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

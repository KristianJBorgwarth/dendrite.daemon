using Microsoft.Data.Sqlite;

namespace Dendrite.Infrastructure.Persistence;

/// <summary>The SQLite index of a single vault.</summary>
internal sealed class VaultIndex : IDisposable
{
    public SqliteConnection Connection { get; }

    public SqliteTransaction? Transaction { get; private set; }

    public VaultIndex(string dbPath)
    {
        Directory.CreateDirectory(Path.GetDirectoryName(dbPath)!);

        var connection = new SqliteConnection($"Data Source={dbPath}");
        connection.Open();
        Migrator.Apply(connection);
        Connection = connection;
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
        Connection.Dispose();
    }
}

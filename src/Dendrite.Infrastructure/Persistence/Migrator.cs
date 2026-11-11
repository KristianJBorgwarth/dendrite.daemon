using System.Reflection;
using Dendrite.Application;
using Microsoft.Data.Sqlite;

namespace Dendrite.Infrastructure.Persistence;

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

using Dendrite.Application.Features.Vault;
using Dendrite.Infrastructure.Persistence;
using Mediator;
using Microsoft.Data.Sqlite;
using Microsoft.Extensions.DependencyInjection;

namespace Dendrite.Daemon.Tests;

/// <summary>An isolated vault directory plus index database, wired through the real service graph.</summary>
internal sealed class TestVault : IDisposable
{
    private readonly ServiceProvider _services = Composition.BuildServices();
    private readonly string _vaultName = $"test-{Guid.NewGuid():N}";

    public string Root { get; } = Directory.CreateTempSubdirectory("dendrite-vault-").FullName;

    public VaultDispatcher Sender => _services.GetRequiredService<VaultDispatcher>();

    public T Get<T>() where T : notnull => _services.GetRequiredService<T>();

    public TestVault(Action<string>? seed = null, Func<VaultConfigDto>? config = null)
    {
        seed?.Invoke(Root);
        Initialize(config?.Invoke());
    }

    public void Initialize(VaultConfigDto? config = null)
    {
        config ??= new VaultConfigDto { TemplatesDir = "/.templates" };
        var withPaths = new VaultConfigDto
        {
            VaultName = _vaultName,
            VaultPath = Root,
            TemplatesDir = config.TemplatesDir,
            ExcludeIndexFiles = config.ExcludeIndexFiles,
            OverrideDefaultIgnores = config.OverrideDefaultIgnores,
        };
        Sender.Send(new InitializeVault(withPaths)).AsTask().GetAwaiter().GetResult();
    }

    public string Write(string relativePath, string content)
    {
        var path = Path.Combine(Root, relativePath);
        Directory.CreateDirectory(Path.GetDirectoryName(path)!);
        File.WriteAllText(path, content);
        return path;
    }

    public long Count(string sql, params (string Name, object Value)[] args)
    {
        var db = Get<Database>();
        using var cmd = db.Connection.CreateCommand();
        cmd.CommandText = sql;
        foreach (var (name, value) in args)
            cmd.Parameters.AddWithValue(name, value);
        return (long)cmd.ExecuteScalar()!;
    }

    public void Dispose()
    {
        _services.Dispose();
        SqliteConnection.ClearAllPools();
        Directory.Delete(Root, recursive: true);
        var dbDir = Path.GetDirectoryName(Database.IndexPathForVault(_vaultName))!;
        if (Directory.Exists(dbDir))
            Directory.Delete(dbDir, recursive: true);
    }
}

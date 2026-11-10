using AwesomeAssertions;
using Dendrite.Application.Features.Notes;
using Dendrite.Application.Features.Vault;
using Dendrite.Infrastructure.Persistence;
using Microsoft.Data.Sqlite;
using Microsoft.Extensions.DependencyInjection;

namespace Dendrite.Daemon.Tests;

public sealed class MultiVaultTests : IDisposable
{
    private readonly ServiceProvider _services = Composition.BuildServices();
    private readonly List<string> _names = [];
    private readonly List<string> _roots = [];

    private VaultDispatcher Sender => _services.GetRequiredService<VaultDispatcher>();

    private async Task<(string Name, string Root)> AddVault(string? name = null, string? root = null)
    {
        name ??= $"mv-{Guid.NewGuid():N}";
        root ??= Directory.CreateTempSubdirectory("dendrite-mv-").FullName;
        if (!_names.Contains(name))
            _names.Add(name);
        if (!_roots.Contains(root))
            _roots.Add(root);
        await Sender.Send(new InitializeVault(new VaultConfigDto
        {
            VaultName = name,
            VaultPath = root,
            TemplatesDir = "/.templates",
        }));
        return (name, root);
    }

    private long Count(string vault, string sql)
    {
        using var cmd = _services.GetRequiredService<Database>().ConnectionFor(vault).CreateCommand();
        cmd.CommandText = sql;
        return (long)cmd.ExecuteScalar()!;
    }

    [Fact]
    public async Task SameSlugInTwoVaults_StaysIsolated()
    {
        var a = await AddVault();
        var b = await AddVault();

        var pathA = await Sender.Send(new CreateNote("Shared", "", "", a.Name));
        var pathB = await Sender.Send(new CreateNote("Shared", "", "", b.Name));

        pathA.Should().Be(Path.Combine(a.Root, "shared.md"));
        pathB.Should().Be(Path.Combine(b.Root, "shared.md"));
        Count(a.Name, "SELECT COUNT(*) FROM note").Should().Be(1);
        Count(b.Name, "SELECT COUNT(*) FROM note").Should().Be(1);

        await Sender.Send(new DeleteNote(pathA));

        Count(a.Name, "SELECT COUNT(*) FROM note").Should().Be(0);
        Count(b.Name, "SELECT COUNT(*) FROM note").Should().Be(1);
    }

    [Fact]
    public async Task NoVaultParam_ResolvesFromPath()
    {
        var a = await AddVault();
        var b = await AddVault();
        var path = await Sender.Send(new CreateNote("Inside B", "", "", b.Name));

        await Sender.Send(new RebuildIndex(a.Name));
        await Sender.Send(new SaveNote(path));

        Count(a.Name, "SELECT COUNT(*) FROM note").Should().Be(0);
        Count(b.Name, "SELECT COUNT(*) FROM note").Should().Be(1);
    }

    [Fact]
    public async Task SeveralVaults_WithoutHint_IsAmbiguous()
    {
        await AddVault();
        await AddVault();

        var act = async () => await Sender.Send(new RebuildIndex());

        await act.Should().ThrowAsync<InvalidOperationException>().WithMessage("*several vaults*");
    }

    [Fact]
    public async Task UnknownVault_Throws()
    {
        await AddVault();

        var act = async () => await Sender.Send(new RebuildIndex("nope"));

        await act.Should().ThrowAsync<InvalidOperationException>().WithMessage("*unknown vault 'nope'*");
    }

    [Fact]
    public async Task ReinitSameVault_IsIdempotent()
    {
        var a = await AddVault();

        await AddVault(a.Name, a.Root);
        await Sender.Send(new CreateNote("Still Works", "", "", a.Name));

        Count(a.Name, "SELECT COUNT(*) FROM note").Should().Be(1);
    }

    [Fact]
    public async Task ReinitSameNameAtOtherPath_IsRejected()
    {
        var a = await AddVault();

        var act = async () => await AddVault(a.Name);

        await act.Should().ThrowAsync<InvalidOperationException>().WithMessage("*already registered*");
    }

    [Theory]
    [InlineData("")]
    [InlineData("a/b")]
    [InlineData("..")]
    public async Task InvalidVaultName_IsRejected(string name)
    {
        var act = async () => await AddVault(name);

        await act.Should().ThrowAsync<InvalidOperationException>().WithMessage("*invalid vault name*");
    }

    public void Dispose()
    {
        _services.Dispose();
        SqliteConnection.ClearAllPools();
        foreach (var root in _roots)
            Directory.Delete(root, recursive: true);
        foreach (var name in _names)
        {
            var dir = Path.GetDirectoryName(Database.IndexPathForVault(name))!;
            if (name.StartsWith("mv-", StringComparison.Ordinal) && Directory.Exists(dir))
                Directory.Delete(dir, recursive: true);
        }
    }
}

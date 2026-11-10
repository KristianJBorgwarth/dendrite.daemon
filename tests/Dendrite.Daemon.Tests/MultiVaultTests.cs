using AwesomeAssertions;
using Dendrite.Application;
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
    public async Task ListVaults_ReturnsNameRootAndNoteCount()
    {
        var a = await AddVault();
        var b = await AddVault();
        await Sender.Send(new CreateNote("One", "", "", b.Name));

        var vaults = await Sender.Send(new ListVaults());

        vaults.Should().BeEquivalentTo([
            new VaultInfoDto(a.Name, a.Root, 0),
            new VaultInfoDto(b.Name, b.Root, 1),
        ]);
    }

    [Fact]
    public async Task RemoveVault_UnregistersAndLeavesIndexOnDisk()
    {
        var a = await AddVault();
        var b = await AddVault();

        await Sender.Send(new RemoveVault(a.Name));

        (await Sender.Send(new ListVaults())).Select(v => v.Name).Should().Equal(b.Name);
        File.Exists(Database.IndexPathForVault(a.Name)).Should().BeTrue();
        var act = async () => await Sender.Send(new RebuildIndex(a.Name));
        await act.Should().ThrowAsync<InvalidOperationException>().WithMessage("*unknown vault*");
    }

    [Fact]
    public async Task RemoveVault_LeavesRemainingVaultAsDefault()
    {
        var a = await AddVault();
        var b = await AddVault();

        await Sender.Send(new RemoveVault(a.Name));
        var path = await Sender.Send(new CreateNote("Now Default", "", ""));

        path.Should().Be(Path.Combine(b.Root, "now-default.md"));
    }

    [Fact]
    public async Task RemoveUnknownVault_Throws()
    {
        var act = async () => await Sender.Send(new RemoveVault("nope"));

        await act.Should().ThrowAsync<InvalidOperationException>().WithMessage("*unknown vault 'nope'*");
    }

    [Fact]
    public async Task CreateVault_ScaffoldsTemplatesRegistersAndIndexes()
    {
        var name = $"mv-{Guid.NewGuid():N}";
        var root = Path.Combine(Directory.CreateTempSubdirectory("dendrite-mv-").FullName, "fresh");
        _names.Add(name);
        _roots.Add(Path.GetDirectoryName(root)!);

        await Sender.Send(new CreateVault(new VaultConfigDto { VaultName = name, VaultPath = root }));

        File.Exists(Path.Combine(root, ".templates", "daily.md")).Should().BeTrue();
        File.Exists(Path.Combine(root, ".templates", "scratch.md")).Should().BeTrue();
        var path = await Sender.Send(new CreateNote("First", "daily", "/daily", name));
        path.Should().Be(Path.Combine(root, "daily", "first.md"));
        Count(name, "SELECT COUNT(*) FROM note").Should().Be(1);
    }

    [Fact]
    public async Task CreateVault_OverExistingVault_KeepsItsFiles()
    {
        var a = await AddVault();
        await Sender.Send(new CreateNote("Keep", "", "", a.Name));
        await Sender.Send(new RemoveVault(a.Name));
        Directory.CreateDirectory(Path.Combine(a.Root, ".templates"));
        File.WriteAllText(Path.Combine(a.Root, ".templates", "daily.md"), "custom");

        await Sender.Send(new CreateVault(new VaultConfigDto { VaultName = a.Name, VaultPath = a.Root }));

        File.ReadAllText(Path.Combine(a.Root, ".templates", "daily.md")).Should().Be("custom");
        Count(a.Name, "SELECT COUNT(*) FROM note").Should().Be(1);
    }

    [Fact]
    public async Task CreateVault_NonEmptyNonVaultDirectory_IsRefused()
    {
        var root = Directory.CreateTempSubdirectory("dendrite-mv-").FullName;
        _roots.Add(root);
        File.WriteAllText(Path.Combine(root, "stuff.txt"), "x");

        var act = async () => await Sender.Send(new CreateVault(new VaultConfigDto { VaultName = "mv-refused", VaultPath = root }));

        await act.Should().ThrowAsync<InvalidOperationException>().WithMessage("*not empty*");
        Directory.Exists(Path.Combine(root, ".templates")).Should().BeFalse();
    }

    [Fact]
    public async Task CreateVault_NameTakenAtOtherPath_IsRefusedWithoutScaffolding()
    {
        var a = await AddVault();
        var other = Directory.CreateTempSubdirectory("dendrite-mv-").FullName;
        _roots.Add(other);

        var act = async () => await Sender.Send(new CreateVault(new VaultConfigDto { VaultName = a.Name, VaultPath = other }));

        await act.Should().ThrowAsync<InvalidOperationException>().WithMessage("*already registered*");
        Directory.Exists(Path.Combine(other, ".templates")).Should().BeFalse();
    }

    [Theory]
    [InlineData("")]
    [InlineData("a/b")]
    [InlineData("..")]
    public async Task CreateVault_InvalidName_IsRefusedWithoutScaffolding(string name)
    {
        var root = Directory.CreateTempSubdirectory("dendrite-mv-").FullName;
        _roots.Add(root);

        var act = async () => await Sender.Send(new CreateVault(new VaultConfigDto { VaultName = name, VaultPath = root }));

        await act.Should().ThrowAsync<InvalidOperationException>().WithMessage("*invalid vault name*");
        Directory.Exists(Path.Combine(root, ".templates")).Should().BeFalse();
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

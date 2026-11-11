using AwesomeAssertions;
using Dendrite.Application;
using Dendrite.Application.Features.Notes;
using Dendrite.Application.Features.Vault;
using Dendrite.Application.Services;
using Dendrite.Application.Vaults;
using Dendrite.Infrastructure.Persistence;
using Microsoft.Data.Sqlite;
using Microsoft.Extensions.DependencyInjection;

namespace Dendrite.Daemon.Tests;

public sealed class MultiVaultTests : IDisposable
{
    private readonly InMemoryVaultCatalog _catalog = new();
    private readonly ServiceProvider _services;
    private readonly List<string> _names = [];
    private readonly List<string> _roots = [];

    public MultiVaultTests()
    {
        _services = Composition.BuildServices(_catalog);
    }

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
            new VaultInfoDto(a.Name, a.Root, 0, false),
            new VaultInfoDto(b.Name, b.Root, 1, false),
        ]);
    }

    [Fact]
    public async Task RemoveVault_UnregistersAndLeavesIndexOnDisk()
    {
        var a = await AddVault();
        var b = await AddVault();

        await Sender.Send(new RemoveVault(a.Name));

        (await Sender.Send(new ListVaults())).Select(v => v.Name).Should().Equal(b.Name);
        File.Exists(Database.IndexPath(a.Root)).Should().BeTrue();
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

    [Fact]
    public async Task RestoredVaults_ComeBackAfterRestart()
    {
        var a = await AddVault();

        using var restarted = Composition.BuildServices(_catalog);
        restarted.GetRequiredService<VaultRestorer>().Restore();

        var vaults = await restarted.GetRequiredService<VaultDispatcher>().Send(new ListVaults());
        vaults.Should().Contain(v => v.Name == a.Name && v.Path == a.Root);
    }

    [Fact]
    public async Task UseVault_MakesItTheFallbackAndSurvivesRestart()
    {
        var a = await AddVault();
        var b = await AddVault();

        await Sender.Send(new UseVault(b.Name));
        var path = await Sender.Send(new CreateNote("Fallback", "", ""));

        path.Should().StartWith(b.Root);
        (await Sender.Send(new ListVaults())).Single(v => v.Active).Name.Should().Be(b.Name);

        using var restarted = Composition.BuildServices(_catalog);
        restarted.GetRequiredService<VaultRestorer>().Restore();

        var vaults = await restarted.GetRequiredService<VaultDispatcher>().Send(new ListVaults());
        vaults.Should().Contain(v => v.Name == a.Name && v.Path == a.Root && !v.Active);
        vaults.Should().Contain(v => v.Name == b.Name && v.Path == b.Root && v.Active);
    }

    [Fact]
    public async Task UseUnknownVault_Throws()
    {
        var act = async () => await Sender.Send(new UseVault("nope"));

        await act.Should().ThrowAsync<InvalidOperationException>().WithMessage("*unknown vault*");
    }

    private string Move(string root)
    {
        var parent = Directory.CreateTempSubdirectory("dendrite-mv-").FullName;
        var moved = Path.Combine(parent, "moved");
        Directory.Move(root, moved);
        _roots.Remove(root);
        _roots.Add(parent);
        return moved;
    }

    private string NotePath(string vault)
    {
        using var cmd = _services.GetRequiredService<Database>().ConnectionFor(vault).CreateCommand();
        cmd.CommandText = "SELECT path FROM note LIMIT 1";
        return (string)cmd.ExecuteScalar()!;
    }

    [Fact]
    public async Task OpenMovedFolder_KeepsTheIndexAndRewritesNotePaths()
    {
        var a = await AddVault();
        await Sender.Send(new CreateNote("Moves With Me", "", "", a.Name));
        var moved = Move(a.Root);

        await Sender.Send(new OpenVault(moved));

        NotePath(a.Name).Should().StartWith(moved + "/");
        Count(a.Name, "SELECT COUNT(*) FROM note").Should().Be(1);
        var info = (await Sender.Send(new ListVaults())).Single(v => v.Name == a.Name);
        info.Path.Should().Be(moved);
        info.Missing.Should().BeFalse();
    }

    [Fact]
    public async Task TouchingAFileInAMovedVault_RepointsItAutomatically()
    {
        var a = await AddVault();
        var moved = Move(a.Root);

        using var restarted = Composition.BuildServices(_catalog);
        restarted.GetRequiredService<VaultRestorer>().Restore();
        var sender = restarted.GetRequiredService<VaultDispatcher>();

        (await sender.Send(new ListVaults())).Single(v => v.Name == a.Name).Missing.Should().BeTrue();

        await sender.Send(new GetBacklinks(Path.Combine(moved, "x.md")));

        (await sender.Send(new ListVaults())).Single(v => v.Name == a.Name).Path.Should().Be(moved);
    }

    [Fact]
    public async Task OpenVault_RegistersACreatedVaultFromItsManifest()
    {
        var root = Directory.CreateTempSubdirectory("dendrite-mv-").FullName;
        _roots.Add(root);
        var name = $"mv-{Guid.NewGuid():N}";
        _names.Add(name);
        await Sender.Send(new CreateVault(new VaultConfigDto
        {
            VaultName = name,
            VaultPath = root,
            DailyNotes = new DailyNotesConfig { Dir = "/logs" },
        }));
        await Sender.Send(new RemoveVault(name));

        await Sender.Send(new OpenVault(root));

        var info = (await Sender.Send(new ListVaults())).Single(v => v.Name == name);
        info.Path.Should().Be(root);
        info.Active.Should().BeTrue();
        _services.GetRequiredService<VaultRegistry>().Resolve(name, null).DailyNotes.Dir.Should().Be("/logs");
    }

    [Fact]
    public async Task OpenVault_AtANewFolder_RepointsTheVaultOfThatName()
    {
        var root = Directory.CreateTempSubdirectory("dendrite-mv-").FullName;
        _roots.Add(root);
        var name = $"mv-{Guid.NewGuid():N}";
        _names.Add(name);
        await Sender.Send(new CreateVault(new VaultConfigDto { VaultName = name, VaultPath = root }));
        var moved = Move(root);

        await Sender.Send(new OpenVault(moved));

        (await Sender.Send(new ListVaults())).Single(v => v.Name == name).Path.Should().Be(moved);
    }

    public void Dispose()
    {
        _services.Dispose();
        SqliteConnection.ClearAllPools();
        foreach (var root in _roots)
            Directory.Delete(root, recursive: true);
    }
}

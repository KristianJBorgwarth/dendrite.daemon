using AwesomeAssertions;

namespace Dendrite.Test.Integration;

public sealed class DaemonProcessTests : IDisposable
{
    private readonly string _dataHome = Directory.CreateTempSubdirectory("dendrite-it-data-").FullName;
    private readonly string _vaultRoot = Directory.CreateTempSubdirectory("dendrite-it-vault-").FullName;
    private string? _movedRoot;

    private static object CreateConfig(string root) => new { config = new { vault_name = "it", vault_path = root } };

    [Fact]
    public async Task CreatedNotesAndBacklinks_WorkOverStdio()
    {
        using var daemon = new DaemonProcess(_dataHome);
        await daemon.Call("vault/create", CreateConfig(_vaultRoot));

        var target = (await daemon.Call("note/create", new { title = "Target", templateName = "", directory = "" })).GetString()!;
        var source = (await daemon.Call("note/create", new { title = "Source", templateName = "", directory = "" })).GetString()!;
        await File.WriteAllTextAsync(source, "---\ntitle: Source\n---\nsee [[target]]");
        await daemon.Call("note/save", new { path = source });

        var backlinks = await daemon.Call("note/backlinks", new { path = target });

        backlinks.GetArrayLength().Should().Be(1);
    }

    [Fact]
    public async Task VaultsAndIndex_SurviveARestart()
    {
        using (var first = new DaemonProcess(_dataHome))
        {
            await first.Call("vault/create", CreateConfig(_vaultRoot));
            await first.Call("note/create", new { title = "Stays", templateName = "", directory = "" });
        }

        using var second = new DaemonProcess(_dataHome);
        var vault = (await second.Call("vault/list")).EnumerateArray().Single();

        vault.GetProperty("name").GetString().Should().Be("it");
        vault.GetProperty("noteCount").GetInt32().Should().Be(1);
        vault.GetProperty("missing").GetBoolean().Should().BeFalse();
    }

    [Fact]
    public async Task MovedVault_IsReportedMissingThenRecoveredByOpen()
    {
        using (var first = new DaemonProcess(_dataHome))
        {
            await first.Call("vault/create", CreateConfig(_vaultRoot));
            await first.Call("note/create", new { title = "Travels", templateName = "", directory = "" });
        }

        _movedRoot = _vaultRoot + "-moved";
        Directory.Move(_vaultRoot, _movedRoot);

        using var second = new DaemonProcess(_dataHome);
        (await second.Call("vault/list")).EnumerateArray().Single().GetProperty("missing").GetBoolean().Should().BeTrue();

        await second.Call("vault/open", new { path = _movedRoot });

        var vault = (await second.Call("vault/list")).EnumerateArray().Single();
        vault.GetProperty("missing").GetBoolean().Should().BeFalse();
        vault.GetProperty("path").GetString().Should().Be(_movedRoot);
        vault.GetProperty("noteCount").GetInt32().Should().Be(1);
    }

    public void Dispose()
    {
        foreach (var dir in new[] { _dataHome, _vaultRoot, _movedRoot })
        {
            if (dir is not null && Directory.Exists(dir))
                Directory.Delete(dir, true);
        }
    }
}

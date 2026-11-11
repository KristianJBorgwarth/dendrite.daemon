using System.Text.Json;
using AwesomeAssertions;
using Dendrite.Daemon.Rpc;
using Microsoft.Extensions.DependencyInjection;

namespace Dendrite.Test.Integration;

public class RpcServerTests
{
    private static async Task<JsonElement> Call(TestVault vault, string line)
    {
        var reply = await vault.Get<RpcServer>().HandleLineAsync(line, CancellationToken.None);
        reply.Should().NotBeNull();
        return JsonDocument.Parse(reply!).RootElement;
    }

    [Fact]
    public async Task InvalidJson_ReturnsParseError()
    {
        using var vault = new TestVault();

        var reply = await Call(vault, "{nope");

        reply.GetProperty("error").GetProperty("code").GetInt32().Should().Be(-32700);
        reply.GetProperty("id").ValueKind.Should().Be(JsonValueKind.Null);
    }

    [Fact]
    public async Task UnknownMethod_ReturnsMethodNotFound()
    {
        using var vault = new TestVault();

        var reply = await Call(vault, """{"jsonrpc":"2.0","id":4,"method":"nope"}""");

        reply.GetProperty("id").GetInt32().Should().Be(4);
        reply.GetProperty("error").GetProperty("code").GetInt32().Should().Be(-32601);
    }

    [Fact]
    public async Task Notification_GetsNoReply()
    {
        using var vault = new TestVault();

        var reply = await vault.Get<RpcServer>().HandleLineAsync("""{"jsonrpc":"2.0","method":"vault/rebuild"}""", CancellationToken.None);

        reply.Should().BeNull();
    }

    [Fact]
    public async Task NoteCreate_UsesCamelCaseParamsAndReturnsPathAsResult()
    {
        using var vault = new TestVault();

        var reply = await Call(vault, """{"jsonrpc":"2.0","id":1,"method":"note/create","params":{"title":"Hi There","templateName":"","directory":""}}""");

        reply.GetProperty("result").GetString().Should().Be(Path.Combine(vault.Root, "hi-there.md"));
    }

    [Fact]
    public async Task VaultInit_AcceptsSnakeCaseConfig()
    {
        using var vault = new TestVault();
        var other = Directory.CreateTempSubdirectory("dendrite-vault-").FullName;
        var name = $"rpc-{Guid.NewGuid():N}";
        try
        {
            var json = JsonSerializer.Serialize(new { vault_name = name, vault_path = other });
            var reply = await Call(vault, "{\"jsonrpc\":\"2.0\",\"id\":1,\"method\":\"vault/init\",\"params\":{\"config\":" + json + "}}");

            reply.TryGetProperty("error", out var error).Should().BeFalse(error.ToString());
            reply.GetProperty("result").ValueKind.Should().Be(JsonValueKind.Null);
        }
        finally
        {
            Directory.Delete(other, true);
        }
    }

    [Fact]
    public async Task GotoUnknownNote_ReturnsNullResult()
    {
        using var vault = new TestVault();

        var reply = await Call(vault, """{"jsonrpc":"2.0","id":2,"method":"note/goto","params":{"link":"[[nope]]"}}""");

        reply.GetProperty("result").ValueKind.Should().Be(JsonValueKind.Null);
    }
}

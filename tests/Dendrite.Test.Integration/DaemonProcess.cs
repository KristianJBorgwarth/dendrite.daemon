using System.Diagnostics;
using System.Text.Json;

namespace Dendrite.Test.Integration;

internal sealed class DaemonProcess : IDisposable
{
    private readonly Process _process;
    private int _id;

    public DaemonProcess(string dataHome)
    {
        var start = new ProcessStartInfo("dotnet", $"\"{Path.Combine(AppContext.BaseDirectory, "dendrite.dll")}\"")
        {
            RedirectStandardInput = true,
            RedirectStandardOutput = true,
            RedirectStandardError = true,
            UseShellExecute = false,
        };
        start.Environment["XDG_DATA_HOME"] = dataHome;

        _process = Process.Start(start)!;
    }

    public async Task<JsonElement> Call(string method, object? parameters = null)
    {
        var request = JsonSerializer.Serialize(new { jsonrpc = "2.0", id = ++_id, method, @params = parameters });
        await _process.StandardInput.WriteLineAsync(request);

        var line = await _process.StandardOutput.ReadLineAsync().WaitAsync(TimeSpan.FromSeconds(30))
            ?? throw new InvalidOperationException("the daemon closed its output");

        var reply = JsonDocument.Parse(line).RootElement;
        if (reply.TryGetProperty("error", out var error))
            throw new InvalidOperationException($"{method} failed: {error}");

        return reply.GetProperty("result");
    }

    public void Dispose()
    {
        _process.StandardInput.Close();
        if (!_process.WaitForExit(5000))
            _process.Kill();

        _process.Dispose();
    }
}

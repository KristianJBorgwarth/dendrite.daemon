namespace Dendrite.Application;

/// <summary>Stdout carries the RPC protocol, so diagnostics go to a sink the host configures (a log file).</summary>
public static class Log
{
    private static Action<string, string>? _sink;

    public static void Configure(Action<string, string> sink) => _sink = sink;

    public static void Debug(string message) => _sink?.Invoke("DEBUG", message);
    public static void Info(string message) => _sink?.Invoke("INFO", message);
    public static void Warn(string message) => _sink?.Invoke("WARN", message);
    public static void Error(string message) => _sink?.Invoke("ERROR", message);
}

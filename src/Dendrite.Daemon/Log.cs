namespace Dendrite.Daemon;

/// <summary>Stdout carries the RPC protocol, so all diagnostics go to a log file.</summary>
public static class Log
{
    private static readonly Lock Gate = new();
    private static string? _path;

    public static void Init() => _path = Path.Combine(Path.GetTempPath(), "dendrite.log");

    public static void Debug(string message) => Write("DEBUG", message);
    public static void Info(string message) => Write("INFO", message);
    public static void Warn(string message) => Write("WARN", message);
    public static void Error(string message) => Write("ERROR", message);

    private static void Write(string level, string message)
    {
        if (_path is null)
            return;

        try
        {
            lock (Gate)
                File.AppendAllText(_path, $"{DateTime.Now:yyyy-MM-ddTHH:mm:ss.fff} {level} {message}\n");
        }
        catch (IOException)
        {
        }
    }
}

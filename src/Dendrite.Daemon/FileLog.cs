namespace Dendrite.Daemon;

/// <summary>Stdout carries the RPC protocol, so all diagnostics go to a log file.</summary>
internal static class FileLog
{
    private static readonly Lock Gate = new();
    private static readonly string Path = System.IO.Path.Combine(System.IO.Path.GetTempPath(), "dendrite.log");

    public static void Write(string level, string message)
    {
        try
        {
            lock (Gate)
                File.AppendAllText(Path, $"{DateTime.Now:yyyy-MM-ddTHH:mm:ss.fff} {level} {message}\n");
        }
        catch (IOException)
        {
        }
    }
}

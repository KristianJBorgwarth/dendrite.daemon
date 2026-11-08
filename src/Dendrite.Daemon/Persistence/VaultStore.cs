namespace Dendrite.Daemon.Persistence;

public sealed record DailyNotes(string Dir, string FilenameFormat, string TemplateName);

public sealed record VaultConfig(
    string VaultName,
    string VaultPath,
    string TemplatesDir,
    IReadOnlyList<string> ExcludeIndexFiles,
    bool OverrideDefaultIgnores,
    DailyNotes DailyNotes)
{
    private static readonly string[] DefaultIgnores = ["index.md", "index", ".templates"];

    public IReadOnlyList<string> EffectiveExcludes =>
        OverrideDefaultIgnores ? ExcludeIndexFiles : [.. DefaultIgnores, .. ExcludeIndexFiles];
}

public sealed class VaultStore
{
    private VaultConfig? _config;

    public VaultConfig Config =>
        _config ?? throw new InvalidOperationException("vault is not initialized — call vault/init first");

    public void Initialize(VaultConfig config) => _config = config;

    public string GetTemplatePath(string templateName)
    {
        if (Path.GetExtension(templateName) != ".md")
            templateName += ".md";

        return Path.GetFullPath(Path.Join(Config.VaultPath, Config.TemplatesDir, templateName));
    }
}

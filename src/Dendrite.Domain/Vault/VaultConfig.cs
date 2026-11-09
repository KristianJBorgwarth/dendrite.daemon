namespace Dendrite.Domain.Vault;

public sealed record VaultConfig
{
    private static readonly string[] DefaultIgnores = ["index.md", "index", ".templates"];

    public string VaultName { get; }
    public string VaultPath { get; }
    public string TemplatesDir { get; }
    public IReadOnlyList<string> ExcludeIndexFiles { get; }
    public bool OverrideDefaultIgnores { get; }
    public DailyNotes DailyNotes { get; }

    public VaultConfig(
        string vaultName,
        string vaultPath,
        string templatesDir,
        IReadOnlyList<string> excludeIndexFiles,
        bool overrideDefaultIgnores,
        DailyNotes dailyNotes)
    {
        VaultName = vaultName;
        VaultPath = vaultPath;
        TemplatesDir = templatesDir;
        ExcludeIndexFiles = excludeIndexFiles;
        OverrideDefaultIgnores = overrideDefaultIgnores;
        DailyNotes = dailyNotes;
    }

    public IReadOnlyList<string> EffectiveExcludes =>
        OverrideDefaultIgnores ? ExcludeIndexFiles : [.. DefaultIgnores, .. ExcludeIndexFiles];
}

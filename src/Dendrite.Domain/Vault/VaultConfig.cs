namespace Dendrite.Domain.Vault;

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

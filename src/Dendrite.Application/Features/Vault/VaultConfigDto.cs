using System.Text.Json.Serialization;

namespace Dendrite.Application.Features.Vault;

public sealed class VaultConfigDto
{
    [JsonPropertyName("vault_name")] public string VaultName { get; init; } = "";
    [JsonPropertyName("vault_path")] public string VaultPath { get; init; } = "";
    [JsonPropertyName("templates_dir")] public string TemplatesDir { get; init; } = "";
    [JsonPropertyName("exclude_index_files")] public string[] ExcludeIndexFiles { get; init; } = [];
    [JsonPropertyName("override_default_ignores")] public bool OverrideDefaultIgnores { get; init; }
    [JsonPropertyName("daily_notes")] public DailyNotesConfig DailyNotes { get; init; } = new();

    public VaultConfigDto WithPath(string vaultPath) => new()
    {
        VaultName = VaultName,
        VaultPath = vaultPath,
        TemplatesDir = TemplatesDir,
        ExcludeIndexFiles = ExcludeIndexFiles,
        OverrideDefaultIgnores = OverrideDefaultIgnores,
        DailyNotes = DailyNotes,
    };

    public VaultConfigDto WithDefaults() => new()
    {
        VaultName = VaultName,
        VaultPath = VaultPath,
        TemplatesDir = string.IsNullOrEmpty(TemplatesDir) ? "/.templates" : TemplatesDir,
        ExcludeIndexFiles = ExcludeIndexFiles,
        OverrideDefaultIgnores = OverrideDefaultIgnores,
        DailyNotes = new DailyNotesConfig
        {
            Dir = string.IsNullOrEmpty(DailyNotes.Dir) ? "/daily" : DailyNotes.Dir,
            FilenameFormat = string.IsNullOrEmpty(DailyNotes.FilenameFormat) ? "%Y-%m-%d.md" : DailyNotes.FilenameFormat,
            TemplateName = string.IsNullOrEmpty(DailyNotes.TemplateName) ? "daily" : DailyNotes.TemplateName,
        },
    };
}

using System.Text.Json.Serialization;
using Dendrite.Daemon.Persistence;
using Dendrite.Daemon.Persistence.Repositories;
using Dendrite.Daemon.Services;
using Mediator;

namespace Dendrite.Daemon.Features.Vault;

public sealed class DailyNotesConfig
{
    [JsonPropertyName("dir")] public string Dir { get; init; } = "";
    [JsonPropertyName("filename_format")] public string FilenameFormat { get; init; } = "";
    [JsonPropertyName("template_name")] public string TemplateName { get; init; } = "";
}

public sealed class VaultConfigDto
{
    [JsonPropertyName("vault_name")] public string VaultName { get; init; } = "";
    [JsonPropertyName("vault_path")] public string VaultPath { get; init; } = "";
    [JsonPropertyName("templates_dir")] public string TemplatesDir { get; init; } = "";
    [JsonPropertyName("exclude_index_files")] public string[] ExcludeIndexFiles { get; init; } = [];
    [JsonPropertyName("override_default_ignores")] public bool OverrideDefaultIgnores { get; init; }
    [JsonPropertyName("daily_notes")] public DailyNotesConfig DailyNotes { get; init; } = new();
}

public sealed record InitializeVault(VaultConfigDto Config) : ICommand<Unit>;

public sealed class InitializeVaultHandler(
    VaultStore vaultStore,
    Database db,
    NoteRepository notes,
    IndexRebuilder rebuilder) : ICommandHandler<InitializeVault, Unit>
{
    public ValueTask<Unit> Handle(InitializeVault command, CancellationToken cancellationToken)
    {
        var c = command.Config;
        if (string.IsNullOrEmpty(c.VaultPath))
            throw new InvalidOperationException("vault_path is empty — check your dendrite.nvim config");

        var daily = c.DailyNotes ?? new DailyNotesConfig();
        var vaultPath = c.VaultPath.Length > 1 ? c.VaultPath.TrimEnd('/', '\\') : c.VaultPath;

        vaultStore.Initialize(new VaultConfig(
            c.VaultName,
            vaultPath,
            c.TemplatesDir,
            c.ExcludeIndexFiles ?? [],
            c.OverrideDefaultIgnores,
            new DailyNotes(daily.Dir, daily.FilenameFormat, daily.TemplateName)));

        db.Initialize(Database.IndexPathForVault(c.VaultName));

        if (notes.GetNoteCount() > 0)
        {
            Log.Debug("Vault already initialized, skipping index rebuild");
            return ValueTask.FromResult(Unit.Value);
        }

        rebuilder.Rebuild(vaultPath);
        return ValueTask.FromResult(Unit.Value);
    }
}

using System.Text.Json;
using Dendrite.Application.Abstractions;
using Dendrite.Application.Features.Vault;
using Dendrite.Application.Services;
using Dendrite.Domain.Vault;

namespace Dendrite.Infrastructure.FileSystem;

/// <summary>Keeps a vault's settings in <c>.dendrite/vault.json</c>. The vault's location is never stored.</summary>
internal sealed class VaultManifest : IVaultManifest
{
    private static string FilePath(string vaultPath) => Path.Combine(vaultPath, ".dendrite", "vault.json");

    public VaultConfig? Read(string vaultPath)
    {
        var file = FilePath(vaultPath);
        if (!File.Exists(file))
            return null;

        var dto = JsonSerializer.Deserialize(File.ReadAllText(file), ManifestJsonContext.Default.VaultConfigDto)
            ?? throw new InvalidOperationException($"'{file}' is empty");

        if (string.IsNullOrWhiteSpace(dto.VaultName))
            throw new InvalidOperationException($"'{file}' has no vault_name");

        return VaultInitializer.ToConfig(dto.WithPath(vaultPath).WithDefaults());
    }

    public string? FindRoot(string path)
    {
        var dir = new DirectoryInfo(Path.GetFullPath(path));
        if (!dir.Exists)
            dir = dir.Parent;

        for (; dir is not null; dir = dir.Parent)
        {
            if (File.Exists(FilePath(dir.FullName)))
                return dir.FullName;
        }

        return null;
    }

    public void Write(VaultConfig config)
    {
        var dto = new VaultConfigDto
        {
            VaultName = config.VaultName,
            TemplatesDir = config.TemplatesDir,
            ExcludeIndexFiles = [.. config.ExcludeIndexFiles],
            OverrideDefaultIgnores = config.OverrideDefaultIgnores,
            DailyNotes = new DailyNotesConfig
            {
                Dir = config.DailyNotes.Dir,
                FilenameFormat = config.DailyNotes.FilenameFormat,
                TemplateName = config.DailyNotes.TemplateName,
            },
        };

        var file = FilePath(config.VaultPath);
        Directory.CreateDirectory(Path.GetDirectoryName(file)!);
        File.WriteAllText(file, JsonSerializer.Serialize(dto, ManifestJsonContext.Default.VaultConfigDto));
    }
}

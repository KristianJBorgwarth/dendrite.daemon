using Dendrite.Application.Abstractions;
using Dendrite.Domain.Vault;

namespace Dendrite.Infrastructure.FileSystem;

internal sealed class VaultScaffolder : IVaultScaffolder
{
    public bool IsEmptyOrMissing(string vaultPath) =>
        !Directory.Exists(vaultPath) || !Directory.EnumerateFileSystemEntries(vaultPath).Any();

    public bool LooksLikeVault(VaultConfig config) => Directory.Exists(TemplatesPath(config));

    public void Scaffold(VaultConfig config)
    {
        var templates = TemplatesPath(config);
        Directory.CreateDirectory(templates);

        WriteIfMissing(Path.Combine(templates, config.DailyNotes.TemplateName + ".md"), Template());
        WriteIfMissing(Path.Combine(templates, "scratch.md"), Template());
    }

    private static string TemplatesPath(VaultConfig config) =>
        Path.GetFullPath(Path.Join(config.VaultPath, config.TemplatesDir));

    private static void WriteIfMissing(string path, string content)
    {
        if (!File.Exists(path))
            File.WriteAllText(path, content);
    }

    private static string Template() => "---\ntitle: {{title}}\ndate: {{date}}\n---\n";
}

using Dendrite.Application;
using Dendrite.Application.Abstractions;
using Dendrite.Application.Vaults;

namespace Dendrite.Infrastructure.FileSystem;

internal sealed class VaultFiles : IVaultFiles
{
    private readonly VaultContext _vault;

    public VaultFiles(VaultContext vault)
    {
        _vault = vault;
    }

    public string ReadText(string path) => File.ReadAllText(path);

    public string ReadTemplate(string templateName)
    {
        if (Path.GetExtension(templateName) != ".md")
            templateName += ".md";

        var config = _vault.Config;
        return File.ReadAllText(Path.GetFullPath(Path.Join(config.VaultPath, config.TemplatesDir, templateName)));
    }

    public IReadOnlyList<string> ListNotePaths()
    {
        var config = _vault.Config;
        var ignored = config.EffectiveExcludes;
        var paths = new List<string>();
        Walk(config.VaultPath);
        return paths;

        void Walk(string dir)
        {
            foreach (var entry in new DirectoryInfo(dir).EnumerateFileSystemInfos().OrderBy(e => e.Name, StringComparer.Ordinal))
            {
                if (entry is DirectoryInfo)
                {
                    if (entry.LinkTarget is not null)
                        continue;

                    if (ignored.Contains(entry.Name))
                    {
                        Log.Debug($"Skipping directory during index rebuild: {entry.FullName}");
                        continue;
                    }
                    Walk(entry.FullName);
                }
                else if (entry.Extension == ".md")
                {
                    paths.Add(entry.FullName);
                }
            }
        }
    }
}
